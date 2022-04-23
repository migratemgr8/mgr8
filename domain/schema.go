package domain

type Schema struct {
	Tables map[string]*Table
	Views  map[string]*View
}

type Table struct {
	Columns map[string]*Column
}

type Column struct {
	Datatype   string
	Parameters map[string]interface{}
	IsNotNull  bool
}

type View struct {
	SQL string
}
