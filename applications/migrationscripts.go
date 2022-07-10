package applications

import (
	"fmt"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"
)

type MigrationFileService interface {
	GetNextMigrationNumber(dir string) (int, error)
	GetSchemaFromFile(filename string) (*domain.Schema, error)
	WriteStatementsToFile(migrationDir string, statements []string, migrationNumber int, migrationType string) error
}

type migrationFileService struct {
	fileService       infrastructure.FileService
	driver            domain.Driver
	fileNameFormatter FileNameFormatter
	logService        infrastructure.LogService
}

func NewMigrationFileService(fService infrastructure.FileService, fileNameFormatter FileNameFormatter, driver domain.Driver, logService infrastructure.LogService) *migrationFileService {
	return &migrationFileService{
		fileService:       fService,
		driver:            driver,
		fileNameFormatter: fileNameFormatter,
		logService:        logService,
	}
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

func (g *migrationFileService) GetSchemaFromFile(filename string) (*domain.Schema, error) {
	content, err := g.fileService.Read(filename)
	if err != nil {
		g.logService.Critical("Could not read from", filename)
		return nil, err
	}

	return g.driver.ParseMigration(content)
}

func (g *migrationFileService) WriteStatementsToFile(migrationDir string, statements []string, migrationNumber int, migrationType string) error {
	filename := g.fileNameFormatter.FormatFilename(migrationNumber, migrationType)
	g.logService.Info("Generating file for", migrationType, "migration:", filename)
	content := g.driver.Deparser().WriteScript(statements)
	return g.fileService.Write(migrationDir, filename, content)
}

type FileNameFormatter interface {
	FormatFilename(int, string) string
}

type fileNameFormatter struct {
	clock infrastructure.Clock
}

func NewFileNameFormatter(clock infrastructure.Clock) *fileNameFormatter {
	return &fileNameFormatter{clock: clock}
}

func (m *fileNameFormatter) FormatFilename(migrationNumber int, migrationType string) string {
	return fmt.Sprintf("%04d_%d.%s.sql", migrationNumber, m.clock.Now().Unix(), migrationType)
}
