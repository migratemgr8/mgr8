package applications

import (
	"fmt"
	"log"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type MigrationFileService interface {
	GetNextMigrationNumber(dir string) (int, error)
	FormatFilename(migrationNumber int, migrationType string) string
	GetSchemaFromFile(filename string) (*domain.Schema, error)
	WriteStatementsToFile(migrationDir string, statements []string, migrationNumber int, migrationType string) error
}

type migrationFileService struct {
	fileService infrastructure.FileService
	clock       infrastructure.Clock
	driver domain.Driver
}

func NewMigrationFileService(fService infrastructure.FileService, clock infrastructure.Clock, driver domain.Driver) *migrationFileService {
	return &migrationFileService{fileService: fService, clock: clock, driver: driver}
}

func (m *migrationFileService) GetNextMigrationNumber(dir string) (int, error) {
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
	return fmt.Sprintf("%04d_%d.%s.sql", migrationNumber,  m.clock.Now().Unix(), migrationType)
}

func (g *migrationFileService) GetSchemaFromFile(filename string) (*domain.Schema, error) {
	content, err := g.fileService.Read(filename)
	if err != nil {
		return nil, err
	}

	return g.driver.ParseMigration(content)
}

func (g *migrationFileService) WriteStatementsToFile(migrationDir string, statements []string, migrationNumber int, migrationType string) error {
	filename := g.FormatFilename(migrationNumber, "up")
	log.Printf("Generating file %s migration %s", migrationType, filename)
	content := g.driver.Deparser().WriteScript(statements)
	return g.fileService.Write(migrationDir, filename, content)
}
