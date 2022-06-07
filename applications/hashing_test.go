package applications

import (
	"errors"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kenji-yamane/mgr8/infrastructure"
)

var _ = Describe("Hashing", func() {
	var (
		subject *hashService
	)

	Context("GetSqlHash", func() {
		var (
			mockFileService *infrastructure.MockFileService
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			mockFileService = infrastructure.NewMockFileService(ctrl)
			subject = NewHashService(mockFileService)
		})
		When("Asked to execute", func() {
			It("Succeeds", func() {
				content := "file_content"
				mockFileService.EXPECT().Read("file_name").Return(content, nil)
				hash, err := subject.GetSqlHash("file_name")
				Expect(err).To(BeNil())
				Expect(hash).To(Equal("7f0b6bb0b7e951b7fd2b2a4a326297e1"))
			})
		})

		When("Fails to read file", func() {
			It("Fails", func() {
				mockFileService.EXPECT().Read("file_name").Return("", errors.New("could not read file"))
				_, err := subject.GetSqlHash("file_name")
				Expect(err).To(Equal(errors.New("could not read file")))
			})
		})

	})

})
