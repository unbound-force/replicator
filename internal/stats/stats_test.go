package stats

import (
	"bytes"
	"strings"
	"testing"

	"github.com/unbound-force/replicator/internal/db"
)

func testStore(t *testing.T) *db.Store {
	t.Helper()
	store, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestRun_EmptyDatabase(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Replicator Stats") {
		t.Error("expected header in output")
	}
	if !strings.Contains(output, "(no events)") {
		t.Error("expected '(no events)' for empty database")
	}
	if !strings.Contains(output, "(no cells)") {
		t.Error("expected '(no cells)' for empty database")
	}
	if !strings.Contains(output, "Recent Activity (24h): 0 events") {
		t.Error("expected 0 recent events")
	}
}

func TestRun_WithEvents(t *testing.T) {
	store := testStore(t)

	// Insert some events.
	for i := 0; i < 3; i++ {
		store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES (?, '{}', 'test')`, "forge_init")
	}
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES (?, '{}', 'test')`, "forge_complete")

	var buf bytes.Buffer
	err := Run(store, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "forge_init") {
		t.Error("expected forge_init in output")
	}
	if !strings.Contains(output, "forge_complete") {
		t.Error("expected forge_complete in output")
	}
}

func TestRun_WithCells(t *testing.T) {
	store := testStore(t)

	// Insert cells with different statuses.
	store.DB.Exec(`INSERT INTO beads (id, title, status) VALUES ('c1', 'Task 1', 'open')`)
	store.DB.Exec(`INSERT INTO beads (id, title, status) VALUES ('c2', 'Task 2', 'open')`)
	store.DB.Exec(`INSERT INTO beads (id, title, status) VALUES ('c3', 'Task 3', 'closed')`)

	var buf bytes.Buffer
	err := Run(store, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "3 total") {
		t.Errorf("expected '3 total' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "open") {
		t.Error("expected 'open' status in output")
	}
	if !strings.Contains(output, "closed") {
		t.Error("expected 'closed' status in output")
	}
}

func TestRun_RecentActivity(t *testing.T) {
	store := testStore(t)

	// Insert events with current timestamp (default).
	for i := 0; i < 5; i++ {
		store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('test', '{}', 'test')`)
	}

	var buf bytes.Buffer
	err := Run(store, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Recent Activity (24h): 5 events") {
		t.Errorf("expected 5 recent events, got:\n%s", output)
	}
}

func TestRun_WritesToWriter(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
