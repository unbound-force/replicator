package main

import (
	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/mcp"
	"github.com/unbound-force/replicator/internal/tools/hive"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

// serveMCP starts the MCP JSON-RPC server on stdio.
func serveMCP() error {
	cfg := config.Load()

	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer store.Close()

	reg := registry.New()
	hive.Register(reg, store)

	server := mcp.NewServer(reg)
	return server.ServeStdio()
}
