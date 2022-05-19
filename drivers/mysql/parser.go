package mysql

import (
	"database/sql"
	"errors"
	"log"

	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"github.com/kenji-yamane/mgr8/domain"
)

type mySqlDriver struct {
	tx *sql.Tx
}

func NewMySqlDriver() *mySqlDriver {
	return &mySqlDriver{}
}

func (d *mySqlDriver) Deparser() domain.Deparser {
	return &deparser{}
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

func (d *mySqlDriver) IsToolInstalled() (bool, error) {
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

func (d *mySqlDriver) InstallTool() error {
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

func (d *mySqlDriver) UninstallTool() error {
	err := d.DropMigrationsLogsTable()
	if err != nil {
		return err
	}
	err = d.DropAppliedMigrationsTable()
	return err
}

func (d *mySqlDriver) HasMigrationLogsTable() (bool, error) {
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

func (d *mySqlDriver) HasAppliedMigrationsTable() (bool, error) {
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

func (d *mySqlDriver) CreateMigrationsLogsTable() error {
	_, err := d.tx.Exec(`CREATE TABLE migration_log (
		migration_number INTEGER,
		type VARCHAR(4),
		username VARCHAR(32),
		date VARCHAR(32)
		)`)
	if err != nil {
		return err
	}
	return err
}

func (d *mySqlDriver) CreateAppliedMigrationsTable() error {
	_, err := d.tx.Exec(`CREATE TABLE applied_migrations (
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

func (d *mySqlDriver) DropMigrationsLogsTable() error {
	_, err := d.tx.Exec(`DROP TABLE migration_log`)
	if err != nil {
		return err
	}
	return err
}

func (d *mySqlDriver) DropAppliedMigrationsTable() error {
	_, err := d.tx.Exec(`DROP TABLE applied_migrations`)
	if err != nil {
		return err
	}
	return err
}

func (d *mySqlDriver) InsertIntoMigrationLog(migrationNum int, migrationType string, username string, currentDate string) error {
	_, err := d.tx.Exec(`INSERT INTO migration_log (
		migration_number,
		type,
		username,
		date
		) VALUES ($1, $2, $3, $4)`, migrationNum, migrationType, username, currentDate)
	return err
}

func (d *mySqlDriver) InsertIntoAppliedMigrations(version int, username string, currentDate string, hash string) error {
	_, err := d.tx.Exec(`INSERT INTO applied_migrations (
		version,
		username,
		date,
		hash
		) VALUES ($1, $2, $3, $4)`, version, username, currentDate, hash)
	return err
}

func (d *mySqlDriver) RemoveAppliedMigration(version int) error {
	_, err := d.tx.Exec(`DELETE FROM applied_migrations WHERE version = $1`, version)
	return err
}

func (d *mySqlDriver) GetLatestMigrationVersion() (int, error) {
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

func (d *mySqlDriver) GetVersionHashing(version int) (string, error) {
	var hash string
	err := d.tx.QueryRow(`SELECT hash FROM applied_migrations WHERE version = $1 LIMIT 1`, version).Scan(&hash)
	if err != nil {
		return ``, err
	}
	return hash, nil
}

type Visitor interface {
	Enter(n ast.Node) (node ast.Node, skipChildren bool)
	Leave(n ast.Node) (node ast.Node, ok bool)
}

type extractor struct {
	tables map[string]*domain.Table
	views  map[string]*domain.View
}

func (x *extractor) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.CreateTableStmt:
		createStmt := in.(*ast.CreateTableStmt)
		x.tables[createStmt.Table.Name.O] = x.parseTable(createStmt.Table.Name.O, createStmt)
	case *ast.CreateViewStmt:
		createStmt := in.(*ast.CreateViewStmt)
		x.views[createStmt.ViewName.Name.O] = x.parseView(createStmt)
	}
	return in, false
}

func (x *extractor) parseTable(tableName string, stmt *ast.CreateTableStmt) *domain.Table {
	columns := make(map[string]*domain.Column)

	for _, c := range stmt.Cols {
		columns[c.Name.Name.O] = x.parseColumn(c)
	}
	return &domain.Table{
		Name:    tableName,
		Columns: columns,
	}
}

func (x *extractor) parseView(stmt *ast.CreateViewStmt) *domain.View {
	// TODO
	return &domain.View{
		SQL: "",
	}
}

func (x *extractor) parseColumn(col *ast.ColumnDef) *domain.Column {
	dt := col.Tp.Tp
	parameters := make(map[string]interface{})

	if dt == mysql.TypeVarchar {
		parameters["size"] = col.Tp.Flen
	}

	isNotNull := false
	for _, opt := range col.Options {
		if opt.Tp == ast.ColumnOptionNotNull {
			isNotNull = true
		}
	}

	return &domain.Column{
		Datatype:   types.TypeStr(col.Tp.Tp),
		Parameters: parameters,
		IsNotNull:  isNotNull,
	}
}

func (x *extractor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func (d *mySqlDriver) ParseMigration(scriptFile string) (*domain.Schema, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(scriptFile, "", "")
	if err != nil {
		return nil, err
	}

	e := &extractor{
		tables: make(map[string]*domain.Table),
		views:  make(map[string]*domain.View),
	}
	for _, n := range stmtNodes {
		n.Accept(e)
	}
	return &domain.Schema{
		Tables: e.tables,
		Views:  e.views,
	}, nil
}
