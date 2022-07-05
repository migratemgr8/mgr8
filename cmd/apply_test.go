package cmd

import (
	"github.com/migratemgr8/mgr8/domain"
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
					[]string{"up"},
					dm.GetConnectionString(global.Postgres),
					testMigrationsFolder,
					postgresDriver,
					infrastructure.CriticalLogLevel,
				)
				Expect(err).To(BeNil())

				var exists bool
				err = postgresDb.QueryRow(`
				SELECT EXISTS (
				    SELECT FROM information_schema.tables WHERE  table_name   = $1
	   			)`, domain.LogsTableName).Scan(&exists)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
	})
})
