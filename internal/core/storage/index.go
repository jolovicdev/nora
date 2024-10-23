package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Index struct {
	rootPath string
}

func NewIndex(rootPath string) *Index {
	return &Index{rootPath: rootPath}
}

func (idx *Index) PrepareFiles(files map[string]string) error {
	data, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(idx.rootPath, "index", "prepared.json"), data, 0644)
}

func (idx *Index) ForgetFiles(paths []string) error {
    prepared := make(map[string]string)
    
    data, err := os.ReadFile(filepath.Join(idx.rootPath, "index", "prepared.json"))
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to read prepared.json: %w", err)
    }
        if len(data) > 0 {
        if err := json.Unmarshal(data, &prepared); err != nil {
            return fmt.Errorf("failed to unmarshal prepared files: %w", err)
        }
    }
        for _, path := range paths {
        delete(prepared, path)
    }
    
    updatedData, err := json.Marshal(prepared)
    if err != nil {
        return fmt.Errorf("failed to marshal updated prepared files: %w", err)
    }
    indexDir := filepath.Join(idx.rootPath, "index")
    if err := os.MkdirAll(indexDir, 0755); err != nil {
        return fmt.Errorf("failed to create index directory: %w", err)
    }
    

    if err := os.WriteFile(filepath.Join(indexDir, "prepared.json"), updatedData, 0644); err != nil {
        return fmt.Errorf("failed to write prepared.json: %w", err)
    }
    
    return nil
}

func (idx *Index) GetPreparedFiles() (map[string]string, error) {
	prepared := make(map[string]string)
	data, err := os.ReadFile(filepath.Join(idx.rootPath, "index", "prepared.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return prepared, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, &prepared)
	return prepared, err
}
