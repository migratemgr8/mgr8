package domain

type Schema struct {
	Tables map[string]*Table
	Views  map[string]*View
}

type Table struct {
	Name    string
	Columns map[string]*Column
}

func NewTable(name string, columns map[string]*Column) *Table {
	return &Table{Name: name, Columns: columns}
}

type Column struct {
	Datatype     string
	Parameters   map[string]interface{}
	IsNotNull    bool
	DefaultValue interface{}
}

type View struct {
	SQL string
}
