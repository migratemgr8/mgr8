package domain

type Driver string

const (
	MySql    Driver = "mysql"
	Postgres Driver = "postgres"

	LogsTableName string = "migration_log"
)
