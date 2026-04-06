package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// NewTable creates a consistently styled lipgloss table.
//
// lipgloss/table uses the global renderer internally; StyleFunc is used
// to apply renderer-aware styles per-cell. The Border style uses the
// Styles.Border color for visual consistency across all CLI commands.
func NewTable(s *Styles, headers []string, rows [][]string) *table.Table {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(s.Border).
		Headers(headers...).
		Rows(rows...)
	return t
}
