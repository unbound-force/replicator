package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupLogger_CreatesLogFile(t *testing.T) {
	// Run setupLogger in a temp directory so it creates
	// .uf/replicator/replicator.log there.
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	logger, closer := setupLogger()
	if closer != nil {
		defer closer.Close()
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	logPath := filepath.Join(dir, ".uf", "replicator", "replicator.log")
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
}

func TestSetupLogger_Truncates(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	// First call: write a marker to the log file.
	logDir := filepath.Join(dir, ".uf", "replicator")
	os.MkdirAll(logDir, 0o755)
	logPath := filepath.Join(logDir, "replicator.log")
	os.WriteFile(logPath, []byte("MARKER_SHOULD_BE_GONE"), 0o644)

	// Second call: setupLogger uses os.Create which truncates.
	_, closer := setupLogger()
	if closer != nil {
		closer.Close()
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) == "MARKER_SHOULD_BE_GONE" {
		t.Error("log file was not truncated; marker still present")
	}
}

func TestSetupLogger_ReadOnlyDir(t *testing.T) {
	// Verify that a read-only directory doesn't cause a panic.
	// setupLogger should fall back to stderr-only logging.
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	// Create the .uf/replicator dir as read-only so file creation fails.
	logDir := filepath.Join(dir, ".uf", "replicator")
	os.MkdirAll(logDir, 0o755)
	os.Chmod(logDir, 0o444)
	t.Cleanup(func() {
		os.Chmod(logDir, 0o755) // restore so TempDir cleanup works
	})

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	// Should not panic — falls back to stderr-only.
	logger, closer := setupLogger()
	if closer != nil {
		defer closer.Close()
	}
	if logger == nil {
		t.Fatal("expected non-nil logger even when log file creation fails")
	}
}
