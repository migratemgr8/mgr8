package domain

type CreateColumnDiff struct {
	tableName string
	column *Column
}

func NewCreateColumnDiff(tableName string, column *Column) *CreateColumnDiff {
	return &CreateColumnDiff{tableName: tableName, column: column}
}


type DropColumnDiff struct {
	tableName string
	columnName string
}

func NewDropColumnDiff(tableName string, columnName string) *DropColumnDiff {
	return &DropColumnDiff{tableName: tableName, columnName: columnName}
}

