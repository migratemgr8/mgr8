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
	verbose bool
	silent bool

	driverName    string
	databaseURL   string
	migrationsDir string

	cmd           CommandExecutor
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	driver, err := drivers.GetDriver(c.driverName)
	if err != nil {
		log.Fatal(err)
	}

	logLevel, err := c.getLogLevel()
	if err != nil {
		log.Fatal(err)
	}

	err = c.cmd.execute(args, c.databaseURL, c.migrationsDir, driver, logLevel)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Command) getLogLevel() (applications.LogLevel, error) {
	if c.silent && c.verbose {
		return "", errors.New("flags silent and verbose are mutually exclusive")
	} else if c.silent {
		return applications.CriticalLogLevel, nil
	} else if c.verbose {
		return applications.DebugLogLevel, nil
	}
	return applications.InfoLogLevel, nil
}