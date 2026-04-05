package hive

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/unbound-force/replicator/internal/db"
)

// SessionStart creates a new session and returns the previous session's
// handoff notes (if any). This enables context preservation across sessions.
func SessionStart(store *db.Store, activeCellID string) (string, error) {
	// Find the most recent ended session's handoff notes.
	// Use rowid for deterministic ordering when timestamps collide.
	var prevNotes string
	err := store.DB.QueryRow(
		"SELECT handoff_notes FROM sessions WHERE ended_at IS NOT NULL ORDER BY rowid DESC LIMIT 1",
	).Scan(&prevNotes)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("query previous session: %w", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return "", fmt.Errorf("generate session ID: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	var cellID *string
	if activeCellID != "" {
		cellID = &activeCellID
	}

	_, err = store.DB.Exec(
		"INSERT INTO sessions (session_id, started_at, active_cell_id) VALUES (?, ?, ?)",
		sessionID, now, cellID,
	)
	if err != nil {
		return "", fmt.Errorf("insert session: %w", err)
	}

	return prevNotes, nil
}

// SessionEnd ends the current (most recent un-ended) session with handoff notes.
func SessionEnd(store *db.Store, handoffNotes string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := store.DB.Exec(`
		UPDATE sessions SET ended_at = ?, handoff_notes = ?
		WHERE session_id = (
			SELECT session_id FROM sessions WHERE ended_at IS NULL ORDER BY started_at DESC LIMIT 1
		)`,
		now, handoffNotes,
	)
	if err != nil {
		return fmt.Errorf("end session: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no active session to end")
	}
	return nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "sess-" + hex.EncodeToString(b), nil
}
