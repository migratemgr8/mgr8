package applications

import (
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
)

var _ = Describe("Generate Command", func() {
	var (
		subject *generateCommand
	)

	Context("GetNextMigrationNumber", func() {
		var (
			driver                   *domain.MockDriver
			deparser                   *domain.MockDeparser
			migrationFileServiceMock *MockMigrationFileService
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			driver = domain.NewMockDriver(ctrl)
			migrationFileServiceMock = NewMockMigrationFileService(ctrl)
			subject = NewGenerateCommand(driver, migrationFileServiceMock)
			deparser = domain.NewMockDeparser(ctrl)

		})
		When("Asked to execute", func(){
			It("Succeeds", func(){
				driver.EXPECT().Deparser().Return(deparser).Times(2)
				migrationFileServiceMock.EXPECT().GetSchemaFromFile("mock_old_path").Return(&domain.Schema{}, nil)
				migrationFileServiceMock.EXPECT().GetSchemaFromFile("mock_new_path").Return(&domain.Schema{}, nil)
				migrationFileServiceMock.EXPECT().GetNextMigrationNumber("mock_dir").Return(3, nil)

				migrationFileServiceMock.EXPECT().
					WriteStatementsToFile("mock_dir", gomock.Len(0), 3, "up").
					Return( nil)

				migrationFileServiceMock.EXPECT().
					WriteStatementsToFile("mock_dir", gomock.Len(0), 3, "down").
					Return( nil)

				err := subject.Execute(&GenerateParameters{
					OldSchemaPath: "mock_old_path",
					NewSchemaPath: "mock_new_path",
					MigrationDir:  "mock_dir",
				})
				Expect(err).To(BeNil())
			})

		})

	})

})
