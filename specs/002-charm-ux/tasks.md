# Tasks: Charm Bracelet CLI UX

**Input**: Design documents from `/specs/002-charm-ux/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, quickstart.md
**CI Gate**: `go vet ./...`, `go test ./... -count=1 -race`, `go build ./...`

**Coverage Strategy**: New packages (`internal/ui/`) must achieve ≥80% line coverage. Modified command files must not regress existing coverage. All formatting functions must have tests that verify both colored (TTY) and plain (piped) output paths.

**TDD Note**: Within each task, follow Red-Green-Refactor: write the test first, then implement the code to pass it. Tasks are listed as implementation+test pairs for readability, but the test is written before the implementation within each task.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

---

## Phase 0: Foundation — `internal/ui/` Package (US1, US2, US3)

**Purpose**: Create the centralized style system that all subsequent phases depend on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T001 Add `charmbracelet/lipgloss@v1.1.0`, `charmbracelet/log@v1.0.0`, `muesli/termenv@v0.16.0` dependencies to `go.mod` via `go get`; run `go mod tidy`
- [x] T002 Create `internal/ui/styles.go`: define `Styles` struct with `Pass`, `Warn`, `Fail`, `Dim`, `Bold`, `Title`, `Box`, `Border` fields (all `lipgloss.Style`), `HasColor bool`, and `Renderer *lipgloss.Renderer`; implement `NewStyles(w io.Writer)` constructor using `lipgloss.NewRenderer(w)` with color palette green(10), yellow(11), red(9), gray(241), pink(212), purple(63) per FR-001, FR-002, FR-012
- [x] T003 Create `internal/ui/styles_test.go`: test that `NewStyles(&bytes.Buffer{})` produces `HasColor == false` and `Render()` output contains no ANSI escape codes (`\x1b[`); test that color constants match FR-012 palette values per FR-005
- [x] T004 [P] Create `internal/ui/table.go`: implement `NewTable(s *Styles, headers []string, rows [][]string) *table.Table` helper wrapping `lipgloss/table` with `NormalBorder()`, purple(63) `BorderStyle`, and renderer-aware `StyleFunc` per FR-004
- [x] T005 [P] Create `internal/ui/table_test.go`: test table renders with borders for TTY-like output; test table renders without ANSI codes when writer is `bytes.Buffer` (non-TTY) per FR-005

**Checkpoint**: `go test ./internal/ui/... -count=1 -race` passes. `go vet ./internal/ui/...` clean.

---

## Phase 1: Doctor Styling (US1) — Priority P1 🎯

**Goal**: Replace raw `fmt.Printf` doctor output with lipgloss-styled rendering matching `uf doctor`.

**Independent Test**: Run `replicator doctor` in terminal → colored output with summary box. Run `replicator doctor | cat` → clean plain text with `[PASS]`/`[WARN]`/`[FAIL]`.

- [x] T006 [US1] Create `internal/doctor/format.go`: implement `FormatText(results []CheckResult, w io.Writer) error` using `ui.NewStyles(w)` — render each result with color-coded indicator (green ✅ pass, yellow ⚠️ warn, red ❌ fail), check name, message, and duration per FR-003
- [x] T007 [US1] Add summary box to `internal/doctor/format.go`: render a `lipgloss.RoundedBorder()` box with purple(63) border showing pass/warn/fail counts and contextual message ("Everything looks good!" vs "Run after fixes") per FR-003
- [x] T008 [US1] Create `internal/doctor/format_test.go`: test `FormatText` with pass result → output contains `[PASS]` indicator (non-TTY buffer); test warn result → `[WARN]`; test fail result → `[FAIL]` per FR-005
- [x] T009 [US1] Add test in `internal/doctor/format_test.go`: verify summary box renders with pass/warn/fail counts for mixed results; verify no ANSI escape codes in `bytes.Buffer` output per FR-003, FR-005
- [x] T010 [US1] Update `cmd/replicator/doctor.go`: replace `fmt.Printf` table and `statusIcon()` calls with `doctor.FormatText(results, os.Stdout)`; remove `statusIcon()` helper function (dead code after migration) per FR-003

**Checkpoint**: `go test ./internal/doctor/... -count=1 -race` passes. `replicator doctor` renders styled output. `replicator doctor | cat` produces clean text.

---

## Phase 2: Cells Table (US2) — Priority P1 🎯

**Goal**: Replace raw JSON dump with a bordered, color-coded table.

**Independent Test**: Create 4 cells with different statuses → `replicator cells` shows styled table. `replicator cells --json` outputs JSON. `replicator cells | cat` shows clean table.

- [x] T011 [US2] Extract formatting to `internal/hive/format.go`: create `FormatCells(cells []Cell, w io.Writer, styles *ui.Styles) error` — renders bordered table with columns ID (truncated to 8 chars), TITLE (truncated to terminal width), STATUS, TYPE, PRIORITY per FR-004, AP-002
- [x] T012 [US2] Add status column coloring in `internal/hive/format.go`: use `StyleFunc` to color-code status — green=open, yellow=in_progress, red=blocked, gray=closed per FR-004
- [x] T012a [US2] Add terminal-width-aware title truncation in `internal/hive/format.go`: use `table.Width(n)` to fit the table within terminal width; truncate titles that would overflow per spec edge case
- [x] T013 [US2] Update `cmd/replicator/cells.go`: add `--json` flag; when set, output `json.MarshalIndent` (backward compatible); when not set, call `hive.FormatCells()` per FR-004
- [x] T014 [US2] Handle empty state in `internal/hive/format.go`: when `len(cells) == 0`, write styled "No cells found" message using `styles.Dim` instead of empty table per FR-004
- [x] T015 [US2] Add test in `internal/hive/format_test.go`: verify table output to `bytes.Buffer` contains no ANSI codes; verify status coloring; verify empty state; verify long title truncation per FR-004, FR-005

**Checkpoint**: `go test ./cmd/replicator/... -count=1 -race` passes. `replicator cells` renders styled table. `replicator cells --json` outputs JSON.

---

## Phase 3: Remaining CLI Commands (US3) — Priority P2

**Goal**: Apply consistent styling to setup, init, version, stats, and query commands.

**Independent Test**: Run each command → verify same color palette and indicator style as doctor/cells.

### Setup & Init (different files, parallelizable)

- [x] T016 [P] [US3] Update `cmd/replicator/setup.go`: replace `\u2713`/`\u2717` Unicode with `ui.NewStyles(os.Stdout)` — green `Pass.Render("✓")` for success, red `Fail.Render("✗")` for failure; style completion message per FR-006
- [x] T017 [P] [US3] Update `cmd/replicator/init.go`: use `ui.NewStyles(os.Stdout)` — green `Pass.Render("initialized .hive/")` for fresh init, dim `Dim.Render("already initialized")` for idempotent case per FR-006

### Version (different file, parallelizable)

- [x] T018 [P] [US3] Update `versionCmd()` in `cmd/replicator/main.go`: use `ui.NewStyles(os.Stdout)` — bold version number via `Bold.Render()`, dim commit and date via `Dim.Render()` per FR-006

### Stats & Query (internal packages, sequential)

- [x] T019 [US3] Update `internal/stats/stats.go`: replace `fmt.Fprintln` headers (`=== Replicator Stats ===`, `Events by Type:`) with `ui.NewStyles(w)` — use `Title.Render()` for section headers, `Dim.Render()` for empty-state messages per FR-006
- [x] T020 [US3] Update `internal/query/presets.go`: replace `fmt.Fprintf` column headers and `---` separators with `ui.NewTable` for tabular presets (`runAgentActivity`, `runCellsByStatus`, `runRecentEvents`); use `ui.NewStyles(w)` for non-tabular output (`runSwarmCompletionRate`) per FR-006
- [x] T021 [US3] Update `cmd/replicator/query.go` `listQueryPresets()`: replace `fmt.Println` with `ui.NewStyles(os.Stdout)` — bold preset names via `Bold.Render()`, dim usage hint via `Dim.Render()` per FR-006

**Checkpoint**: `go test ./... -count=1 -race` passes. All 8 styled CLI commands (doctor, cells, setup, init, version, stats, query, serve) use `ui.Styles` — `docs` is excluded (markdown output).

---

## Phase 4: Structured Logging (US4) — Priority P2

**Goal**: Add `charmbracelet/log` with per-repo log file for MCP server sessions.

**Independent Test**: Start `replicator serve` → `.uf/replicator/replicator.log` created with structured entries. CLI commands do not create log file.

- [x] T022 [US4] Update `cmd/replicator/serve.go`: on `serveMCP()` entry, create `.uf/replicator/` directory with `os.MkdirAll(0o755)`, open `.uf/replicator/replicator.log` with `os.Create` (truncate), configure `charmbracelet/log` with `io.MultiWriter(os.Stderr, logFile)` per FR-007, FR-008
- [x] T023 [US4] Handle log file creation failure in `cmd/replicator/serve.go`: if `os.Create` or `os.MkdirAll` fails, emit warning to stderr via `fmt.Fprintf(os.Stderr, ...)` and continue with stderr-only logger — do not crash per FR-011. Bootstrap exception to CS-008: `charmbracelet/log` cannot be used here because the logger itself is what failed to initialize. Add `defer logFile.Close()` for explicit cleanup.
- [x] T024 [US4] Update `internal/mcp/server.go`: add `Logger` field to `Server` struct (interface or `*log.Logger`); update `NewServer` to accept logger; wrap `handleToolsCall` to log tool name, `time.Since(start)` duration, and success/error status per FR-009
- [x] T025 [US4] Update `cmd/replicator/serve.go`: pass configured logger to `mcp.NewServer(reg, Version, logger)` per FR-009
- [x] T026 [US4] Create `cmd/replicator/serve_test.go`: test that `serveMCP`-style setup in `t.TempDir()` creates `.uf/replicator/replicator.log`; test truncation by writing marker, re-creating file, verifying marker absent per FR-007, FR-008
- [x] T027 [US4] Add test in `cmd/replicator/serve_test.go`: test that log file creation failure (read-only directory) does not panic — logger falls back to stderr-only per FR-011
- [x] T028 [US4] Add test in `internal/mcp/server_test.go`: test that `handleToolsCall` with a mock logger records tool name and duration for a successful call and an error call per FR-009

**Checkpoint**: `go test ./... -count=1 -race` passes. `replicator serve` creates log file. CLI commands (doctor, cells, etc.) do not create `.uf/replicator/replicator.log`.

---

## Phase 5: Polish & Verification

**Purpose**: Final integration testing, documentation updates, CI parity gate.

- [x] T029 [P] Verify zero ANSI codes in piped output: run `replicator doctor | cat`, `replicator cells | cat`, `replicator setup | cat` and confirm no `\x1b[` sequences per FR-005, SC-004
- [x] T030 [P] Verify `NO_COLOR=1` produces plain text: run `NO_COLOR=1 replicator doctor` and confirm text indicators `[PASS]`/`[WARN]`/`[FAIL]` with no escape codes per FR-005
- [x] T031 [P] Verify zero raw `fmt.Printf` for user-facing output: `grep -rn 'fmt.Printf\|fmt.Println' cmd/replicator/ --include='*.go'` returns only non-user-facing uses (error paths, docs.go) per SC-003
- [x] T032 Update `AGENTS.md`: add `internal/ui/` to Project Structure with description "Centralized lipgloss styles + table helpers"; add `charmbracelet/lipgloss`, `charmbracelet/log`, `muesli/termenv` to Active Technologies per documentation gate
- [x] T033 Update `README.md` if CLI output examples exist: update any doctor/cells output samples to reflect new styled format per documentation gate
- [x] T034 Run `make check` (full CI parity gate): `go vet ./...` + `go test ./... -count=1 -race` + `go build ./...` — all must pass per SC-001 through SC-006

**Checkpoint**: `make check` passes. All success criteria (SC-001 through SC-006) verified. quickstart.md steps 1–18 validated.

---

## Dependencies & Execution Order

### Phase Dependencies

```text
Phase 0 (ui package) ──BLOCKS──┬──> Phase 1 (doctor)  ──┐
                                ├──> Phase 2 (cells)   ──┼──> Phase 5 (polish)
                                └──> Phase 3 (commands)──┤
                                                         │
                          Phase 4 (logging) ─────────────┘
```

- **Phase 0**: No dependencies — start immediately
- **Phases 1, 2, 3**: All depend on Phase 0 completion. Can run in parallel (different files, read-only dependency on `internal/ui/`)
- **Phase 4**: Independent of Phases 1–3 (only touches `serve.go` and `server.go`). Sequenced after Phase 3 for clean diffs.
- **Phase 5**: Depends on all previous phases

### Parallel Opportunities

| Tasks | Can Parallel? | Reason |
|-------|--------------|--------|
| T002, T003 | No | T003 tests T002's output |
| T004, T005 | No | T005 tests T004's output |
| T002+T003, T004+T005 | Yes | Different files (`styles.go` vs `table.go`) |
| T006–T010 (Phase 1), T011–T015 (Phase 2), T016–T021 (Phase 3) | Yes | Different files, shared `ui/` is read-only |
| T016, T017, T018 | Yes | Different files (`setup.go`, `init.go`, `main.go`) |
| T019, T020 | No | Both modify internal packages consumed by same commands |
| T029, T030, T031 | Yes | Independent verification scripts |

### Within Each Phase

- Tests SHOULD be written alongside implementation (same task or immediately after)
- Phase gate must pass before proceeding to next phase
- Mark `- [ ]` to `- [x]` immediately after each task completion

<!-- spec-review: passed -->
