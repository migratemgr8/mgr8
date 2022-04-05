package postgres

import (
	"database/sql"
	"fmt"
	"log"
 _ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type postgresDriver struct {
	tx *sql.Tx
}

func NewPostgresDriver() *postgresDriver{
	return &postgresDriver{}
}

func (p *postgresDriver) Execute(statements []string) error{
	for _, stmt := range statements {
		_, err := p.tx.Exec(stmt)
		if err != nil {
			 err2 := p.tx.Rollback()
			if err2 != nil {
				return err2
			}
			return err
		}
	}
	 return nil
}

func (p *postgresDriver) Begin(url string) error {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	p.tx = tx
	return nil
}

func (p *postgresDriver)  Commit() error {
	if p.tx == nil {
		return fmt.Errorf("no transaction running")
	}
	return p.tx.Commit()
}