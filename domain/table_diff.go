package domain

type CreateTableDiff struct {
	table *Table
}

func NewCreateTableDiff(table *Table) *CreateTableDiff {
	return &CreateTableDiff{table: table}
}

func (d *CreateTableDiff) Up(deparser Deparser) string{
	return deparser.CreateTable(d.table)
}

func (d *CreateTableDiff) Down(deparser Deparser) string{
	return deparser.DropTable(d.table.Name)
}

type DropTableDiff struct {
	table *Table
}

func NewDropTableDiff(table *Table) *DropTableDiff {
	return &DropTableDiff{table: table}
}

func (d *DropTableDiff) Up(deparser Deparser) string{
	return deparser.DropTable(d.table.Name)
}

func (d *DropTableDiff) Down(deparser Deparser) string{
	return deparser.CreateTable(d.table)
}

