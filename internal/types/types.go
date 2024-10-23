package types

import "os"

type Snapshot struct {
	ID        string            `json:"id"`
	Timestamp int64            `json:"timestamp"`
	Message   string            `json:"message"`
	Files     map[string]string `json:"files"`
	Parent    string            `json:"parent"`
}

type FileChange struct {
    Path string
    State string
}

type Timeline struct {
	Name      string   `json:"name"`
	Current   string   `json:"current"`
	Snapshots []string `json:"snapshots"`
}
type FileStatus struct {
    Path     string
    State    string
    Mode     os.FileMode
    IsSymlink bool
}

type Config struct {
	CurrentTimeline string            `json:"current_timeline"`
	Timelines      map[string]string `json:"timelines"`
}

type Status struct {
	CurrentTimeline string            `json:"current_timeline"`
	Timelines       map[string]string `json:"timelines"`
}

type DiffStep struct {
	Type     string
	Content  string
	Position int
}
