package cmd

import (
	"fmt"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type diff struct{}

func (g *diff) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver, verbosity infrastructure.LogLevel) error {
	newSchemaPath := args[0]

	fileService := infrastructure.NewFileService()
	clock := infrastructure.NewClock()
	logService, err := infrastructure.NewLogService(verbosity)
	if err != nil {
		return fmt.Errorf("could not start logger, error: %w", err)
	}

	generateCommand := applications.NewGenerateCommand(
		driver,
		applications.NewMigrationFileService(fileService, applications.NewFileNameFormatter(clock), driver, logService),
		fileService,
		logService,
	)

	err = generateCommand.Execute(&applications.GenerateParameters{
		OldSchemaPath: ".mgr8/reference.sql",
		NewSchemaPath: newSchemaPath,
		MigrationDir:  migrationsDir,
	})
	if err != nil {
		logService.Critical("Failed generating diff migration")
	}
	return nil
}
