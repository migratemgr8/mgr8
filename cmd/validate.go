package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/drivers"
)

type validate struct{}

func (v *validate) execute(folderName, database string, driver drivers.Driver) error {
	return driver.ExecuteTransaction(database, func() error {
		previousMigrationNumber, err := applications.GetPreviousMigrationNumber(driver)
		if err != nil {
			return err
		}

		_, err = validateFolderMigrations(folderName, previousMigrationNumber, driver)

		return err
	})
}

func validateFolderMigrations(folderName string, previousMigrationNumber int, driver drivers.Driver) (int, error) {
	items, err := ioutil.ReadDir(folderName)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		fileName := item.Name()
		fullName := path.Join(folderName, fileName)

		version, err := applications.GetMigrationNumber(fileName)
		if err != nil {
			return 0, err
		}

		valid, err := validateFileMigration(version, fullName, driver)
		if err != nil {
			return 0, err
		}

		if valid {
			fmt.Printf("✅  %s\n", item.Name())
		} else {
			fmt.Printf("❌  %s\n", item.Name())
		}
	}

	return 0, nil
}

func validateFileMigration(version int, fileName string, driver drivers.Driver) (bool, error) {
	hash_file, err := applications.GetSqlHash(fileName)
	if err != nil {
		return false, err
	}

	hash_db, err := driver.GetVersionHashing(version)
	if err != nil {
		return false, err
	}

	return hash_file == hash_db, nil
}
