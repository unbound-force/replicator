// Package swarmmail provides agent messaging and file reservation for swarm coordination.
//
// Agents register themselves, send messages to each other, and reserve files
// for exclusive editing to prevent conflicts during parallel work.
package comms

import (
	"fmt"
	"time"

	"github.com/unbound-force/replicator/internal/db"
)

// Agent represents a registered swarm agent.
//
// TaskDescription is always serialized (no omitempty) to maintain shape
// parity with the TypeScript cyborg-swarm responses.
type Agent struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ProjectPath     string `json:"project_path"`
	TaskDescription string `json:"task_description"`
	Status          string `json:"status"`
}

// Init registers or updates an agent in the swarm.
// Uses upsert semantics -- if the agent already exists, updates its metadata.
func Init(store *db.Store, agentName, projectPath, taskDescription string) (*Agent, error) {
	if agentName == "" {
		return nil, fmt.Errorf("agent name is required")
	}

	now := time.Now().UTC().Format(time.RFC3339)

	// Upsert: insert or update on conflict.
	_, err := store.DB.Exec(`
		INSERT INTO agents (name, project_path, task_description, last_seen_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			project_path = excluded.project_path,
			task_description = excluded.task_description,
			last_seen_at = excluded.last_seen_at`,
		agentName, projectPath, taskDescription, now,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert agent: %w", err)
	}

	// Read back the agent.
	var agent Agent
	err = store.DB.QueryRow(
		"SELECT id, name, project_path, COALESCE(task_description, ''), status FROM agents WHERE name = ?",
		agentName,
	).Scan(&agent.ID, &agent.Name, &agent.ProjectPath, &agent.TaskDescription, &agent.Status)
	if err != nil {
		return nil, fmt.Errorf("read agent: %w", err)
	}

	return &agent, nil
}
