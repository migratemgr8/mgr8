package postgres

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
)

var _ = Describe("Postgres Driver", func() {
	var (
		subject *postgresDriver
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
					size INT
				);`
				schema, err := subject.ParseMigration(migration)
				Expect(err).To(BeNil())
				Expect(schema).To(Equal(&domain.Schema{
					Tables: map[string]*domain.Table{
						"users": {
							Columns: map[string]*domain.Column{
								"social_number": {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(9)}},
								"phone":         {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(11)}},
								"name":          {Datatype: "varchar", IsNotNull: false, Parameters: map[string]interface{}{"size": int32(15)}},
								"age":           {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
								"size":          {Datatype: "int4", IsNotNull: false, Parameters: map[string]interface{}{}},
							},
						},
					},
					Views: nil,
				}))
			})
		})
	})
})
