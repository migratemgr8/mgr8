package drivers

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/drivers/mysql"
	"github.com/kenji-yamane/mgr8/drivers/postgres"
)

type Driver interface {
	ExecuteTransaction(url string, f func() error) error

	Execute(statements []string) error
	GetLatestMigration() (int, error)
	GetVersionHashing(version int) (string, error)
	InsertLatestMigration(int, string, string, string) error
	CreateBaseTable() error
	HasBaseTable() (bool, error)

	ParseMigration(scriptFile string) (*domain.Schema, error)
}

func GetDriver(driverName string) (Driver, error) {
	driver := domain.Driver(driverName)
	switch driver {
	case domain.Postgres:
		return postgres.NewPostgresDriver(), nil
	case domain.MySql:
		return mysql.NewMySqlDriver(), nil
	}
	return nil, fmt.Errorf("inexistent driver %s", driverName)
}
