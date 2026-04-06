package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit_FreshDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	hiveDir := filepath.Join(dir, ".uf", "replicator")
	info, err := os.Stat(hiveDir)
	if err != nil {
		t.Fatalf(".uf/replicator/ not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal(".uf/replicator/ is not a directory")
	}

	cellsPath := filepath.Join(hiveDir, "cells.json")
	data, err := os.ReadFile(cellsPath)
	if err != nil {
		t.Fatalf("cells.json not created: %v", err)
	}
	if string(data) != "[]\n" {
		t.Errorf("cells.json content = %q, want %q", string(data), "[]\n")
	}
}

func TestRunInit_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()

	// First init.
	if err := runInit(dir); err != nil {
		t.Fatalf("first runInit: %v", err)
	}

	// Write something to cells.json to verify it's not overwritten.
	cellsPath := filepath.Join(dir, ".uf", "replicator", "cells.json")
	os.WriteFile(cellsPath, []byte(`[{"id":"test"}]`), 0o644)

	// Second init — should be idempotent.
	if err := runInit(dir); err != nil {
		t.Fatalf("second runInit: %v", err)
	}

	// Verify cells.json was NOT overwritten.
	data, _ := os.ReadFile(cellsPath)
	if string(data) != `[{"id":"test"}]` {
		t.Errorf("cells.json was overwritten: got %q", string(data))
	}
}

func TestRunInit_CustomPath(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "myproject")
	os.MkdirAll(target, 0o755)

	if err := runInit(target); err != nil {
		t.Fatalf("runInit with custom path: %v", err)
	}

	cellsPath := filepath.Join(target, ".uf", "replicator", "cells.json")
	if _, err := os.Stat(cellsPath); err != nil {
		t.Fatalf("cells.json not created at custom path: %v", err)
	}
}

func TestRunInit_InvalidPath(t *testing.T) {
	err := runInit("/nonexistent/path/that/cannot/exist")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
