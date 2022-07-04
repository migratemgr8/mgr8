package cmd

import (
	"log"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type empty struct {
	emptyCommand applications.EmptyCommand
}

func (c *empty) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver, verbosity infrastructure.LogLevel) error {
	if c.emptyCommand == nil {
		fileService := infrastructure.NewFileService()
		clock := infrastructure.NewClock()
		logService, err := infrastructure.NewLogService(verbosity)
		if err != nil {
			log.Fatal(err)
		}

		migrationFileService := applications.NewMigrationFileService(fileService, applications.NewFileNameFormatter(clock), driver, logService)
		c.emptyCommand = applications.NewEmptyCommand(migrationFileService)
	}

	err := c.emptyCommand.Execute(migrationsDir)
	if err != nil {
		log.Print(err)
	}
	return err
}
