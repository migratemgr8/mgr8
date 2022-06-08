package cmd

import (
	"errors"
	"strconv"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type apply struct { }

type CommandArgs struct {
	migrationType string
	numMigrations int
}

func (a *apply) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	commandArgs, err := parseArgs(args)
	if err != nil {
		return err
	}

	applyCommand := applications.NewApplyCommand(driver, applications.NewHashService(infrastructure.NewFileService()))

	return applyCommand.Execute(&applications.ApplyCommandParameters{
		MigrationsDir: migrationsDir,
		DatabaseURL:   databaseURL,
		NumMigrations: commandArgs.numMigrations,
		MigrationType: commandArgs.migrationType,
	})
}

func parseArgs(args []string) (CommandArgs, error) {
	var commandArgs CommandArgs

	if len(args) == 0 {
		return commandArgs, errors.New("arguments missing")
	}

	migrationType := args[0]
	if migrationType != "up" && migrationType != "down" {
		return commandArgs, errors.New("apply's first argument should be either up/down")
	}

	numMigrations := 1
	if len(args) == 2 {
		var err error
		numMigrations, err = strconv.Atoi(args[1])
		if err != nil {
			return commandArgs, err
		}
		if numMigrations == 0 {
			return commandArgs, errors.New("can't run 0 migrations")
		}
	}

	commandArgs.migrationType = migrationType
	commandArgs.numMigrations = numMigrations

	return commandArgs, nil
}
