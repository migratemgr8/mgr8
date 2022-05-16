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
							"old_column": &Column{},
							"kept_column": &Column{},
						}),
					},
					Views:  nil,
				}
				newSchema := &Schema{
					Tables: map[string]*Table{
						"new_table": NewTable("new_table", map[string]*Column{}),
						"kept_table": NewTable("kept_table", map[string]*Column{
							"kept_column": &Column{},
							"new_column": &Column{},
						}),
					},
					Views:  nil,
				}

				diffQueue := newSchema.Diff(oldSchema)
				Expect(diffQueue).To(HaveLen(4))
				Expect(diffQueue).To(ContainElements(
					NewDropTableDiff(NewTable("old_table", map[string]*Column{})),
					NewDropColumnDiff("kept_table", "old_column"),
					NewCreateTableDiff( NewTable("new_table", map[string]*Column{})),
					NewCreateColumnDiff("kept_table", &Column{}),
					))
			})
		})
	})
})
