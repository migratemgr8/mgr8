package postgres

import (
	"database/sql"
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

func (p *postgresDriver) ExecuteTransaction(url string, f func() error) error {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	p.tx = tx

	err = f()
	if err != nil {
		err2 := p.tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}

	return p.tx.Commit()
}

func (p *postgresDriver) GetLatestMigration() (int, error) {
	var version int
	err := p.tx.QueryRow(`SELECT version FROM migration_log ORDER BY version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (p *postgresDriver) InsertLatestMigration(version int, username string, hash string) error {
	_, err := p.tx.Exec(`INSERT INTO migration_log (version, username, date, hash) VALUES ($1, $2, NOW(), $3)`, version, username, hash)
	return err
}

func (p *postgresDriver) HasBaseTable() (bool, error) {
	var installed bool
	err := p.tx.QueryRow(`SELECT EXISTS (
	   SELECT FROM information_schema.tables 
	   WHERE  table_name   = 'migration_log'
	   )`).Scan(&installed)
	if err != nil {
		return false, err
	}
	return installed, err
}

func (p *postgresDriver) CreateBaseTable() error {
	_, err := p.tx.Exec(`CREATE TABLE migration_log( version INTEGER, username VARCHAR(32), date TIMESTAMPTZ, hash VARCHAR(32) )`)
	if err != nil {
		return err
	}
	
	_, err = p.tx.Exec(`INSERT INTO migration_log (version, username, date, hash) VALUES (0, 'InitialSetup', NOW(), '-') `)
	if err != nil {
		return err
	}
	return err
}
