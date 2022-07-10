package cmd

import (
	"errors"
	"log"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/drivers"
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/infrastructure"

	"github.com/spf13/cobra"
)

var defaultDriverName = "postgres"

type CommandExecutor interface {
	execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver, verbosity infrastructure.LogLevel) error
}

type Command struct {
	driverName    string
	databaseURL   string
	migrationsDir string

	cmd CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	verbose, err := cmd.Flags().GetBool(verboseFlag)
	if err != nil {
		panic(err)
	}
	silent, err := cmd.Flags().GetBool(silentFlag)
	if err != nil {
		panic(err)
	}

	var d global.Database
	d.FromStr(c.driverName)
	driver, err := drivers.GetDriver(d)
	if err != nil {
		log.Fatal(err)
	}

	logLevel, err := c.getLogLevel(verbose, silent)
	if err != nil {
		log.Fatal(err)
	}

	err = c.cmd.execute(args, c.databaseURL, c.migrationsDir, driver, logLevel)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Command) getLogLevel(verbose, silent bool) (infrastructure.LogLevel, error) {
	if silent && verbose {
		return "", errors.New("flags silent and verbose are mutually exclusive")
	}

	if silent {
		return infrastructure.CriticalLogLevel, nil
	}

	if verbose {
		return infrastructure.DebugLogLevel, nil
	}

	return infrastructure.InfoLogLevel, nil
}
