package main

import (
	"fmt"
	"time"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/doctor"
)

// runDoctor executes health checks and prints results as a table.
func runDoctor(cfg *config.Config) error {
	store, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()

	results, err := doctor.Run(store, cfg)
	if err != nil {
		return fmt.Errorf("run checks: %w", err)
	}

	// Print header.
	fmt.Printf("%-12s %-6s %s\n", "CHECK", "STATUS", "MESSAGE")
	fmt.Printf("%-12s %-6s %s\n", "-----", "------", "-------")

	hasFailure := false
	for _, r := range results {
		icon := statusIcon(r.Status)
		fmt.Printf("%-12s %s %-4s %s (%s)\n", r.Name, icon, r.Status, r.Message, r.Duration.Round(time.Millisecond))
		if r.Status == "fail" {
			hasFailure = true
		}
	}

	if hasFailure {
		fmt.Println("\nSome checks failed. Run 'replicator setup' to fix common issues.")
	} else {
		fmt.Println("\nAll checks passed.")
	}

	return nil
}

func statusIcon(status string) string {
	switch status {
	case "pass":
		return "\u2713" // checkmark
	case "fail":
		return "\u2717" // X mark
	case "warn":
		return "!" // warning
	default:
		return "?"
	}
}
