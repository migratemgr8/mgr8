package drivers

import (
	"fmt"

	"github.com/kenji-yamane/mgr8/drivers/postgres"
)

type Driver interface {
	Execute(db string, statements []string) error
}

func GetDriver(driverName string) (Driver, error) {
	switch driverName {
	case "postgres":
		return postgres.NewPostgresDriver(), nil
	}
	return nil, fmt.Errorf("inexistent driver %s", driverName)
}