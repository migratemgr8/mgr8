package domain

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Table Diff", func() {
	var (
		table    *Table
		deparser *MockDeparser

		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(_t)
		table = &Table{Name: "fake_table_name"}
		deparser = NewMockDeparser(ctrl)
	})

	Context("Create Table", func() {
		var (
			subject *CreateTableDiff
		)
		When("Asked to go up", func() {
			It("Calls Create Table deparser", func() {
				subject = NewCreateTableDiff(table)

				deparser.EXPECT().CreateTable(table).Return("FAKE CREATE TABLE")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE CREATE TABLE"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Table deparser", func() {
				subject = NewCreateTableDiff(table)

				deparser.EXPECT().DropTable("fake_table_name").Return("FAKE DROP TABLE")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE DROP TABLE"))
			})
		})
	})

	Context("Drop Table", func() {
		var (
			subject *DropTableDiff
		)
		When("Asked to go up", func() {
			It("Calls Drop Table deparser", func() {
				subject = NewDropTableDiff(table)

				deparser.EXPECT().DropTable("fake_table_name").Return("FAKE DROP TABLE")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE DROP TABLE"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Table deparser", func() {
				subject = NewDropTableDiff(table)

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
