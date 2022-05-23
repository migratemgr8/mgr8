package cmd

import (
	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
)

type generate struct{}

func (g *generate) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	newSchemaPath := args[0]

	// TODO: get this from schemas folder
	oldSchemaPath := args[1]

	// TODO: get these as "(highest migration + 1)_timestamp.up/down.sql
	upMigrationOutputFile := args[2]
	downMigrationOutputFile := args[3]

	fileService := applications.NewFileService()

	generateCommand := applications.NewGenerateCommand(driver, fileService)

	return generateCommand.Execute(&applications.GenerateParameters{
		OldSchemaPath:         oldSchemaPath,
		NewSchemaPath:         newSchemaPath,
		UpMigrationFilename:   upMigrationOutputFile,
		DownMigrationFilename: downMigrationOutputFile,
	})
}
