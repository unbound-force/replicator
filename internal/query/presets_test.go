package query

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

func TestListPresets(t *testing.T) {
	presets := ListPresets()
	if len(presets) != 4 {
		t.Fatalf("expected 4 presets, got %d", len(presets))
	}

	expected := map[string]bool{
		AgentActivity24h:    true,
		CellsByStatus:       true,
		ForgeCompletionRate: true,
		RecentEvents:        true,
	}
	for _, p := range presets {
		if !expected[p] {
			t.Errorf("unexpected preset: %q", p)
		}
	}
}

func TestRun_UnknownPreset(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, "nonexistent", &buf)
	if err == nil {
		t.Fatal("expected error for unknown preset")
	}
	if !strings.Contains(err.Error(), "unknown preset") {
		t.Errorf("error = %q, want 'unknown preset'", err.Error())
	}
}

func TestRun_AgentActivity_Empty(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, AgentActivity24h, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "no activity") {
		t.Error("expected empty message")
	}
}

func TestRun_AgentActivity_WithData(t *testing.T) {
	store := testStore(t)

	// Insert events with agent_name in payload.
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_init', '{"agent_name": "worker-1"}', 'test')`)
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_init', '{"agent_name": "worker-1"}', 'test')`)
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_init', '{"agent_name": "worker-2"}', 'test')`)

	var buf bytes.Buffer
	err := Run(store, AgentActivity24h, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "worker-1") {
		t.Error("expected worker-1 in output")
	}
}

func TestRun_CellsByStatus_Empty(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, CellsByStatus, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "no cells") {
		t.Error("expected empty message")
	}
}

func TestRun_CellsByStatus_WithData(t *testing.T) {
	store := testStore(t)

	store.DB.Exec(`INSERT INTO beads (id, title, status, type) VALUES ('c1', 'Task 1', 'open', 'task')`)
	store.DB.Exec(`INSERT INTO beads (id, title, status, type) VALUES ('c2', 'Task 2', 'closed', 'task')`)
	store.DB.Exec(`INSERT INTO beads (id, title, status, type) VALUES ('c3', 'Bug 1', 'open', 'bug')`)

	var buf bytes.Buffer
	err := Run(store, CellsByStatus, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "open") {
		t.Error("expected 'open' in output")
	}
	if !strings.Contains(output, "closed") {
		t.Error("expected 'closed' in output")
	}
}

func TestRun_ForgeCompletionRate_Empty(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, ForgeCompletionRate, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Forge Completion Rate") {
		t.Error("expected header")
	}
	if !strings.Contains(output, "N/A") {
		t.Error("expected N/A for empty data")
	}
}

func TestRun_ForgeCompletionRate_WithData(t *testing.T) {
	store := testStore(t)

	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_init', '{}', 'test')`)
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_progress', '{}', 'test')`)
	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('forge_complete', '{}', 'test')`)

	var buf bytes.Buffer
	err := Run(store, ForgeCompletionRate, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Total forge events") {
		t.Error("expected total count")
	}
	if !strings.Contains(output, "Completed") {
		t.Error("expected completed count")
	}
}

func TestRun_RecentEvents_Empty(t *testing.T) {
	store := testStore(t)
	var buf bytes.Buffer

	err := Run(store, RecentEvents, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "no events") {
		t.Error("expected empty message")
	}
}

func TestRun_RecentEvents_WithData(t *testing.T) {
	store := testStore(t)

	store.DB.Exec(`INSERT INTO events (type, payload, project_key) VALUES ('test_event', '{}', 'my-project')`)

	var buf bytes.Buffer
	err := Run(store, RecentEvents, &buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test_event") {
		t.Error("expected test_event in output")
	}
	if !strings.Contains(output, "my-project") {
		t.Error("expected my-project in output")
	}
}

func TestRun_AllPresets(t *testing.T) {
	store := testStore(t)

	// Verify all presets run without error on an empty database.
	for _, preset := range ListPresets() {
		t.Run(preset, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Run(store, preset, &buf); err != nil {
				t.Fatalf("Run(%q): %v", preset, err)
			}
			if buf.Len() == 0 {
				t.Errorf("preset %q produced no output", preset)
			}
		})
	}
}
