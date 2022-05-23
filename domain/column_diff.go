package domain

type CreateColumnDiff struct {
	tableName  string
	columnName string
	column     *Column
}

func NewCreateColumnDiff(tableName, columnName string, column *Column) *CreateColumnDiff {
	return &CreateColumnDiff{tableName: tableName, column: column, columnName: columnName}
}

func (d *CreateColumnDiff) Up(deparser Deparser) string {
	return deparser.AddColumn()
}

func (d *CreateColumnDiff) Down(deparser Deparser) string {
	return deparser.DropColumn(d.tableName, d.columnName)
}

type DropColumnDiff struct {
	tableName  string
	columnName string
}

func NewDropColumnDiff(tableName string, columnName string) *DropColumnDiff {
	return &DropColumnDiff{tableName: tableName, columnName: columnName}
}

func (d *DropColumnDiff) Up(deparser Deparser) string {
	return deparser.DropColumn(d.tableName, d.columnName)
}

func (d *DropColumnDiff) Down(deparser Deparser) string {
	return deparser.AddColumn()
}

type MakeColumnNotNullDiff struct {
	tableName  string
	columnName string
}

func NewMakeColumnNotNullDiff(tableName string, columnName string) *MakeColumnNotNullDiff {
	return &MakeColumnNotNullDiff{tableName: tableName, columnName: columnName}
}

func (m *MakeColumnNotNullDiff) Up(deparser Deparser) string {
	return deparser.MakeColumnNotNull(m.tableName, m.columnName)
}

func (m *MakeColumnNotNullDiff) Down(deparser Deparser) string {
	return deparser.UnmakeColumnNotNull(m.tableName, m.columnName)
}

type UnmakeColumnNotNullDiff struct {
	tableName  string
	columnName string
}

func NewUnmakeColumnNotNullDiff(tableName string, columnName string) *UnmakeColumnNotNullDiff {
	return &UnmakeColumnNotNullDiff{tableName: tableName, columnName: columnName}
}

func (m *UnmakeColumnNotNullDiff) Up(deparser Deparser) string {
	return deparser.UnmakeColumnNotNull(m.tableName, m.columnName)
}

func (m *UnmakeColumnNotNullDiff) Down(deparser Deparser) string {
	return deparser.MakeColumnNotNull(m.tableName, m.columnName)
}
