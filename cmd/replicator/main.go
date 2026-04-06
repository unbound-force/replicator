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
	"github.com/unbound-force/replicator/internal/ui"
)

// Build-time variables set via ldflags.
var (
	Version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

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
	root.AddCommand(doctorCmd())
	root.AddCommand(statsCmd())
	root.AddCommand(queryCmd())
	root.AddCommand(setupCmd())
	root.AddCommand(initCmd())
	root.AddCommand(docsCmd())

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
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON instead of a styled table")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			styles := ui.NewStyles(os.Stdout)
			fmt.Printf("replicator %s\n", styles.Bold.Render(Version))
			if commit != "unknown" {
				fmt.Printf("  commit: %s\n", styles.Dim.Render(commit))
			}
			if date != "unknown" {
				fmt.Printf("  built:  %s\n", styles.Dim.Render(date))
			}
		},
	}
}

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run health checks on the replicator environment",
		Long:  "Checks git, database, Dewey connectivity, and config directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			return runDoctor(cfg)
		},
	}
}

func statsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show database statistics",
		Long:  "Displays event counts, recent activity, and cell status summary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			return runStats(cfg)
		},
	}
}

func queryCmd() *cobra.Command {
	var listFlag bool

	cmd := &cobra.Command{
		Use:   "query [preset]",
		Short: "Run a preset database query",
		Long:  "Executes a named preset query against the database. Use --list to see available presets.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag || len(args) == 0 {
				listQueryPresets()
				return nil
			}
			cfg := config.Load()
			return runQuery(cfg, args[0])
		},
	}
	cmd.Flags().BoolVar(&listFlag, "list", false, "List available query presets")
	return cmd
}

func setupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Initialize replicator environment",
		Long:  "Creates config directory, initializes database, and verifies git.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetup()
		},
	}
}
