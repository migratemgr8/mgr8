package applications

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

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type ApplyCommand interface {

}

type applyCommand struct {
	driver domain.Driver
	hashService HashService
}

func NewApplyCommand(driver domain.Driver, hashService HashService) *applyCommand {
	return &applyCommand{driver: driver, hashService: hashService}
}

type Migrations struct {
	files    []infrastructure.MigrationFile
	isUpType bool
}

type ApplyCommandParameters struct {
	MigrationsDir string
	DatabaseURL string
	NumMigrations int
	MigrationType string
}

func (a *applyCommand) Execute(parameters *ApplyCommandParameters) error {
	migrationFiles, err := a.getMigrationsFiles(parameters.MigrationsDir)
	if err != nil {
		return err
	}

	return a.driver.ExecuteTransaction(parameters.DatabaseURL, func() error {
		err := CheckAndInstallTool(a.driver)

		version, err := a.driver.GetLatestMigrationVersion()
		if err != nil {
			return err
		}

		migrationsToRun, err := a.getMigrationsToRun(migrationFiles, version, parameters.NumMigrations, parameters.MigrationType)
		if err != nil {
			return err
		}

		_, err = a.runMigrations(migrationsToRun, version, a.driver)
		if err != nil {
			return err
		}

		return err
	})
}

// reads directory and returns an array containing full paths of files inside
func (a *applyCommand) getMigrationsFiles(dir string) ([]infrastructure.MigrationFile, error) {
	migrationFiles := []infrastructure.MigrationFile{}

	dirPath, err := filepath.Abs(dir)
	if err != nil {
		return []infrastructure.MigrationFile{}, err
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return []infrastructure.MigrationFile{}, err
	}

	for _, fileInfo := range fileInfos {
		var migrationFile infrastructure.MigrationFile
		migrationFile.Name = fileInfo.Name()
		migrationFile.FullPath = filepath.Join(dirPath, fileInfo.Name())
		migrationFiles = append(migrationFiles, migrationFile)
	}

	return migrationFiles, err
}

// returns sorted migration files
// if migration of type up orders ascending, descending otherwise
func (a *applyCommand) sortMigrationFiles(files []infrastructure.MigrationFile, isUpType bool) []infrastructure.MigrationFile {
	if isUpType {
		// sort by ascending
		sort.Slice(files, func(i, j int) bool {
			iNum, _ := GetMigrationNumber(files[i].Name)
			jNum, _ := GetMigrationNumber(files[j].Name)
			return iNum < jNum
		})
	} else {
		// sort by descending
		sort.Slice(files, func(i, j int) bool {
			iNum, _ := GetMigrationNumber(files[i].Name)
			jNum, _ := GetMigrationNumber(files[j].Name)
			return iNum >= jNum
		})
	}
	return files
}

// returns migrations files in folder that match type specified (up/down)
func (a *applyCommand) getMigrationsToRun(migrationFiles []infrastructure.MigrationFile, currentVersion int, numMigrations int, migrationType string) (Migrations, error) {
	var migrations Migrations

	isUpType := migrationType == "up"
	var files []infrastructure.MigrationFile

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
		migrationNum, err := GetMigrationNumber(file.Name)
		if err != nil {
			return migrations, err
		}

		fileMigrationType, err := GetMigrationType(file.Name)
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

	files = a.sortMigrationFiles(files, isUpType)

	migrations.files = files
	migrations.isUpType = isUpType

	return migrations, nil
}

func (a *applyCommand) runMigrations(migrations Migrations, version int, driver domain.Driver) (int, error) {
	usernameService := NewUserNameService()
	username, err := usernameService.GetUserName()
	if err != nil {
		return 0, err
	}

	migrationType := "up"
	if !migrations.isUpType {
		migrationType = "down"
	}

	for _, file := range migrations.files {
		migrationNum, err := GetMigrationNumber(file.Name)
		if err != nil {
			return 0, err
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")

		hash, err := a.hashService.GetSqlHash(file.FullPath)
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
				valid, err := a.hashService.ValidateFileMigration(migrationNum, file.FullPath, driver)
				if err != nil {
					return 0, err
				}
				if !valid {
					return 0, fmt.Errorf("❌ invalid migration file %s", file.Name)
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

func (a *applyCommand) applyMigration(driver domain.Driver, migration infrastructure.MigrationFile) error {
	fmt.Printf("Applying file %s\n", migration.Name)
	content, err := os.ReadFile(migration.FullPath)
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
