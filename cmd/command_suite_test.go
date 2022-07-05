package cmd

import (
	"github.com/jmoiron/sqlx"
	"github.com/migratemgr8/mgr8/infrastructure"
	"github.com/migratemgr8/mgr8/testing/fixtures"
	"testing"

	"github.com/migratemgr8/mgr8/domain"
	"github.com/migratemgr8/mgr8/drivers"
	"github.com/migratemgr8/mgr8/global"
	mgr8testing "github.com/migratemgr8/mgr8/testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _t *testing.T
var dm *mgr8testing.DockerManager
var postgresDb *sqlx.DB
var mySqlDb *sqlx.DB

var postgresDriver domain.Driver
var testMigrationsFolder = "apply-test-migrations"
var postgresMigrations fixtures.MigrationsFixture
var userFixture0001 *fixtures.Fixture
var firstNewColumnFixture0002 fixtures.VarcharFixture
var userViewFixture0002 *fixtures.ViewFixture
var secondNewColumnFixture0003 fixtures.VarcharFixture

func TestCommand(t *testing.T) {
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Test Suite")
}

var _ = BeforeSuite(func() {
	dm = mgr8testing.NewDockerManager()
	var err error
	postgresDb, err = sqlx.Connect(global.Postgres.String(), dm.GetConnectionString(global.Postgres))
	Expect(err).To(BeNil())
	mySqlDb, err = sqlx.Connect(global.MySql.String(), dm.GetConnectionString(global.MySql))
	Expect(err).To(BeNil())

	postgresDriver = getDriverSuccessfully(global.Postgres)
	postgresMigrations = fixtures.NewMigrationsFixture(testMigrationsFolder,
		infrastructure.NewFileService(),
		postgresDriver.Deparser(),
	)
	userFixture0001 = postgresMigrations.AddMigration0001()
	firstNewColumnFixture0002, userViewFixture0002 = postgresMigrations.AddMigration0002()
	secondNewColumnFixture0003 = postgresMigrations.AddMigration0003()
})

var _ = AfterSuite(func() {
	err := postgresDb.Close()
	Expect(err).To(BeNil())
	err = mySqlDb.Close()
	Expect(err).To(BeNil())
	err = dm.CloseAll()
	Expect(err).To(BeNil())
})

func getDriverSuccessfully(d global.Database) domain.Driver {
	driver, err := drivers.GetDriver(d)
	Expect(err).To(BeNil())
	return driver
}
