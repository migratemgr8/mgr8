package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type apply struct {
	hashService applications.HashService
}

// TODO: replace file usage with infrastructure.FileService
type MigrationFile struct {
	fullPath string
	name     string
}

type CommandArgs struct {
	migrationType string
	numMigrations int
}

type Migrations struct {
	files    []MigrationFile
	isUpType bool
}

func (a *apply) execute(args []string, databaseURL string, migrationsDir string, driver domain.Driver) error {
	a.hashService = applications.NewHashService(infrastructure.NewFileService())
	dir := migrationsDir
	migrationFiles, err := getMigrationsFiles(dir)
	if err != nil {
		return err
	}

	commandArgs, err := parseArgs(args)
	if err != nil {
		return err
	}

	return driver.ExecuteTransaction(databaseURL, func() error {
		err := applications.CheckAndInstallTool(driver)

		version, err := driver.GetLatestMigrationVersion()
		if err != nil {
			return err
		}

		migrationsToRun, err := getMigrationsToRun(migrationFiles, version, commandArgs.numMigrations, commandArgs.migrationType)
		if err != nil {
			return err
		}

		_, err = a.runMigrations(migrationsToRun, version, driver)
		if err != nil {
			return err
		}

		return err
	})
}

func parseArgs(args []string) (CommandArgs, error) {
	var commandArgs CommandArgs

	if len(args) == 0 {
		return commandArgs, errors.New("arguments missing")
	}

	migrationType := args[0]
	if migrationType != "up" && migrationType != "down" {
		return commandArgs, errors.New("apply's first argument should be either up/down")
	}

	numMigrations := 1
	if len(args) == 2 {
		var err error
		numMigrations, err = strconv.Atoi(args[1])
		if err != nil {
			return commandArgs, err
		}
		if numMigrations == 0 {
			return commandArgs, errors.New("can't run 0 migrations")
		}
	}

	commandArgs.migrationType = migrationType
	commandArgs.numMigrations = numMigrations

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
func getMigrationsToRun(migrationFiles []MigrationFile, currentVersion int, numMigrations int, migrationType string) (Migrations, error) {
	var migrations Migrations

	isUpType := migrationType == "up"
	var files []MigrationFile

	migrationsToBeIncluded := map[int]bool{}

	var firstNum int
	var lastNum int

	// set range of migrations
	if isUpType {
		firstNum = currentVersion + 1
		lastNum = currentVersion + numMigrations
	} else {
		firstNum = currentVersion - numMigrations + 1
		lastNum = currentVersion
	}

	if firstNum <= 0 {
		return migrations, errors.New("migrations would exceed current migration version")
	}

	for i := firstNum; i <= lastNum; i++ {
		migrationsToBeIncluded[i] = true
	}

	for _, file := range migrationFiles {
		migrationNum, err := applications.GetMigrationNumber(file.name)
		if err != nil {
			return migrations, err
		}

		fileMigrationType, err := applications.GetMigrationType(file.name)
		if err != nil {
			return migrations, err
		}

		shouldInclude, ok := migrationsToBeIncluded[migrationNum]

		if migrationType == fileMigrationType && ok && shouldInclude {
			files = append(files, file)
			migrationsToBeIncluded[migrationNum] = false
		}
	}

	// check if all migrations needed were found
	for key, element := range migrationsToBeIncluded {
		if element == true {
			return migrations, errors.New("missing migration number " + strconv.Itoa(key))
		}
	}

	files = sortMigrationFiles(files, isUpType)

	migrations.files = files
	migrations.isUpType = isUpType

	return migrations, nil
}

func (a *apply) runMigrations(migrations Migrations, version int, driver domain.Driver) (int, error) {
	username_service := applications.NewUserNameService()
	username, err := username_service.GetUserName()
	if err != nil {
		return 0, err
	}

	migrationType := "up"
	if !migrations.isUpType {
		migrationType = "down"
	}

	for _, file := range migrations.files {
		migrationNum, err := applications.GetMigrationNumber(file.name)
		if err != nil {
			return 0, err
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")

		hash, err := a.hashService.GetSqlHash(file.fullPath)
		if err != nil {
			return 0, err
		}

		if migrations.isUpType {
			if migrationNum == version+1 {
				err = a.applyMigration(driver, file)
				if err != nil {
					return 0, err
				}

				version = version + 1
				err = driver.InsertIntoMigrationLog(migrationNum, migrationType, username, currentDate)
				if err != nil {
					return 0, err
				}
				err = driver.InsertIntoAppliedMigrations(version, username, currentDate, hash)
				if err != nil {
					return 0, err
				}
			} else {
				valid, err := validateFileMigration(migrationNum, file.fullPath, driver, a.hashService)
				if err != nil {
					return 0, err
				}
				if !valid {
					return 0, fmt.Errorf("❌ invalid migration file %s", file.name)
				}
			}
		} else if !migrations.isUpType && migrationNum == version {
			err = a.applyMigration(driver, file)
			if err != nil {
				return 0, err
			}

			err = driver.InsertIntoMigrationLog(migrationNum, migrationType, username, currentDate)
			if err != nil {
				return 0, err
			}

			err = driver.RemoveAppliedMigration(version)
			if err != nil {
				return 0, err
			}
			version = version - 1
		}
	}

	return version, nil
}

func (a *apply) applyMigration(driver domain.Driver, migration MigrationFile) error {
	fmt.Printf("Applying file %s\n", migration.name)
	content, err := os.ReadFile(migration.fullPath)
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
