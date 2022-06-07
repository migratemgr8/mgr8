package domain

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Column Diff", func() {
	var (
		tableName string
		columnName string
		column *Column
		deparser *MockDeparser

		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(_t)
		tableName = "fake_table_name"
		columnName = "fake_column_name"
		column = &Column{Datatype: "fake_datatype", Parameters: map[string]interface{}{"fake_param": 10}, IsNotNull: false}
		deparser = NewMockDeparser(ctrl)
	})

	Context("Create Column", func() {
		var (
			subject *CreateColumnDiff
		)
		When("Asked to go up", func() {
			It("Calls Create Column deparser", func() {
				subject = NewCreateColumnDiff(tableName, columnName, column)

				deparser.EXPECT().AddColumn(tableName, columnName, column).Return("FAKE CREATE COLUMN")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE CREATE COLUMN"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Create Column deparser", func() {
				subject = NewCreateColumnDiff(tableName, columnName, column)

				deparser.EXPECT().DropColumn(tableName, columnName).Return("FAKE DROP COLUMN")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE DROP COLUMN"))
			})
		})
	})

	Context("Drop Column", func() {
		var (
			subject *DropColumnDiff
		)
		When("Asked to go up", func() {
			It("Calls Drop Column deparser", func() {
				subject = NewDropColumnDiff(tableName, columnName)

				deparser.EXPECT().DropColumn(tableName, columnName).Return("FAKE DROP COLUMN")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE DROP COLUMN"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Column deparser", func() {
				subject = NewDropColumnDiff(tableName, columnName)

				deparser.EXPECT().AddColumn(tableName, columnName, nil).Return("FAKE CREATE COLUMN")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE CREATE COLUMN"))
			})
		})
	})

	Context("Column Not Null", func() {
		var (
			subject *MakeColumnNotNullDiff
		)
		When("Asked to go up", func() {
			It("Calls Column Not Null deparser", func() {
				subject = NewMakeColumnNotNullDiff(tableName, columnName)

				deparser.EXPECT().MakeColumnNotNull(tableName, columnName, nil).Return("FAKE COLUMN NOT NULL")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE COLUMN NOT NULL"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Column Not Null deparser", func() {
				subject = NewMakeColumnNotNullDiff(tableName, columnName)

				deparser.EXPECT().MakeColumnNullable(tableName, columnName, nil).Return("FAKE COLUMN NULL")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE COLUMN NULL"))
			})
		})
	})

	Context("Column Nullable", func() {
		var (
			subject *MakeColumnNullableDiff
		)
		When("Asked to go up", func() {
			It("Calls Column Nullable deparser", func() {
				subject = NewMakeColumnNullableDiff(tableName, columnName)

				deparser.EXPECT().MakeColumnNullable(tableName, columnName, nil).Return("FAKE COLUMN NULL")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE COLUMN NULL"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Column Nullable deparser", func() {
				subject = NewMakeColumnNullableDiff(tableName, columnName)

				deparser.EXPECT().MakeColumnNotNull(tableName, columnName, nil).Return("FAKE COLUMN NOT NULL")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE COLUMN NOT NULL"))
			})
		})
	})

	AfterEach(func(){
		ctrl.Finish()
	})
})
