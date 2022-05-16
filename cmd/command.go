package cmd

import (
	"fmt"
	"log"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/drivers"
	"github.com/spf13/cobra"
)

var defaultDriverName = string(drivers.Postgres)

type CommandExecutor interface {
	execute(args []string, databaseURL string, driver domain.Driver) error
}

type Command struct {
	driverName  string
	databaseURL string
	cmd         CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	driver, err := drivers.GetDriver(c.driverName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Driver %s started\n", c.driverName)

	err = c.cmd.execute(args, c.databaseURL, driver)
	if err != nil {
		log.Fatal(err)
	}
}
