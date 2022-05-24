package cmd

import (
	"log"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type generate struct{}

func (g *generate) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	newSchemaPath := args[0]

	// TODO: get this from schemas folder
	oldSchemaPath := args[1]

	fileService := infrastructure.NewFileService()
	clock := infrastructure.NewClock()

	generateCommand := applications.NewGenerateCommand(
		driver,
		applications.NewMigrationFileService(fileService, clock, driver),
	)

	err := generateCommand.Execute(&applications.GenerateParameters{
		OldSchemaPath: oldSchemaPath,
		NewSchemaPath: newSchemaPath,
		MigrationDir:  migrationsDir,
	})
	if err != nil {
		log.Print(err)
	}
	return err
}
