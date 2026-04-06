package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unbound-force/replicator/internal/ui"
)

func initCmd() *cobra.Command {
	var pathFlag string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a project directory for swarm operations",
		Long: `Creates a .uf/replicator/ directory with an empty cells.json in the target
directory. Idempotent — safe to run multiple times.

This is the per-repo initialization command. It does not require the
global database (replicator setup) or any external services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(pathFlag)
		},
	}
	cmd.Flags().StringVar(&pathFlag, "path", ".", "Target directory for .uf/replicator/ initialization")
	return cmd
}

// runInit creates the .uf/replicator/ directory and seeds cells.json.
// Uses styled output: green for success, dim for already-initialized.
func runInit(targetDir string) error {
	styles := ui.NewStyles(os.Stdout)
	hiveDir := filepath.Join(targetDir, ".uf", "replicator")

	// Check if already initialized.
	if info, err := os.Stat(hiveDir); err == nil && info.IsDir() {
		fmt.Println(styles.Dim.Render("already initialized"))
		return nil
	}

	// Create .uf/replicator/ directory.
	if err := os.MkdirAll(hiveDir, 0o755); err != nil {
		return fmt.Errorf("create .uf/replicator directory: %w", err)
	}

	// Write empty cells.json.
	cellsPath := filepath.Join(hiveDir, "cells.json")
	if err := os.WriteFile(cellsPath, []byte("[]\n"), 0o644); err != nil {
		return fmt.Errorf("write cells.json: %w", err)
	}

	fmt.Println(styles.Pass.Render("initialized .uf/replicator/"))
	return nil
}
