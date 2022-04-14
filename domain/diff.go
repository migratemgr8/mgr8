package domain

type SchemaDiff struct { }
type TableDiff struct { }
type ColumnDiff struct { }

func (s *Schema) Diff(originalSchema *Schema) *SchemaDiff{
	return nil
}

func (t *Table) Diff(originalTable *Schema) *TableDiff{
	return nil
}

func (t *Column) Diff(originalColumn *Schema) *ColumnDiff{
	return nil
}