package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/drivers"
)

var defaultDriver = "postgres"

type apply struct {
	Database string
}

func (a *apply) Execute(cmd *cobra.Command, args []string){
	folderName := args[0]

	driverName := defaultDriver
	if len(args) > 1 {
		driverName = args[1]
	}

	driver, err := drivers.GetDriver(driverName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Driver %s started", driverName)

	err = driver.Begin(a.Database)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	items, _ := ioutil.ReadDir(folderName)

	for _, item := range items {
		a.applyMigrationScript(driver, path.Join(folderName, item.Name()))
	}

	err = driver.Commit()
	if err != nil {
		log.Fatalf("could not commit transaction: %s", err)
	}
}

func (a *apply) applyMigrationScript(driver drivers.Driver, scriptName string){
	log.Printf("Applying file %s", scriptName)
	content, err := os.ReadFile(scriptName)
	if err != nil {
		log.Fatal("could not read from file")
	}

	statements := strings.Split(string(content), ";")
	err = driver.Execute(statements)
	if err != nil {
		log.Fatalf("could not execute transaction: %s", err)
	}
}