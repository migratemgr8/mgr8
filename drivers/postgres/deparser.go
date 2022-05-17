package postgres

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/domain"
)

type deparser struct { }

func inStringList(stringList []string, needle string) bool {
	isIn := false
	for _, s := range stringList {
		if needle == s {
			isIn = true
		}
	}
	return isIn
}

func hasSingleArg(datatype string) bool {
	singleArgTypes := []string{"char", "varchar", "bit", "varbit", "time", "timestamp"}
	if inStringList(singleArgTypes, datatype) {
		return true
	} else {
		return false
	}
}

func hasDoubleArg(datatype string) bool {
	doubleArgTypes := []string{"decimal"}
	if inStringList(doubleArgTypes, datatype) {
		return true
	}	else {
		return false
	}
}

func (d *deparser) CreateTable(table *domain.Table) string {
	string := fmt.Sprintf("CREATE TABLE %s (", table.Name)

	for columnName, column := range table.Columns {
		string = string + fmt.Sprintf("%s %s", columnName, column.Datatype)

		if hasSingleArg(column.Datatype) {
			string = string + fmt.Sprintf("(%s)", column.Datatype, column.Parameters["size"])
		} else if hasDoubleArg(column.Datatype) {
			string = string + fmt.Sprintf("(%s,%s)", column.Datatype, column.Parameters["precision"], column.Parameters["scale"])
		}

		if column.IsNotNull {
			string = string + fmt.Sprintf(" NOT NULL")
		}

		string = string + fmt.Sprintf(",")
	}

	string = string[0:len(string) - 2]
	string = string + fmt.Sprintf(")")
	return string
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
