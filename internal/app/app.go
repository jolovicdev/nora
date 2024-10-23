package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jolovicdev/nora/internal/core/diff"
	"github.com/jolovicdev/nora/internal/core/snapshot"
	"github.com/jolovicdev/nora/internal/core/storage"
	"github.com/jolovicdev/nora/internal/core/timeline"
	"github.com/jolovicdev/nora/internal/types"
	"github.com/jolovicdev/nora/internal/utils"
)

const(
	Reset = "\033[0m" 
	Red = "\033[31m" 
	Green = "\033[32m" 
	Yellow = "\033[33m" 
	Blue = "\033[34m" 
	Magenta = "\033[35m" 
	Cyan = "\033[36m" 
	Gray = "\033[37m" 
	White = "\033[97m"
)


type App struct {
    contentStore *storage.ContentStore
    index       *storage.Index
    snapshots   *snapshot.Store
    timelines   *timeline.Manager
}
func (app *App) PrepareFiles(paths []string) error {
    prepared := make(map[string]string)
    
    existing, err := app.index.GetPreparedFiles()
    if err == nil {
        for k, v := range existing {
            prepared[k] = v
        }
    }
    ignorePatterns, err := app.loadIgnorePatterns()
    if err != nil {
        return fmt.Errorf("failed to load ignore patterns: %w", err)
    }

    for _, path := range paths {
        if path == "." {
            err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }


                if strings.HasPrefix(path, ".nora") {
                    return filepath.SkipDir
                }


                if shouldIgnore(path, ignorePatterns) {
                    if info.IsDir() {
                        return filepath.SkipDir
                    }
                    return nil
                }

                if !info.IsDir() {
                    status, err := app.getFileStatus(path, info)
                    if err != nil {
                        return fmt.Errorf("failed to get status for %s: %w", path, err)
                    }

                    if status.State != "unchanged" {
                        if err := app.prepareFile(path, prepared); err != nil {
                            return fmt.Errorf("failed to prepare %s: %w", path, err)
                        }
                    }
                }
                return nil
            })
            if err != nil {
                return fmt.Errorf("failed to walk directory: %w", err)
            }
        } else {

            _, err := os.Lstat(path)
            if err != nil {
                return fmt.Errorf("failed to stat %s: %w", path, err)
            }

            if shouldIgnore(path, ignorePatterns) {
                continue
            }

            if err := app.prepareFile(path, prepared); err != nil {
                return fmt.Errorf("failed to prepare %s: %w", path, err)
            }
        }
    }

    return app.index.PrepareFiles(prepared)
}

func (app *App) prepareFile(path string, prepared map[string]string) error {
    info, err := os.Lstat(path)
    if err != nil {
        return fmt.Errorf("failed to stat %s: %w", path, err)
    }

    var content []byte
    if info.Mode()&os.ModeSymlink != 0 {

        target, err := os.Readlink(path)
        if err != nil {
            return fmt.Errorf("failed to read symlink %s: %w", path, err)
        }
        content = []byte(target)
    } else {

        content, err = os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("failed to read file %s: %w", path, err)
        }
    }
    hash, err := app.contentStore.Store(content)
    if err != nil {
        return fmt.Errorf("failed to store content for %s: %w", path, err)
    }

    prepared[path] = hash
    fmt.Printf("Prepared: %s\n", path)
    return nil
}
func (app *App) GetStatus() error {

    timeline, err := app.timelines.GetCurrent()
    if err != nil {
        return fmt.Errorf("failed to get current timeline: %w", err)
    }


    prepared, err := app.index.GetPreparedFiles()
    if err != nil {
        prepared = make(map[string]string)
    }


    var snapshotFiles map[string]string
    if timeline.Current != "" {
        snapshot, err := app.snapshots.Get(timeline.Current)
        if err != nil {
            return fmt.Errorf("failed to get current snapshot: %w", err)
        }
        snapshotFiles = snapshot.Files
    } else {
        snapshotFiles = make(map[string]string)
    }


    changes := make(map[string]types.FileChange)


    err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }


        if strings.HasPrefix(path, ".nora") || shouldIgnore(path, []string{}) {
            if info.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }


        if info.IsDir() {
            return nil
        }


        content, err := os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("failed to read file %s: %w", path, err)
        }

        currentHash, err := app.contentStore.Store(content)
        if err != nil {
            return fmt.Errorf("failed to hash content for %s: %w", path, err)
        }


        preparedHash, isPrepared := prepared[path]
        

        snapshotHash, inSnapshot := snapshotFiles[path]

        var state string
        switch {
        case isPrepared:
            if inSnapshot {
                if preparedHash != snapshotHash {
                    state = "modified (prepared)"
                } else {
                    state = "unchanged"
                }
            } else {
                state = "added (prepared)"
            }
        case inSnapshot:
            if currentHash != snapshotHash {
                state = "modified"
            } else {
                state = "unchanged"
            }
        default:
            state = "untracked"
        }

        if state != "unchanged" {
            changes[path] = types.FileChange{Path: path, State: state}
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("failed to walk directory: %w", err)
    }


    fmt.Printf("\nOn timeline: %s\n", timeline.Name)
    

    hasPrepared := false
    fmt.Println("\nChanges prepared for snapshot:")
    for path, change := range changes {
        if strings.Contains(change.State, "prepared") {
            hasPrepared = true
            fmt.Printf("%s%s: %s%s\n", Green, path, change.State, Reset)
        }
    }
    if !hasPrepared {
        fmt.Println("  no changes prepared")
    }


    hasUnprepared := false
    fmt.Println("\nChanges not prepared for snapshot:")
    for path, change := range changes {
        if !strings.Contains(change.State, "prepared") && change.State != "unchanged" {
            hasUnprepared = true
            switch change.State {
            case "modified":
                fmt.Printf("%s%s: %s%s\n", Red, path, change.State, Reset)
            case "untracked":
                fmt.Printf("%s%s: %s%s\n", Blue, path, change.State, Reset)
            }
        }
    }
    if !hasUnprepared {
        fmt.Println("  working directory clean")
    }

    fmt.Println()
    return nil
}

func (app *App) getFileStatus(path string, info os.FileInfo) (types.FileStatus, error) {
    status := types.FileStatus{
        Path:     path,
        Mode:     info.Mode(),
        IsSymlink: info.Mode()&os.ModeSymlink != 0,
    }


    timeline, err := app.timelines.GetCurrent()
    if err != nil {
        return status, fmt.Errorf("failed to get current timeline: %w", err)
    }


    if timeline.Current == "" {
        status.State = "added"
        return status, nil
    }


    snapshot, err := app.snapshots.Get(timeline.Current)
    if err != nil {
        return status, fmt.Errorf("failed to get current snapshot: %w", err)
    }


    oldHash, exists := snapshot.Files[path]
    if !exists {
        status.State = "added"
        return status, nil
    }


    content, err := os.ReadFile(path)
    if err != nil {
        return status, fmt.Errorf("failed to read file: %w", err)
    }


    newHash, err := app.contentStore.Store(content)
    if err != nil {
        return status, fmt.Errorf("failed to hash content: %w", err)
    }

    if newHash != oldHash {
        status.State = "modified"
    } else {
        status.State = "unchanged"
    }

    return status, nil
}


func (app *App) loadIgnorePatterns() ([]string, error) {
    patterns := []string{
        ".nora",
        ".noraignore",
        ".git",
        ".svn",
        "node_modules",
        "*.tmp",
        "*.swp",
	  "nora",
    }

    data, err := os.ReadFile(".noraignore")
    if err == nil {
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line != "" && !strings.HasPrefix(line, "#") {
                patterns = append(patterns, line)
            }
        }
    }

    return patterns, nil
}


func shouldIgnore(path string, patterns []string) bool {
    for _, pattern := range patterns {
        matched, err := filepath.Match(pattern, filepath.Base(path))
        if err == nil && matched {
            return true
        }

        if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern) {
            return true
        }
    }
    return false
}

func (app *App) CreateSnapshot(message string) error {
    prepared, err := app.index.GetPreparedFiles()
    if err != nil {
        return fmt.Errorf("failed to get prepared files: %v", err)
    }

    if len(prepared) == 0 {
        return fmt.Errorf("no files prepared for snapshot")
    }

    timeline, err := app.timelines.GetCurrent()
    if err != nil {
        return fmt.Errorf("failed to get current timeline: %v", err)
    }

    snap, err := app.snapshots.Create(message, prepared, timeline.Current)
    if err != nil {
        return fmt.Errorf("failed to create snapshot: %v", err)
    }

    timeline.Current = snap.ID
    timeline.Snapshots = append(timeline.Snapshots, snap.ID)
    if err := app.timelines.Update(timeline); err != nil {
        return fmt.Errorf("failed to update timeline: %v", err)
    }

    if err := app.index.PrepareFiles(make(map[string]string)); err != nil {
        return fmt.Errorf("failed to clear prepared files: %v", err)
    }

    fmt.Printf("Created snapshot: %s\n", snap.ID)
    return nil
}

func (app *App) RecallSnapshot(id string) error {
    snapshot, err := app.snapshots.Get(id)
    if err != nil {
        return fmt.Errorf("failed to get snapshot: %v", err)
    }

    fmt.Printf("Snapshot: %s\n", snapshot.ID)
    fmt.Printf("Message: %s\n", snapshot.Message)
    fmt.Printf("Files:\n")
    
    for path, hash := range snapshot.Files {
        content, err := app.contentStore.Get(hash)
        if err != nil {
            fmt.Printf("  %s: [error reading content: %v]\n", path, err)
            continue
        }
        fmt.Printf("  %s: %d bytes\n", path, len(content))
    }

    return nil
}

func (app *App) ShowDiff(path string) error {

    prepared, err := app.index.GetPreparedFiles()
    if err != nil {
        return fmt.Errorf("failed to get prepared files: %w", err)
    }


    newHash, exists := prepared[path]
    if !exists {
        return fmt.Errorf("file not prepared: %s", path)
    }


    timeline, err := app.timelines.GetCurrent()
    if err != nil {
        return fmt.Errorf("failed to get current timeline: %w", err)
    }


    if timeline.Current == "" {

        newContent, err := app.contentStore.Get(newHash)
        if err != nil {
            return fmt.Errorf("failed to get new content: %w", err)
        }
        
        lines := strings.Split(string(newContent), "\n")
        for _, line := range lines {
            if line != "" {
                fmt.Printf("%s+ %s%s\n", Green, line, Reset)
            }
        }
        return nil
    }


    snapshot, err := app.snapshots.Get(timeline.Current)
    if err != nil {
        return fmt.Errorf("failed to get current snapshot: %w", err)
    }


    oldHash, exists := snapshot.Files[path]
    if !exists {

        newContent, err := app.contentStore.Get(newHash)
        if err != nil {
            return fmt.Errorf("failed to get new content: %w", err)
        }
        
        lines := strings.Split(string(newContent), "\n")
        for _, line := range lines {
            if line != "" {
                fmt.Printf("%s+ %s%s\n", Green, line, Reset)
            }
        }
        return nil
    }


    oldContent, err := app.contentStore.Get(oldHash)
    if err != nil {
        return fmt.Errorf("failed to get old content: %w", err)
    }

    newContent, err := app.contentStore.Get(newHash)
    if err != nil {
        return fmt.Errorf("failed to get new content: %w", err)
    }


    oldLines := strings.Split(string(oldContent), "\n")
    newLines := strings.Split(string(newContent), "\n")


    steps := diff.SimpleMyers(oldLines, newLines)
    if steps == nil {
        return fmt.Errorf("failed to calculate diff")
    }


    var oldIdx, newIdx int
    for _, step := range steps {
        switch step.Type {
        case "keep":
            if oldIdx < len(oldLines) {
                fmt.Printf("  %s\n", oldLines[oldIdx])
            }
            oldIdx++
            newIdx++
        case "delete":
            if oldIdx < len(oldLines) {
                fmt.Printf("%s- %s%s\n", Red, oldLines[oldIdx], Reset)
            }
            oldIdx++
        case "add":
            if newIdx < len(newLines) {
                fmt.Printf("%s+ %s%s\n", Green, newLines[newIdx], Reset)
            }
            newIdx++
        }
    }

    return nil
}

type Step struct {
    Type string
}
func (app *App) Forget(files []string) error {
	existing, err := app.index.GetPreparedFiles()
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		app.index.ForgetFiles(files)
	}
	return nil
}
func (app *App) Initialize() error {
    dirs := []string{
        ".nora",
        ".nora/memories",
        ".nora/timelines",
        ".nora/index",
        ".nora/config",
        ".nora/meta",
        ".nora/objects",
    }

    for _, dir := range dirs {
        if err := utils.CreateDirIfNotExists(dir); err != nil {
            return fmt.Errorf("failed to create directory %s: %v", dir, err)
        }
    }

    if err := app.timelines.Create("main"); err != nil {
        return err
    }


    cwd, err := os.Getwd()
    if err != nil {
        cwd = "current directory"
    }

    utils.PrintInitMessage(cwd)

    return nil
}
func New(rootPath string) *App {
    return &App{
        contentStore: storage.NewContentStore(rootPath),
        index:       storage.NewIndex(rootPath),
        snapshots:   snapshot.NewStore(rootPath),
        timelines:   timeline.NewManager(rootPath),
    }
}
