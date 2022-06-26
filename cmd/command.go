package cmd

import (
	"log"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/drivers"
	"github.com/migratemgr8/mgr8/global"
	"github.com/spf13/cobra"
)

var defaultDriverName = "postgres"

type CommandExecutor interface {
	execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error
}

type Command struct {
	driverName    string
	databaseURL   string
	migrationsDir string
	cmd           CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	var d global.Database
	d.FromStr(c.driverName)
	driver, err := drivers.GetDriver(d)
	if err != nil {
		log.Fatal(err)
	}

	err = c.cmd.execute(args, c.databaseURL, c.migrationsDir, driver)
	if err != nil {
		log.Fatal(err)
	}
}
