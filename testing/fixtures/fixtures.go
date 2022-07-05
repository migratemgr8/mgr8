package fixtures

import "github.com/migratemgr8/mgr8/domain"

type VarcharFixture struct {
	name string
	cap  int64
}

func (f *VarcharFixture) ToDomainColumn() *domain.Column {
	return &domain.Column{
		Parameters: map[string]interface{}{
			"size": f.cap,
		},
		Datatype: "varchar",
	}
}

type Fixture struct {
	tableName      string
	varcharColumns []VarcharFixture
	textColumns    []string
}

func (f *Fixture) ToDomainTable() *domain.Table {
	t := domain.NewTable(f.tableName, map[string]*domain.Column{})
	for _, varchar := range f.varcharColumns {
		t.Columns[varchar.name] = varchar.ToDomainColumn()
	}
	for _, text := range f.textColumns {
		t.Columns[text] = &domain.Column{Datatype: "text"}
	}
	return t
}

type ViewFixture struct {
	viewName string
	columns  []string
}
