package applications

import (
	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type GenerateCommand interface {
	Execute(parameters *GenerateParameters) error
}

type GenerateParameters struct {
	OldSchemaPath string
	NewSchemaPath string
	MigrationDir  string
}

type generateCommand struct {
	driver            domain.Driver
	fService          infrastructure.FileService
	migrationFService MigrationFileService
}

func NewGenerateCommand(driver domain.Driver, migrationFService MigrationFileService) *generateCommand {
	return &generateCommand{driver: driver, migrationFService: migrationFService}
}

func (g *generateCommand) Execute(parameters *GenerateParameters) error {
	oldSchema, err := g.migrationFService.GetSchemaFromFile(parameters.OldSchemaPath)
	if err != nil {
		return err
	}

	newSchema, err := g.migrationFService.GetSchemaFromFile(parameters.NewSchemaPath)
	if err != nil {
		return err
	}

	diffQueue := newSchema.Diff(oldSchema)
	upStatements := diffQueue.GetUpStatements(g.driver.Deparser())
	downStatements := diffQueue.GetDownStatements(g.driver.Deparser())

	nextMigration, err := g.migrationFService.GetNextMigrationNumber(parameters.MigrationDir)
	if err != nil {
		return err
	}

	err = g.migrationFService.WriteStatementsToFile(parameters.MigrationDir, upStatements, nextMigration, "up")
	if err != nil {
		return err
	}

	err = g.migrationFService.WriteStatementsToFile(parameters.MigrationDir, downStatements, nextMigration, "down")
	if err != nil {
		return err
	}

	return nil
}