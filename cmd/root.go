package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultMigrationDir = "migrations"

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "mgr8",
		Short: "mgr8 is a framework-agnostic database migrations tool",
		Long:  `Long version: mgr8 is an agnostic tool that abstracts database migration operations`,
	}

	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose")
	rootCmd.PersistentFlags().Bool("silent", false, "Silent")

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate creates migration scripts",
	}

	diffCommand := Command{cmd: &diff{}}
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "diff creates migration script based on the diff between schema versions",
		Run:   diffCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	diffCmd.Flags().StringVar(&diffCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	diffCmd.Flags().StringVar(&diffCommand.driverName, "driver", defaultDriverName, "Driver Name")
	diffCmd.Flags().StringVar(&diffCommand.migrationsDir, "dir", defaultMigrationDir, "Migrations Directory")

	emptyCommand := Command{cmd: &empty{}}
	emptyCmd := &cobra.Command{
		Use:   "empty",
		Short: "empty creates empty migration",
		Run:   emptyCommand.Execute,
		Args:  cobra.NoArgs,
	}
	emptyCmd.Flags().StringVar(&emptyCommand.migrationsDir, "dir", defaultMigrationDir, "Migrations Directory")
	emptyCmd.Flags().StringVar(&emptyCommand.driverName, "driver", defaultDriverName, "Driver Name")

	initCommand := &InitCommand{}
	initCmd := &cobra.Command{
		Use:   "init file",
		Short: "init sets the schema as reference",
		Run:   initCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}

	checkCommand := &CheckCommand{}
	checkCmd := &cobra.Command{
		Use:   "check file",
		Short: "check returns 0 if files match",
		Run:   checkCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}

	generateCmd.AddCommand(emptyCmd, diffCmd, initCmd, checkCmd)

	applyCommand := Command{cmd: &apply{}}
	applyCmd := &cobra.Command{
		Use:   "apply n",
		Short: "apply runs migrations in the selected database",
		Run:   applyCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	applyCmd.Flags().StringVar(&applyCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	applyCmd.Flags().StringVar(&applyCommand.driverName, "driver", defaultDriverName, "Driver Name")
	applyCmd.Flags().StringVar(&applyCommand.migrationsDir, "dir", "", "Migrations Directory")

	validateCommand := Command{cmd: &validate{}}
	validateCmd := &cobra.Command{
		Use:   "validate file",
		Short: "validate compares migrations sql scripts against hashing ",
		Run:   validateCommand.Execute,
		Args:  cobra.MinimumNArgs(1),
	}
	validateCmd.Flags().StringVar(&validateCommand.databaseURL, "database", os.Getenv("DB_HOST"), "Database URL")
	validateCmd.Flags().StringVar(&validateCommand.driverName, "driver", defaultDriverName, "Driver Name")
	validateCmd.Flags().StringVar(&validateCommand.migrationsDir, "dir", "", "Migrations Directory")

	rootCmd.AddCommand(applyCmd, generateCmd, validateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
