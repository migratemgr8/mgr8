package domain

type Diff interface{
	Up(driver Deparser) string
	Down(driver Deparser) string
}

func (s *Schema) Diff(originalSchema *Schema) []Diff {
	diffsQueue := []Diff{}

	for tableName, table := range s.Tables {
		originalTable, originalHasTable := originalSchema.Tables[tableName]
		if !originalHasTable {
			diffsQueue = append(diffsQueue, NewCreateTableDiff(table))
		} else {
			diffsQueue = append(diffsQueue, table.Diff(originalTable)...)
		}
	}

	for tableName, table := range originalSchema.Tables {
		if _, ok := s.Tables[tableName]; !ok {
			diffsQueue = append(diffsQueue, NewDropTableDiff(table))
		}
	}

	return diffsQueue
}

func (t *Table) Diff(originalTable *Table) []Diff {
	diffsQueue := []Diff{}

	for columnName, column := range t.Columns {
		originalColumn, originalHasColumn := originalTable.Columns[columnName]
		if !originalHasColumn {
			diffsQueue = append(diffsQueue, NewCreateColumnDiff(t.Name, column))
		} else {
			diffsQueue = append(diffsQueue, column.Diff(t, columnName, originalColumn)...)
		}
	}

	for columnName := range originalTable.Columns {
		if _, ok := t.Columns[columnName]; !ok {
			diffsQueue = append(diffsQueue, NewDropColumnDiff(t.Name, columnName))
		}
	}

	return diffsQueue
}

func (c *Column) Diff(table *Table, columnName string, originalColumn *Column) []Diff {
	diffsQueue := []Diff{}
	column := table.Columns[columnName]
	if column.IsNotNull != originalColumn.IsNotNull {
		if column.IsNotNull {
			diffsQueue = append(diffsQueue, NewMakeColumnNotNullDiff(table.Name, columnName))
		} else {
		 	diffsQueue = append(diffsQueue, NewUnmakeColumnNotNullDiff(table.Name, columnName))
		}
	}
	return diffsQueue
}
