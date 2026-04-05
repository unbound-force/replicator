// Package stats provides database statistics for the replicator CLI.
//
// Queries the events and cells tables to produce a human-readable summary
// of system activity and work item status.
package stats

import (
	"fmt"
	"io"

	"github.com/unbound-force/replicator/internal/db"
)

// eventCount holds a type and its count from the events table.
type eventCount struct {
	Type  string
	Count int
}

// Run queries the database for statistics and writes a formatted report.
func Run(store *db.Store, w io.Writer) error {
	// Events by type.
	eventCounts, err := queryEventCounts(store)
	if err != nil {
		return fmt.Errorf("query event counts: %w", err)
	}

	// Recent events (last 24h).
	recentCount, err := queryRecentEvents(store)
	if err != nil {
		return fmt.Errorf("query recent events: %w", err)
	}

	// Cells by status.
	cellCounts, err := queryCellCounts(store)
	if err != nil {
		return fmt.Errorf("query cell counts: %w", err)
	}

	// Total cells.
	totalCells := 0
	for _, c := range cellCounts {
		totalCells += c.Count
	}

	// Print report.
	fmt.Fprintln(w, "=== Replicator Stats ===")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Events by Type:")
	if len(eventCounts) == 0 {
		fmt.Fprintln(w, "  (no events)")
	}
	for _, ec := range eventCounts {
		fmt.Fprintf(w, "  %-30s %d\n", ec.Type, ec.Count)
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Recent Activity (24h): %d events\n", recentCount)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Cells (%d total):\n", totalCells)
	if len(cellCounts) == 0 {
		fmt.Fprintln(w, "  (no cells)")
	}
	for _, cc := range cellCounts {
		fmt.Fprintf(w, "  %-15s %d\n", cc.Type, cc.Count)
	}

	return nil
}

func queryEventCounts(store *db.Store) ([]eventCount, error) {
	rows, err := store.DB.Query(`
		SELECT type, COUNT(*) as count
		FROM events
		GROUP BY type
		ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []eventCount
	for rows.Next() {
		var ec eventCount
		if err := rows.Scan(&ec.Type, &ec.Count); err != nil {
			return nil, err
		}
		counts = append(counts, ec)
	}
	return counts, rows.Err()
}

func queryRecentEvents(store *db.Store) (int, error) {
	var count int
	err := store.DB.QueryRow(`
		SELECT COUNT(*)
		FROM events
		WHERE created_at >= datetime('now', '-24 hours')`).Scan(&count)
	return count, err
}

func queryCellCounts(store *db.Store) ([]eventCount, error) {
	rows, err := store.DB.Query(`
		SELECT status, COUNT(*) as count
		FROM beads
		GROUP BY status
		ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []eventCount
	for rows.Next() {
		var ec eventCount
		if err := rows.Scan(&ec.Type, &ec.Count); err != nil {
			return nil, err
		}
		counts = append(counts, ec)
	}
	return counts, rows.Err()
}
