package mysql

import (
	"database/sql"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
	"log"

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
		Name: tableName,
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
