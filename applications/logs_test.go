package applications

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	domain_mock "github.com/migratemgr8/mgr8/mock/domain"
)

var _ = Describe("Logs", func() {
	Context("CheckAndInstallTool", func() {
		var (
			driver *domain_mock.MockDriver
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			driver = domain_mock.NewMockDriver(ctrl)
		})
		When("Driver has error", func() {
			It("Returns error", func() {
				driverError := errors.New("driver error")
				driver.EXPECT().IsToolInstalled().Return(false, driverError)
				err := CheckAndInstallTool(driver)
				Expect(err).To(Equal(driverError))
			})
		})

		When("Tool is not installed", func() {
			It("Installs and returns nil", func() {
				driver.EXPECT().IsToolInstalled().Return(false, nil)
				driver.EXPECT().InstallTool().Return(nil)
				err := CheckAndInstallTool(driver)
				Expect(err).To(BeNil())
			})
		})

		When("Tool is installed", func() {
			It("Returns nil", func() {
				driver.EXPECT().IsToolInstalled().Return(true, nil)
				err := CheckAndInstallTool(driver)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("GetMigrationNumber", func() {
		testCases := []struct {
			name           string
			expectedOutput int
		}{
			{expectedOutput: 1, name: "0001_test.up.sql"},
			{expectedOutput: 14, name: "0014_test"},
			{expectedOutput: 114, name: "0114_test.down.sql"},
			{expectedOutput: 1234, name: "1234_.sql"},
		}
		for _, testCase := range testCases {
			When(fmt.Sprintf("Asked with name %s", testCase.name), func() {
				It("Returns expected number", func() {
					migrationNumber, err := GetMigrationNumber(testCase.name)
					Expect(migrationNumber).To(Equal(testCase.expectedOutput))
					Expect(err).To(BeNil())
				})
			})
		}
	})

	Context("GetMigrationType", func() {
		testCases := []struct {
			name           string
			expectedOutput string
		}{
			{expectedOutput: "up", name: "0001_test.up.sql"},
			{expectedOutput: "down", name: "0114_test.down.sql"},
		}
		When(fmt.Sprintf("Asked with valid name"), func() {
			It("Returns expected type", func() {
				for _, testCase := range testCases {
					migrationType, err := GetMigrationType(testCase.name)
					Expect(migrationType).To(Equal(testCase.expectedOutput))
					Expect(err).To(BeNil())
				}

			})
		})
		When("Invalid migration type", func() {
			It("Returns error", func() {
				_, err := GetMigrationType("0014_test.left.sql")
				Expect(err).To(Equal(ErrInvalidMigrationType))
			})
		})

		When("Invalid migration format", func() {
			It("Returns error", func() {
				_, err := GetMigrationType("0014_test.sql")
				Expect(err).To(Equal(ErrInvalidMigrationName))
			})
		})
	})

})
