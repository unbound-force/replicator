package doctor

import (
	"fmt"
	"io"
	"time"

	"github.com/unbound-force/replicator/internal/ui"
)

// FormatText renders health check results as styled terminal output.
//
// Uses lipgloss for styling with automatic NO_COLOR and pipe detection
// following the UF doctor formatting pattern. When output is not a TTY,
// falls back to plain-text indicators ([PASS], [WARN], [FAIL]).
func FormatText(results []CheckResult, w io.Writer) error {
	styles := ui.NewStyles(w)

	// Header with stethoscope emoji.
	fmt.Fprintln(w, styles.Title.Render("🩺 Replicator Doctor"))
	fmt.Fprintln(w)

	// Tally counters for the summary box.
	var passed, warned, failed int

	for _, r := range results {
		indicator := styles.Indicator(r.Status)
		name := fmt.Sprintf("%-14s", r.Name)
		duration := styles.Dim.Render(fmt.Sprintf("(%s)", r.Duration.Round(time.Millisecond)))

		fmt.Fprintf(w, "  %s %s %s %s\n", indicator, name, r.Message, duration)

		switch r.Status {
		case "pass":
			passed++
		case "warn":
			warned++
		case "fail":
			failed++
		}
	}

	fmt.Fprintln(w)

	// Boxed summary with emoji counters.
	summaryContent := fmt.Sprintf("  ✅ %d passed  ⚠️  %d warnings  ❌ %d failed",
		passed, warned, failed)
	fmt.Fprintln(w, styles.Box.Render(summaryContent))

	// Contextual completion message.
	if failed == 0 && warned == 0 {
		fmt.Fprintln(w, styles.Pass.Render("🎉 Everything looks good!"))
	} else if failed > 0 {
		fmt.Fprintln(w, styles.Dim.Render("  Run 'replicator setup' to fix common issues."))
	} else {
		fmt.Fprintln(w, styles.Dim.Render("  All critical checks passed."))
	}

	return nil
}
