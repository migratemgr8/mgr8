package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/drivers"
)

type apply struct{}

func (a *apply) execute(args []string, database string, driver drivers.Driver) error {
	folderName := args[0]
	return driver.ExecuteTransaction(database, func() error {
		previousMigrationNumber, err := applications.GetPreviousMigrationNumber(driver)
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

	for _, item := range items {
		fileName := item.Name()
		fullName := path.Join(folderName, fileName)

		itemMigrationNumber, err := applications.GetMigrationNumber(fileName)
		if err != nil {
			continue
		}
		if itemMigrationNumber > latestMigrationNumber {
			latestMigrationNumber = itemMigrationNumber
		}
		if itemMigrationNumber <= previousMigrationNumber {
			valid, err := validateFileMigration(itemMigrationNumber, fullName, driver)
			if err != nil {
				return 0, err
			}
			if !valid {
				return 0, fmt.Errorf("❌ invalid migration file %s", fileName)
			}
			continue
		}
		err = a.applyMigrationScript(driver, fullName)
		if err != nil {
			return 0, err
		}
		currentDate := time.Now().Format("2006-01-02 15:04:05")

		hash, err := applications.GetSqlHash(fullName)
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

func (a *apply) applyMigrationScript(driver drivers.Driver, scriptName string) error {
	fmt.Printf("Applying file %s\n", scriptName)
	content, err := os.ReadFile(scriptName)
	if err != nil {
		return fmt.Errorf("could not read from file: %s", err)
	}

	statements := FilterNonEmpty(strings.Split(string(content), ";"))
	err = driver.Execute(statements)
	if err != nil {
		return fmt.Errorf("could not execute transaction: %s", err)
	}
	return nil
}

func FilterNonEmpty(statements []string) []string {
	filtered := make([]string, 0)
	for _, s := range statements {
		if strings.TrimSpace(s) != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
