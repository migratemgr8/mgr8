package drivers

import (
	"fmt"
	"github.com/migratemgr8/mgr8/global"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/drivers/mysql"
	"github.com/migratemgr8/mgr8/drivers/postgres"
)

func GetDriver(d global.Database) (domain.Driver, error) {
	switch d {
	case global.Postgres:
		return postgres.NewPostgresDriver(), nil
	case global.MySql:
		return mysql.NewMySqlDriver(), nil
	}
	return nil, fmt.Errorf("inexistent driver %s", d)
}
