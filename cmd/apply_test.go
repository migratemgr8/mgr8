package cmd

import (
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/infrastructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply integration test", func() {
	var (
		subject CommandExecutor = &apply{}
	)

	executeApply := func(args []string) {
		err := subject.execute(
			args,
			dm.GetConnectionString(global.Postgres),
			testMigrationsFolder,
			postgresDriver,
			infrastructure.CriticalLogLevel,
		)
		Expect(err).To(BeNil())
	}

	Context("execute", func() {
		When("up with 3", func() {
			It("executes all three files that increase migration in folder", func() {
				AssertStateBeforeAllMigrations()
				executeApply([]string{"up", "3"})
				AssertStateAfterAllMigrations()
			})
		})
		When("down with no number specified", func() {
			It("decreases one", func() {
				executeApply([]string{"down"})
				AssertStateAfterMigration0002AndBefore0003()
				executeApply([]string{"down"})
				AssertStateAfterMigration0001AndBefore0002()
				executeApply([]string{"down"})
				AssertStateBeforeAllMigrations()
			})
		})
	})
})

func AssertStateBeforeAllMigrations() {
	exists, err := postgresTestDriver.AssertTableExistence(userFixture0001.TableName)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
	exists, err = postgresTestDriver.AssertViewExistence(userViewFixture0002.ViewName)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func AssertStateAfterMigration0001AndBefore0002() {
	exists, err := postgresTestDriver.AssertFixtureExistence(userFixture0001)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, firstNewColumnFixture0002)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
	exists, err = postgresTestDriver.AssertViewExistence(userViewFixture0002.ViewName)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, secondNewColumnFixture0003)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func AssertStateAfterMigration0002AndBefore0003() {
	exists, err := postgresTestDriver.AssertFixtureExistence(userFixture0001)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, firstNewColumnFixture0002)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertViewFixtureExistence(userViewFixture0002)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, secondNewColumnFixture0003)
	Expect(err).To(BeNil())
	Expect(exists).To(BeFalse())
}

func AssertStateAfterAllMigrations() {
	exists, err := postgresTestDriver.AssertFixtureExistence(userFixture0001)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, firstNewColumnFixture0002)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertViewFixtureExistence(userViewFixture0002)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
	exists, err = postgresTestDriver.AssertVarcharExistence(userFixture0001.TableName, secondNewColumnFixture0003)
	Expect(err).To(BeNil())
	Expect(exists).To(BeTrue())
}
