package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/kenji-yamane/mgr8/domain"
)

type mySqlDriver struct {
	tx *sqlx.Tx
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
	return nil
}

func (d *mySqlDriver) GetLatestMigration() (int, error) {
	return 0, nil
}

func (d *mySqlDriver) InsertLatestMigration(version int, username string, currentDate string, hash string) error {
	return nil
}

func (d *mySqlDriver) CreateBaseTable() error {
	return nil
}

func (d *mySqlDriver) HasBaseTable() (bool, error) {
	return false, nil
}

func (d *mySqlDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	return &domain.Schema{}, nil
}
