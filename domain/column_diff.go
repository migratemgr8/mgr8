package domain

type CreateColumnDiff struct {
	tableName string
	columnName string
	column *Column
}

func NewCreateColumnDiff(tableName string, column *Column) *CreateColumnDiff {
	return &CreateColumnDiff{tableName: tableName, column: column}
}

func (d *CreateColumnDiff) Up(deparser Deparser) string{
	return deparser.AddColumn()
}

func (d *CreateColumnDiff) Down(deparser Deparser) string{
	return deparser.DropColumn(d.tableName, d.columnName)
}

type DropColumnDiff struct {
	tableName string
	columnName string
}

func NewDropColumnDiff(tableName string, columnName string) *DropColumnDiff {
	return &DropColumnDiff{tableName: tableName, columnName: columnName}
}

func (d *DropColumnDiff) Up(deparser Deparser) string{
	return deparser.DropColumn(d.tableName, d.columnName)
}

func (d *DropColumnDiff) Down(deparser Deparser) string{
	return deparser.AddColumn()
}