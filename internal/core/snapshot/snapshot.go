package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jolovicdev/nora/internal/types"
	"github.com/jolovicdev/nora/internal/utils"
)

type Store struct {
	rootPath string
}

func NewStore(rootPath string) *Store {
	return &Store{rootPath: rootPath}
}

func (s *Store) Create(message string, files map[string]string, parent string) (*types.Snapshot, error) {
	snapshot := &types.Snapshot{
		ID:        utils.GenerateID(),
		Timestamp: time.Now().Unix(),
		Message:   message,
		Files:     files,
		Parent:    parent,
	}

	if err := s.Save(snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (s *Store) Save(snapshot *types.Snapshot) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.rootPath, "memories", snapshot.ID+".json"), data, 0644)
}

func (s *Store) Get(id string) (*types.Snapshot, error) {
	data, err := os.ReadFile(filepath.Join(s.rootPath, "memories", id+".json"))
	if err != nil {
		return nil, err
	}

	var snapshot types.Snapshot
	err = json.Unmarshal(data, &snapshot)
	return &snapshot, err
}