package cmd

import (
	"log"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type empty struct{
	emptyCommand applications.EmptyCommand
}

func (c *empty) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	if c.emptyCommand == nil {
		fileService := infrastructure.NewFileService()
		clock := infrastructure.NewClock()
		migrationFileService := applications.NewMigrationFileService(fileService, applications.NewFileNameFormatter(clock), driver)
		c.emptyCommand = applications.NewEmptyCommand(migrationFileService)
	}

	err := c.emptyCommand.Execute(migrationsDir)
	if err != nil {
		log.Print(err)
	}
	return err
}
