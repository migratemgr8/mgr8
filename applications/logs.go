package applications

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kenji-yamane/mgr8/domain"
)

func CheckAndInstallTool(driver domain.Driver) error {
	hasTables, err := driver.HasBaseTable()
	if err != nil {
		return err
	}
	if !hasTables {
		fmt.Printf("Installing mgr8 into the database...\n")
		return driver.CreateBaseTable()
	}
	return nil
}

func GetPreviousMigrationNumber(driver domain.Driver) (int, error) {
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

func GetMigrationType(fileName string) (string, error) {
	re := regexp.MustCompile(`\.(.*?)\.`) // gets string in between dots
	match := re.FindStringSubmatch(fileName)
	if len(match) > 1 {
		migrationType := match[1]
		if migrationType != "up" && migrationType != "down" {
			return "", errors.New("migration type should be either up/down")
		}
		return migrationType, nil
	}
	return "", errors.New("migration file name in wrong format")
}
