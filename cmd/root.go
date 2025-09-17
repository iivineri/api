package cmd

import (
	"iivineri/cmd/migration"
	"iivineri/cmd/serve"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "CLI",
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true

	rootCmd.AddCommand(serve.ServeCmd)
	rootCmd.AddCommand(migration.MigrationCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
