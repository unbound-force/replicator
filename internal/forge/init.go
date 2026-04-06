// Package forge implements multi-agent coordination for parallel task execution.
//
// The forge package provides session initialization, task decomposition (prompt
// generation), subtask spawning, progress tracking, worktree management, code
// review prompts, and historical insights. Most functions are either prompt
// generators (returning strings) or event recorders (writing to the events table).
package forge

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/unbound-force/replicator/internal/db"
)

// Init initializes a forge session by recording an event and returning session info.
// The isolation parameter specifies the isolation strategy ("worktree" or "reservation").
func Init(store *db.Store, projectPath, isolation string) (map[string]any, error) {
	if projectPath == "" {
		return nil, fmt.Errorf("project_path is required")
	}
	if isolation == "" {
		isolation = "reservation"
	}

	payload := map[string]string{
		"project_path": projectPath,
		"isolation":    isolation,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	_, err = store.DB.Exec(
		"INSERT INTO events (type, payload, project_key) VALUES (?, ?, ?)",
		"forge_init", string(payloadJSON), projectPath,
	)
	if err != nil {
		return nil, fmt.Errorf("record forge_init event: %w", err)
	}

	return map[string]any{
		"status":       "initialized",
		"project_path": projectPath,
		"isolation":    isolation,
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
	}, nil
}
