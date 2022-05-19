package postgres

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
)

var _ = Describe("Postgres Driver", func() {
	var (
		subject *postgresDriver
		dp      *deparser
	)

	Context("Parse Migration", func() {
		BeforeEach(func() {
			subject = NewPostgresDriver()
		})

		When("Table has all data types", func() {
			It("Parses each of them", func() {
				migration := `
				CREATE TABLE users (
					social_number VARCHAR(9) PRIMARY KEY,
					phone VARCHAR(11),
					name VARCHAR(15),
					age INTEGER,
					size INT,
					ddi VARCHAR(3)
				);

				CREATE VIEW user_phones AS
				SELECT name, CONCAT(ddi, phone) AS full_phone FROM users;`

				schema, err := subject.ParseMigration(migration)
				Expect(err).To(BeNil())
				Expect(schema).To(Equal(&domain.Schema{
					Tables: map[string]*domain.Table{
						"users": {
							Name: "users",
							Columns: map[string]*domain.Column{
								"social_number": {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(9)}},
								"phone":         {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(11)}},
								"name":          {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(15)}},
								"ddi":           {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(3)}},
								"age":           {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
								"size":          {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
							},
						},
					},
					Views: map[string]*domain.View{
						"user_phones": {SQL: ""},
					},
				}))
			})
		})
	})

	Context("Deparse Migration", func() {
		BeforeEach(func() {
			dp = &deparser{}
		})

		When("New schema has a new column", func() {
			It("Creates alter table statement for column", func() {
				column := &domain.Column{
					Datatype:   "int",
					IsNotNull:  false,
					Parameters: map[string]interface{}{},
				}
				stmt := dp.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col int"))
			})

			It("Identifies not null property", func() {
				column := &domain.Column{
					Datatype:   "char",
					IsNotNull:  true,
					Parameters: map[string]interface{}{},
				}
				stmt := dp.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col char not null"))
			})

			It("Places correct parameters in column definition", func() {
				column := &domain.Column{
					Datatype:   "varchar",
					IsNotNull:  false,
					Parameters: map[string]interface{}{"size": 10},
				}
				stmt := dp.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col varchar(10)"))
			})
		})

		When("New schema doesn't have a column", func() {
			It("Drops the column completly", func() {
				columnName := "col"
				tableName := "tbl"
				stmt := dp.DropColumn(tableName, columnName)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl drop column col"))
			})
		})

		When("A column changes its null property", func() {
			It("Makes a int null column become not null", func() {
				columnName := "col"
				tableName := "tbl"
				stmt := dp.MakeColumnNotNull(tableName, columnName)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl alter column col set not null"))
			})

			It("Makes a int not null column become null", func() {
				columnName := "col"
				tableName := "tbl"
				stmt := dp.UnmakeColumnNotNull(tableName, columnName)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl alter column col drop not null"))
			})
		})

		When("New schema doesn't have a table", func() {
			It("Drops the table", func() {
				tableName := "tbl"
				stmt := dp.DropTable(tableName)
				Expect(strings.ToLower(stmt)).To(Equal("drop table if exists tbl"))
			})
		})
	})
})
