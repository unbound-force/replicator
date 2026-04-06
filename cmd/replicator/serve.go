package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	charmlog "github.com/charmbracelet/log"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/mcp"
	"github.com/unbound-force/replicator/internal/memory"
	commstools "github.com/unbound-force/replicator/internal/tools/comms"
	forgetools "github.com/unbound-force/replicator/internal/tools/forge"
	memorytools "github.com/unbound-force/replicator/internal/tools/memory"
	"github.com/unbound-force/replicator/internal/tools/org"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

// serveMCP starts the MCP JSON-RPC server on stdio.
func serveMCP() error {
	cfg := config.Load()

	// Set up structured logging to file (and stderr).
	// Bootstrap exception: use fmt.Fprintf for errors before the logger exists.
	logger, logCloser := setupLogger()
	if logCloser != nil {
		defer logCloser.Close()
	}

	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer store.Close()

	reg := registry.New()
	org.Register(reg, store)
	commstools.Register(reg, store)
	forgetools.Register(reg, store)

	// Memory tools proxy to Dewey for semantic search.
	memClient := memory.NewClient(cfg.DeweyURL)
	memorytools.Register(reg, memClient)

	server := mcp.NewServer(reg, Version, logger)
	return server.ServeStdio()
}

// setupLogger creates a charmbracelet/log logger that writes to both
// stderr and .uf/replicator/replicator.log. If the log file cannot be
// created, logging falls back to stderr only and a warning is printed.
// The returned io.Closer should be deferred by the caller; it may be nil.
func setupLogger() (*charmlog.Logger, io.Closer) {
	logDir := filepath.Join(".", ".uf", "replicator")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot create log directory: %v (logging to stderr only)\n", err)
		return charmlog.NewWithOptions(os.Stderr, charmlog.Options{
			ReportTimestamp: true,
			Level:           charmlog.InfoLevel,
		}), nil
	}

	logFile, err := os.Create(filepath.Join(logDir, "replicator.log"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot create log file: %v (logging to stderr only)\n", err)
		return charmlog.NewWithOptions(os.Stderr, charmlog.Options{
			ReportTimestamp: true,
			Level:           charmlog.InfoLevel,
		}), nil
	}

	logWriter := io.MultiWriter(os.Stderr, logFile)
	logger := charmlog.NewWithOptions(logWriter, charmlog.Options{
		ReportTimestamp: true,
		Level:           charmlog.InfoLevel,
	})
	return logger, logFile
}
