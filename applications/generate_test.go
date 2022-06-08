package applications

import (
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
	applications_mock "github.com/kenji-yamane/mgr8/mock/applications"
	domain_mock "github.com/kenji-yamane/mgr8/mock/domain"
	infrastructure_mock "github.com/kenji-yamane/mgr8/mock/infrastructure"
)

var _ = Describe("Generate Command", func() {
	var (
		subject *generateCommand
	)

	Context("GetNextMigrationNumber", func() {
		var (
			driver                   *domain_mock.MockDriver
			deparser                 *domain_mock.MockDeparser
			migrationFileServiceMock *applications_mock.MockMigrationFileService
			fileService              *infrastructure_mock.MockFileService
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			driver = domain_mock.NewMockDriver(ctrl)
			migrationFileServiceMock = applications_mock.NewMockMigrationFileService(ctrl)
			fileService = infrastructure_mock.NewMockFileService(ctrl)
			subject = NewGenerateCommand(driver, migrationFileServiceMock, fileService)
			deparser = domain_mock.NewMockDeparser(ctrl)

		})
		When("Asked to execute", func() {
			It("Succeeds", func() {
				fileService.EXPECT().Read("mock_new_path").Return("content", nil)
				fileService.EXPECT().Write(".mgr8", "reference.sql", "content")
				driver.EXPECT().Deparser().Return(deparser).Times(2)
				migrationFileServiceMock.EXPECT().GetSchemaFromFile("mock_old_path").Return(&domain.Schema{}, nil)
				migrationFileServiceMock.EXPECT().GetSchemaFromFile("mock_new_path").Return(&domain.Schema{}, nil)
				migrationFileServiceMock.EXPECT().GetNextMigrationNumber("mock_dir").Return(3, nil)

				migrationFileServiceMock.EXPECT().
					WriteStatementsToFile("mock_dir", gomock.Len(0), 3, "up").
					Return(nil)

				migrationFileServiceMock.EXPECT().
					WriteStatementsToFile("mock_dir", gomock.Len(0), 3, "down").
					Return(nil)

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
