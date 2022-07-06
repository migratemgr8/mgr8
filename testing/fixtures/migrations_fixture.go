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
		usersTable:  &Fixture{TableName: "users"},
		deparser:    deparser,
	}
}

func (f *migrationsFixture) AddRawFile(filename, content string) {
	err := f.fileService.Write(f.folderPath, filename, content)
	Expect(err).To(BeNil())
}

func (f *migrationsFixture) AddMigration0001() *Fixture {
	f.usersTable.VarcharColumns = []VarcharFixture{
		{"social_number", 9},
		{"name", 15},
		{"phone", 11},
	}
	upStatements := []string{
		f.deparser.CreateTable(f.usersTable.ToDomainTable()),
	}
	f.AddRawFile(f.migrationUpMockName(1), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		f.deparser.DropTable(f.usersTable.TableName),
	}
	f.AddRawFile(f.migrationDownMockName(1), f.deparser.WriteScript(downStatements))
	return f.usersTable
}

func (f *migrationsFixture) AddMigration0002() (VarcharFixture, *ViewFixture) {
	f.usersView = &ViewFixture{
		ViewName:       "user_phones",
		VarcharColumns: []VarcharFixture{f.usersTable.VarcharColumns[2]},
		TextColumns:    []string{"full_phone"},
	}
	newVarcharFixture := VarcharFixture{Name: "ddi", Cap: 3}
	f.usersTable.VarcharColumns = append(f.usersTable.VarcharColumns, newVarcharFixture)
	f.usersView.Statement = fmt.Sprintf(`SELECT %s, CONCAT(%s, %s) AS %s FROM %s`,
		f.usersView.VarcharColumns[0].Name,
		f.usersTable.VarcharColumns[2].Name, f.usersTable.VarcharColumns[3].Name,
		f.usersView.TextColumns[0], f.usersTable.TableName)
	upStatements := []string{
		f.deparser.AddColumn(f.usersTable.TableName, newVarcharFixture.Name, newVarcharFixture.ToDomainColumn()),
		fmt.Sprintf(`CREATE VIEW %s AS %s`, f.usersView.ViewName, f.usersView.Statement),
	}
	f.AddRawFile(f.migrationUpMockName(2), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		fmt.Sprintf("DROP VIEW IF EXISTS %s", f.usersView.ViewName),
		f.deparser.DropColumn(f.usersTable.TableName, newVarcharFixture.Name),
	}
	f.AddRawFile(f.migrationDownMockName(2), f.deparser.WriteScript(downStatements))
	return newVarcharFixture, f.usersView
}

func (f *migrationsFixture) AddMigration0003() VarcharFixture {
	newVarcharFixture := VarcharFixture{Name: "abc", Cap: 3}
	f.usersTable.VarcharColumns = append(f.usersTable.VarcharColumns, newVarcharFixture)
	upStatements := []string{
		f.deparser.AddColumn(f.usersTable.TableName, newVarcharFixture.Name, newVarcharFixture.ToDomainColumn()),
	}
	f.AddRawFile(f.migrationUpMockName(3), f.deparser.WriteScript(upStatements))
	downStatements := []string{
		f.deparser.DropColumn(f.usersTable.TableName, newVarcharFixture.Name),
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
