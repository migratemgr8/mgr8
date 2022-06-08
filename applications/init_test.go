package applications

import (
	"errors"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	infrastructure_mock "github.com/kenji-yamane/mgr8/mock/infrastructure"
)

var _ = Describe("Init Command", func() {
	var (
		subject *initCommand
	)

	Context("Execute", func() {
		var (
			mockFileService *infrastructure_mock.MockFileService
			content = "file_content"
			fileName = "file_name"
			applicationFolder = "app_folder"
			referenceFile = "ref_file"
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(_t)
			mockFileService = infrastructure_mock.NewMockFileService(ctrl)
			subject = NewInitCommand(mockFileService)
		})
		When("Asked to execute", func() {
			It("Succeeds", func() {
				mockFileService.EXPECT().CreateFolderIfNotExists(applicationFolder).Return(nil)
				mockFileService.EXPECT().Read(fileName).Return(content, nil)
				mockFileService.EXPECT().Write(applicationFolder, referenceFile, content).Return(nil)
				err := subject.Execute(applicationFolder, referenceFile, fileName)
				Expect(err).To(BeNil())
			})
		})

		When("Fails to read file", func() {
			It("Fails", func() {
				expectedError := errors.New("could not read file")
				mockFileService.EXPECT().Read("file_name").Return("", expectedError)
				err := subject.Execute(applicationFolder, referenceFile,"file_name")
				Expect(err).To(Equal(expectedError))
			})
		})

		When("Fails to write file", func() {
			It("Fails", func() {
				expectedError := errors.New("could not read file")
				mockFileService.EXPECT().CreateFolderIfNotExists(applicationFolder).Return(nil)
				mockFileService.EXPECT().Read(fileName).Return(content, nil)
				mockFileService.EXPECT().Write(applicationFolder, referenceFile, content).Return(expectedError)
				err := subject.Execute(applicationFolder, referenceFile,"file_name")
				Expect(err).To(Equal(expectedError))
			})
		})

		When("Fails to create folder", func() {
			It("Fails", func() {
				expectedError := errors.New("could not create folder")
				mockFileService.EXPECT().CreateFolderIfNotExists(applicationFolder).Return(expectedError)
				mockFileService.EXPECT().Read(fileName).Return(content, nil)
				err := subject.Execute(applicationFolder, referenceFile,"file_name")
				Expect(err).To(Equal(expectedError))
			})
		})

	})

})
