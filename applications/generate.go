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
	deparser := g.driver.Deparser()

	var upStatements []string
	for _, diff := range diffQueue {
		upStatements = append(upStatements, diff.Up(deparser))
	}

	var downStatements []string
	for i := len(diffQueue) - 1; i >= 0; i-- {
		downStatements = append(downStatements, diffQueue[i].Down(deparser))
	}

	nextMigration, err := g.migrationFService.GetNextMigration(parameters.MigrationDir)
	if err != nil {
		return err
	}

	upMigrationFilename := g.migrationFService.FormatFilename(nextMigration, "up")
	downMigrationFilename := g.migrationFService.FormatFilename(nextMigration, "down")

	log.Printf("Generating files %s and %s", upMigrationFilename, downMigrationFilename)

	err = g.writeStatementsToFile(parameters.MigrationDir, upMigrationFilename, upStatements)
	if err != nil {
		return err
	}

	err = g.writeStatementsToFile(parameters.MigrationDir, downMigrationFilename, downStatements)
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

func (g *generateCommand) writeStatementsToFile(migrationDir, filename string, statements []string) error {
	content := g.driver.Deparser().WriteScript(statements)
	return g.fService.Write(migrationDir, filename, content)
}
