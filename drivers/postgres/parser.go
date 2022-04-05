package postgres

import (
	"log"
 _ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type postgresDriver struct {

}

func NewPostgresDriver() *postgresDriver{
	return &postgresDriver{}
}

func (p *postgresDriver) Execute(statements []string) error{
	db, err := sqlx.Connect("postgres", "user=root dbname=core password=root sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, stmt := range statements {
		_, err := tx.Exec(stmt)
		if err != nil {
			 err2 := tx.Rollback()
			if err2 != nil {
				return err2
			}
			return err
		}
	}

	 return tx.Commit()
}