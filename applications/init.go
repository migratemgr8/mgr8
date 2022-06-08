package applications

import (
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type InitCommand interface {
	Execute(string, string, string) error
}

type initCommand struct {
	fileService infrastructure.FileService
}

func NewInitCommand(fileService infrastructure.FileService) *initCommand {
	return &initCommand{fileService: fileService}
}

func (g *initCommand) Execute(applicationFolder, referenceFile, initialFile string) error {
	content, err := g.fileService.Read(initialFile)
	if err != nil {
		return err
	}

	err = g.fileService.CreateFolderIfNotExists(applicationFolder)
	if err != nil {
		return err
	}

	return g.fileService.Write(applicationFolder, referenceFile, content)
}
