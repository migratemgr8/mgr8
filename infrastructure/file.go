package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type MigrationFile struct {
	FullPath string
	Name     string
}

type FileService interface {
	Write(migrationDir, filename, content string) error
	Read(filename string) (string, error)
	List(fileDirectory string) ([]MigrationFile, error)
}

type fileService struct{}

func NewFileService() *fileService {
	return &fileService{}
}

func (f *fileService) Write(migrationDir, filename, content string) error {
	file, err := os.Create(path.Join(migrationDir, filename))
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

func (f *fileService) List(dir string) ([]MigrationFile, error) {
	var migrationFiles []MigrationFile
	{
	}

	dirPath, err := filepath.Abs(dir)
	if err != nil {
		return []MigrationFile{}, err
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return []MigrationFile{}, err
	}

	for _, fileInfo := range fileInfos {
		var migrationFile MigrationFile
		migrationFile.Name = fileInfo.Name()
		migrationFile.FullPath = filepath.Join(dirPath, fileInfo.Name())
		migrationFiles = append(migrationFiles, migrationFile)
	}

	return migrationFiles, err
}
