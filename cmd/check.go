package cmd

import (
	"log"
	"path"

	"github.com/spf13/cobra"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type CheckCommand struct {
	driverName    string
	databaseURL   string
	migrationsDir string
	cmd           CommandExecutor
}

func (c *CheckCommand) Execute(cmd *cobra.Command, args []string) {
	fileService := infrastructure.NewFileService()
	checkCommand := applications.NewCheckCommand(fileService)

	initialFile := args[0]
	referenceFile := path.Join(global.ApplicationFolder, global.ReferenceFile)

	err := checkCommand.Execute(referenceFile, initialFile)
	if err != nil {
		log.Fatal(err)
	}
}
