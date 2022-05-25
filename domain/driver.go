package domain

type Driver interface {
	ExecuteTransaction(url string, f func() error) error

	Execute(statements []string) error
	GetLatestMigration() (int, error)
	GetVersionHashing(version int) (string, error)
	InsertLatestMigration(int, string, string, string) error
	RemoveMigration(int) error
	CreateBaseTable() error
	HasBaseTable() (bool, error)

	ParseMigration(scriptFile string) (*Schema, error)
	Deparser() Deparser
}

type Deparser interface {
	CreateTable(table *Table) string
	DropTable(tableName string) string

	AddColumn(tableName, columnName string, column *Column) string
	DropColumn(tableName, columnName string) string

	MakeColumnNotNull(tableName, columnName string, column *Column) string
	UnmakeColumnNotNull(tableName, columnName string, column *Column) string
}
