package global

import "sync"

type Database int

const (
	Postgres Database = 0
	MySql    Database = 1
)

var toStr = map[Database]string{
	Postgres: "postgres",
	MySql:    "mysql",
}

var fromStr map[string]Database

var initFromStrOnce sync.Once

func initFromStr() {
	fromStr = make(map[string]Database, len(toStr))
	for k, v := range toStr {
		fromStr[v] = k
	}
}

func (d Database) String() string {
	return toStr[d]
}

func (d *Database) FromStr(s string) {
	initFromStrOnce.Do(initFromStr)
	*d = fromStr[s]
}
