package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type InitCommand struct {
	driverName    string
	databaseURL   string
	migrationsDir string
	cmd           CommandExecutor
}

func (c *InitCommand) Execute(cmd *cobra.Command, args []string) {
	fileService := infrastructure.NewFileService()
	initCommand := applications.NewInitCommand(fileService)

	initialFile := args[0]
	err := initCommand.Execute(initialFile)
	if err != nil {
		log.Fatal(err)
	}
}

