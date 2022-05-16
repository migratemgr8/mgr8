package postgres

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/domain"
)

type deparser struct { }

func (d *deparser) CreateTable(table *domain.Table) string {
	// TODO: how to mount this string?
	return ""
}

func (d *deparser) DropTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
}

func (d *deparser) AddColumn() string {
	// TODO
	return ""
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
