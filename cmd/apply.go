package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kenji-yamane/mgr8/applications"
	"github.com/kenji-yamane/mgr8/domain"
)

type apply struct{}

type MigrationFile struct {
	fullPath string
	name     string
}

type CommandArgs struct {
	migrationFiles []MigrationFile
	migrationType  string
}

type Migrations struct {
	files    []MigrationFile
	isUpType bool
}

func (a *apply) execute(args []string, databaseURL string, driver domain.Driver) error {
	commandArgs, err := parseArgs(args)
	if err != nil {
		return err
	}

	return driver.ExecuteTransaction(databaseURL, func() error {
		previousMigrationNumber, err := applications.GetPreviousMigrationNumber(driver)
		if err != nil {
			return err
		}

		migrationsToRun, err := getMigrationsToRun(commandArgs)
		if err != nil {
			return err
		}

		latestMigrationNumber, err := a.runMigrations(migrationsToRun, previousMigrationNumber, driver)
		if err != nil {
			return err
		}

		if latestMigrationNumber <= previousMigrationNumber {
			return nil
		}

		return err
	})
}

func parseArgs(args []string) (CommandArgs, error) {
	var commandArgs CommandArgs

	if len(args) < 2 {
		return commandArgs, errors.New("arguments missing")
	}

	migrationType := args[0]
	if migrationType != "up" && migrationType != "down" {
		return commandArgs, errors.New("apply's first argument should be either up/down")
	}

	dir := args[1]
	migrationFiles, err := getMigrationsFiles(dir)
	if err != nil {
		return commandArgs, err
	}

	commandArgs.migrationType = migrationType
	commandArgs.migrationFiles = migrationFiles

	return commandArgs, nil
}

// reads directory and returns an array containing full paths of files inside
func getMigrationsFiles(dir string) ([]MigrationFile, error) {
	migrationFiles := []MigrationFile{}

	dirPath, err := filepath.Abs(dir)
	if err != nil {
		return []MigrationFile{}, err
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return []MigrationFile{}, err
	}

	for _, fileInfo := range fileInfos {
		var migrationFile MigrationFile
		migrationFile.name = fileInfo.Name()
		migrationFile.fullPath = filepath.Join(dirPath, fileInfo.Name())
		migrationFiles = append(migrationFiles, migrationFile)
	}

	return migrationFiles, err
}

// returns sorted migration files
// if migration of type up orders ascending, descending otherwise
func sortMigrationFiles(files []MigrationFile, isUpType bool) []MigrationFile {
	if isUpType {
		// sort by ascending
		sort.Slice(files, func(i, j int) bool {
			iNum, _ := applications.GetMigrationNumber(files[i].name)
			jNum, _ := applications.GetMigrationNumber(files[j].name)
			return iNum < jNum
		})
	} else {
		// sort by descending
		sort.Slice(files, func(i, j int) bool {
			iNum, _ := applications.GetMigrationNumber(files[i].name)
			jNum, _ := applications.GetMigrationNumber(files[j].name)
			return iNum >= jNum
		})
	}
	return files
}

// returns migrations files in folder that match type specified (up/down)
func getMigrationsToRun(args CommandArgs) (Migrations, error) {
	var migrations Migrations

	isUpType := args.migrationType == "up"
	var files []MigrationFile

	for _, file := range args.migrationFiles {
		migrationType, err := applications.GetMigrationType(file.name)
		if err != nil {
			return migrations, err
		}

		if migrationType == args.migrationType {
			files = append(files, file)
		}
	}

	migrations.files = sortMigrationFiles(files, isUpType)
	migrations.isUpType = isUpType

	return migrations, nil
}

func (a *apply) runMigrations(migrations Migrations, previousMigrationNumber int, driver domain.Driver) (int, error) {
	version := previousMigrationNumber

	username_service := applications.NewUserNameService()
	username, err := username_service.GetUserName()
	if err != nil {
		return 0, err
	}
	fmt.Println("User detected: " + username)

	for _, file := range migrations.files {
		migrationNum, err := applications.GetMigrationNumber(file.name)
		if err != nil {
			return 0, err
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")

		hash, err := applications.GetSqlHash(file.fullPath)
		if err != nil {
			return 0, err
		}

		if migrations.isUpType {
			if migrationNum == version+1 {
				err = a.applyMigrationScript(driver, file.fullPath)
				if err != nil {
					return 0, err
				}

				version = version + 1
				err = driver.InsertLatestMigration(version, username, currentDate, hash)
				if err != nil {
					return 0, err
				}
			} else {
				valid, err := validateFileMigration(migrationNum, file.fullPath, driver)
				if err != nil {
					return 0, err
				}
				if !valid {
					return 0, fmt.Errorf("❌ invalid migration file %s", file.name)
				}
			}
		} else if !migrations.isUpType && migrationNum == version {
			err = a.applyMigrationScript(driver, file.fullPath)
			if err != nil {
				return 0, err
			}

			err = driver.RemoveMigration(version)
			if err != nil {
				return 0, err
			}
			version = version - 1
		}
	}

	return version, nil
}

func (a *apply) applyMigrationScript(driver domain.Driver, scriptName string) error {
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
