package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type postgresDriver struct {
	tx *sql.Tx
}

func NewPostgresDriver() *postgresDriver {
	return &postgresDriver{}
}

func (p *postgresDriver) Execute(statements []string) error {
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

func (p *postgresDriver) Commit() error {
	if p.tx == nil {
		return fmt.Errorf("no transaction running")
	}
	return p.tx.Commit()
}

func (p *postgresDriver) GetLatestMigration() (int, error) {
	var version int
	err := p.tx.QueryRow(`SELECT version FROM migration_version`).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (p *postgresDriver) UpdateLatestMigration(version int) error {
	_, err := p.tx.Exec(`UPDATE migration_version SET version = $1`, version)
	return err
}

func (p *postgresDriver) HasBaseTable() (bool, error) {
	var installed bool
	err := p.tx.QueryRow(`SELECT EXISTS (
	   SELECT FROM information_schema.tables 
	   WHERE  table_name   = 'migration_version'
	   )`).Scan(&installed)
	if err != nil {
		return false, err
	}
	return installed, err
}

func (p *postgresDriver) CreateBaseTable() error {
	_, err := p.tx.Exec(`CREATE TABLE migration_version( version INTEGER )`)
	if err != nil {
		return err
	}
	_, err = p.tx.Exec(`INSERT INTO migration_version (version) VALUES (0) `)
	if err != nil {
		return err
	}
	return err
}
