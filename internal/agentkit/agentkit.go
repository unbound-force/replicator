// Package agentkit provides embedded agent kit content for scaffolding
// new project directories. The kit includes command definitions, skill
// files, and agent role descriptions that are written to .opencode/
// during `replicator init`.
package agentkit

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed content/*
var content embed.FS

// ScaffoldResult describes the outcome of writing a single agent kit file.
type ScaffoldResult struct {
	Path   string `json:"path"`
	Action string `json:"action"` // "created", "skipped", "overwritten"
}

// Scaffold writes the embedded agent kit files to targetDir/.opencode/.
// If force is false, existing files are skipped. If force is true,
// existing files are overwritten.
func Scaffold(targetDir string, force bool) ([]ScaffoldResult, error) {
	var results []ScaffoldResult

	err := fs.WalkDir(content, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Strip "content/" prefix to get relative path under .opencode/.
		relPath, _ := filepath.Rel("content", path)
		destPath := filepath.Join(targetDir, ".opencode", relPath)

		// Check if file exists.
		if _, statErr := os.Stat(destPath); statErr == nil {
			if !force {
				results = append(results, ScaffoldResult{Path: relPath, Action: "skipped"})
				return nil
			}
			results = append(results, ScaffoldResult{Path: relPath, Action: "overwritten"})
		} else {
			results = append(results, ScaffoldResult{Path: relPath, Action: "created"})
		}

		// Create parent directories.
		if mkErr := os.MkdirAll(filepath.Dir(destPath), 0o755); mkErr != nil {
			return fmt.Errorf("create directory for %s: %w", relPath, mkErr)
		}

		// Read from embedded FS and write to disk.
		data, readErr := content.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read embedded %s: %w", path, readErr)
		}

		return os.WriteFile(destPath, data, 0o644)
	})

	return results, err
}
