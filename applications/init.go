package applications

import (
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type InitCommand interface {
	Execute(string) error
}

type initCommand struct {
	fileService infrastructure.FileService
}

func NewInitCommand(fileService infrastructure.FileService) *initCommand {
	return &initCommand{fileService: fileService}
}

func (g *initCommand) Execute(initialFile string) error {
	content, err := g.fileService.Read(initialFile)
	if err != nil {
		return err
	}

	err = g.fileService.CreateFolderIfNotExists(".mgr8")
	if err != nil {
		return err
	}

	return g.fileService.Write(".mgr8", "reference.sql", content)
}
