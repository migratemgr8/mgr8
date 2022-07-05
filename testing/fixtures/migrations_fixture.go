package fixtures

import (
	"fmt"
	"os"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/infrastructure"

	. "github.com/onsi/gomega"
)

type MigrationsFixture interface {
	AddRawFile(filename, content string)
	AddMigration0001() *Fixture
	AddMigration0002() (VarcharFixture, *ViewFixture)
	AddMigration0003() VarcharFixture
	TearDown()
}

type migrationsFixture struct {
	folderPath  string
	fileService infrastructure.FileService
	usersTable  *Fixture
	usersView   *ViewFixture
	deparser    domain.Deparser
}

func NewMigrationsFixture(folderPath string, fileService infrastructure.FileService, deparser domain.Deparser) MigrationsFixture {
	err := fileService.CreateFolderIfNotExists(folderPath)
	Expect(err).To(BeNil())
	return &migrationsFixture{
		folderPath:  folderPath,
		fileService: fileService,
		usersTable:  &Fixture{tableName: "users"},
		deparser:    deparser,
	}
}

func (f *migrationsFixture) AddRawFile(filename, content string) {
	err := f.fileService.Write(f.folderPath, filename, content)
	Expect(err).To(BeNil())
}

func (f *migrationsFixture) AddMigration0001() *Fixture {
	f.usersTable.varcharColumns = []VarcharFixture{
		{"social_number", 9},
		{"name", 15},
		{"phone", 11},
	}
	upStatements := []string{
		f.deparser.CreateTable(f.usersTable.ToDomainTable()),
	}
	f.AddRawFile(f.migrationUpMockName(1), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		f.deparser.DropTable(f.usersTable.tableName),
	}
	f.AddRawFile(f.migrationDownMockName(1), f.deparser.WriteScript(downStatements))
	return f.usersTable
}

func (f *migrationsFixture) AddMigration0002() (VarcharFixture, *ViewFixture) {
	f.usersView = &ViewFixture{viewName: "user_phones", columns: []string{"name", "full_phone"}}
	newVarcharFixture := VarcharFixture{name: "ddi", cap: 3}
	f.usersTable.varcharColumns = append(f.usersTable.varcharColumns, newVarcharFixture)
	upStatements := []string{
		f.deparser.AddColumn(f.usersTable.tableName, newVarcharFixture.name, newVarcharFixture.ToDomainColumn()),
		fmt.Sprintf(`
		CREATE VIEW %s AS
			SELECT %s,
			CONCAT(%s, %s) AS %s
			FROM %s`, f.usersView.viewName,
			f.usersView.columns[0],
			f.usersTable.varcharColumns[2].name, f.usersTable.varcharColumns[3].name, f.usersView.columns[1],
			f.usersTable.tableName),
	}
	f.AddRawFile(f.migrationUpMockName(2), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		fmt.Sprintf("DROP VIEW IF EXISTS %s", f.usersView.viewName),
		f.deparser.DropColumn(f.usersTable.tableName, newVarcharFixture.name),
	}
	f.AddRawFile(f.migrationDownMockName(2), f.deparser.WriteScript(downStatements))
	return newVarcharFixture, f.usersView
}

func (f *migrationsFixture) AddMigration0003() VarcharFixture {
	newVarcharFixture := VarcharFixture{name: "abc", cap: 3}
	f.usersTable.varcharColumns = append(f.usersTable.varcharColumns, newVarcharFixture)
	upStatements := []string{
		f.deparser.AddColumn(f.usersTable.tableName, newVarcharFixture.name, newVarcharFixture.ToDomainColumn()),
	}
	f.AddRawFile(f.migrationUpMockName(3), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		f.deparser.DropColumn(f.usersTable.tableName, newVarcharFixture.name),
	}
	f.AddRawFile(f.migrationDownMockName(3), f.deparser.WriteScript(downStatements))
	return newVarcharFixture
}

func (f *migrationsFixture) TearDown() {
	err := os.RemoveAll(f.folderPath)
	Expect(err).To(BeNil())
}

func (f *migrationsFixture) migrationUpMockName(n int) string {
	return fmt.Sprintf("%04d_test_migration.up.sql", n)
}

func (f *migrationsFixture) migrationDownMockName(n int) string {
	return fmt.Sprintf("%04d_test_migration.down.sql", n)
}
