# Research: Charm Bracelet CLI UX

**Branch**: `002-charm-ux` | **Date**: 2026-04-06

## Library Analysis

### charmbracelet/lipgloss v1.1.0

**Purpose**: Declarative terminal styling with automatic TTY/pipe/NO_COLOR
detection.

**Key API surface**:

```go
// Renderer-aware construction (REQUIRED pattern — never use lipgloss.NewStyle() globally)
renderer := lipgloss.NewRenderer(w)  // w is the output writer (os.Stdout, etc.)
style := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))

// Color profiles — detected automatically from writer
renderer.ColorProfile()  // returns termenv.TrueColor, ANSI256, ANSI, Ascii

// Rendering — writes styled string
style.Render("hello")  // returns string with ANSI codes (or plain if no color)

// Border styles
lipgloss.RoundedBorder()  // ╭╮╰╯│─
lipgloss.NormalBorder()   // ┌┐└┘│─
lipgloss.ThickBorder()    // ┏┓┗┛┃━
```

**Critical pattern — renderer-aware styles**:

The `uf doctor` format.go and `gaze` styles.go both demonstrate the
correct pattern: create a `lipgloss.NewRenderer(w)` from the output
writer, then derive all styles from that renderer. This ensures:

1. Pipe detection: when `w` is not a TTY, the renderer's color profile
   is `termenv.Ascii`, and all `Render()` calls produce plain text.
2. `NO_COLOR` respect: the renderer checks the environment variable.
3. No global state: each command creates its own renderer from its writer.

**Anti-pattern — global styles**:

```go
// WRONG: lipgloss.NewStyle() uses a global renderer that may not match the output writer
var passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
```

This fails when output is piped because the global renderer was
initialized before the pipe context was known.

### charmbracelet/lipgloss/table

**Purpose**: Structured table rendering with borders and column styling.

**Key API surface**:

```go
import "github.com/charmbracelet/lipgloss/table"

t := table.New().
    Border(lipgloss.NormalBorder()).
    BorderStyle(renderer.NewStyle().Foreground(lipgloss.Color("63"))).
    Headers("ID", "TITLE", "STATUS", "TYPE").
    Rows(rows...).                    // [][]string
    StyleFunc(func(row, col int) lipgloss.Style {
        // row 0 = header (handled by HeaderStyle)
        // row 1+ = data rows
        if col == 2 { // status column
            return statusStyle(rows[row-1][2])
        }
        return lipgloss.NewStyle()
    })

fmt.Fprintln(w, t.Render())
```

**Key behaviors**:
- `table.New()` does NOT take a renderer — it uses the global renderer
  by default. To make it pipe-safe, set `BorderStyle` from a
  renderer-aware style and use `StyleFunc` with renderer-aware styles.
- Column widths auto-calculated from content.
- `Width(n)` constrains total table width (for terminal fitting).

### charmbracelet/log v1.0.0

**Purpose**: Structured, leveled logging with human-readable output.

**Key API surface**:

```go
import "github.com/charmbracelet/log"

// Create logger with specific writer
logger := log.NewWithOptions(w, log.Options{
    Level:           log.DebugLevel,
    ReportTimestamp: true,
    ReportCaller:    false,
    Prefix:          "replicator",
})

// Structured logging
logger.Info("tool call", "tool", "hive_cells", "duration", "12ms")
logger.Error("tool failed", "tool", "hive_create", "error", err)

// Multi-writer for stderr + file
multiWriter := io.MultiWriter(os.Stderr, logFile)
logger := log.NewWithOptions(multiWriter, opts)
```

**Key behaviors**:
- Automatically styles output when writer is a TTY (colors, bold levels).
- Falls back to plain text when piped or `NO_COLOR` is set.
- Thread-safe for concurrent MCP tool calls.
- `log.Options.TimeFormat` defaults to `time.Kitchen` — consider
  `time.RFC3339` for log files.

### muesli/termenv v0.16.0

**Purpose**: Terminal environment detection. Used internally by lipgloss.

**Key constants for replicator**:

```go
import "github.com/muesli/termenv"

// Check if color is supported
if renderer.ColorProfile() == termenv.Ascii {
    // No color support — use plain text indicators
}

// Color profiles (from lowest to highest fidelity):
// termenv.Ascii      — no color (pipe, NO_COLOR, dumb terminal)
// termenv.ANSI       — 16 colors
// termenv.ANSI256    — 256 colors (our palette uses this)
// termenv.TrueColor  — 16M colors
```

## Reference Implementation Analysis

### `uf doctor` format.go Pattern

The `unbound-force/internal/doctor/format.go` is the canonical reference
for replicator's doctor output. Key patterns to replicate:

1. **Renderer from writer**: `renderer := lipgloss.NewRenderer(w)`
2. **Named styles**: `passStyle`, `warnStyle`, `failStyle`, `dimStyle`,
   `titleStyle`, `boxStyle` — all derived from the renderer.
3. **Color palette**: green(10), yellow(11), red(9), gray(241),
   pink-bold(212), purple-border(63).
4. **Plain-text fallback**: `hasColor := renderer.ColorProfile() != termenv.Ascii`
   then `[PASS]`/`[WARN]`/`[FAIL]` text indicators.
5. **Summary box**: `lipgloss.RoundedBorder()` with purple(63) border,
   emoji counters (✅, ⚠️, ❌).
6. **Contextual message**: "Everything looks good!" / "Run after fixes."

### `gaze` styles.go Pattern

The `gaze/internal/report/styles.go` demonstrates the centralized style
struct pattern:

1. **Styles struct**: Named fields for each semantic style (Header,
   Pass, Fail, Muted, Border, TableHeader, TableCell, etc.).
2. **DefaultStyles() constructor**: Returns a fully populated `Styles`
   with the color palette.
3. **Method dispatchers**: `TierStyle(tier)`, `ClassificationStyle(label)`
   — map domain values to styles.

**Difference from uf doctor**: Gaze uses `lipgloss.NewStyle()` (global
renderer) because it always writes to stdout. Replicator must use
`renderer.NewStyle()` because the MCP server uses stdout for JSON-RPC.

### Replicator-Specific Adaptation

Replicator's `internal/ui/styles.go` should combine both patterns:

```go
// Styles holds the visual theme for CLI output.
// Created via NewStyles(w) which detects TTY/pipe/NO_COLOR from the writer.
type Styles struct {
    Pass     lipgloss.Style  // green(10) — success indicators
    Warn     lipgloss.Style  // yellow(11) — warning indicators
    Fail     lipgloss.Style  // red(9) — failure indicators
    Dim      lipgloss.Style  // gray(241) — de-emphasized text
    Bold     lipgloss.Style  // bold — emphasis
    Title    lipgloss.Style  // pink(212) bold — section headers
    Box      lipgloss.Style  // purple(63) rounded border — summary boxes
    Border   lipgloss.Style  // purple(63) — table borders

    // HasColor indicates whether the output supports ANSI colors.
    // When false, use text indicators ([PASS], [WARN], [FAIL]).
    HasColor bool

    // Renderer is the lipgloss renderer for this output context.
    Renderer *lipgloss.Renderer
}
```

## Color Palette Reference

| Name | ANSI Code | Usage | Hex (approx) |
|------|-----------|-------|-------------|
| Green | 10 | Pass, success, open status | #00ff00 |
| Yellow | 11 | Warn, in_progress status | #ffff00 |
| Red | 9 | Fail, blocked status | #ff0000 |
| Gray | 241 | Dim, closed status, hints | #626262 |
| Pink | 212 | Title, bold headers | #ff87d7 |
| Purple | 63 | Borders, table frames | #5f5fff |

This palette matches the Unbound Force ecosystem (uf doctor, gaze) per
FR-012.

## Existing Output Audit

Commands that currently use raw `fmt.Printf` and need migration:

| Command | Current Pattern | Migration Target |
|---------|----------------|-----------------|
| `doctor` | `fmt.Printf` with `\u2713`/`\u2717` Unicode | `doctor.FormatText` with lipgloss |
| `cells` | `json.MarshalIndent` dump | `ui.Table` with status coloring |
| `setup` | `fmt.Printf` with `\u2713`/`\u2717` | `ui.Styles.Pass`/`Fail` indicators |
| `init` | `fmt.Println` plain text | `ui.Styles.Pass`/`Dim` messages |
| `version` | `fmt.Printf` plain text | Bold version, dim commit/date |
| `stats` | `fmt.Fprintf` with `===` headers | `ui.Styles.Title` + `ui.Table` |
| `query` | `fmt.Fprintf` with `---` separators | `ui.Table` for tabular output |
| `query --list` | `fmt.Println` plain list | `ui.Styles.Bold` + `Dim` |
| `docs` | Markdown output | UNCHANGED (markdown, not terminal) |
| `serve` | No user-facing output (stdio JSON-RPC) | Add logging only |

## Testing Strategy

### Output Capture Pattern

All formatters accept `io.Writer`. Tests use `bytes.Buffer`:

```go
func TestFormatText_PassResult(t *testing.T) {
    var buf bytes.Buffer
    results := []doctor.CheckResult{{Name: "git", Status: "pass", Message: "ok"}}
    if err := doctor.FormatText(results, &buf); err != nil {
        t.Fatalf("FormatText: %v", err)
    }
    // Buffer is non-TTY → plain text fallback
    if !strings.Contains(buf.String(), "[PASS]") {
        t.Errorf("expected [PASS] indicator in non-TTY output")
    }
}
```

### ANSI Detection Pattern

To verify no ANSI codes in piped output:

```go
func hasANSI(s string) bool {
    return strings.Contains(s, "\x1b[")
}

func TestFormatText_NoANSI_WhenPiped(t *testing.T) {
    var buf bytes.Buffer
    // bytes.Buffer is not a TTY → renderer uses Ascii profile
    doctor.FormatText(results, &buf)
    if hasANSI(buf.String()) {
        t.Error("piped output contains ANSI escape codes")
    }
}
```

### Log File Testing Pattern

```go
func TestServeLogging_CreatesLogFile(t *testing.T) {
    dir := t.TempDir()
    logPath := filepath.Join(dir, ".unbound-force", "replicator.log")
    // ... start server with dir as working directory
    // ... send tool call
    // ... verify log file exists and contains tool name
}
```

## MCP Server Logging Architecture

```text
┌─────────────────────────────────────────────┐
│ replicator serve                            │
│                                             │
│  stdin ──→ [JSON-RPC parser] ──→ stdout     │
│                 │                           │
│                 ▼                           │
│         [tool handler]                      │
│                 │                           │
│                 ▼                           │
│         [charmbracelet/log]                 │
│              │       │                      │
│              ▼       ▼                      │
│          stderr    .unbound-force/          │
│                    replicator.log           │
└─────────────────────────────────────────────┘
```

**Critical constraint**: The MCP server uses stdout exclusively for
JSON-RPC responses. Logging MUST NOT write to stdout. The logger writes
to `io.MultiWriter(os.Stderr, logFile)`.

## Dependencies to Add

```bash
go get github.com/charmbracelet/lipgloss@v1.1.0
go get github.com/charmbracelet/log@v1.0.0
go get github.com/muesli/termenv@v0.16.0
```

Note: `lipgloss/table` is a sub-package of `lipgloss` — no separate
`go get` needed. Import as `github.com/charmbracelet/lipgloss/table`.
