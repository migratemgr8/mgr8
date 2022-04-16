package drivers

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/drivers/postgres"
)

type Driver interface {
	ExecuteTransaction(url string, f func() error) error

	Execute(statements []string) error
	GetLatestMigration() (int, error)
	InsertLatestMigration(int, string, string) error
	CreateBaseTable() error
	HasBaseTable() (bool, error)
}

func GetDriver(driverName string) (Driver, error) {
	switch driverName {
	case "postgres":
		return postgres.NewPostgresDriver(), nil
	}
	return nil, fmt.Errorf("inexistent driver %s", driverName)
}
