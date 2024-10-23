package timeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jolovicdev/nora/internal/types"
)

type Manager struct {
    rootPath string
    config   *types.Config
}

func NewManager(rootPath string) *Manager {
    return &Manager{
        rootPath: rootPath,
        config: &types.Config{
            CurrentTimeline: "main",
            Timelines:      make(map[string]string),
        },
    }
}

func (m *Manager) loadConfig() error {
    configPath := filepath.Join(m.rootPath, ".nora", "config", "config.json")
    

    configDir := filepath.Dir(configPath)
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        if os.IsNotExist(err) {

            m.config = &types.Config{
                CurrentTimeline: "main",
                Timelines: map[string]string{"main": "main"},
            }
            return m.saveConfig()
        }
        return fmt.Errorf("failed to read config: %w", err)
    }

    if err := json.Unmarshal(data, &m.config); err != nil {
        return fmt.Errorf("failed to parse config: %w", err)
    }


    if m.config.Timelines == nil {
        m.config.Timelines = make(map[string]string)
    }

    return nil
}

func (m *Manager) saveConfig() error {
    configPath := filepath.Join(m.rootPath, ".nora", "config", "config.json")
    
    data, err := json.MarshalIndent(m.config, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    if err := os.WriteFile(configPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    return nil
}

func (m *Manager) Create(name string) error {
    timeline := types.Timeline{
        Name:      name,
        Current:   "",
        Snapshots: []string{},
    }


    timelinePath := filepath.Join(m.rootPath, ".nora", "timelines", name+".json")
    

    timelineDir := filepath.Dir(timelinePath)
    if err := os.MkdirAll(timelineDir, 0755); err != nil {
        return fmt.Errorf("failed to create timelines directory: %w", err)
    }

    data, err := json.MarshalIndent(timeline, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal timeline: %w", err)
    }

    if err := os.WriteFile(timelinePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write timeline: %w", err)
    }


    if err := m.loadConfig(); err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    m.config.Timelines[name] = name
    m.config.CurrentTimeline = name

    return m.saveConfig()
}

func (m *Manager) GetCurrent() (*types.Timeline, error) {
    if err := m.loadConfig(); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    if m.config.CurrentTimeline == "" {
        return nil, fmt.Errorf("no current timeline set")
    }

    timelinePath := filepath.Join(m.rootPath, ".nora", "timelines", m.config.CurrentTimeline+".json")
    data, err := os.ReadFile(timelinePath)
    if err != nil {
        if os.IsNotExist(err) {

            timeline := &types.Timeline{
                Name:      m.config.CurrentTimeline,
                Current:   "",
                Snapshots: []string{},
            }
            return timeline, m.Update(timeline)
        }
        return nil, fmt.Errorf("failed to read timeline %s: %w", m.config.CurrentTimeline, err)
    }

    var timeline types.Timeline
    if err := json.Unmarshal(data, &timeline); err != nil {
        return nil, fmt.Errorf("failed to parse timeline: %w", err)
    }

    return &timeline, nil
}

func (m *Manager) Update(timeline *types.Timeline) error {
    timelinePath := filepath.Join(m.rootPath, ".nora", "timelines", timeline.Name+".json")
    
    data, err := json.MarshalIndent(timeline, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal timeline: %w", err)
    }

    if err := os.WriteFile(timelinePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write timeline: %w", err)
    }

    return nil
}