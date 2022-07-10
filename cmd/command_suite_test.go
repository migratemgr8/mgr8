package cmd

import (
	"github.com/migratemgr8/mgr8/applications"
	"github.com/migratemgr8/mgr8/infrastructure"
	"github.com/migratemgr8/mgr8/testing/fixtures"
	"os"
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

var (
	postgresTestDriver         mgr8testing.TestDriver
	postgresDriver             domain.Driver
	postgresMigrations         fixtures.MigrationsFixture
	userFixture0001            fixtures.Fixture
	firstNewColumnFixture0002  fixtures.VarcharFixture
	userViewFixture0002        fixtures.ViewFixture
	secondNewColumnFixture0003 fixtures.VarcharFixture
)

var (
	testMigrationsFolder = "apply-test-migrations"
	mockUser             = "mock-user"
)

func TestCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Test Suite")
}

var _ = BeforeSuite(func() {
	createConfigFileIfNotExists()
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

func createConfigFileIfNotExists() {
	configPath, err := applications.GetConfigFilePath()
	Expect(err).To(BeNil())
	config, err := os.Open(configPath)
	if err == nil {
		return
	}
	config, err = os.Create(configPath)
	Expect(err).To(BeNil())
	username := mockUser
	hostname, err := os.Hostname()
	Expect(err).To(BeNil())
	err = applications.InsertUserDetails(username, hostname, config)
	Expect(err).To(BeNil())
	err = config.Close()
	Expect(err).To(BeNil())
}

func getDriverSuccessfully(d global.Database) domain.Driver {
	driver, err := drivers.GetDriver(d)
	Expect(err).To(BeNil())
	return driver
}
