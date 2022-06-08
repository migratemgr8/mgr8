package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type validate struct{}

func (v *validate) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
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

		valid, err := hashService.ValidateFileMigration(version, fullName, driver)
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
