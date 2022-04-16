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
	"github.com/kenji-yamane/mgr8/applications"
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

	err := a.execute(folderName, driverName)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *apply) execute(folderName, driverName string) error {
	driver, err := drivers.GetDriver(driverName)
	if err != nil {
		return err
	}

	fmt.Printf("Driver %s started\n", driverName)

	return driver.ExecuteTransaction(a.Database, func() error {
		previousMigrationNumber, err := a.getPreviousMigrationNumber(driver)
		if err != nil {
			return err
		}

		latestMigrationNumber, err := a.runFolderMigrations(folderName, previousMigrationNumber, driver)
		if err != nil {
			return err
		}

		if latestMigrationNumber <= previousMigrationNumber {
			return nil
		}

		return err
	})
}

func (a *apply) runFolderMigrations(folderName string, previousMigrationNumber int, driver drivers.Driver) (int, error) {
	latestMigrationNumber := 0
	items, err := ioutil.ReadDir(folderName)
	if err != nil {
		return 0, err
	}

	username_service := applications.NewUserNameService()
	username, err := username_service.GetUserName()
	if err != nil {
		return 0, err
	}
	fmt.Println("User detected: " + username)

	hash_service := applications.NewHashService()

	for _, item := range items {
		itemMigrationNumber, err := a.getMigrationNumber(item.Name())
		if err != nil {
			continue
		}
		if itemMigrationNumber > latestMigrationNumber {
			latestMigrationNumber = itemMigrationNumber
		}
		if itemMigrationNumber <= previousMigrationNumber {
			continue
		}
		err = a.applyMigrationScript(driver, path.Join(folderName, item.Name()))
		if err != nil {
			return 0, err
		}
		hash, err := hash_service.GetSqlHash(path.Join(folderName, item.Name()))
		if err != nil {
			return 0, err
		}
		err = driver.InsertLatestMigration(latestMigrationNumber, username, hash)
		if err != nil {
			return 0, err
		}
	}
	return latestMigrationNumber, nil
}

func (a *apply) getPreviousMigrationNumber(driver drivers.Driver) (int, error) {
	hasTables, err := driver.HasBaseTable()
	if err != nil {
		return 0, err
	}
	if hasTables {
		return driver.GetLatestMigration()
	}
	fmt.Printf("Installing mgr8 into the database...\n")
	return 0, driver.CreateBaseTable()
}

func (a *apply) applyMigrationScript(driver drivers.Driver, scriptName string) error {
	fmt.Printf("Applying file %s\n", scriptName)
	content, err := os.ReadFile(scriptName)
	if err != nil {
		return fmt.Errorf("could not read from file: %s", err)
	}

	statements := strings.Split(string(content), ";")
	err = driver.Execute(statements)
	if err != nil {
		return fmt.Errorf("could not execute transaction: %s", err)
	}
	return nil
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
