package domain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema Diff", func() {
	Context("Generate Diff", func() {
		When("Table has all data types", func() {
			It("Parses each of them", func() {
				oldSchema := &Schema{
					Tables: map[string]*Table{
						"old_table": NewTable("old_table", map[string]*Column{}),
						"kept_table": NewTable("kept_table", map[string]*Column{
							"old_column":  &Column{},
							"kept_column": &Column{},
						}),
					},
					Views: nil,
				}
				newSchema := &Schema{
					Tables: map[string]*Table{
						"new_table": NewTable("new_table", map[string]*Column{}),
						"kept_table": NewTable("kept_table", map[string]*Column{
							"kept_column": &Column{},
							"new_column":  &Column{},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(4))
				Expect(diffQueue).To(ContainElements(
					NewDropTableDiff(NewTable("old_table", map[string]*Column{})),
					NewDropColumnDiff("kept_table", "old_column"),
					NewCreateTableDiff(NewTable("new_table", map[string]*Column{})),
					NewCreateColumnDiff("kept_table", "new_column", &Column{}),
				))
			})
		})

		When("Column switches to not null", func() {
			It("Returns MakeColumnNotNull", func() {
				oldSchema := &Schema{
					Tables: map[string]*Table{
						"table": NewTable("table", map[string]*Column{
							"column": &Column{IsNotNull: false},
						}),
					},
					Views: nil,
				}
				newSchema := &Schema{
					Tables: map[string]*Table{
						"table": NewTable("table", map[string]*Column{
							"column": &Column{IsNotNull: true},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(1))
				Expect(diffQueue).To(ContainElements(
					NewMakeColumnNotNullDiff("table", "column"),
				))
			})
		})

		When("Column switches to nullable", func() {
			It("Returns MakeColumnNullable", func() {
				oldSchema := &Schema{
					Tables: map[string]*Table{
						"table": NewTable("table", map[string]*Column{
							"column": {IsNotNull: true},
						}),
					},
					Views: nil,
				}
				newSchema := &Schema{
					Tables: map[string]*Table{
						"table": NewTable("table", map[string]*Column{
							"column": {IsNotNull: false},
						}),
					},
					Views: nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(1))
				Expect(diffQueue).To(ContainElements(
					NewMakeColumnNullableDiff("table", "column"),
				))
			})
		})

		When("New column is added", func() {
			It("Returns CreateColumn", func() {
				column := &Column{Datatype: "integer"}

				oldSchema := &Schema{
					Tables: map[string]*Table{
						"tableName": NewTable("tableName", map[string]*Column{}),
					},
				}
				newSchema := &Schema{
					Tables: map[string]*Table{
						"tableName": NewTable("tableName", map[string]*Column{
							"newColumn": column,
						}),
					},
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(1))
				Expect(diffQueue).To(ContainElements(
					NewCreateColumnDiff("tableName", "newColumn", column),
				))
			})
		})
	})
})
