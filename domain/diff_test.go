package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
)

var _ = Describe("Schema Diff", func() {
	Context("Generate Diff", func() {
		When("Table has all data types", func() {
			It("Parses each of them", func() {
				oldSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"old_table": domain.NewTable("old_table", map[string]*domain.Column{}),
						"kept_table": domain.NewTable("kept_table", map[string]*domain.Column{
							"old_column":  &domain.Column{},
							"kept_column": &domain.Column{},
						}),
					},
					Views: nil,
				}
				newSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"new_table": domain.NewTable("new_table", map[string]*domain.Column{}),
						"kept_table": domain.NewTable("kept_table", map[string]*domain.Column{
							"kept_column": &domain.Column{},
							"new_column":  &domain.Column{},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(4))
				Expect(diffQueue).To(ContainElements(
					domain.NewDropTableDiff(domain.NewTable("old_table", map[string]*domain.Column{})),
					domain.NewDropColumnDiff("kept_table", "old_column"),
					domain.NewCreateTableDiff(domain.NewTable("new_table", map[string]*domain.Column{})),
					domain.NewCreateColumnDiff("kept_table", "new_column", &domain.Column{}),
				))
			})
		})

		When("Column switches to not null", func() {
			It("Returns MakeColumnNotNull", func() {
				oldSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"table": domain.NewTable("table", map[string]*domain.Column{
							"column": &domain.Column{IsNotNull: false},
						}),
					},
					Views: nil,
				}
				newSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"table": domain.NewTable("table", map[string]*domain.Column{
							"column": &domain.Column{IsNotNull: true},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(1))
				Expect(diffQueue).To(ContainElements(
					domain.NewMakeColumnNotNullDiff("table", "column"),
				))
			})
		})

		When("Column switches to nullable", func() {
			It("Returns UnmakeColumnNotNull", func() {
				oldSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"table": domain.NewTable("table", map[string]*domain.Column{
							"column": &domain.Column{IsNotNull: true},
						}),
					},
					Views: nil,
				}
				newSchema := &domain.Schema{
					Tables: map[string]*domain.Table{
						"table": domain.NewTable("table", map[string]*domain.Column{
							"column": &domain.Column{IsNotNull: false},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(1))
				Expect(diffQueue).To(ContainElements(
					domain.NewMakeColumnNullableDiff("table", "column"),
				))
			})
		})
	})
})
