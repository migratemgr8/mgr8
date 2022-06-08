package applications

import (
	"errors"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	applications_mock "github.com/kenji-yamane/mgr8/mock/applications"
)

var _ = Describe("Empty Command", func() {
	var (
		subject *emptyCommand
	)

	Context("Execute", func() {
		var (
			mockFileService *applications_mock.MockMigrationFileService
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			mockFileService = applications_mock.NewMockMigrationFileService(ctrl)
			subject = NewEmptyCommand(mockFileService)
		})
		When("Asked to execute", func() {
			It("Succeeds", func() {
				mockFileService.EXPECT().GetNextMigrationNumber("migrations_dir").Return(2, nil)
				mockFileService.EXPECT().WriteStatementsToFile("migrations_dir", []string{}, 2, "up").Return(nil)
				mockFileService.EXPECT().WriteStatementsToFile("migrations_dir", []string{}, 2, "down").Return(nil)
				err := subject.Execute("migrations_dir")
				Expect(err).To(BeNil())
			})
		})

		When("Fails to read file", func() {
			It("Fails", func() {
				expectedError := errors.New("could not get next migration")
				mockFileService.EXPECT().GetNextMigrationNumber("migrations_dir").Return(0, expectedError)
				err := subject.Execute("migrations_dir")
				Expect(err).To(Equal(expectedError))
			})
		})

	})

})
