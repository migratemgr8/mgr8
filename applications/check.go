package applications

import (
	"errors"
	"log"

	"github.com/kenji-yamane/mgr8/infrastructure"
)

type CheckCommand interface {
	Execute(string) error
}

type checkCommand struct {
	fileService infrastructure.FileService
}

var ErrFilesDoNotMatch = errors.New("reference and schema dont match")

func NewCheckCommand(fileService infrastructure.FileService) *checkCommand {
	return &checkCommand{fileService: fileService}
}

func (g *checkCommand) Execute(referenceFile, initialFile string) error {
	schemaContent, err := g.fileService.Read(initialFile)
	if err != nil {
		return err
	}

	referenceContent, err := g.fileService.Read(referenceFile)
	if err != nil {
		return err
	}

	if schemaContent != referenceContent {
		return ErrFilesDoNotMatch
	}
	log.Print("Files match")

	return nil
}
