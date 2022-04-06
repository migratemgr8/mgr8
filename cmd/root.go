package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "mgr8",
		Short: "mgr8 is a framework-agnostic database migrations tool",
		Long: `Lonog version: MGR8 a framework agnostic database migrations tool
                sbrubbles
                sbrubbles`,
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate creates migration script based on the diff between schema versions",
	}

	applyCommand := apply{}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply runs migrations in the selected database",
		Run:   applyCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	applyCmd.Flags().StringVar(&applyCommand.Database, "database", "", "Database URL")
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(applyCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
