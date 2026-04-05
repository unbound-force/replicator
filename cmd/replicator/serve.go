package main

import (
	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/mcp"
	"github.com/unbound-force/replicator/internal/memory"
	"github.com/unbound-force/replicator/internal/tools/hive"
	memorytools "github.com/unbound-force/replicator/internal/tools/memory"
	"github.com/unbound-force/replicator/internal/tools/registry"
	swarmtools "github.com/unbound-force/replicator/internal/tools/swarm"
	swarmmailtools "github.com/unbound-force/replicator/internal/tools/swarmmail"
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
	swarmmailtools.Register(reg, store)
	swarmtools.Register(reg, store)

	// Memory tools proxy to Dewey for semantic search.
	memClient := memory.NewClient(cfg.DeweyURL)
	memorytools.Register(reg, memClient)

	server := mcp.NewServer(reg, Version)
	return server.ServeStdio()
}
