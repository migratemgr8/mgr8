package cmd

import (
	"errors"
	"log"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/drivers"
	"github.com/spf13/cobra"
)

var defaultDriverName = string(drivers.Postgres)

type CommandExecutor interface {
	execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver, verbosity applications.LogLevel) error
}

type Command struct {
	driverName    string
	databaseURL   string
	migrationsDir string

	cmd           CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		panic(err)
	}
	silent, err := 	cmd.Flags().GetBool("silent")
	if err != nil {
		panic(err)
	}

	driver, err := drivers.GetDriver(c.driverName)
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

func (c *Command) getLogLevel(verbose, silent bool) (applications.LogLevel, error) {
	if silent && verbose {
		return "", errors.New("flags silent and verbose are mutually exclusive")
	} else if silent {
		return applications.CriticalLogLevel, nil
	} else if verbose {
		return applications.DebugLogLevel, nil
	}
	return applications.InfoLogLevel, nil
}