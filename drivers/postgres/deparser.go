package postgres

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/domain"
)

type deparser struct{}

func (d *deparser) CreateTable(table *domain.Table) string {
	// TODO: how to mount this string?
	return ""
}

func (d *deparser) DropTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
}

func (d *deparser) AddColumn(tableName, columnName string, column *domain.Column) string {
	columnDatatype := column.Datatype
	if size, ok := column.Parameters["size"]; ok {
		columnDatatype = fmt.Sprintf("%s(%d)", column.Datatype, size)
	}
	if column.IsNotNull {
		columnDatatype = fmt.Sprintf("%s NOT NULL", columnDatatype)
	}
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %v", tableName, columnName, columnDatatype)
}

func (d *deparser) DropColumn(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName)
}
func (d *deparser) MakeColumnNotNull(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL", tableName, columnName)
}

func (d *deparser) UnmakeColumnNotNull(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL", tableName, columnName)
}
