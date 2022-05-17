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
		Args:  cobra.MinimumNArgs(1),
	}
	generateCmd.Flags().StringVar(&generateCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	generateCmd.Flags().StringVar(&generateCommand.driverName, "driver", defaultDriverName, "Driver Name")

	applyCommand := Command{cmd: &apply{}}
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply runs migrations in the selected database",
		Run:   applyCommand.Execute,
		Args:  cobra.MinimumNArgs(2),
	}
	applyCmd.Flags().StringVar(&applyCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	applyCmd.Flags().StringVar(&applyCommand.driverName, "driver", defaultDriverName, "Driver Name")

	validateCommand := Command{cmd: &validate{}}
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "validate compares migrations sql scripts against hashing ",
		Run:   validateCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	validateCmd.Flags().StringVar(&validateCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	validateCmd.Flags().StringVar(&validateCommand.driverName, "driver", defaultDriverName, "Driver Name")

	rootCmd.AddCommand(applyCmd, generateCmd, validateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
