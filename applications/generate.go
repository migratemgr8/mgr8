package applications

import (
	"path"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/infrastructure"
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
	logService        infrastructure.LogService
}

func NewGenerateCommand(driver domain.Driver, migrationFService MigrationFileService, fService infrastructure.FileService, logService infrastructure.LogService) *generateCommand {
	return &generateCommand{driver: driver, migrationFService: migrationFService, fService: fService, logService: logService}
}

func (g *generateCommand) Execute(parameters *GenerateParameters) error {
	newSchemaContent, err := g.fService.Read(parameters.NewSchemaPath)
	if err != nil {
		g.logService.Critical("Could not read from", parameters.NewSchemaPath)
		return err
	}

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
	g.logService.Debug("Latest migration found:", nextMigration)

	err = g.migrationFService.WriteStatementsToFile(parameters.MigrationDir, upStatements, nextMigration, "up")
	if err != nil {
		return err
	}

	err = g.migrationFService.WriteStatementsToFile(parameters.MigrationDir, downStatements, nextMigration, "down")
	if err != nil {
		return err
	}

	g.logService.Debug("Updating reference file at", path.Join(global.ApplicationFolder, global.ReferenceFile))
	return g.fService.Write(global.ApplicationFolder, global.ReferenceFile, newSchemaContent)
}
