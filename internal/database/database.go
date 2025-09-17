package database

import (
	"context"
	"fmt"
	"iivineri/internal/config"
	"iivineri/internal/logger"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type DatabaseInterface interface {
	Connect(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) error
	Close()
}

type Database struct {
	config *config.DatabaseConfig
	logger *logger.Logger
	pool   *pgxpool.Pool
	once   sync.Once
	err    error
}

type pgxLogger struct {
	logger   *logger.Logger
	dbLogLevel string
}

func (l *pgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	// Check if we should log based on DB_LOG_LEVEL
	shouldLog := l.shouldLogLevel(level)
	if !shouldLog {
		return
	}

	sql, sqlExists := data["sql"]
	args, argsExists := data["args"]
	duration, durationExists := data["time"]

	// Create base log entry
	logEntry := l.logger.WithField("database", "postgres")
	
	if sqlExists {
		logEntry = logEntry.WithField("sql", sql)
	}
	if argsExists && args != nil {
		logEntry = logEntry.WithField("args", args)
	}
	if durationExists {
		logEntry = logEntry.WithField("duration", duration)
	}

	// Add other fields
	for key, value := range data {
		if key != "sql" && key != "args" && key != "time" {
			logEntry = logEntry.WithField(key, value)
		}
	}

	// Log based on the original level but with database context
	if sqlExists {
		// For SQL queries, add a clear prefix
		logEntry.Info("SQL Query: " + msg)
	} else {
		// For non-SQL messages, use original level
		switch level {
		case tracelog.LogLevelTrace:
			logEntry.Trace(msg)
		case tracelog.LogLevelDebug:
			logEntry.Debug(msg)
		case tracelog.LogLevelInfo:
			logEntry.Info(msg)
		case tracelog.LogLevelWarn:
			logEntry.Warn(msg)
		case tracelog.LogLevelError:
			logEntry.Error(msg)
		}
	}
}

func (l *pgxLogger) shouldLogLevel(level tracelog.LogLevel) bool {
	dbLogLevel := strings.ToLower(l.dbLogLevel)
	
	switch dbLogLevel {
	case "trace":
		return true // Log everything
	case "debug":
		return level >= tracelog.LogLevelDebug
	case "info":
		return level >= tracelog.LogLevelInfo
	case "warn":
		return level >= tracelog.LogLevelWarn
	case "error":
		return level >= tracelog.LogLevelError
	case "none", "off":
		return false // Log nothing
	default:
		return level >= tracelog.LogLevelInfo // Default to info and above
	}
}

func NewDatabase(cfg *config.DatabaseConfig, logger *logger.Logger) DatabaseInterface {
	return &Database{
		config: cfg,
		logger: logger,
	}
}

func (db *Database) Connect(ctx context.Context) error {
	db.once.Do(func() {
		db.err = db.connect(ctx)
	})
	return db.err
}

func (db *Database) connect(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(db.config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(db.config.DBMaxConns())
	poolConfig.MinConns = int32(db.config.DBMinConns())
	poolConfig.MinIdleConns = int32(db.config.DBMinIdleConns())
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute * 5
	poolConfig.ConnConfig.ConnectTimeout = time.Second * 30
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"statement_timeout":                   "30000",
		"idle_in_transaction_session_timeout": "60000",
	}

	// Configure SQL logging based on DB_LOG_LEVEL
	dbLogLevel := strings.ToLower(db.config.LogLevel())
	var traceLogLevel tracelog.LogLevel
	
	switch dbLogLevel {
	case "trace":
		traceLogLevel = tracelog.LogLevelTrace
	case "debug":
		traceLogLevel = tracelog.LogLevelDebug
	case "info":
		traceLogLevel = tracelog.LogLevelInfo
	case "warn":
		traceLogLevel = tracelog.LogLevelWarn
	case "error":
		traceLogLevel = tracelog.LogLevelError
	case "none", "off":
		traceLogLevel = tracelog.LogLevelNone
	default:
		traceLogLevel = tracelog.LogLevelInfo // Default to info level
	}

	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgxLogger{
			logger:     db.logger,
			dbLogLevel: db.config.LogLevel(),
		},
		LogLevel: traceLogLevel,
	}

	db.logger.Trace("Creating database connection pool")
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}

	db.pool = pool
	db.logger.Trace("Database connection pool created successfully")
	return nil
}

func (db *Database) HealthCheck(ctx context.Context) error {
	if db.pool == nil {
		return fmt.Errorf("database not connected")
	}

	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if db.pool == nil {
		if err := db.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return db.pool.Query(ctx, sql, args...)
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if db.pool == nil {
		if err := db.Connect(ctx); err != nil {
			return &errorRow{err: err}
		}
	}
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *Database) Exec(ctx context.Context, sql string, args ...any) error {
	if db.pool == nil {
		if err := db.Connect(ctx); err != nil {
			return err
		}
	}
	_, err := db.pool.Exec(ctx, sql, args...)
	return err
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.logger.Trace("Database connection pool closed")
	}
}

type errorRow struct {
	err error
}

func (er *errorRow) Scan(dest ...any) error {
	return er.err
}
