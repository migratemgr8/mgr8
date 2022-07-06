package cmd

import (
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

var postgresTestDriver mgr8testing.TestDriver
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

	postgresTestDriver = mgr8testing.NewTestDriver(global.Postgres)
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
	err := dm.CloseAll()
	Expect(err).To(BeNil())
	postgresMigrations.TearDown()
})

func getDriverSuccessfully(d global.Database) domain.Driver {
	driver, err := drivers.GetDriver(d)
	Expect(err).To(BeNil())
	return driver
}
