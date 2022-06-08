package applications

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kenji-yamane/mgr8/domain"
)

var ErrInvalidMigrationType = errors.New("migration type should be either up/down")
var ErrInvalidMigrationName = errors.New("migration file name in wrong format")

func CheckAndInstallTool(driver domain.Driver) error {
	isToolInstalled, err := driver.IsToolInstalled()
	if err != nil {
		return err
	}

	if !isToolInstalled {
		fmt.Printf("Installing mgr8 into the database...\n")
		return driver.InstallTool()
	}

	return nil
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
			return "", ErrInvalidMigrationType
		}
		return migrationType, nil
	}
	return "", ErrInvalidMigrationName
}
