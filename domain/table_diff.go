package domain

type CreateTableDiff struct {
	table *Table
}

func NewCreateTableDiff(table *Table) *CreateTableDiff {
	return &CreateTableDiff{table: table}
}

type DropTableDiff struct {
	tableName string
}

func NewDropTableDiff(tableName string) *DropTableDiff {
	return &DropTableDiff{tableName: tableName}
}


