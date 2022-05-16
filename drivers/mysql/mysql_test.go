package mysql

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
)

var _ = Describe("MySql Driver", func() {
	var (
		subject *mySqlDriver
	)

	Context("Parse Migration", func() {
		BeforeEach(func() {
			subject = NewMySqlDriver()
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
				Expect(schema).To(Not(BeNil()))
				Expect(schema).To(Equal(&domain.Schema{
					Tables: map[string]*domain.Table{
						"users": {
							Name: "users",
							Columns: map[string]*domain.Column{
								"phone":         {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": 11}},
								"social_number": {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": 9}},
								"name":          {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": 15}},
								"age":           {Datatype: "int", IsNotNull: false, Parameters: map[string]interface{}{}},
								"size":          {Datatype: "int", IsNotNull: false, Parameters: map[string]interface{}{}},
								"ddi":           {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": 3}},
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
})
