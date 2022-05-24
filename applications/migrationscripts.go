package applications

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/infrastructure"
)

type MigrationFileService interface {
	GetNextMigration(dir string) (int, error)
	FormatFilename(migrationNumber int, migrationType string) string
}

type migrationFileService struct {
	fileService infrastructure.FileService
	clock       infrastructure.Clock
}

func NewMigrationFileService(fService infrastructure.FileService, clock infrastructure.Clock) *migrationFileService {
	return &migrationFileService{fileService: fService, clock: clock}
}

func (m *migrationFileService) GetNextMigration(dir string) (int, error) {
	migrationFiles, err := m.fileService.List(dir)
	if err != nil {
		return 0, err
	}
	maxMigration := 0
	for _, file := range migrationFiles {
		mNumber, _ := GetMigrationNumber(file.Name)
		if mNumber > maxMigration {
			maxMigration = mNumber
		}
	}
	return maxMigration + 1, nil
}

func (m *migrationFileService) FormatFilename(migrationNumber int, migrationType string) string {
	return fmt.Sprintf("%04d_%d.%s.sql", migrationNumber,  m.clock.Now().Nanosecond(), migrationType)
}