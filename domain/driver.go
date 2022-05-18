package domain

type Driver interface {
	ExecuteTransaction(url string, f func() error) error

	Execute(statements []string) error

	IsToolInstalled() (bool, error)
	InstallTool() error
	UninstallTool() error

	HasMigrationLogsTable() (bool, error)
	CreateMigrationsLogsTable() error
	DropMigrationsLogsTable() error

	HasAppliedMigrationsTable() (bool, error)
	CreateAppliedMigrationsTable() error
	DropAppliedMigrationsTable() error

	InsertIntoMigrationLog(migrationNum int, migrationType string, username string, currentDate string) error
	InsertIntoAppliedMigrations(version int, username string, currentDate string, hash string) error
	RemoveAppliedMigration(version int) error

	GetLatestMigrationVersion() (int, error)
	GetVersionHashing(version int) (string, error)

	ParseMigration(scriptFile string) (*Schema, error)
	Deparser() Deparser
}

type Deparser interface {
	CreateTable(table *Table) string
	DropTable(tableName string) string

	AddColumn() string
	DropColumn(tableName, columnName string) string

	MakeColumnNotNull(tableName, columnName string) string
	UnmakeColumnNotNull(tableName, columnName string) string
}
