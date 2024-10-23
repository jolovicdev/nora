package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jolovicdev/nora/internal/utils"
)

type ContentStore struct {
	rootPath string
}

func NewContentStore(rootPath string) *ContentStore {
	return &ContentStore{rootPath: rootPath}
}

func (cs *ContentStore) Store(content []byte) (string, error) {
	hash := sha1.Sum(content)
	hashStr := hex.EncodeToString(hash[:])
	
	objPath := filepath.Join(cs.rootPath, "objects", hashStr[:2], hashStr[2:])
	if err := utils.CreateDirIfNotExists(filepath.Dir(objPath)); err != nil {
		return "", err
	}

	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		if err := os.WriteFile(objPath, content, 0644); err != nil {
			return "", fmt.Errorf("failed to store content: %v", err)
		}
	}

	return hashStr, nil
}

func (cs *ContentStore) Get(hash string) ([]byte, error) {
	objPath := filepath.Join(cs.rootPath, "objects", hash[:2], hash[2:])
	return os.ReadFile(objPath)
}