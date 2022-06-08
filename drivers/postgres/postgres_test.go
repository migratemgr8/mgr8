package postgres

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/migratemgr8/mgr8/domain"
)

var _ = Describe("Postgres Driver", func() {

	Context("Parse Migration", func() {
		var (
			subject *postgresDriver
		)

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

		var (
			subject *deparser
		)

		BeforeEach(func() {
			subject = &deparser{}
		})

		When("New schema has a new column", func() {
			It("Creates alter table statement for column", func() {
				column := &domain.Column{
					Datatype:   "int",
					IsNotNull:  false,
					Parameters: map[string]interface{}{},
				}
				stmt := subject.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col int"))
			})

			It("Identifies not null property", func() {
				column := &domain.Column{
					Datatype:   "char",
					IsNotNull:  true,
					Parameters: map[string]interface{}{},
				}
				stmt := subject.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col char not null"))
			})

			It("Places correct parameters in column definition", func() {
				column := &domain.Column{
					Datatype:   "varchar",
					IsNotNull:  false,
					Parameters: map[string]interface{}{"size": 10},
				}
				stmt := subject.AddColumn("tbl", "col", column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl add column col varchar(10)"))
			})
		})

		When("New schema doesn't have a column", func() {
			It("Drops the column completly", func() {
				columnName := "col"
				tableName := "tbl"
				stmt := subject.DropColumn(tableName, columnName)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl drop column col"))
			})
		})

		When("A column changes its null property", func() {
			It("Makes a int null column become not null", func() {
				columnName := "col"
				tableName := "tbl"
				column := &domain.Column{Datatype: "int", IsNotNull: false}
				stmt := subject.MakeColumnNotNull(tableName, columnName, column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl alter column col set not null"))
			})

			It("Makes a int not null column become null", func() {
				columnName := "col"
				tableName := "tbl"
				column := &domain.Column{Datatype: "int", IsNotNull: false}
				stmt := subject.MakeColumnNullable(tableName, columnName, column)
				Expect(strings.ToLower(stmt)).To(Equal("alter table tbl alter column col drop not null"))
			})
		})

		When("New schema doesn't have a table", func() {
			It("Drops the table", func() {
				tableName := "tbl"
				stmt := subject.DropTable(tableName)
				Expect(strings.ToLower(stmt)).To(Equal("drop table if exists tbl"))
			})
		})

		When("Table has 1 as maximum argument in data type", func() {
			It("Generate CREATE TABLE statement", func() {
				table := &domain.Table{
					Name: "users",
					Columns: map[string]*domain.Column{
						"social_number":   {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(9)}},
						"favorite_number": {Datatype: "bit", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(6)}},
						"size":            {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
					},
				}

				statement := subject.CreateTable(table)
				answer := `CREATE TABLE users (
favorite_number bit(6),
size int4,
social_number varchar(9)
)`
				Expect(statement).To(Equal(answer))
			})
		})

		When("Table has 2 as maximum argument in data type", func() {
			It("Generate CREATE TABLE statement", func() {
				table := &domain.Table{
					Name: "users",
					Columns: map[string]*domain.Column{
						"area":      {Datatype: "decimal", IsNotNull: false, Parameters: map[string]interface{}{"precision": int32(9), "scale": int32(1)}},
						"perimeter": {Datatype: "numeric", IsNotNull: false, Parameters: map[string]interface{}{"precision": int32(6), "scale": int32(4)}},
						"name":      {Datatype: "char", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(10)}},
					},
				}

				statement := subject.CreateTable(table)
				answer := `CREATE TABLE users (
area decimal(9,1),
name char(10),
perimeter numeric(6,4)
)`
				Expect(statement).To(Equal(answer))
			})
		})
	})
})
