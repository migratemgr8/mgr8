package cmd

import (
	"fmt"
	"log"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/drivers"
	"github.com/spf13/cobra"
)

var defaultDriver = domain.Postgres

type CommandExecutor interface {
	execute(pathName, database string, driver drivers.Driver) error
}

type Command struct {
	driverName string `default:"postgres"`
	Database   string
	cmd        CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	pathName := args[0]

	c.driverName = string(defaultDriver)
	if len(args) > 1 {
		c.driverName = args[1]
	}

	driver, err := drivers.GetDriver(c.driverName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Driver %s started\n", c.driverName)

	err = c.cmd.execute(pathName, c.Database, driver)
	if err != nil {
		log.Fatal(err)
	}
}
