package domain

type CreateColumnDiff struct {
	tableName  string
	columnName string
	column     *Column
}

func NewCreateColumnDiff(tableName, columnName string, column *Column) *CreateColumnDiff {
	return &CreateColumnDiff{tableName: tableName, columnName: columnName, column: column}
}

func (d *CreateColumnDiff) Up(deparser Deparser) string {
	return deparser.AddColumn(d.tableName, d.columnName, d.column)
}

func (d *CreateColumnDiff) Down(deparser Deparser) string {
	return deparser.DropColumn(d.tableName, d.columnName)
}

type DropColumnDiff struct {
	tableName  string
	columnName string
	column     *Column
}

func NewDropColumnDiff(tableName string, columnName string) *DropColumnDiff {
	return &DropColumnDiff{tableName: tableName, columnName: columnName}
}

func (d *DropColumnDiff) Up(deparser Deparser) string {
	return deparser.DropColumn(d.tableName, d.columnName)
}

func (d *DropColumnDiff) Down(deparser Deparser) string {
	return deparser.AddColumn(d.tableName, d.columnName, d.column)
}

type MakeColumnNotNullDiff struct {
	tableName  string
	columnName string
	column     *Column
}

func NewMakeColumnNotNullDiff(tableName string, columnName string) *MakeColumnNotNullDiff {
	return &MakeColumnNotNullDiff{tableName: tableName, columnName: columnName}
}

func (m *MakeColumnNotNullDiff) Up(deparser Deparser) string {
	return deparser.MakeColumnNotNull(m.tableName, m.columnName, m.column)
}

func (m *MakeColumnNotNullDiff) Down(deparser Deparser) string {
	return deparser.MakeColumnNullable(m.tableName, m.columnName, m.column)
}

type MakeColumnNullableDiff struct {
	tableName  string
	columnName string
	column     *Column
}

func NewMakeColumnNullableDiff(tableName string, columnName string) *MakeColumnNullableDiff {
	return &MakeColumnNullableDiff{tableName: tableName, columnName: columnName}
}

func (m *MakeColumnNullableDiff) Up(deparser Deparser) string {
	return deparser.MakeColumnNullable(m.tableName, m.columnName, m.column)
}

func (m *MakeColumnNullableDiff) Down(deparser Deparser) string {
	return deparser.MakeColumnNotNull(m.tableName, m.columnName, m.column)
}
