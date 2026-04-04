// Package main is the entry point for the replicator CLI.
//
// Replicator provides multi-agent coordination for AI coding agents.
// It exposes tools via the MCP (Model Context Protocol) JSON-RPC interface
// and a CLI for setup, querying, and observability.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unbound-force/replicator/internal/config"
)

// Version is set at build time via ldflags.
var Version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "replicator",
		Short: "Multi-agent coordination for AI coding agents",
		Long: `Replicator provides multi-agent coordination tools for AI coding agents.

It exposes tools via the MCP protocol (JSON-RPC over stdio) and a CLI
for setup, querying, and observability.

Tools include: hive (work items), swarm mail (messaging), swarm
(orchestration), worktrees (isolation), and memory (semantic search
via Dewey).`,
		Version: Version,
	}

	root.AddCommand(serveCmd())
	root.AddCommand(cellsCmd())
	root.AddCommand(versionCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start MCP JSON-RPC server on stdio",
		Long:  "Starts the MCP server that AI coding agents connect to for tool access.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serveMCP()
		},
	}
}

func cellsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "List hive cells (work items)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			return listCells(cfg)
		},
	}
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("replicator %s\n", Version)
		},
	}
}
