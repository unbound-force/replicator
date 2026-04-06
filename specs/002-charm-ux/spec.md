# Feature Specification: Charm Bracelet CLI UX

**Feature Branch**: `002-charm-ux`  
**Created**: 2026-04-06  
**Status**: Ready  
**Input**: User description: "Upgrade CLI UX with Charm Bracelet libraries (lipgloss styling, lipgloss/table, charmbracelet/log, per-repo log file)"  
**References**: [unbound-force doctor format](https://github.com/unbound-force/unbound-force/blob/main/internal/doctor/format.go), [gaze styles](https://github.com/unbound-force/gaze/blob/main/internal/report/styles.go)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Styled Doctor Output (Priority: P1)

A developer runs `replicator doctor` and sees a visually clear, color-coded health report. Pass results appear in green, warnings in yellow, failures in red. A rounded summary box shows the total counts. When the output is piped to a file or another command, colors are automatically stripped and plain-text indicators (`[PASS]`, `[WARN]`, `[FAIL]`) are used instead.

**Why this priority**: Doctor is the first command developers run after installation. Its visual quality sets the tone for the entire tool. The `uf doctor` output is the benchmark the ecosystem expects.

**Independent Test**: Run `replicator doctor` in a terminal and verify colored output with a summary box. Run `replicator doctor | cat` and verify clean plain-text fallback with no escape codes.

**Acceptance Scenarios**:

1. **Given** a terminal with color support, **When** `replicator doctor` is run, **Then** pass results display in green, warnings in yellow, failures in red, and a rounded summary box appears at the bottom.
2. **Given** output is piped (not a TTY), **When** `replicator doctor | cat` is run, **Then** results use `[PASS]`, `[WARN]`, `[FAIL]` text indicators with no ANSI escape codes.
3. **Given** the `NO_COLOR` environment variable is set, **When** `replicator doctor` is run, **Then** plain-text indicators are used regardless of terminal capabilities.

---

### User Story 2 - Styled Cells Table (Priority: P1)

A developer runs `replicator cells` and sees a formatted table with borders, styled headers, and color-coded status badges. Open cells appear in green, in-progress in yellow, blocked in red, and closed in gray. The table fits within the terminal width.

**Why this priority**: Cells is the most frequently used CLI command for inspecting work items. A styled table makes scanning many cells fast and intuitive.

**Independent Test**: Create 4 cells with different statuses, run `replicator cells`, and verify the table renders with correct colors per status. Pipe output to verify clean fallback.

**Acceptance Scenarios**:

1. **Given** cells exist with varying statuses, **When** `replicator cells` is run in a terminal, **Then** a bordered table displays with color-coded status columns.
2. **Given** no cells exist, **When** `replicator cells` is run, **Then** a styled message indicates "No cells found" instead of raw `[]`.
3. **Given** output is piped, **When** `replicator cells | cat` is run, **Then** the table renders without ANSI codes but retains its structure.

---

### User Story 3 - Consistent Styling Across CLI Commands (Priority: P2)

A developer uses setup, init, version, stats, and query commands and sees visually consistent output that matches the doctor and cells styling. Indicators, headers, and emphasis use the same color palette and typography across all commands.

**Why this priority**: Visual consistency makes the tool feel polished and professional. Inconsistent styling (some commands colored, others plain) looks unfinished.

**Independent Test**: Run each CLI command in sequence and verify the same color palette and indicator style is used throughout.

**Acceptance Scenarios**:

1. **Given** a terminal, **When** `replicator setup` is run, **Then** success indicators are green, failures red, and the same style as doctor output.
2. **Given** a terminal, **When** `replicator version` is run, **Then** the version number is bold, commit and date are dimmed.
3. **Given** a terminal, **When** `replicator init` is run, **Then** "initialized .hive/" is green and "already initialized" is dimmed.
4. **Given** a terminal, **When** `replicator stats` and `replicator query` are run, **Then** table headers are styled consistently with the cells table.

---

### User Story 4 - Structured Logging with Per-Repo Log File (Priority: P2)

When the MCP server runs (`replicator serve`), structured log messages are written to both stderr and a per-repo log file at `[repo]/.unbound-force/replicator.log`. The log file is truncated on each server startup so it only contains the current session's logs. CLI commands (doctor, cells, etc.) log to stderr only -- no log file.

**Why this priority**: Debugging MCP tool calls is difficult without a persistent log. The per-repo log captures a full session's activity alongside the project it serves, without interleaving with other sessions.

**Independent Test**: Start `replicator serve`, send several MCP requests, stop the server, and verify `[repo]/.unbound-force/replicator.log` contains structured log entries for each tool call. Restart the server and verify the log file is truncated.

**Acceptance Scenarios**:

1. **Given** a project directory, **When** `replicator serve` starts, **Then** `.unbound-force/replicator.log` is created (truncating any existing file) and structured log entries begin writing.
2. **Given** the MCP server is running, **When** a `tools/call` request is processed, **Then** a log entry appears in both stderr and the log file with the tool name, duration, and success/error status.
3. **Given** the MCP server is stopped and restarted, **When** the log file is inspected, **Then** it contains only entries from the most recent session (truncated on startup).
4. **Given** a CLI command is run (not `serve`), **When** `replicator doctor` runs, **Then** no `.unbound-force/replicator.log` file is created or modified.

---

### Edge Cases

- What happens when the terminal does not support colors? The renderer detects the color profile and falls back to plain text automatically.
- What happens when the cells table has very long titles? Titles are truncated to fit the terminal width, preserving the table structure.
- What happens when `.unbound-force/` directory does not exist for logging? It is created with `0o755` permissions on `serve` startup.
- What happens when the log file cannot be written (permissions)? A warning is emitted to stderr (via `fmt.Fprintf`, not `charmbracelet/log`, since the logger itself failed to initialize -- bootstrap exception) and the server continues without file logging.
- What happens when `NO_COLOR=1` is set? All lipgloss output uses the ASCII color profile, producing no escape codes. This is handled automatically by the renderer.
- What happens when two `replicator serve` instances run in the same repo? Only one instance per repository is supported. Concurrent instances will produce interleaved log output in the shared log file. This is a known limitation.
- What happens during a long-running MCP server session? Log file size is bounded by session duration (truncated on restart). No rotation is needed for typical sessions. Long-running sessions may produce large log files.

## Requirements *(mandatory)*

### Functional Requirements

#### Styling (US1, US2, US3)
- **FR-001**: All CLI output that uses color MUST use a renderer-aware style system that automatically detects TTY, `NO_COLOR`, and pipe contexts.
- **FR-002**: A centralized style definition MUST be shared across all CLI commands to ensure consistent colors and typography.
- **FR-003**: The doctor command MUST display results with color-coded indicators (green=pass, yellow=warn, red=fail) and a rounded summary box.
- **FR-004**: The cells command MUST display a bordered table with styled headers and per-row status coloring (green=open, yellow=in_progress, red=blocked, gray=closed).
- **FR-005**: When output is piped or `NO_COLOR` is set, all commands MUST produce clean text with no ANSI escape codes.
- **FR-006**: The setup, init, version, stats, and query commands MUST use the shared style system for indicators and headers.

#### Logging (US4)
- **FR-007**: The MCP server MUST write structured log messages to both stderr and a per-repo log file at `[repo]/.unbound-force/replicator.log`.
- **FR-008**: The log file MUST be truncated (not appended) on each `replicator serve` startup.
- **FR-009**: Each MCP `tools/call` invocation MUST be logged with at minimum: tool name, duration, and success/error status.
- **FR-010**: CLI commands (not `serve`) MUST log to stderr only -- no log file creation.
- **FR-011**: If the log file cannot be created or written, the server MUST continue operating with stderr-only logging and emit a warning.

#### Color Palette
- **FR-012**: The color palette MUST match the Unbound Force ecosystem: green (ANSI 10), yellow (ANSI 11), red (ANSI 9), gray (ANSI 241), pink bold (ANSI 212), purple border (ANSI 63).

### Key Entities

- **Styles**: A centralized set of named style definitions (pass, warn, fail, dim, bold, title, box) that all CLI commands share. Created from a writer-aware renderer to support pipe and NO_COLOR detection.
- **Log Entry**: A structured log message with timestamp, level, message, and key-value metadata (tool name, duration, error). Written as human-readable text to stderr and log file.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The doctor command output matches the visual style of `uf doctor` when run in a color-capable terminal (same indicator colors, same summary box style).
- **SC-002**: The cells command renders a bordered table with per-row status coloring for 4+ cells in under 50 milliseconds.
- **SC-003**: All 9 CLI commands use the shared style system -- zero commands produce raw `fmt.Printf` output in the final implementation.
- **SC-004**: When piped (`replicator doctor | cat`), the output contains zero ANSI escape code sequences (verified by `grep -P '\x1b\['`).
- **SC-005**: The MCP server log file at `.unbound-force/replicator.log` contains at least one structured entry per tool call, with tool name and duration.
- **SC-006**: Restarting `replicator serve` truncates the log file (file size resets to 0 before new entries).

## Assumptions

- The existing `internal/doctor/checks.go` separation of health check logic from formatting is preserved -- only the formatting layer changes.
- The `cells` command currently dumps raw JSON; it will switch to a human-readable table for TTY and retain JSON output for piped contexts (or with a `--json` flag).
- The `docs` command output is markdown (not terminal-styled) and is excluded from lipgloss styling.
- The log file path `.unbound-force/replicator.log` is relative to the working directory where `replicator serve` is launched (typically the project root).

## Dependencies

- **charmbracelet/lipgloss v1.1.0**: Style rendering with TTY/pipe detection.
- **charmbracelet/log v1.0.0**: Structured logging with leveled output.
- **muesli/termenv v0.16.0**: Terminal environment detection (used by lipgloss internally, referenced for `termenv.Ascii` constant).
