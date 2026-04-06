package comms

import "testing"

func TestReserve(t *testing.T) {
	store := testStore(t)

	reservations, err := Reserve(store, "worker-1", []string{"foo.go", "bar.go"}, true, "implementing feature", 300)
	if err != nil {
		t.Fatalf("Reserve: %v", err)
	}
	if len(reservations) != 2 {
		t.Fatalf("expected 2 reservations, got %d", len(reservations))
	}
	if reservations[0].Path != "foo.go" {
		t.Errorf("path = %q, want %q", reservations[0].Path, "foo.go")
	}
	if reservations[0].AgentName != "worker-1" {
		t.Errorf("agent_name = %q, want %q", reservations[0].AgentName, "worker-1")
	}
	if !reservations[0].Exclusive {
		t.Error("exclusive should be true")
	}
	if reservations[0].TTLSeconds != 300 {
		t.Errorf("ttl_seconds = %d, want 300", reservations[0].TTLSeconds)
	}
}

func TestReserve_DefaultTTL(t *testing.T) {
	store := testStore(t)

	reservations, err := Reserve(store, "worker-1", []string{"foo.go"}, true, "test", 0)
	if err != nil {
		t.Fatalf("Reserve: %v", err)
	}
	if reservations[0].TTLSeconds != 300 {
		t.Errorf("default ttl = %d, want 300", reservations[0].TTLSeconds)
	}
}

func TestReserve_ExclusiveConflict(t *testing.T) {
	store := testStore(t)

	// Worker 1 reserves foo.go exclusively.
	_, err := Reserve(store, "worker-1", []string{"foo.go"}, true, "first", 300)
	if err != nil {
		t.Fatalf("first Reserve: %v", err)
	}

	// Worker 2 tries to reserve foo.go -- should fail.
	_, err = Reserve(store, "worker-2", []string{"foo.go"}, true, "second", 300)
	if err == nil {
		t.Error("expected conflict error for exclusive reservation")
	}
}

func TestReserve_SameAgentNoConflict(t *testing.T) {
	store := testStore(t)

	// Same agent can reserve the same path again.
	_, err := Reserve(store, "worker-1", []string{"foo.go"}, true, "first", 300)
	if err != nil {
		t.Fatalf("first Reserve: %v", err)
	}

	_, err = Reserve(store, "worker-1", []string{"foo.go"}, true, "second", 300)
	if err != nil {
		t.Fatalf("same agent should not conflict: %v", err)
	}
}

func TestReserve_NonExclusiveNoConflict(t *testing.T) {
	store := testStore(t)

	// Non-exclusive reservation should not block others.
	_, err := Reserve(store, "worker-1", []string{"foo.go"}, false, "reading", 300)
	if err != nil {
		t.Fatalf("first Reserve: %v", err)
	}

	_, err = Reserve(store, "worker-2", []string{"foo.go"}, false, "also reading", 300)
	if err != nil {
		t.Fatalf("non-exclusive should not conflict: %v", err)
	}
}

func TestRelease_ByPath(t *testing.T) {
	store := testStore(t)

	Reserve(store, "worker-1", []string{"foo.go"}, true, "test", 300)

	if err := Release(store, []string{"foo.go"}, nil); err != nil {
		t.Fatalf("Release: %v", err)
	}

	// Should be able to reserve again by another agent.
	_, err := Reserve(store, "worker-2", []string{"foo.go"}, true, "after release", 300)
	if err != nil {
		t.Fatalf("reserve after release: %v", err)
	}
}

func TestRelease_ByID(t *testing.T) {
	store := testStore(t)

	reservations, _ := Reserve(store, "worker-1", []string{"foo.go"}, true, "test", 300)

	if err := Release(store, nil, []int{reservations[0].ID}); err != nil {
		t.Fatalf("Release by ID: %v", err)
	}

	// Should be able to reserve again.
	_, err := Reserve(store, "worker-2", []string{"foo.go"}, true, "after release", 300)
	if err != nil {
		t.Fatalf("reserve after release by ID: %v", err)
	}
}

func TestRelease_NothingSpecified(t *testing.T) {
	store := testStore(t)

	err := Release(store, nil, nil)
	if err == nil {
		t.Error("expected error when nothing specified")
	}
}

func TestReleaseAll(t *testing.T) {
	store := testStore(t)

	Reserve(store, "worker-1", []string{"a.go"}, true, "test", 300)
	Reserve(store, "worker-2", []string{"b.go"}, true, "test", 300)

	if err := ReleaseAll(store, ""); err != nil {
		t.Fatalf("ReleaseAll: %v", err)
	}

	// Both should be available now.
	_, err := Reserve(store, "worker-3", []string{"a.go", "b.go"}, true, "after release all", 300)
	if err != nil {
		t.Fatalf("reserve after release all: %v", err)
	}
}

func TestReleaseAgent(t *testing.T) {
	store := testStore(t)

	Reserve(store, "worker-1", []string{"a.go", "b.go"}, true, "test", 300)
	Reserve(store, "worker-2", []string{"c.go"}, true, "test", 300)

	if err := ReleaseAgent(store, "worker-1"); err != nil {
		t.Fatalf("ReleaseAgent: %v", err)
	}

	// worker-1's paths should be available.
	_, err := Reserve(store, "worker-3", []string{"a.go"}, true, "after agent release", 300)
	if err != nil {
		t.Fatalf("reserve a.go after agent release: %v", err)
	}

	// worker-2's path should still be reserved.
	_, err = Reserve(store, "worker-3", []string{"c.go"}, true, "conflict", 300)
	if err == nil {
		t.Error("expected conflict -- worker-2's reservation should still exist")
	}
}
