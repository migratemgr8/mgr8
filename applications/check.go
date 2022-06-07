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

func NewCheckCommand(fileService infrastructure.FileService) *checkCommand {
	return &checkCommand{fileService: fileService}
}

func (g *checkCommand) Execute(initialFile string) error {
	schemaContent, err := g.fileService.Read(initialFile)
	if err != nil {
		return err
	}

	referenceContent, err := g.fileService.Read(".mgr8/reference.sql")
	if err != nil {
		return err
	}

	if schemaContent != referenceContent {
		return errors.New("reference and schema dont match")
	}
	log.Print("Files match")

	return nil
}
