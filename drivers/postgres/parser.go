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

func (p *postgresDriver) InsertLatestMigration(version int, username string, currentDate string, hash string) error {
	_, err := p.tx.Exec(`INSERT INTO migration_log (version, username, date, hash) VALUES ($1, $2, $3, $4)`, version, username, currentDate, hash)
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
	_, err := p.tx.Exec(`CREATE TABLE migration_log( version INTEGER, username VARCHAR(32), date VARCHAR(32), hash VARCHAR(32) )`)
	if err != nil {
		return err
	}
	return err
}

func (p *postgresDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	result, err := pg_query.Parse(scriptFile)
	if err != nil {
		return nil, err
	}

	tables := make(map[string]*domain.Table)
	views := make(map[string]*domain.View)
	for _, statement := range result.Stmts {
		switch statement.Stmt.Node.(type){
		case *pg_query.Node_CreateStmt:
			parsedStatement := statement.Stmt.GetCreateStmt()
			tableName := parsedStatement.Relation.Relname
			tables[tableName] = p.parseTable(parsedStatement)
		case *pg_query.Node_ViewStmt:
			parsedStatement := statement.Stmt.GetViewStmt()
			viewName := parsedStatement.View.Relname
			views[viewName] = p.parseView(parsedStatement)
		default:
			return nil, fmt.Errorf("found an unsuported statement:\n %s", statement.Stmt.String())
		}
	}

	return &domain.Schema{
		Tables: tables,
		Views:  views,
	}, nil
}

func (p *postgresDriver) parseTable(parsedStatement *pg_query.CreateStmt) *domain.Table {
	columns := make(map[string]*domain.Column)
	for _, elts := range parsedStatement.TableElts {
		columnDefinition := elts.GetColumnDef()
		columns[columnDefinition.Colname] = p.parseColumn(columnDefinition)
	}
	return &domain.Table{
		Columns: columns,
	}
}
func (p *postgresDriver) parseView(parsedStatement *pg_query.ViewStmt) *domain.View {
	// TODO
	return &domain.View{
		SQL: "",
	}
}

func (p *postgresDriver) parseColumn(columnDefinition *pg_query.ColumnDef) *domain.Column {
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
