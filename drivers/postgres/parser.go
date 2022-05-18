package postgres

import (
	"database/sql"
	"errors"
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

func (d *postgresDriver) Deparser() domain.Deparser {
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

func (d *postgresDriver) IsToolInstalled() (bool, error) {
	hasMigrationLogsTable, err := d.HasMigrationLogsTable()
	if err != nil {
		return false, err
	}
	hasAppliedMigrationsTable, err := d.HasAppliedMigrationsTable()
	if err != nil {
		return false, err
	}
	if hasAppliedMigrationsTable != hasMigrationLogsTable { // Has just one of them (XOR)
		return false, errors.New("database in dirty state, tool is partially installed")
	}
	return hasAppliedMigrationsTable && hasMigrationLogsTable, err
}

func (d *postgresDriver) InstallTool() error {
	err := d.CreateMigrationsLogsTable()
	if err != nil {
		return err
	}
	err = d.CreateAppliedMigrationsTable()
	if err != nil {
		return err
	}
	return err
}

func (d *postgresDriver) UninstallTool() error {
	err := d.DropMigrationsLogsTable()
	if err != nil {
		return err
	}
	err = d.DropAppliedMigrationsTable()
	return err
}

func (d *postgresDriver) HasMigrationLogsTable() (bool, error) {
	var hasMigrationLogsTable bool
	err := d.tx.QueryRow(`SELECT EXISTS (
	   SELECT FROM information_schema.tables
	   WHERE  table_name   = $1
	   )`, domain.LogsTableName).Scan(&hasMigrationLogsTable)
	if err != nil {
		return false, err
	}
	return hasMigrationLogsTable, err
}

func (d *postgresDriver) HasAppliedMigrationsTable() (bool, error) {
	var hasAppliedMigrationsTable bool
	err := d.tx.QueryRow(`SELECT EXISTS (
	   SELECT FROM information_schema.tables
	   WHERE  table_name   = $1
	   )`, domain.AppliedTableName).Scan(&hasAppliedMigrationsTable)
	if err != nil {
		return false, err
	}
	return hasAppliedMigrationsTable, err
}

func (d *postgresDriver) CreateMigrationsLogsTable() error {
	_, err := d.tx.Exec(`CREATE TABLE migration_log(
		num INTEGER,
		type VARCHAR(4),
		username VARCHAR(32),
		date VARCHAR(32)
		)`)
	if err != nil {
		return err
	}
	return err
}

func (d *postgresDriver) CreateAppliedMigrationsTable() error {
	_, err := d.tx.Exec(`CREATE TABLE applied_migrations(
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

func (d *postgresDriver) DropMigrationsLogsTable() error {
	_, err := d.tx.Exec(`DROP TABLE migration_log`)
	if err != nil {
		return err
	}
	return err
}

func (d *postgresDriver) DropAppliedMigrationsTable() error {
	_, err := d.tx.Exec(`DROP TABLE applied_migrations`)
	if err != nil {
		return err
	}
	return err
}

func (d *postgresDriver) InsertIntoMigrationLog(migrationNum int, migrationType string, username string, currentDate string) error {
	_, err := d.tx.Exec(`INSERT INTO migration_log (
		num,
		type,
		username,
		date
		) VALUES ($1, $2, $3, $4)`, migrationNum, migrationType, username, currentDate)
	return err
}

func (d *postgresDriver) InsertIntoAppliedMigrations(version int, username string, currentDate string, hash string) error {
	_, err := d.tx.Exec(`INSERT INTO applied_migrations (
		version,
		username,
		date,
		hash
		) VALUES ($1, $2, $3, $4)`, version, username, currentDate, hash)
	return err
}

func (d *postgresDriver) RemoveAppliedMigration(version int) error {
	_, err := d.tx.Exec(`DELETE FROM applied_migrations WHERE version = $1`, version)
	return err
}

func (d *postgresDriver) GetLatestMigrationVersion() (int, error) {
	var version int
	err := d.tx.QueryRow(`SELECT version FROM applied_migrations ORDER BY version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

func (d *postgresDriver) GetVersionHashing(version int) (string, error) {
	var hash string
	err := d.tx.QueryRow(`SELECT hash FROM applied_migrations WHERE version = $1 LIMIT 1`, version).Scan(&hash)
	if err != nil {
		return ``, err
	}
	return hash, nil
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
		Name:    tableName,
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
