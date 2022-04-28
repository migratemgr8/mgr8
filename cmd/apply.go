package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/drivers"
)

type apply struct{}

func (a *apply) execute(folderName, database string, driver drivers.Driver) error {
	return driver.ExecuteTransaction(database, func() error {
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
		currentDate := time.Now().Format("2006-01-02 15:04:05")

		hash, err := hash_service.GetSqlHash(path.Join(folderName, item.Name()))
		if err != nil {
			return 0, err
		}
		err = driver.InsertLatestMigration(latestMigrationNumber, username, currentDate, hash)
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

	statements := a.filterNonEmpty(strings.Split(string(content), ";"))
	err = driver.Execute(statements)
	if err != nil {
		return fmt.Errorf("could not execute transaction: %s", err)
	}
	return nil
}

func (a *apply) filterNonEmpty(statements []string) []string {
	filtered := make([]string, 0)
	for _, s := range statements {
		if strings.TrimSpace(s) != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
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
