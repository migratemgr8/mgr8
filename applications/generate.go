package applications

import (
	"log"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type GenerateCommand interface {
	Execute(parameters *GenerateParameters) error
}

type GenerateParameters struct {
	OldSchemaPath string
	NewSchemaPath string
	MigrationDir string
}

type generateCommand struct {
	driver   domain.Driver
	fService infrastructure.FileService
	migrationFService MigrationFileService
}

func NewGenerateCommand(driver domain.Driver, fileService infrastructure.FileService, migrationFService MigrationFileService) *generateCommand {
	return &generateCommand{driver: driver, fService: fileService, migrationFService: migrationFService}
}

func (g *generateCommand) Execute(parameters *GenerateParameters) error {
	oldSchema, err := g.getSchemaFromFile(parameters.OldSchemaPath)
	if err != nil {
		return err
	}

	newSchema, err := g.getSchemaFromFile(parameters.NewSchemaPath)
	if err != nil {
		return err
	}

	diffQueue := newSchema.Diff(oldSchema)
	upStatements := diffQueue.GetUpStatements(g.driver.Deparser())
	downStatements := diffQueue.GetDownStatements(g.driver.Deparser())

	nextMigration, err := g.migrationFService.GetNextMigration(parameters.MigrationDir)
	if err != nil {
		return err
	}

	err = g.writeStatementsToFile(parameters.MigrationDir, upStatements, nextMigration, "up")
	if err != nil {
		return err
	}

	err = g.writeStatementsToFile(parameters.MigrationDir, downStatements, nextMigration,"down")
	if err != nil {
		return err
	}

	return nil
}

func (g *generateCommand) getSchemaFromFile(filename string) (*domain.Schema, error) {
	content, err := g.fService.Read(filename)
	if err != nil {
		return nil, err
	}

	return g.driver.ParseMigration(content)
}

func (g *generateCommand) writeStatementsToFile(migrationDir string, statements []string, nextMigration int, migrationType string) error {
	filename := g.migrationFService.FormatFilename(nextMigration, "up")
	log.Printf("Generating file %s migration %s", migrationType, filename)
	content := g.driver.Deparser().WriteScript(statements)
	return g.fService.Write(migrationDir, filename, content)
}
