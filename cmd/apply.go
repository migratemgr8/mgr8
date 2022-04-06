package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenji-yamane/mgr8/drivers"
)

var defaultDriver = "postgres"

type apply struct {
	Database string
}

func (a *apply) Execute(cmd *cobra.Command, args []string) {
	folderName := args[0]

	driverName := defaultDriver
	if len(args) > 1 {
		driverName = args[1]
	}

	driver, err := drivers.GetDriver(driverName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("Driver %s started\n", driverName)

	err = driver.Begin(a.Database)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	hasTables, err := driver.HasBaseTable()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if !hasTables {
		fmt.Printf("Installing mgr8 into the database...\n")
		err := driver.CreateBaseTable()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	previousMigrationNumber, err := driver.GetLatestMigration()
	if err != nil {
		log.Fatalf("%s", err)
	}

	latestMigrationNumber := 0
	items, _ := ioutil.ReadDir(folderName)
	for _, item := range items {
		itemMigrationNumber, err := a.getMigrationNumber(item.Name())
		if err != nil {
			continue
		}
		if itemMigrationNumber <= previousMigrationNumber {
			continue
		}
		if itemMigrationNumber > latestMigrationNumber {
			latestMigrationNumber = itemMigrationNumber
		}
		a.applyMigrationScript(driver, path.Join(folderName, item.Name()))
	}

	err = driver.UpdateLatestMigration(latestMigrationNumber)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = driver.Commit()
	if err != nil {
		log.Fatalf("could not commit transaction: %s", err)
	}
}

func (a *apply) getMigrationNumber(itemName string) (int, error) {
	itemNameParts := strings.Split(itemName, "_")
	migrationVersionStr := itemNameParts[0]
	migrationVersion, err := strconv.Atoi(migrationVersionStr)
	if err != nil {
		return 0, err
	}
	return migrationVersion, nil
}

func (a *apply) applyMigrationScript(driver drivers.Driver, scriptName string) {
	fmt.Printf("Applying file %s\n", scriptName)
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
