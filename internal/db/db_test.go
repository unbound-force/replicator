package db

import "testing"

func TestOpenMemory(t *testing.T) {
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer store.Close()

	// Verify tables exist.
	tables := []string{"events", "agents", "beads", "cell_events"}
	for _, table := range tables {
		var name string
		err := store.DB.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestOpenMemory_Idempotent(t *testing.T) {
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("first open: %v", err)
	}

	// Running migrate again should not fail.
	if err := store.migrate(); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	store.Close()
}
