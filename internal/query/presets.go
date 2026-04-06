// Package query provides preset database queries for the replicator CLI.
//
// Each preset is a named SQL query that produces a styled table.
// Presets cover common observability needs: agent activity, cell status,
// swarm completion rates, and recent events.
package query

import (
	"fmt"
	"io"
	"strconv"

	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/ui"
)

// Preset names.
const (
	AgentActivity24h    = "agent_activity_24h"
	CellsByStatus       = "cells_by_status"
	SwarmCompletionRate = "swarm_completion_rate"
	RecentEvents        = "recent_events"
)

// ListPresets returns all available preset names.
func ListPresets() []string {
	return []string{
		AgentActivity24h,
		CellsByStatus,
		SwarmCompletionRate,
		RecentEvents,
	}
}

// Run executes a preset query and writes the results to w.
func Run(store *db.Store, presetName string, w io.Writer) error {
	switch presetName {
	case AgentActivity24h:
		return runAgentActivity(store, w)
	case CellsByStatus:
		return runCellsByStatus(store, w)
	case SwarmCompletionRate:
		return runSwarmCompletionRate(store, w)
	case RecentEvents:
		return runRecentEvents(store, w)
	default:
		return fmt.Errorf("unknown preset: %q (use --list to see available presets)", presetName)
	}
}

func runAgentActivity(store *db.Store, w io.Writer) error {
	rows, err := store.DB.Query(`
		SELECT COALESCE(agent_name, '(unknown)') as agent, COUNT(*) as events
		FROM events
		WHERE created_at >= datetime('now', '-24 hours')
		GROUP BY agent_name
		ORDER BY events DESC
		LIMIT 20`)
	if err != nil {
		return fmt.Errorf("query agent activity: %w", err)
	}
	defer rows.Close()

	styles := ui.NewStyles(w)
	var tableRows [][]string

	for rows.Next() {
		var agent string
		var events int
		if err := rows.Scan(&agent, &events); err != nil {
			return err
		}
		tableRows = append(tableRows, []string{agent, strconv.Itoa(events)})
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(tableRows) == 0 {
		fmt.Fprintln(w, styles.Dim.Render("(no activity in last 24 hours)"))
		return nil
	}

	t := ui.NewTable(styles, []string{"AGENT", "EVENTS (24h)"}, tableRows)
	fmt.Fprintln(w, t.String())
	return nil
}

func runCellsByStatus(store *db.Store, w io.Writer) error {
	rows, err := store.DB.Query(`
		SELECT status, type, COUNT(*) as count
		FROM beads
		GROUP BY status, type
		ORDER BY status, type`)
	if err != nil {
		return fmt.Errorf("query cells by status: %w", err)
	}
	defer rows.Close()

	styles := ui.NewStyles(w)
	var tableRows [][]string

	for rows.Next() {
		var status, cellType string
		var n int
		if err := rows.Scan(&status, &cellType, &n); err != nil {
			return err
		}
		tableRows = append(tableRows, []string{status, cellType, strconv.Itoa(n)})
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(tableRows) == 0 {
		fmt.Fprintln(w, styles.Dim.Render("(no cells)"))
		return nil
	}

	t := ui.NewTable(styles, []string{"STATUS", "TYPE", "COUNT"}, tableRows)
	fmt.Fprintln(w, t.String())
	return nil
}

func runSwarmCompletionRate(store *db.Store, w io.Writer) error {
	styles := ui.NewStyles(w)

	// Count completed vs total swarm events.
	var total, completed int
	store.DB.QueryRow(`SELECT COUNT(*) FROM events WHERE type LIKE 'swarm_%'`).Scan(&total)
	store.DB.QueryRow(`SELECT COUNT(*) FROM events WHERE type = 'swarm_complete'`).Scan(&completed)

	fmt.Fprintln(w, styles.Bold.Render("Swarm Completion Rate:"))
	fmt.Fprintf(w, "  Total swarm events:     %d\n", total)
	fmt.Fprintf(w, "  Completed:              %d\n", completed)
	if total > 0 {
		rate := float64(completed) / float64(total) * 100
		fmt.Fprintf(w, "  Completion rate:        %.1f%%\n", rate)
	} else {
		fmt.Fprintln(w, styles.Dim.Render("  Completion rate:        N/A (no swarm events)"))
	}
	return nil
}

func runRecentEvents(store *db.Store, w io.Writer) error {
	rows, err := store.DB.Query(`
		SELECT id, type, project_key, created_at
		FROM events
		ORDER BY created_at DESC
		LIMIT 20`)
	if err != nil {
		return fmt.Errorf("query recent events: %w", err)
	}
	defer rows.Close()

	styles := ui.NewStyles(w)
	var tableRows [][]string

	for rows.Next() {
		var id int
		var eventType, projectKey, createdAt string
		if err := rows.Scan(&id, &eventType, &projectKey, &createdAt); err != nil {
			return err
		}
		tableRows = append(tableRows, []string{strconv.Itoa(id), eventType, projectKey, createdAt})
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(tableRows) == 0 {
		fmt.Fprintln(w, styles.Dim.Render("(no events)"))
		return nil
	}

	t := ui.NewTable(styles, []string{"ID", "TYPE", "PROJECT", "CREATED"}, tableRows)
	fmt.Fprintln(w, t.String())
	return nil
}
