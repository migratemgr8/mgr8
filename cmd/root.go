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

	generateCommand := Command{cmd: &generate{}}
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate creates migration script based on the diff between schema versions",
		Run:   generateCommand.Execute,
	}

	applyCommand := Command{cmd: &apply{}}
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply runs migrations in the selected database",
		Run:   applyCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	applyCmd.Flags().StringVar(&applyCommand.Database, "database", os.Getenv("DB_HOST"), "Database URL")

	validateCommand := Command{cmd: &validate{}}
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "validate compares migrations sql scripts against hashing ",
		Run:   validateCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	validateCmd.Flags().StringVar(&validateCommand.Database, "database", os.Getenv("DB_HOST"), "Database URL")

	rootCmd.AddCommand(applyCmd, generateCmd, validateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
