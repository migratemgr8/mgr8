package applications

type EmptyCommand interface {
	Execute(string) error
}

type emptyCommand struct {
	migrationFService MigrationFileService
}

func NewEmptyCommand(migrationFService MigrationFileService) *emptyCommand {
	return &emptyCommand{migrationFService: migrationFService}
}

func (g *emptyCommand) Execute(migrationDir  string) error {
	nextMigration, err := g.migrationFService.GetNextMigrationNumber(migrationDir)
	if err != nil {
		return err
	}

	err = g.migrationFService.WriteStatementsToFile(migrationDir, []string{}, nextMigration, "up")
	if err != nil {
		return err
	}

	err = g.migrationFService.WriteStatementsToFile(migrationDir, []string{}, nextMigration, "down")
	if err != nil {
		return err
	}

	return nil
}
