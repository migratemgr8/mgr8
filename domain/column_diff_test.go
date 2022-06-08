package domain_test

import (
	"github.com/golang/mock/gomock"
	"github.com/kenji-yamane/mgr8/domain"
	domain_mock "github.com/kenji-yamane/mgr8/mock/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Column Diff", func() {
	var (
		tableName  string
		columnName string
		column     *domain.Column
		deparser   *domain_mock.MockDeparser

		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(_t)
		tableName = "fake_table_name"
		columnName = "fake_column_name"
		column = &domain.Column{Datatype: "fake_datatype", Parameters: map[string]interface{}{"fake_param": 10}, IsNotNull: false}
		deparser = domain_mock.NewMockDeparser(ctrl)
	})

	Context("Create Column", func() {
		var (
			subject *domain.CreateColumnDiff
		)
		When("Asked to go up", func() {
			It("Calls Create Column deparser", func() {
				subject = domain.NewCreateColumnDiff(tableName, columnName, column)

				deparser.EXPECT().AddColumn(tableName, columnName, column).Return("FAKE CREATE COLUMN")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE CREATE COLUMN"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Create Column deparser", func() {
				subject = domain.NewCreateColumnDiff(tableName, columnName, column)

				deparser.EXPECT().DropColumn(tableName, columnName).Return("FAKE DROP COLUMN")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE DROP COLUMN"))
			})
		})
	})

	Context("Drop Column", func() {
		var (
			subject *domain.DropColumnDiff
		)
		When("Asked to go up", func() {
			It("Calls Drop Column deparser", func() {
				subject = domain.NewDropColumnDiff(tableName, columnName, column)

				deparser.EXPECT().DropColumn(tableName, columnName).Return("FAKE DROP COLUMN")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE DROP COLUMN"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Drop Column deparser", func() {
				subject = domain.NewDropColumnDiff(tableName, columnName, column)

				deparser.EXPECT().AddColumn(tableName, columnName, column).Return("FAKE CREATE COLUMN")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE CREATE COLUMN"))
			})
		})
	})

	Context("Column Not Null", func() {
		var (
			subject *domain.MakeColumnNotNullDiff
		)
		When("Asked to go up", func() {
			It("Calls Column Not Null deparser", func() {
				subject = domain.NewMakeColumnNotNullDiff(tableName, columnName, column)

				deparser.EXPECT().MakeColumnNotNull(tableName, columnName, column).Return("FAKE COLUMN NOT NULL")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE COLUMN NOT NULL"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Column Not Null deparser", func() {
				subject = domain.NewMakeColumnNotNullDiff(tableName, columnName, column)

				deparser.EXPECT().MakeColumnNullable(tableName, columnName, column).Return("FAKE COLUMN NULL")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE COLUMN NULL"))
			})
		})
	})

	Context("Column Nullable", func() {
		var (
			subject *domain.MakeColumnNullableDiff
		)
		When("Asked to go up", func() {
			It("Calls Column Nullable deparser", func() {
				subject = domain.NewMakeColumnNullableDiff(tableName, columnName, column)

				deparser.EXPECT().MakeColumnNullable(tableName, columnName, column).Return("FAKE COLUMN NULL")
				result := subject.Up(deparser)
				Expect(result).To(Equal("FAKE COLUMN NULL"))
			})
		})
		When("Asked to go down", func() {
			It("Calls Column Nullable deparser", func() {
				subject = domain.NewMakeColumnNullableDiff(tableName, columnName, column)

				deparser.EXPECT().MakeColumnNotNull(tableName, columnName, column).Return("FAKE COLUMN NOT NULL")
				result := subject.Down(deparser)
				Expect(result).To(Equal("FAKE COLUMN NOT NULL"))
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
