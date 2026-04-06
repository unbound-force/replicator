// Package ui provides shared terminal styling for the replicator CLI.
//
// All CLI commands use NewStyles to get a renderer-aware style set that
// automatically degrades to plain text when output is piped or NO_COLOR
// is set. This follows the UF doctor formatting pattern.
package ui

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Styles holds renderer-aware lipgloss styles for consistent CLI output.
// Created via NewStyles with an io.Writer so the renderer can detect
// color support (TTY vs pipe, NO_COLOR, etc.).
type Styles struct {
	Pass     lipgloss.Style
	Warn     lipgloss.Style
	Fail     lipgloss.Style
	Dim      lipgloss.Style
	Bold     lipgloss.Style
	Title    lipgloss.Style
	Box      lipgloss.Style
	Border   lipgloss.Style
	HasColor bool
	Renderer *lipgloss.Renderer
}

// NewStyles creates a Styles set bound to the given writer.
// Color detection is automatic — a bytes.Buffer or pipe gets plain text,
// a real TTY gets full color.
func NewStyles(w io.Writer) *Styles {
	r := lipgloss.NewRenderer(w)
	hasColor := r.ColorProfile() != termenv.Ascii

	return &Styles{
		Pass:     r.NewStyle().Foreground(lipgloss.Color("10")),
		Warn:     r.NewStyle().Foreground(lipgloss.Color("11")),
		Fail:     r.NewStyle().Foreground(lipgloss.Color("9")),
		Dim:      r.NewStyle().Foreground(lipgloss.Color("241")),
		Bold:     r.NewStyle().Bold(true),
		Title:    r.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		Box:      r.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1),
		Border:   r.NewStyle().Foreground(lipgloss.Color("63")),
		HasColor: hasColor,
		Renderer: r,
	}
}

// Indicator returns a status icon appropriate for the color mode.
// In color mode, returns styled emoji. In plain mode, returns bracketed text.
func (s *Styles) Indicator(status string) string {
	switch status {
	case "pass":
		if s.HasColor {
			return s.Pass.Render("✅")
		}
		return "[PASS]"
	case "warn":
		if s.HasColor {
			return s.Warn.Render("⚠️")
		}
		return "[WARN]"
	case "fail":
		if s.HasColor {
			return s.Fail.Render("❌")
		}
		return "[FAIL]"
	default:
		return status
	}
}
