package applications

import "github.com/kenji-yamane/mgr8/domain"

type GenerateCommand interface {
	Execute(parameters *GenerateParameters) error
}

type GenerateParameters struct {
	OldSchemaPath string
	NewSchemaPath string

	UpMigrationFilename   string
	DownMigrationFilename string
}

type generateCommand struct {
	driver   domain.Driver
	fService FileService
}

func NewGenerateCommand(driver domain.Driver, fileService FileService) *generateCommand {
	return &generateCommand{driver: driver, fService: fileService}
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

	err = g.writeStatementsToFile(parameters.UpMigrationFilename, upStatements)
	if err != nil {
		return err
	}

	err = g.writeStatementsToFile(parameters.DownMigrationFilename, downStatements)
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

func (g *generateCommand) writeStatementsToFile(filename string, statements []string) error {
	content := g.driver.Deparser().WriteScript(statements)
	return g.fService.Write(filename, content)
}
