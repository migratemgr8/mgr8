package applications

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenji-yamane/mgr8/drivers"
)

func GetPreviousMigrationNumber(driver drivers.Driver) (int, error) {
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

func GetMigrationNumber(itemName string) (int, error) {
	itemNameParts := strings.Split(itemName, "_")
	migrationVersionStr := itemNameParts[0]
	migrationVersion, err := strconv.Atoi(migrationVersionStr)
	if err != nil {
		return 0, err
	}
	return migrationVersion, nil
}
