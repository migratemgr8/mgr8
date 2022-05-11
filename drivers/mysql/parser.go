package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kenji-yamane/mgr8/domain"

	"database/sql"
	"log"
)

type mySqlDriver struct {
	tx *sql.Tx
}

func NewMySqlDriver() *mySqlDriver {
	return &mySqlDriver{}
}

func (d *mySqlDriver) Execute(statements []string) error {
	for _, stmt := range statements {
		_, err := d.tx.Exec(stmt)
		if err != nil {
			err2 := d.tx.Rollback()
			if err2 != nil {
				return err2
			}
			return err
		}
	}
	return nil
}

func (d *mySqlDriver) ExecuteTransaction(url string, f func() error) error {
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	d.tx = tx

	err = f()
	if err != nil {
		err2 := d.tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}

	return d.tx.Commit()
}

func (d *mySqlDriver) GetLatestMigration() (int, error) {
	var version int
	err := d.tx.QueryRow(`SELECT version FROM migration_log ORDER BY version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (d *mySqlDriver) GetVersionHashing(version int) (string, error) {
	var hash string
	err := d.tx.QueryRow(`SELECT hash FROM migration_log WHERE version = ? LIMIT 1`, version).Scan(&hash)
	if err != nil {
		return ``, err
	}
	return hash, nil
}

func (d *mySqlDriver) InsertLatestMigration(version int, username string, currentDate string, hash string) error {
	_, err := d.tx.Exec(
		`INSERT INTO migration_log (version, username, date, hash) VALUES (?, ?, ?, ?)`,
		version, username, currentDate, hash)
	return err
}

func (d *mySqlDriver) HasBaseTable() (bool, error) {
	var installed bool
	err := d.tx.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
	    WHERE table_name = ?`, domain.LogsTableName).Scan(&installed)
	if err != nil {
		return false, err
	}
	return installed, err
}

func (d *mySqlDriver) CreateBaseTable() error {
	_, err := d.tx.Exec(`
		CREATE TABLE migration_log(
			version INTEGER,
			username VARCHAR(32),
			date VARCHAR(32),
			hash VARCHAR(32)
	  	)`)
	if err != nil {
		return err
	}
	return err
}

func (d *mySqlDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	return &domain.Schema{}, nil
}
