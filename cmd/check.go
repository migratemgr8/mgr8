package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/infrastructure"
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
	err := checkCommand.Execute(initialFile)
	if err != nil {
		log.Fatal(err)
	}
}

