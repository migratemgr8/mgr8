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

func NewDropColumnDiff(tableName string, columnName string, column *Column) *DropColumnDiff {
	return &DropColumnDiff{tableName: tableName, columnName: columnName, column: column}
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

func NewMakeColumnNotNullDiff(tableName string, columnName string, column *Column) *MakeColumnNotNullDiff {
	return &MakeColumnNotNullDiff{tableName: tableName, columnName: columnName, column: column}
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

func NewMakeColumnNullableDiff(tableName string, columnName string, column *Column) *MakeColumnNullableDiff {
	return &MakeColumnNullableDiff{tableName: tableName, columnName: columnName, column: column}
}

func (m *MakeColumnNullableDiff) Up(deparser Deparser) string {
	return deparser.MakeColumnNullable(m.tableName, m.columnName, m.column)
}

func (m *MakeColumnNullableDiff) Down(deparser Deparser) string {
	return deparser.MakeColumnNotNull(m.tableName, m.columnName, m.column)
}

type ChangeColumnDataTypeParameterDiff struct {
	tableName      string
	columnName     string
	originalColumn *Column
	column         *Column
}

func NewChangeColumnParameterDiff(tableName string, columnName string, originalColumn *Column, column *Column) *ChangeColumnDataTypeParameterDiff {
	return &ChangeColumnDataTypeParameterDiff{tableName: tableName, columnName: columnName, originalColumn: originalColumn, column: column}
}

func (m *ChangeColumnDataTypeParameterDiff) Up(deparser Deparser) string {
	return deparser.ChangeDataTypeParameters(m.tableName, m.columnName, m.column)
}

func (m *ChangeColumnDataTypeParameterDiff) Down(deparser Deparser) string {
	return deparser.ChangeDataTypeParameters(m.tableName, m.columnName, m.originalColumn)
}

type SetDefaultValueDiff struct {
	tableName      string
	columnName     string
	newDefaultValue interface{}
	originalDefaultValue         interface{}
}

func NewSetDefaultValueDiff(tableName string, columnName string, newDefaultValue, originalDefaultValue interface{}) *SetDefaultValueDiff {
	return &SetDefaultValueDiff{tableName: tableName, columnName: columnName, newDefaultValue: newDefaultValue, originalDefaultValue: originalDefaultValue}
}

func (m *SetDefaultValueDiff) Up(deparser Deparser) string {
	return deparser.SetColumnDefault(m.tableName, m.columnName, m.newDefaultValue)
}

func (m *SetDefaultValueDiff) Down(deparser Deparser) string {
	return deparser.SetColumnDefault(m.tableName, m.columnName, m.originalDefaultValue)
}
