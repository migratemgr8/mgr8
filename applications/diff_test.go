package applications

//
//var _ = Describe("Migration Scripts", func() {
//	Context("Format Filename", func(){
//		When("Asked with one number", func() {
//			It("Parses each of them", func() {
//				migrationNumber := 1
//				migrationType := "up"
//
//
//				diffQueue := newSchema.Diff(oldSchema)
//				Expect(diffQueue).To(HaveLen(4))
//				Expect(diffQueue).To(ContainElements(
//					NewDropTableDiff(NewTable("old_table", map[string]*Column{})),
//					NewDropColumnDiff("kept_table", "old_column"),
//					NewCreateTableDiff(NewTable("new_table", map[string]*Column{})),
//					NewCreateColumnDiff("kept_table", "new_column", &Column{}),
//				))
//			})
//		})
//	})
//
//})