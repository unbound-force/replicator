package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
)

// runSetup creates the config directory, initializes the database, and
// verifies git is available.
func runSetup() error {
	// 1. Create config directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("determine home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "swarm-tools")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	fmt.Printf("\u2713 Config directory: %s\n", configDir)

	// 2. Initialize database.
	cfg := config.Load()
	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	store.Close()
	fmt.Printf("\u2713 Database: %s\n", cfg.DatabasePath)

	// 3. Verify git.
	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("\u2717 Git: not found (%v)\n", err)
		fmt.Println("  Install git: https://git-scm.com/downloads")
	} else {
		fmt.Printf("\u2713 Git: %s\n", strings.TrimSpace(string(out)))
	}

	fmt.Println()
	fmt.Println("Setup complete. Run 'replicator doctor' to verify all checks pass.")
	return nil
}
