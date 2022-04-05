package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/drivers"
)

var defaultDriver = "postgres"

type apply struct {
	Database string
}

func (a *apply) Execute(cmd *cobra.Command, args []string){
	fileName := args[0]

	driverName := defaultDriver
	if len(args) > 1 {
		driverName = args[1]
	}

	log.Printf("Running migrations from file %s with driver %s", fileName, driverName)
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("could not read from file")
	}

	driver, err := drivers.GetDriver(driverName)
	if err != nil {
		log.Fatal("could not read from file")
	}

	statements := strings.Split(string(content), ";")
	err = driver.Execute(a.Database, statements)
	if err != nil {
		log.Fatal("could not execute transaction")
	}

}

func (a *apply) applyMigrationScript(scriptName string){

}