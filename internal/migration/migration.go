package migration

import (
	"database/sql"
	"fmt"
	"iivineri/internal/config"
	"iivineri/internal/logger"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationInterface interface {
	Up(steps int) error
	Down(steps int) error
	Drop() error
	Force(version int) error
	Status() error
	CreateMigration(name string) error
	Close() error
}

type Migration struct {
	config *config.DatabaseConfig
	logger *logger.Logger
	m      *migrate.Migrate
	db     *sql.DB
}

func NewMigration(
	cfg *config.DatabaseConfig,
	log *logger.Logger,
) MigrationInterface {
	return &Migration{
		config: cfg,
		logger: log,
	}
}

func (m *Migration) connect() error {
	if m.m != nil {
		return nil
	}

	m.logger.Info("Connecting to database for migrations...")
	
	db, err := sql.Open("postgres", m.config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	m.db = db
	m.m = migrator
	m.logger.Info("Migration connection established successfully")
	return nil
}

func (m *Migration) CreateMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join("migrations", filename+".up.sql")
	downFile := filepath.Join("migrations", filename+".down.sql")

	if err := createMigrationFile(upFile, "-- Migration up"); err != nil {
		return fmt.Errorf("error creating up migration file: %w", err)
	}

	if err := createMigrationFile(downFile, "-- Migration down"); err != nil {
		return fmt.Errorf("error creating down migration file: %w", err)
	}

	m.logger.Infof("Created migrations:\n\t%s\n\t%s", upFile, downFile)
	return nil
}

func (m *Migration) Up(steps int) error {
	if err := m.connect(); err != nil {
		return err
	}

	if steps > 0 {
		if err := m.m.Steps(steps); err != nil {
			if err.Error() == "no change" {
				m.logger.Info("No migrations to run")
				return nil
			}
			return fmt.Errorf("failed to run %d migration steps: %w", steps, err)
		}
		m.logger.Infof("Successfully ran %d migration steps", steps)
	} else {
		if err := m.m.Up(); err != nil {
			if err.Error() == "no change" {
				m.logger.Info("No migrations to run")
				return nil
			}
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		m.logger.Info("Migrations completed successfully")
	}

	return nil
}

func (m *Migration) Down(steps int) error {
	if err := m.connect(); err != nil {
		return err
	}

	if steps > 0 {
		if err := m.m.Steps(-steps); err != nil {
			if err.Error() == "no change" {
				m.logger.Info("No migrations to rollback")
				return nil
			}
			return fmt.Errorf("failed to rollback %d migration steps: %w", steps, err)
		}
		m.logger.Infof("Successfully rolled back %d migration steps", steps)
	} else {
		if err := m.m.Down(); err != nil {
			if err.Error() == "no change" {
				m.logger.Info("No migrations to rollback")
				return nil
			}
			return fmt.Errorf("failed to rollback all migrations: %w", err)
		}
		m.logger.Info("All migrations rolled back successfully")
	}

	return nil
}

func (m *Migration) Drop() error {
	if err := m.connect(); err != nil {
		return err
	}

	if err := m.m.Drop(); err != nil {
		return fmt.Errorf("failed to drop database schema: %w", err)
	}

	m.logger.Info("Database schema dropped successfully")
	return nil
}

func (m *Migration) Force(version int) error {
	if err := m.connect(); err != nil {
		return err
	}

	if err := m.m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version to %d: %w", version, err)
	}

	m.logger.Infof("Forced migration version to %d", version)
	return nil
}

func (m *Migration) Status() error {
	if err := m.connect(); err != nil {
		return err
	}

	version, dirty, err := m.m.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	m.logger.Infof("Current migration version: %d", version)
	if dirty {
		m.logger.Warn("Status: DIRTY (migration failed)")
	} else {
		m.logger.Info("Status: CLEAN")
	}

	return nil
}

func (m *Migration) Close() error {
	if m.db != nil {
		err := m.db.Close()
		m.logger.Info("Migration database connection closed")
		return err
	}
	return nil
}