package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/query"
)

// runQuery executes a preset query and prints results.
func runQuery(cfg *config.Config, presetName string) error {
	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()

	return query.Run(store, presetName, os.Stdout)
}

// listQueryPresets prints available preset names.
func listQueryPresets() {
	fmt.Println("Available query presets:")
	for _, p := range query.ListPresets() {
		fmt.Printf("  %s\n", p)
	}
	fmt.Println()
	fmt.Println("Usage: replicator query <preset>")
	fmt.Printf("Example: replicator query %s\n", strings.Join(query.ListPresets()[:1], ""))
}
