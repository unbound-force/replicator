package parity

import (
	"fmt"
	"io"
)

// ToolResult captures the parity test outcome for a single MCP tool.
type ToolResult struct {
	Name        string
	Match       bool
	Differences []Difference
}

// GenerateReport writes a human-readable parity report to w.
//
// Output format:
//
//	Parity Report
//	=============
//	hive_cells      ✓
//	hive_create     ✓
//	hivemind_find   ✗  $.content[0].text: expected object, got string
//	---
//	22/23 tools match (95.7%)
func GenerateReport(results []ToolResult, w io.Writer) {
	fmt.Fprintln(w, "Parity Report")
	fmt.Fprintln(w, "=============")

	matched := 0
	for _, r := range results {
		if r.Match {
			fmt.Fprintf(w, "%-30s ✓\n", r.Name)
			matched++
		} else {
			// Show first difference inline, rest on subsequent lines.
			for i, d := range r.Differences {
				if i == 0 {
					fmt.Fprintf(w, "%-30s ✗  %s: expected %s, got %s\n",
						r.Name, d.Path, d.ExpectedType, d.ActualType)
				} else {
					fmt.Fprintf(w, "%-30s    %s: expected %s, got %s\n",
						"", d.Path, d.ExpectedType, d.ActualType)
				}
			}
		}
	}

	fmt.Fprintln(w, "---")

	total := len(results)
	if total == 0 {
		fmt.Fprintln(w, "0/0 tools match (0.0%)")
		return
	}

	pct := float64(matched) / float64(total) * 100
	fmt.Fprintf(w, "%d/%d tools match (%.1f%%)\n", matched, total, pct)
}
