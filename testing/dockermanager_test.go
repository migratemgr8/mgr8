package testing

import (
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/migratemgr8/mgr8/global"
)

var _t *testing.T

func TestTestingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testing Test Suite")
}

var _ = Describe("Check Command", func() {
	var (
		subject *DockerManager
	)

	BeforeSuite(func() {
		subject = NewDockerManager()
	})
	AfterSuite(func() {
		err := subject.CloseAll()
		Expect(err).To(BeNil())
	})

	Context("GetConnectionString", func() {
		When("postgres requested", func() {
			It("should return viable connection", func() {
				postgresConn := subject.GetConnectionString(global.Postgres)
				db, err := sqlx.Connect(global.Postgres.String(), postgresConn)
				Expect(err).To(BeNil())
				Expect(db).To(Not(BeNil()))
				_, err = db.Exec(`SELECT 1`)
				Expect(err).To(BeNil())
				err = db.Close()
				Expect(err).To(BeNil())
			})
		})

		When("mysql requested", func() {
			It("should return viable connection", func() {
				mySqlConn := subject.GetConnectionString(global.MySql)
				db, err := sqlx.Connect(global.MySql.String(), mySqlConn)
				Expect(err).To(BeNil())
				Expect(db).To(Not(BeNil()))
				_, err = db.Exec(`SELECT 1`)
				Expect(err).To(BeNil())
				err = db.Close()
				Expect(err).To(BeNil())
			})
		})
	})
})
