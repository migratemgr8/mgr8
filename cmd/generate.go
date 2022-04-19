package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/drivers"
)

type generate struct {
	Database string
}

func (g *generate) Execute(cmd *cobra.Command, args []string) {
	fileName := args[0]
	driverName := defaultDriver
	if len(args) > 1 {
		driverName = args[1]
	}

	err := g.execute(fileName, driverName)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *generate) execute(fileName, driverName string) error {
	driver, err := drivers.GetDriver(driverName)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not read from file: %s", err)
	}

	_, err = driver.ParseMigration(string(content))
	if err != nil {
		return err
	}

	return nil
}
