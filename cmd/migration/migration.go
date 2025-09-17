package migration

import (
	"fmt"
	"iivineri/internal/wire"
	"strconv"

	"github.com/spf13/cobra"
)

var MigrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Database migration commands",
}

var upCmd = &cobra.Command{
	Use:   "up [steps]",
	Short: "Run pending migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		steps := 0
		if len(args) > 0 {
			steps, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}

		return container.Migration.Up(steps)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [steps]",
	Short: "Rollback migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		steps := 0
		if len(args) > 0 {
			steps, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}

		return container.Migration.Down(steps)
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		return container.Migration.CreateMigration(args[0])
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		return container.Migration.Status()
	},
}

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all database objects",
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		return container.Migration.Drop()
	},
}

var forceCmd = &cobra.Command{
	Use:   "force <version>",
	Short: "Force migration version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container, err := wire.InitializeContainer()
		if err != nil {
			return fmt.Errorf("failed to initialize container: %w", err)
		}
		defer func() {
			if err := container.Shutdown(); err != nil {
				container.Logger.WithError(err).Error("Error during shutdown")
			}
		}()

		version, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return container.Migration.Force(version)
	},
}

func init() {
	MigrationCmd.AddCommand(upCmd)
	MigrationCmd.AddCommand(downCmd)
	MigrationCmd.AddCommand(createCmd)
	MigrationCmd.AddCommand(statusCmd)
	MigrationCmd.AddCommand(dropCmd)
	MigrationCmd.AddCommand(forceCmd)
}