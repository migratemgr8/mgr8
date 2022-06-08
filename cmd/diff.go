package cmd

import (
	"log"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type diff struct{}

func (g *diff) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	newSchemaPath := args[0]

	fileService := infrastructure.NewFileService()
	clock := infrastructure.NewClock()

	generateCommand := applications.NewGenerateCommand(
		driver,
		applications.NewMigrationFileService(fileService, applications.NewFileNameFormatter(clock), driver),
		fileService,
	)

	err := generateCommand.Execute(&applications.GenerateParameters{
		OldSchemaPath: ".mgr8/reference.sql",
		NewSchemaPath: newSchemaPath,
		MigrationDir:  migrationsDir,
	})
	if err != nil {
		log.Print(err)
	}
	return err
}
