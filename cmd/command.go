package cmd

import (
	"fmt"
	"log"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/drivers"
	"github.com/spf13/cobra"
)

var defaultDriverName = string(domain.DefaultDriver)

type CommandExecutor interface {
	execute(pathName, databaseURL string, driver drivers.Driver) error
}

type Command struct {
	driverName  string
	databaseURL string
	cmd         CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	pathName := args[0]

	driver, err := drivers.GetDriver(c.driverName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Driver %s started\n", c.driverName)

	err = c.cmd.execute(pathName, c.databaseURL, driver)
	if err != nil {
		log.Fatal(err)
	}
}
