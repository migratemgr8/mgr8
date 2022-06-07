package cmd

import (
	"log"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type empty struct{ }

func (c *empty) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	fileService := infrastructure.NewFileService()
	clock := infrastructure.NewClock()

	emptyCommand := applications.NewEmptyCommand(
		applications.NewMigrationFileService(fileService, applications.NewFileNameFormatter(clock), driver),
	)

	err := emptyCommand.Execute(migrationsDir)
	if err != nil {
		log.Print(err)
	}
	return err
}
