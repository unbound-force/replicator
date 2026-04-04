package main

import (
	"encoding/json"
	"fmt"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/hive"
)

// listCells queries and prints hive cells as JSON.
func listCells(cfg *config.Config) error {
	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()

	cells, err := hive.QueryCells(store, hive.CellQuery{})
	if err != nil {
		return fmt.Errorf("query cells: %w", err)
	}

	out, err := json.MarshalIndent(cells, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
