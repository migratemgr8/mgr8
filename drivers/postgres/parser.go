package postgres

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pg_query "github.com/pganalyze/pg_query_go/v2"

	"fmt"

	"github.com/kenji-yamane/mgr8/domain"
)

type postgresDriver struct {
	tx *sql.Tx
}

func NewPostgresDriver() *postgresDriver {
	return &postgresDriver{}
}

func (d *postgresDriver) Deparser() domain.Deparser{
	return &deparser{}
}

func (d *postgresDriver) Execute(statements []string) error {
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

func (d *postgresDriver) ExecuteTransaction(url string, f func() error) error {
	db, err := sqlx.Connect("postgres", url)
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

func (d *postgresDriver) GetLatestMigration() (int, error) {
	var version int
	err := d.tx.QueryRow(`SELECT version FROM migration_log ORDER BY version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (d *postgresDriver) GetVersionHashing(version int) (string, error) {
	var hash string
	err := d.tx.QueryRow(`SELECT hash FROM migration_log WHERE version = $1 LIMIT 1`, version).Scan(&hash)
	if err != nil {
		return ``, err
	}
	return hash, nil
}

func (d *postgresDriver) InsertLatestMigration(version int, username string, currentDate string, hash string) error {
	_, err := d.tx.Exec(`INSERT INTO migration_log (version, username, date, hash) VALUES ($1, $2, $3, $4)`, version, username, currentDate, hash)
	return err
}

func (d *postgresDriver) HasBaseTable() (bool, error) {
	var installed bool
	err := d.tx.QueryRow(`SELECT EXISTS (
	   SELECT FROM information_schema.tables 
	   WHERE  table_name   = $1
	   )`, domain.LogsTableName).Scan(&installed)
	if err != nil {
		return false, err
	}
	return installed, err
}

func (d *postgresDriver) CreateBaseTable() error {
	_, err := d.tx.Exec(`CREATE TABLE migration_log( version INTEGER, username VARCHAR(32), date VARCHAR(32), hash VARCHAR(32) )`)
	if err != nil {
		return err
	}
	return err
}

func (d *postgresDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	result, err := pg_query.Parse(scriptFile)
	if err != nil {
		return nil, err
	}

	tables := make(map[string]*domain.Table)
	views := make(map[string]*domain.View)
	for _, statement := range result.Stmts {
		switch statement.Stmt.Node.(type) {
		case *pg_query.Node_CreateStmt:
			parsedStatement := statement.Stmt.GetCreateStmt()
			tableName := parsedStatement.Relation.Relname
			tables[tableName] = d.parseTable(tableName, parsedStatement)
		case *pg_query.Node_ViewStmt:
			parsedStatement := statement.Stmt.GetViewStmt()
			viewName := parsedStatement.View.Relname
			views[viewName] = d.parseView(parsedStatement)
		default:
			return nil, fmt.Errorf("found an unsuported statement:\n %s", statement.Stmt.String())
		}
	}

	return &domain.Schema{
		Tables: tables,
		Views:  views,
	}, nil
}

func (d *postgresDriver) parseTable(tableName string, parsedStatement *pg_query.CreateStmt) *domain.Table {
	columns := make(map[string]*domain.Column)
	for _, elts := range parsedStatement.TableElts {
		columnDefinition := elts.GetColumnDef()
		columns[columnDefinition.Colname] = d.parseColumn(columnDefinition)
	}
	return &domain.Table{
		Name: tableName,
		Columns: columns,
	}
}
func (d *postgresDriver) parseView(parsedStatement *pg_query.ViewStmt) *domain.View {
	// TODO
	return &domain.View{
		SQL: "",
	}
}

func (d *postgresDriver) parseColumn(columnDefinition *pg_query.ColumnDef) *domain.Column {
	datatype := columnDefinition.TypeName.Names[1].GetString_().GetStr()
	parameters := make(map[string]interface{})

	if datatype == "varchar" {
		parameters["size"] = columnDefinition.TypeName.Typmods[0].GetAConst().Val.GetInteger().Ival
	}

	return &domain.Column{
		Datatype:   datatype,
		Parameters: parameters,
		IsNotNull:  columnDefinition.IsNotNull,
	}
}
