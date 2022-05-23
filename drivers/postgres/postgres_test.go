package postgres

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
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
								"ddi":          {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(3)}},
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
			subject deparser
		)

		BeforeEach(func() {
			subject = deparser{}
		})

		When("Table has 1 as maximum argument in data type", func() {
			It("Generate CREATE TABLE statement", func() {
				table := &domain.Table{
					Name: "users",
					Columns: map[string]*domain.Column{
						"social_number": {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(9)}},
						"favorite_number": {Datatype: "bit", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(6)}},
						"size":          {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
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
						"area": {Datatype: "decimal", IsNotNull: false, Parameters: map[string]interface{}{"precision": int32(9), "scale": int32(1)}},
						"perimeter": {Datatype: "numeric", IsNotNull: false, Parameters: map[string]interface{}{"precision": int32(6), "scale": int32(4)}},
						"name":          {Datatype: "char", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(10)}},
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
