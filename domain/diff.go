package domain

import "reflect"

type Diff interface {
	Up(driver Deparser) string
	Down(driver Deparser) string
}

func (s *Schema) Diff(originalSchema *Schema) DiffDeque {
	diffsQueue := NewDiffDeque()

	for tableName, table := range s.Tables {
		originalTable, originalHasTable := originalSchema.Tables[tableName]
		if !originalHasTable {
			diffsQueue.Add(NewCreateTableDiff(table))
		} else {
			diffsQueue.Extend(table.Diff(originalTable))
		}
	}

	for tableName, table := range originalSchema.Tables {
		if _, ok := s.Tables[tableName]; !ok {
			diffsQueue.Add(NewDropTableDiff(table))
		}
	}

	return diffsQueue
}

func (t *Table) Diff(originalTable *Table) DiffDeque {
	diffsQueue := NewDiffDeque()

	for columnName, column := range t.Columns {
		originalColumn, originalHasColumn := originalTable.Columns[columnName]
		if !originalHasColumn {
			diffsQueue.Add(NewCreateColumnDiff(t.Name, columnName, column))
		} else {
			diffsQueue.Extend(column.Diff(t, columnName, originalColumn))
		}
	}

	for columnName, column := range originalTable.Columns {
		if _, ok := t.Columns[columnName]; !ok {
			diffsQueue.Add(NewDropColumnDiff(t.Name, columnName, column))
		}
	}

	return diffsQueue
}

func (c *Column) Diff(table *Table, columnName string, originalColumn *Column) DiffDeque {
	diffsQueue := NewDiffDeque()
	column := table.Columns[columnName]

	if !reflect.DeepEqual(column.Parameters, originalColumn.Parameters) {
		diffsQueue.Add(NewChangeColumnParameterDiff(table.Name, columnName, originalColumn, column))
	}

	if column.IsNotNull != originalColumn.IsNotNull {
		if column.IsNotNull {
			diffsQueue.Add(NewMakeColumnNotNullDiff(table.Name, columnName, column))
		} else {
			diffsQueue.Add(NewMakeColumnNullableDiff(table.Name, columnName, column))
		}
	}

	if column.DefaultValue != originalColumn.DefaultValue {
		diffsQueue.Add(NewSetDefaultValueDiff(table.Name, columnName, column.DefaultValue, originalColumn.DefaultValue))
	}

	return diffsQueue
}
