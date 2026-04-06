package main

import (
	"fmt"
	"os"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/query"
	"github.com/unbound-force/replicator/internal/ui"
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

// listQueryPresets prints available preset names with styled output.
func listQueryPresets() {
	styles := ui.NewStyles(os.Stdout)

	fmt.Println(styles.Bold.Render("Available query presets:"))
	for _, p := range query.ListPresets() {
		fmt.Printf("  %s\n", p)
	}
	fmt.Println()
	fmt.Println(styles.Dim.Render("Usage: replicator query <preset>"))
	fmt.Printf("%s replicator query %s\n",
		styles.Dim.Render("Example:"),
		query.ListPresets()[0])
}
