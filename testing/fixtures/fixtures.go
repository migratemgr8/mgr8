package fixtures

import "github.com/migratemgr8/mgr8/domain"

type VarcharFixture struct {
	Name string
	Cap  int64
}

func (f *VarcharFixture) ToDomainColumn() *domain.Column {
	return &domain.Column{
		Parameters: map[string]interface{}{
			"size": f.Cap,
		},
		Datatype: "varchar",
	}
}

type Fixture struct {
	TableName      string
	VarcharColumns []VarcharFixture
	TextColumns    []string
}

func (f *Fixture) ToDomainTable() *domain.Table {
	t := domain.NewTable(f.TableName, map[string]*domain.Column{})
	for _, varchar := range f.VarcharColumns {
		t.Columns[varchar.Name] = varchar.ToDomainColumn()
	}
	for _, text := range f.TextColumns {
		t.Columns[text] = &domain.Column{Datatype: "text"}
	}
	return t
}

type ViewFixture struct {
	ViewName       string
	TextColumns    []string
	VarcharColumns []VarcharFixture
	Statement      string
}
