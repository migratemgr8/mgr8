package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type validate struct{}

func (v *validate) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver, verbosity applications.LogLevel) error {
	dir := args[0]
	return driver.ExecuteTransaction(databaseURL, func() error {
		err := applications.CheckAndInstallTool(driver)

		_, err = validateDirMigrations(dir, driver, applications.NewHashService(infrastructure.NewFileService()))

		return err
	})
}

func validateDirMigrations(dir string, driver domain.Driver, hashService applications.HashService) (int, error) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		fileName := item.Name()
		fullName := path.Join(dir, fileName)

		version, err := applications.GetMigrationNumber(fileName)
		if err != nil {
			return 0, err
		}

		valid, err := validateFileMigration(version, fullName, driver, hashService)
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

func validateFileMigration(version int, filePath string, driver domain.Driver, hashService applications.HashService) (bool, error) {
	hash_file, err := hashService.GetSqlHash(filePath)
	if err != nil {
		return false, err
	}

	hash_db, err := driver.GetVersionHashing(version)
	if err != nil {
		return false, err
	}

	return hash_file == hash_db, nil
}
