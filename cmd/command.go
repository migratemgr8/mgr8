package cmd

import (
	"fmt"
	"log"

	"github.com/kenji-yamane/mgr8/drivers"
	"github.com/spf13/cobra"
)

var defaultDriver = "postgres"

type DatabaseCommand interface {
	execute(pathName, database string, driver drivers.Driver) error
}

type DatabaseCmd struct {
	driverName string `default:"postgres"`
	Database   string
	cmd        DatabaseCommand
}

func (dcmd *DatabaseCmd) Execute(cmd *cobra.Command, args []string) {
	pathName := args[0]

	dcmd.driverName = defaultDriver
	if len(args) > 1 {
		dcmd.driverName = args[1]
	}

	driver, err := drivers.GetDriver(dcmd.driverName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Driver %s started\n", dcmd.driverName)

	err = dcmd.cmd.execute(pathName, dcmd.Database, driver)
	if err != nil {
		log.Fatal(err)
	}
}
