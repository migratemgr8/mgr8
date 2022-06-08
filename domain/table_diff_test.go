package domain_test

import (
	"github.com/golang/mock/gomock"
	"github.com/migratemgr8/mgr8/domain"
	domain_mock "github.com/migratemgr8/mgr8/mock/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Table Diff", func() {
	var (
		table    *domain.Table
		deparser *domain_mock.MockDeparser

		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(_t)
		table = &domain.Table{Name: "fake_table_name"}
		deparser = domain_mock.NewMockDeparser(ctrl)
	})

	Context("Create Table", func() {
		var (
			subject *domain.CreateTableDiff
		)
		When("Asked to go up", func() {
			It("Calls Create Table deparser", func() {
				subject = domain.NewCreateTableDiff(table)

				deparser.EXPECT().CreateTable(table).Return("FAKE CREATE TABLE")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE CREATE TABLE"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Table deparser", func() {
				subject = domain.NewCreateTableDiff(table)

				deparser.EXPECT().DropTable("fake_table_name").Return("FAKE DROP TABLE")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE DROP TABLE"))
			})
		})
	})

	Context("Drop Table", func() {
		var (
			subject *domain.DropTableDiff
		)
		When("Asked to go up", func() {
			It("Calls Drop Table deparser", func() {
				subject = domain.NewDropTableDiff(table)

				deparser.EXPECT().DropTable("fake_table_name").Return("FAKE DROP TABLE")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE DROP TABLE"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Table deparser", func() {
				subject = domain.NewDropTableDiff(table)

				deparser.EXPECT().CreateTable(table).Return("FAKE CREATE TABLE")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE CREATE TABLE"))
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
