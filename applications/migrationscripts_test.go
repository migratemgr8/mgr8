package applications

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

var _ = Describe("Migration Scripts", func() {
	var (
		subject *migrationFileService
	)
	Context("Format Filename", func() {
		var mockTime time.Time
		testCases := []struct {
			migrationNumber int
			migrationType   string
			expectedOutput  string
		}{
			{migrationType: "up", migrationNumber: 1, expectedOutput: "0001_1640995200.up.sql"},
			{migrationType: "down", migrationNumber: 14, expectedOutput: "0014_1640995200.down.sql"},
			{migrationType: "down", migrationNumber: 114, expectedOutput: "0114_1640995200.down.sql"},
			{migrationType: "up", migrationNumber: 1234, expectedOutput: "1234_1640995200.up.sql"},
		}
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			clock := infrastructure.NewMockClock(ctrl)
			mockTime = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

			clock.EXPECT().Now().Return(mockTime)
			fileService := infrastructure.NewMockFileService(ctrl)
			subject = NewMigrationFileService(fileService, clock, domain.NewMockDriver(ctrl))
		})
		for _, testCase := range testCases {
			When(fmt.Sprintf("Asked %s with number %d", testCase.migrationType, testCase.migrationNumber), func() {
				It("Generates expected filename", func() {
					fileName := subject.FormatFilename(testCase.migrationNumber, testCase.migrationType)
					Expect(fileName).To(Equal(testCase.expectedOutput))
				})
			})
		}
	})

	Context("GetNextMigrationNumber", func() {
		var fileService *infrastructure.MockFileService
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			clock := infrastructure.NewMockClock(ctrl)
			fileService = infrastructure.NewMockFileService(ctrl)
			subject = NewMigrationFileService(fileService, clock, domain.NewMockDriver(ctrl))
		})
		When("Has two migration files", func() {
			It("Next migration number is 3", func() {
				fileService.EXPECT().List("dir").Return([]infrastructure.MigrationFile{
					{FullPath: "/migrations/0001_name.up.sql", Name: "0001_name.up.sql"},
					{FullPath: "/migrations/0001_name.down.sql", Name: "0001_name.down.sql"},
					{FullPath: "/migrations/0002_name.up.sql", Name: "0002_name.uo.sql"},
					{FullPath: "/migrations/0002_name.down.sql", Name: "0002_name.down.sql"},
				}, nil)
				migrationNumber, err := subject.GetNextMigrationNumber("dir")
				Expect(err).To(BeNil())
				Expect(migrationNumber).To(Equal(3))
			})
		})
		When("Has one migration file and random file", func() {
			It("Next migration number is 2", func() {
				fileService.EXPECT().List("dir").Return([]infrastructure.MigrationFile{
					{FullPath: "/migrations/0001_name.up.sql", Name: "0001_name.up.sql"},
					{FullPath: "/migrations/0001_name.down.sql", Name: "0001_name.down.sql"},
					{FullPath: "/migrations/random_file.json", Name: "random_file.json"},
				}, nil)
				migrationNumber, err := subject.GetNextMigrationNumber("dir")
				Expect(err).To(BeNil())
				Expect(migrationNumber).To(Equal(2))
			})
		})
		When("Receives error", func() {
			It("Returns err", func() {
				expectedErr := fmt.Errorf("could not list dir")
				fileService.EXPECT().List("dir").Return(nil, expectedErr)
				migrationNumber, err := subject.GetNextMigrationNumber("dir")
				Expect(err).To(Equal(expectedErr))
				Expect(migrationNumber).To(Equal(0))
			})
		})
	})

	Context("GetSchemaFromFile", func() {
		var (
			fileService *infrastructure.MockFileService
			driver *domain.MockDriver
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			clock := infrastructure.NewMockClock(ctrl)
 			fileService = infrastructure.NewMockFileService(ctrl)
			driver = domain.NewMockDriver(ctrl)
							subject = NewMigrationFileService(fileService, clock,driver )
		})
		When("Reads file successfully", func() {
			It("Generates expected filename", func() {
				fileService.EXPECT().Read("filename").Return("content", nil)
				driver.EXPECT().ParseMigration("content").Return(&domain.Schema{ },nil)
				schema, err := subject.GetSchemaFromFile("filename")
				Expect(err).To(BeNil())
				Expect(schema).To(Equal(&domain.Schema{ }))
			})
		})
		When("Reads file returns error", func() {
			It("Returns the error", func() {
				expectedError := fmt.Errorf("expected error")
				fileService.EXPECT().Read("filename").Return("", expectedError)
				schema, err := subject.GetSchemaFromFile("filename")
				Expect(err).To(Equal(expectedError))
				Expect(schema).To(BeNil())
			})
		})
	})


})
