package domain

type Driver interface {
	ExecuteTransaction(url string, f func() error) error

	Execute(statements []string) error
	GetLatestMigration() (int, error)
	GetVersionHashing(version int) (string, error)
	InsertLatestMigration(int, string, string, string) error
	CreateBaseTable() error
	HasBaseTable() (bool, error)

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
