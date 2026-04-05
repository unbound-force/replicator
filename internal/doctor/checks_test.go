package doctor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unbound-force/replicator/internal/config"
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

func TestRun_AllChecks(t *testing.T) {
	store := testStore(t)

	// Mock Dewey as healthy.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		DeweyURL: srv.URL,
	}

	results, err := Run(store, cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(results) != 4 {
		t.Fatalf("expected 4 checks, got %d", len(results))
	}

	// Verify check names.
	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	for _, expected := range []string{"git", "database", "dewey", "config_dir"} {
		if !names[expected] {
			t.Errorf("missing check: %s", expected)
		}
	}
}

func TestCheckGit(t *testing.T) {
	if testing.Short() {
		t.Skip("requires git")
	}

	result := checkGit()
	if result.Name != "git" {
		t.Errorf("name = %q, want %q", result.Name, "git")
	}
	if result.Status != "pass" {
		t.Errorf("status = %q, want %q (message: %s)", result.Status, "pass", result.Message)
	}
	if result.Duration <= 0 {
		t.Error("duration should be positive")
	}
}

func TestCheckDatabase_Healthy(t *testing.T) {
	store := testStore(t)

	result := checkDatabase(store)
	if result.Status != "pass" {
		t.Errorf("status = %q, want %q (message: %s)", result.Status, "pass", result.Message)
	}
}

func TestCheckDatabase_Closed(t *testing.T) {
	store := testStore(t)
	store.Close()

	result := checkDatabase(store)
	if result.Status != "fail" {
		t.Errorf("status = %q, want %q for closed database", result.Status, "fail")
	}
}

func TestCheckDewey_Healthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	result := checkDewey(srv.URL)
	if result.Status != "pass" {
		t.Errorf("status = %q, want %q", result.Status, "pass")
	}
}

func TestCheckDewey_Unreachable(t *testing.T) {
	result := checkDewey("http://127.0.0.1:1")
	if result.Status != "warn" {
		t.Errorf("status = %q, want %q for unreachable Dewey", result.Status, "warn")
	}
}

func TestCheckDewey_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	result := checkDewey(srv.URL)
	if result.Status != "warn" {
		t.Errorf("status = %q, want %q for HTTP 500", result.Status, "warn")
	}
}

func TestCheckConfigDir(t *testing.T) {
	result := checkConfigDir()
	// The config dir may or may not exist in CI, but the check should not panic.
	if result.Name != "config_dir" {
		t.Errorf("name = %q, want %q", result.Name, "config_dir")
	}
	if result.Status != "pass" && result.Status != "fail" {
		t.Errorf("status = %q, want pass or fail", result.Status)
	}
	if result.Duration <= 0 {
		t.Error("duration should be positive")
	}
}

func TestCheckResult_StatusValues(t *testing.T) {
	// Verify that all results use valid status values.
	store := testStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := &config.Config{DeweyURL: srv.URL}
	results, _ := Run(store, cfg)

	validStatuses := map[string]bool{"pass": true, "fail": true, "warn": true}
	for _, r := range results {
		if !validStatuses[r.Status] {
			t.Errorf("check %q has invalid status %q", r.Name, r.Status)
		}
	}
}
