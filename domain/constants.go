package domain

type Driver string

const (
	MySql         Driver = "mysql"
	Postgres      Driver = "postgres"
	DefaultDriver Driver = Postgres

	LogsTableName string = "migration_log"
)
