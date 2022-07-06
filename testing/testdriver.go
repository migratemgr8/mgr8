package testing

import (
	"github.com/jmoiron/sqlx"
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/testing/fixtures"
	"log"
)

type TestDriver interface {
	AssertTableExistence(tableName string) (bool, error)
	AssertViewExistence(viewName string) (bool, error)
	AssertTextExistence(tableName string, columnName string) (bool, error)
	AssertVarcharExistence(tableName string, varchar fixtures.VarcharFixture) (bool, error)
	AssertFixtureExistence(fixture *fixtures.Fixture) (bool, error)
	AssertViewFixtureExistence(fixture *fixtures.ViewFixture) (bool, error)
}

type testDriver struct {
	db *sqlx.DB
}

func NewTestDriver(d global.Database) TestDriver {
	if d == global.MySql {
		// TODO
		log.Fatalf("not implemented")
	}
	dm := NewDockerManager()
	url := dm.GetConnectionString(d)
	conn, err := sqlx.Connect(d.String(), url)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &testDriver{
		db: conn,
	}
}

func (d *testDriver) AssertTableExistence(tableName string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`
				SELECT EXISTS (
				    SELECT FROM information_schema.tables WHERE table_name = $1
	   			)`, tableName).Scan(&exists)
	return exists, err
}

func (d *testDriver) AssertViewExistence(viewName string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`
				SELECT EXISTS (
				    SELECT FROM information_schema.views WHERE table_name = $1
	   			)`, viewName).Scan(&exists)
	return exists, err
}

func (d *testDriver) AssertTextExistence(tableName string, columnName string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`
				SELECT EXISTS (
				    SELECT FROM information_schema.columns WHERE table_name = $1
				    AND column_name = $2 AND data_type = 'text'
	   			)`, tableName, columnName).Scan(&exists)
	return exists, err
}

func (d *testDriver) AssertVarcharExistence(tableName string, varchar fixtures.VarcharFixture) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`
				SELECT EXISTS (
				    SELECT FROM information_schema.columns WHERE table_name = $1
				    AND column_name = $2 AND data_type = 'character varying'
				    AND character_maximum_length = $3
	   			)`, tableName, varchar.Name, varchar.Cap).Scan(&exists)
	return exists, err
}

func (d *testDriver) AssertFixtureExistence(fixture *fixtures.Fixture) (bool, error) {
	exists, err := d.AssertTableExistence(fixture.TableName)
	if err != nil || !exists {
		return false, err
	}
	for _, text := range fixture.TextColumns {
		exists, err = d.AssertTextExistence(fixture.TableName, text)
		if err != nil || !exists {
			return false, err
		}
	}
	for _, varchar := range fixture.VarcharColumns {
		exists, err = d.AssertVarcharExistence(fixture.TableName, varchar)
		if err != nil || !exists {
			return false, err
		}
	}
	return true, nil
}

func (d *testDriver) AssertViewFixtureExistence(fixture *fixtures.ViewFixture) (bool, error) {
	exists, err := d.AssertViewExistence(fixture.ViewName)
	if err != nil || !exists {
		return false, err
	}
	for _, text := range fixture.TextColumns {
		exists, err = d.AssertTextExistence(fixture.ViewName, text)
		if err != nil || !exists {
			return false, err
		}
	}
	for _, varchar := range fixture.VarcharColumns {
		exists, err = d.AssertVarcharExistence(fixture.ViewName, varchar)
		if err != nil || !exists {
			return false, err
		}
	}
	return true, nil
}
