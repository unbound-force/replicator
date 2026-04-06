# Implementation Plan: Charm Bracelet CLI UX

**Branch**: `002-charm-ux` | **Date**: 2026-04-06 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-charm-ux/spec.md`

## Summary

Upgrade all CLI output from raw `fmt.Printf` with Unicode codepoints to
lipgloss-styled rendering with automatic TTY/pipe/NO_COLOR detection.
Add a centralized `internal/ui/` package that defines the Unbound Force
color palette and reusable style constructors. Introduce structured
logging via `charmbracelet/log` with a per-repo log file for MCP server
sessions.

The approach follows the proven pattern from `uf doctor` (format.go) and
`gaze` (styles.go): a renderer-aware style struct created per-command,
with `lipgloss.NewRenderer(w)` for pipe-safe output.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: `charmbracelet/lipgloss v1.1.0`, `charmbracelet/log v1.0.0`, `muesli/termenv v0.16.0`, `charmbracelet/lipgloss/table` (sub-package of lipgloss)
**Storage**: SQLite via `modernc.org/sqlite` (unchanged)
**Testing**: `go test` (stdlib), in-memory SQLite, `bytes.Buffer` for output capture
**Target Platform**: macOS, Linux (CLI + MCP server over stdio)
**Project Type**: CLI + MCP server
**Performance Goals**: Table rendering < 50ms for 100 cells (SC-002)
**Constraints**: Zero ANSI codes when piped (SC-004), MCP stdout must remain clean JSON-RPC
**Scale/Scope**: 9 CLI commands, 1 MCP server, 1 new package

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Autonomous Collaboration** | PASS | UI styling is CLI-only. MCP tool responses remain JSON — no lipgloss in tool output. Logging adds observability without coupling. |
| **II. Composability First** | PASS | lipgloss/log are compile-time dependencies with zero external services. Graceful degradation: no-color terminals get plain text automatically. Log file failure degrades to stderr-only (FR-011). |
| **III. Observable Quality** | PASS | MCP responses unchanged (JSON content blocks). Doctor output gains structured formatting. Log entries are structured with tool name + duration (FR-009). |
| **IV. Testability** | PASS | All formatters accept `io.Writer` — testable with `bytes.Buffer`. Renderer created from writer enables deterministic no-color testing. No external services needed. |

No violations. No complexity tracking needed.

## Project Structure

### Documentation (this feature)

```text
specs/002-charm-ux/
├── plan.md              # This file
├── research.md          # Charm library patterns + reference analysis
├── quickstart.md        # Verification steps
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── ui/                  # NEW — centralized styles + table helpers
│   ├── styles.go        # Styles struct, NewStyles(io.Writer), color constants
│   ├── styles_test.go   # Verify no-color fallback, color values
│   ├── table.go         # Table builder wrapping lipgloss/table
│   └── table_test.go    # Table rendering tests
├── doctor/
│   ├── checks.go        # UNCHANGED — health check logic
│   ├── checks_test.go   # UNCHANGED
│   ├── format.go        # NEW — lipgloss formatter (extracted from cmd)
│   └── format_test.go   # NEW — formatter tests
├── mcp/
│   └── server.go        # MODIFIED — add logging to tools/call handler
├── config/
│   └── config.go        # MODIFIED — add LogFilePath field
├── stats/
│   └── stats.go         # MODIFIED — accept Styles, use table helpers
├── query/
│   └── presets.go        # MODIFIED — accept Styles, use table helpers
└── ...

cmd/replicator/
├── main.go              # MODIFIED — pass os.Stdout to style constructors
├── doctor.go            # MODIFIED — use doctor.FormatText
├── cells.go             # MODIFIED — use ui.Table for TTY, JSON for pipe
├── setup.go             # MODIFIED — use ui.Styles indicators
├── init.go              # MODIFIED — use ui.Styles indicators
├── serve.go             # MODIFIED — configure charmbracelet/log + log file
├── stats.go             # MINOR — pass writer through
├── query.go             # MINOR — pass writer through
└── docs.go              # UNCHANGED — markdown output, no styling
```

**Structure Decision**: New `internal/ui/` package centralizes all style
definitions. This mirrors the Gaze pattern (`internal/report/styles.go`)
where a `Styles` struct holds named styles and is constructed from a
renderer. The `ui` package is imported by `cmd/replicator/` and by
internal packages that format output (doctor, stats, query).

## Phased Implementation

### Phase 0: Foundation — `internal/ui/` Package (US1, US2, US3)

**Goal**: Create the centralized style system that all commands will use.

**Rationale**: Every subsequent phase depends on having a shared style
definition. Building this first with tests establishes the color palette
contract and pipe-safe rendering pattern.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 0.1 | Add lipgloss, log, termenv dependencies | `go.mod`, `go.sum` | — | `go mod tidy` succeeds |
| 0.2 | Create `internal/ui/styles.go` with `Styles` struct and `NewStyles(io.Writer)` constructor | `internal/ui/styles.go` | FR-001, FR-002, FR-012 | Unit: verify color values match palette |
| 0.3 | Test pipe-safe rendering: `NewStyles` with non-TTY writer produces no ANSI codes | `internal/ui/styles_test.go` | FR-005 | Unit: render to `bytes.Buffer`, assert no `\x1b[` |
| 0.4 | Create `internal/ui/table.go` with `Table` builder wrapping `lipgloss/table` | `internal/ui/table.go` | FR-004 | Unit: render table, verify borders + alignment |
| 0.5 | Test table pipe fallback: table renders without ANSI when writer is not TTY | `internal/ui/table_test.go` | FR-005 | Unit: render to buffer, grep for escape codes |

**Phase Gate**: `go test ./internal/ui/...` passes. `go vet ./...` clean.

### Phase 1: Doctor Styling (US1)

**Goal**: Replace the raw `fmt.Printf` doctor output with lipgloss-styled
rendering matching `uf doctor`.

**Rationale**: Doctor is the first command developers run (spec: "sets the
tone for the entire tool"). It has the most complex formatting (grouped
results, indicators, summary box) and serves as the template for all
other commands.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 1.1 | Create `internal/doctor/format.go` with `FormatText(results, io.Writer)` | `internal/doctor/format.go` | FR-003 | Unit: render known results, verify indicators |
| 1.2 | Test color output: pass results -> green indicator, warn -> yellow, fail -> red | `internal/doctor/format_test.go` | FR-003 | Unit: render to color-aware buffer, check ANSI codes |
| 1.3 | Test plain-text fallback: pipe writer -> `[PASS]`, `[WARN]`, `[FAIL]` text indicators | `internal/doctor/format_test.go` | FR-005 | Unit: render to `bytes.Buffer`, assert text indicators |
| 1.4 | Test summary box: verify rounded border and pass/warn/fail counts | `internal/doctor/format_test.go` | FR-003 | Unit: check box characters in output |
| 1.5 | Update `cmd/replicator/doctor.go` to call `doctor.FormatText` instead of `fmt.Printf` | `cmd/replicator/doctor.go` | FR-003 | Integration: `go build` + manual verify |
| 1.6 | Remove `statusIcon()` helper from `cmd/replicator/doctor.go` (dead code) | `cmd/replicator/doctor.go` | — | Build succeeds, no references |

**Phase Gate**: `go test ./internal/doctor/...` passes. `replicator doctor` renders styled output. `replicator doctor | cat` produces clean text.

### Phase 2: Cells Table (US2)

**Goal**: Replace raw JSON dump with a bordered, color-coded table.

**Rationale**: Cells is the most frequently used command. The table format
makes scanning many work items fast and intuitive.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 2.1 | Update `cmd/replicator/cells.go` to render a `ui.Table` for TTY output | `cmd/replicator/cells.go` | FR-004 | Unit: render cells to buffer, verify table structure |
| 2.2 | Add `--json` flag to `cells` command for programmatic output | `cmd/replicator/cells.go`, `cmd/replicator/main.go` | FR-004 | Unit: verify JSON output with flag |
| 2.3 | Color-code status column: green=open, yellow=in_progress, red=blocked, gray=closed | `cmd/replicator/cells.go` | FR-004 | Unit: verify ANSI codes per status |
| 2.4 | Handle empty state: "No cells found" styled message instead of `[]` | `cmd/replicator/cells.go` | FR-004 | Unit: empty slice -> styled message |
| 2.5 | Test pipe fallback: table renders without ANSI codes when piped | `cmd/replicator/cells.go` | FR-005 | Unit: render to buffer, no escape codes |

**Phase Gate**: `go test ./cmd/replicator/...` passes. `replicator cells` renders styled table. `replicator cells --json` outputs JSON.

### Phase 3: Remaining CLI Commands (US3)

**Goal**: Apply consistent styling to setup, init, version, stats, and query.

**Rationale**: Visual consistency across all commands makes the tool feel
polished. These commands have simpler output than doctor/cells, so they
can be updated quickly using the established patterns.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 3.1 | Update `cmd/replicator/setup.go`: green checkmark / red X via `ui.Styles` | `cmd/replicator/setup.go` | FR-006 | Unit: verify styled indicators |
| 3.2 | Update `cmd/replicator/init.go`: green "initialized", dim "already initialized" | `cmd/replicator/init.go` | FR-006 | Unit: verify styled messages |
| 3.3 | Update `cmd/replicator/main.go` `versionCmd`: bold version, dim commit/date | `cmd/replicator/main.go` | FR-006 | Unit: verify styled version output |
| 3.4 | Update `internal/stats/stats.go`: accept `io.Writer`, use `ui.Styles` for headers | `internal/stats/stats.go` | FR-006 | Unit: existing tests pass with styled output |
| 3.5 | Update `internal/query/presets.go`: use `ui.Table` for tabular presets | `internal/query/presets.go` | FR-006 | Unit: existing tests pass with styled output |
| 3.6 | Update `cmd/replicator/query.go`: styled preset list | `cmd/replicator/query.go` | FR-006 | Unit: verify styled list output |

**Phase Gate**: `go test ./...` passes. All 9 CLI commands use `ui.Styles` — zero raw `fmt.Printf` for user-facing output.

### Phase 4: Structured Logging (US4)

**Goal**: Add `charmbracelet/log` with per-repo log file for MCP server.

**Rationale**: Logging is isolated from styling — it only affects the
`serve` command and MCP server internals. Implementing it last avoids
entangling log setup with the style refactoring.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 4.1 | Create log setup in `cmd/replicator/serve.go`: configure `charmbracelet/log` with stderr + file multi-writer | `cmd/replicator/serve.go` | FR-007, FR-008 | Unit: verify log file creation + truncation |
| 4.2 | Create `.uf/replicator/` directory on serve startup (0o755) | `cmd/replicator/serve.go` | FR-007 | Unit: verify directory creation in `t.TempDir()` |
| 4.3 | Handle log file creation failure: warn to stderr, continue without file | `cmd/replicator/serve.go` | FR-011 | Unit: read-only dir -> no crash, warning emitted |
| 4.4 | Add tool call logging to `internal/mcp/server.go`: log tool name, duration, success/error | `internal/mcp/server.go` | FR-009 | Unit: mock logger, verify log entries per tool call |
| 4.5 | Verify CLI commands do NOT create log file | `cmd/replicator/doctor.go` (no change) | FR-010 | Unit: run doctor, verify no `.uf/replicator/replicator.log` |
| 4.6 | Test log truncation: restart serve, verify file contains only new session entries | — | FR-008 | Integration: write marker, restart, verify marker absent |

**Phase Gate**: `go test ./...` passes. `replicator serve` creates log file. CLI commands do not.

### Phase 5: Polish & Verification

**Goal**: Final integration testing, documentation updates, CI validation.

| # | Task | Files | FR | Test Strategy |
|---|------|-------|----|---------------|
| 5.1 | Verify zero ANSI codes in piped output for all commands | — | FR-005, SC-004 | Script: `replicator doctor \| grep -P '\x1b\['` returns empty |
| 5.2 | Verify `NO_COLOR=1` produces plain text for all commands | — | FR-005 | Script: `NO_COLOR=1 replicator doctor` has no escape codes |
| 5.3 | Update `AGENTS.md` with `internal/ui/` package description | `AGENTS.md` | — | Review |
| 5.4 | Update `README.md` if CLI output examples exist | `README.md` | — | Review |
| 5.5 | Run `make check` (full CI parity gate) | — | — | CI parity: `go vet ./... && go test ./...` |

**Phase Gate**: `make check` passes. All success criteria (SC-001 through SC-006) verified.

## Dependency Graph

```text
Phase 0 (ui package)
  ├──> Phase 1 (doctor)
  ├──> Phase 2 (cells)
  └──> Phase 3 (other commands)
         └──> Phase 4 (logging)  <-- independent, but sequenced for clean diffs
                └──> Phase 5 (polish)
```

Phases 1, 2, and 3 can be parallelized after Phase 0 completes — they
touch different files and share only the `internal/ui/` package (read-only
dependency). Phase 4 is independent but sequenced last to avoid merge
conflicts in `serve.go`.

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| lipgloss output leaks into MCP stdout | Medium | HIGH | MCP server writes to stdout only via `json.Marshal`. CLI commands write to `os.Stdout` via lipgloss. Logging writes to stderr + file. These are separate writers — no shared state. |
| Table rendering slow for large cell counts | Low | Medium | lipgloss/table is string-based, not I/O-bound. Benchmark with 1000 cells in Phase 2 gate. |
| `NO_COLOR` not respected by lipgloss | Low | Medium | lipgloss v1.1.0 respects `NO_COLOR` via termenv. Verified in Phase 5 gate. |
| Log file permissions on shared servers | Low | Low | File created with default umask. Document in README. |
| Breaking existing tests that assert on raw output | Medium | Medium | Tests that check `fmt.Printf` output will need updating. Identify all such tests in Phase 0 research. |

## Success Criteria Mapping

| Criterion | Phase | Verification |
|-----------|-------|-------------|
| SC-001: Doctor matches `uf doctor` style | Phase 1 gate | Visual comparison |
| SC-002: Cells table < 50ms for 4+ cells | Phase 2 gate | Benchmark test |
| SC-003: Zero raw `fmt.Printf` in final | Phase 5 gate | `grep -r 'fmt.Printf' cmd/replicator/` returns only non-user-facing |
| SC-004: Zero ANSI in piped output | Phase 5 gate | `replicator doctor \| grep -cP '\x1b\['` = 0 |
| SC-005: Log entry per tool call | Phase 4 gate | Inspect log file after tool calls |
| SC-006: Log truncation on restart | Phase 4 gate | Restart test |
