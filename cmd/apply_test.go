package cmd

import (
	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/infrastructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply integration test", func() {
	var (
		subject CommandExecutor
	)

	Context("execute", func() {
		When("apply up with no arguments and on postgres", func() {
			It("executes all files in folder", func() {
				subject = &apply{}
				err := subject.execute(
					[]string{"up", "3"},
					dm.GetConnectionString(global.Postgres),
					testMigrationsFolder,
					postgresDriver,
					infrastructure.CriticalLogLevel,
				)
				Expect(err).To(BeNil())

				finalUser := *userFixture0001
				finalUser.VarcharColumns = append(finalUser.VarcharColumns, firstNewColumnFixture0002)
				finalUser.VarcharColumns = append(finalUser.VarcharColumns, secondNewColumnFixture0003)

				exists, err := postgresTestDriver.AssertFixtureExistence(&finalUser)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
				exists, err = postgresTestDriver.AssertViewFixtureExistence(userViewFixture0002)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
	})
})
