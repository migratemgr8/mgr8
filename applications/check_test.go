package applications

import (
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	infrastructure_mock "github.com/migratemgr8/mgr8/mock/infrastructure"
)

var _ = Describe("Check Command", func() {
	var (
		subject *checkCommand
	)

	Context("Execute", func() {
		var (
			mockFileService *infrastructure_mock.MockFileService
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			mockFileService = infrastructure_mock.NewMockFileService(ctrl)
			subject = NewCheckCommand(mockFileService)
		})
		When("Files match", func() {
			It("Succeeds", func() {
				mockFileService.EXPECT().Read("new_schema").Return("content", nil)
				mockFileService.EXPECT().Read("reference_schema").Return("content", nil)
				err := subject.Execute("reference_schema", "new_schema")
				Expect(err).To(BeNil())
			})
		})

		When("Files do not match", func() {
			It("Fails", func() {
				mockFileService.EXPECT().Read("new_schema").Return("new_content", nil)
				mockFileService.EXPECT().Read("reference_schema").Return("old_content", nil)
				err := subject.Execute("reference_schema", "new_schema")
				Expect(err).To(Equal(ErrFilesDoNotMatch))
			})
		})

	})

})
