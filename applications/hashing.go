package applications

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/kenji-yamane/mgr8/domain"
	"github.com/kenji-yamane/mgr8/infrastructure"
)

type HashService interface {
	GetSqlHash(sqlFilePath string) (string, error)
	ValidateFileMigration(version int, filePath string, driver domain.Driver) (bool, error)
}

type hashService struct {
	fileService infrastructure.FileService
}

func NewHashService(fileService infrastructure.FileService) *hashService {
	return &hashService{fileService: fileService}
}

func (h *hashService) GetSqlHash(sqlFilePath string) (string, error) {
	content, err := h.fileService.Read(sqlFilePath) // ioutil close file after reading
	if err != nil {
		return "", err
	}
	hash_md5 := md5.Sum([]byte(content))
	string_hash_md5 := hex.EncodeToString(hash_md5[:])
	return string_hash_md5, nil
}

func (h *hashService) ValidateFileMigration(version int, filePath string, driver domain.Driver) (bool, error) {
	hashFile, err := h.GetSqlHash(filePath)
	if err != nil {
		return false, err
	}

	hashDb, err := driver.GetVersionHashing(version)
	if err != nil {
		return false, err
	}

	return hashFile == hashDb, nil
}
