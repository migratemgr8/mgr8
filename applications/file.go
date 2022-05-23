package applications

import (
	"fmt"
	"os"
)

type FileService interface {
	Write(filename, content string) error
	Read(filename string) (string, error)
}

type fileService struct {
}

func NewFileService() *fileService {
	return &fileService{}
}

func (f *fileService) Write(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func (f *fileService) Read(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read from file with path: %s", err)
	}
	return string(content), nil
}
