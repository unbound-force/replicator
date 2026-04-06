# Implementation Plan: Terminology Rename + Agent Kit

**Branch**: `003-rename-terminology` | **Date**: 2026-04-06 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/003-rename-terminology/spec.md`

## Summary

Rename the three core MCP tool families (hive→org, swarmmail→comms, swarm→forge) across all Go packages, tool registrations, and documentation. Simultaneously ship an agent kit (5 commands, 7 skills, 3 agents) embedded in the binary and scaffolded by `replicator init`.

The rename is a mechanical refactor touching 6 package directories, 45 tool `Name:` strings, 4 parity fixture files, and all import paths. The agent kit introduces `embed.FS` for 15 template files and extends the `init` command with `--force` flag support.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: cobra (CLI), modernc.org/sqlite (pure Go SQLite), embed (stdlib)
**Storage**: SQLite at `~/.config/uf/replicator/replicator.db` (WAL mode)
**Testing**: `go test` (stdlib only), parity tests with `//go:build parity` tag
**Target Platform**: macOS, Linux (cross-compiled via GoReleaser)
**Project Type**: CLI + MCP server
**Performance Goals**: N/A (rename is behavior-preserving)
**Constraints**: Zero behavioral change to tool responses; all 53 tools must pass existing tests
**Scale/Scope**: 53 MCP tools, ~20 Go files with import path changes, 15 new embedded files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Autonomous Collaboration** | ✅ PASS | Tools remain independently callable via MCP. Only names change, not behavior or response shapes. |
| **II. Composability First** | ✅ PASS | Binary remains standalone. Agent kit is scaffolded locally, no external service dependency. `embed.FS` is stdlib — no new dependencies. |
| **III. Observable Quality** | ✅ PASS | All tool responses remain JSON. Parity test fixtures updated to new names. `replicator docs` updated to new categories. |
| **IV. Testability** | ✅ PASS | All tests continue using `db.OpenMemory()`, `t.TempDir()`, `httptest`. New init tests verify agent kit scaffolding with `t.TempDir()`. |

No violations. No complexity tracking needed.

## Project Structure

### Documentation (this feature)

```text
specs/003-rename-terminology/
├── plan.md              # This file
├── research.md          # Rename impact analysis
├── quickstart.md        # Implementation quick-reference
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code — Before → After

```text
# Package renames (git mv)
internal/hive/              → internal/org/
internal/swarmmail/         → internal/comms/
internal/swarm/             → internal/forge/
internal/tools/hive/        → internal/tools/org/
internal/tools/swarmmail/   → internal/tools/comms/
internal/tools/swarm/       → internal/tools/forge/

# New files (agent kit)
internal/agentkit/
├── agentkit.go             # embed.FS + scaffold logic
├── agentkit_test.go        # Tests for scaffold, skip, force
└── content/                # Embedded template files
    ├── command/
    │   ├── forge.md
    │   ├── org.md
    │   ├── inbox.md
    │   ├── forge-status.md
    │   └── handoff.md
    ├── skills/
    │   ├── always-on-guidance/SKILL.md
    │   ├── forge-coordination/SKILL.md
    │   ├── replicator-cli/SKILL.md
    │   ├── testing-patterns/SKILL.md
    │   ├── system-design/SKILL.md
    │   ├── learning-systems/SKILL.md
    │   └── forge-global/SKILL.md
    └── agents/
        ├── coordinator.md
        ├── worker.md
        └── background-worker.md

# Modified files (import paths + tool names)
cmd/replicator/serve.go     # Import path updates
cmd/replicator/docs.go      # Import paths + category map
cmd/replicator/init.go      # Agent kit integration + --force flag
cmd/replicator/init_test.go # Agent kit test cases
test/parity/parity_test.go  # Import paths + fixture tool names
test/parity/fixtures/*.json # Tool name keys updated
```

**Structure Decision**: The existing single-project Go structure is preserved. The `internal/agentkit/` package is added for embed.FS content and scaffold logic, following the existing pattern of domain packages under `internal/`.

## Implementation Phases

### Phase 1: Directory Renames [US1]

Rename 6 Go package directories using `git mv` to preserve history. This is the foundation — all subsequent phases depend on these directories existing at their new paths.

**Files touched**:
- `internal/hive/` → `internal/org/` (11 files)
- `internal/swarmmail/` → `internal/comms/` (6 files)
- `internal/swarm/` → `internal/forge/` (14 files)
- `internal/tools/hive/` → `internal/tools/org/` (1 file)
- `internal/tools/swarmmail/` → `internal/tools/comms/` (1 file)
- `internal/tools/swarm/` → `internal/tools/forge/` (1 file)

**Risk**: LOW. `git mv` preserves history. The build will break until Phase 2 completes (expected).

**Checkpoint**: `ls internal/org/ internal/comms/ internal/forge/` succeeds. Old directories do not exist.

---

### Phase 2: Import Path Updates [US1]

Update all Go import paths from old package names to new ones. Also update package declarations in renamed files.

**Files touched** (all `.go` files importing renamed packages):
- `cmd/replicator/serve.go` — 4 import paths
- `cmd/replicator/docs.go` — 4 import paths + import aliases
- `cmd/replicator/cells.go` — if it imports hive
- `test/parity/parity_test.go` — 4 import paths + aliases
- `internal/tools/org/tools.go` — package declaration + imports
- `internal/tools/comms/tools.go` — package declaration + imports
- `internal/tools/forge/tools.go` — package declaration + imports
- `internal/org/*.go` — package declarations (hive → org)
- `internal/comms/*.go` — package declarations (swarmmail → comms)
- `internal/forge/*.go` — package declarations (swarm → forge)
- `internal/query/presets.go` — if it references hive/swarm packages
- `internal/stats/stats_test.go` — if it references hive/swarm packages

**Import path mapping**:
```
github.com/unbound-force/replicator/internal/hive      → .../internal/org
github.com/unbound-force/replicator/internal/swarmmail  → .../internal/comms
github.com/unbound-force/replicator/internal/swarm      → .../internal/forge
github.com/unbound-force/replicator/internal/tools/hive      → .../internal/tools/org
github.com/unbound-force/replicator/internal/tools/swarmmail  → .../internal/tools/comms
github.com/unbound-force/replicator/internal/tools/swarm      → .../internal/tools/forge
```

**Risk**: MEDIUM. Many files touched. Compiler will catch all missed imports.

**Checkpoint**: `go vet ./...` passes. `go build ./cmd/replicator` succeeds.

---

### Phase 3: Tool Name String Updates [US1]

Update the 45 MCP tool `Name:` string literals from old prefixes to new prefixes. Also update Go function names for consistency (e.g., `hiveCells` → `orgCells`).

**Tool name mapping** (45 tools):

| Old Name | New Name |
|----------|----------|
| `hive_cells` | `org_cells` |
| `hive_create` | `org_create` |
| `hive_close` | `org_close` |
| `hive_update` | `org_update` |
| `hive_create_epic` | `org_create_epic` |
| `hive_query` | `org_query` |
| `hive_start` | `org_start` |
| `hive_ready` | `org_ready` |
| `hive_sync` | `org_sync` |
| `hive_session_start` | `org_session_start` |
| `hive_session_end` | `org_session_end` |
| `swarmmail_init` | `comms_init` |
| `swarmmail_send` | `comms_send` |
| `swarmmail_inbox` | `comms_inbox` |
| `swarmmail_read_message` | `comms_read_message` |
| `swarmmail_reserve` | `comms_reserve` |
| `swarmmail_release` | `comms_release` |
| `swarmmail_release_all` | `comms_release_all` |
| `swarmmail_release_agent` | `comms_release_agent` |
| `swarmmail_ack` | `comms_ack` |
| `swarmmail_health` | `comms_health` |
| `swarm_init` | `forge_init` |
| `swarm_select_strategy` | `forge_select_strategy` |
| `swarm_plan_prompt` | `forge_plan_prompt` |
| `swarm_decompose` | `forge_decompose` |
| `swarm_validate_decomposition` | `forge_validate_decomposition` |
| `swarm_subtask_prompt` | `forge_subtask_prompt` |
| `swarm_spawn_subtask` | `forge_spawn_subtask` |
| `swarm_complete_subtask` | `forge_complete_subtask` |
| `swarm_progress` | `forge_progress` |
| `swarm_complete` | `forge_complete` |
| `swarm_status` | `forge_status` |
| `swarm_record_outcome` | `forge_record_outcome` |
| `swarm_worktree_create` | `forge_worktree_create` |
| `swarm_worktree_merge` | `forge_worktree_merge` |
| `swarm_worktree_cleanup` | `forge_worktree_cleanup` |
| `swarm_worktree_list` | `forge_worktree_list` |
| `swarm_review` | `forge_review` |
| `swarm_review_feedback` | `forge_review_feedback` |
| `swarm_adversarial_review` | `forge_adversarial_review` |
| `swarm_evaluation_prompt` | `forge_evaluation_prompt` |
| `swarm_broadcast` | `forge_broadcast` |
| `swarm_get_strategy_insights` | `forge_get_strategy_insights` |
| `swarm_get_file_insights` | `forge_get_file_insights` |
| `swarm_get_pattern_insights` | `forge_get_pattern_insights` |

**Unchanged** (8 tools): `hivemind_store`, `hivemind_find`, `hivemind_get`, `hivemind_remove`, `hivemind_validate`, `hivemind_stats`, `hivemind_index`, `hivemind_sync`.

**Files touched**:
- `internal/tools/org/tools.go` — 11 Name strings + 11 function names
- `internal/tools/comms/tools.go` — 10 Name strings + 10 function names
- `internal/tools/forge/tools.go` — 24 Name strings + 24 function names
- Tool description strings that reference old names (e.g., "alias for hive_cells")

**Risk**: LOW. Mechanical find-and-replace. Parity tests will catch any missed renames.

**Checkpoint**: `go test ./internal/tools/... -count=1` passes. `grep -r '"hive_\|"swarmmail_\|"swarm_' internal/tools/` returns zero matches (excluding `hivemind_`).

---

### Phase 4: Parity Fixtures + Docs Category Map [US1, US4]

Update parity test fixtures to use new tool names. Update the `categories` slice in `docs.go` to map new prefixes to new display names.

**Files touched**:
- `test/parity/fixtures/hive.json` — rename tool name keys to `org_*`
- `test/parity/fixtures/swarmmail.json` — rename tool name keys to `comms_*`
- `test/parity/fixtures/swarm.json` — rename tool name keys to `forge_*`
- `test/parity/parity_test.go` — update hardcoded tool name strings (e.g., `"hive_create"` → `"org_create"`, `"hive_cells"` → `"org_cells"`)
- `cmd/replicator/docs.go` — update `categories` slice:
  ```go
  {"org_", "Org"},
  {"comms_", "Comms"},
  {"forge_", "Forge"},
  {"hivemind_", "Memory"},
  ```
- `cmd/replicator/docs.go` — update Long description text
- `internal/query/presets.go` — update any hardcoded tool name references
- `internal/stats/stats_test.go` — update any hardcoded tool name references

**Risk**: LOW. Fixture files are JSON with tool name keys. Mechanical rename.

**Checkpoint**: `go test -tags parity ./test/parity/ -count=1` passes. `replicator docs` output shows Org/Comms/Forge/Memory categories.

---

### Phase 5: Write Agent Kit Content [US2, US3]

Create the `internal/agentkit/` package with `embed.FS` and 15 template files. The content is adapted from the upstream `joelhooks/swarm-tools` Claude Code plugin, rewritten for OpenCode conventions and the new org/comms/forge terminology.

**Files created**:
- `internal/agentkit/content/command/forge.md` — Forge orchestrator command (adapted from upstream `/swarm:swarm`)
- `internal/agentkit/content/command/org.md` — Org management command
- `internal/agentkit/content/command/inbox.md` — Comms inbox command
- `internal/agentkit/content/command/forge-status.md` — Forge status command
- `internal/agentkit/content/command/handoff.md` — Session handoff command
- `internal/agentkit/content/skills/always-on-guidance/SKILL.md`
- `internal/agentkit/content/skills/forge-coordination/SKILL.md`
- `internal/agentkit/content/skills/replicator-cli/SKILL.md`
- `internal/agentkit/content/skills/testing-patterns/SKILL.md`
- `internal/agentkit/content/skills/system-design/SKILL.md`
- `internal/agentkit/content/skills/learning-systems/SKILL.md`
- `internal/agentkit/content/skills/forge-global/SKILL.md`
- `internal/agentkit/content/agents/coordinator.md`
- `internal/agentkit/content/agents/worker.md`
- `internal/agentkit/content/agents/background-worker.md`
- `internal/agentkit/agentkit.go` — `Scaffold(targetDir string, force bool) ([]string, error)` function + `embed.FS`
- `internal/agentkit/agentkit_test.go` — Tests for scaffold, skip-existing, force-overwrite

**Design decisions**:
- Use `embed.FS` (stdlib) — no external dependency, files compiled into binary
- `Scaffold()` returns `[]string` of created/skipped files for CLI output
- Files are written to `.opencode/` under the target directory
- Skip logic: check `os.Stat()` before writing; skip if exists and `force=false`

**Risk**: MEDIUM. Content authoring is the largest creative effort. Template quality matters for agent usability.

**Checkpoint**: `go test ./internal/agentkit/ -count=1` passes. `go build ./cmd/replicator` succeeds.

---

### Phase 6: Enhance `replicator init` [US2]

Integrate the agent kit scaffold into the `init` command. Add `--force` flag.

**Files touched**:
- `cmd/replicator/init.go` — import `agentkit`, call `Scaffold()`, add `--force` flag
- `cmd/replicator/init_test.go` — add tests for agent kit creation, skip-existing, force-overwrite

**Behavior changes**:
1. `replicator init` now creates `.uf/replicator/cells.json` AND scaffolds 15 agent kit files to `.opencode/`
2. `--force` flag overwrites existing agent kit files
3. Without `--force`, existing files are skipped with a "skipped" message
4. Output lists each file created/skipped using styled output

**Risk**: LOW. The `agentkit.Scaffold()` function encapsulates all complexity. Init just calls it.

**Checkpoint**: `go test ./cmd/replicator/ -count=1 -run TestRunInit` passes. Manual smoke test: `replicator init` in temp dir creates 16 files.

---

### Phase 7: Documentation Updates [US4]

Update all documentation to use the new terminology consistently.

**Files touched**:
- `AGENTS.md` — naming convention table, project structure, MCP protocol example, CLI commands table, constitution references
- `README.md` — if it references old tool names
- `.specify/memory/constitution.md` — update "swarm mail" references to "comms", "hive" to "org"
- `specs/003-rename-terminology/spec.md` — mark as implemented

**Risk**: LOW. Text changes only.

**Checkpoint**: `grep -r 'hive_\|swarmmail_\|swarm_' AGENTS.md README.md .specify/memory/constitution.md` returns zero matches (excluding historical attribution and `hivemind_`).

---

### Phase 8: Verification [ALL]

Full CI-equivalent verification pass.

**Commands**:
```bash
# CI parity gate (from .github/workflows/ci.yml)
go vet ./...
go test ./... -count=1 -race
go build -o bin/replicator ./cmd/replicator

# Parity tests
go test -tags parity ./test/parity/ -count=1 -v

# Grep verification (SC-005)
grep -rn '"hive_\|"swarmmail_\|"swarm_' --include='*.go' | grep -v hivemind_
# Expected: zero matches

# Init smoke test (SC-002)
tmpdir=$(mktemp -d)
./bin/replicator init --path "$tmpdir"
find "$tmpdir" -type f | wc -l
# Expected: 16 files

# Tool count verification (SC-001)
# Run replicator serve, call tools/list, count by prefix
```

**Risk**: LOW. Verification only, no code changes.

**Checkpoint**: All commands pass with zero failures.

## Dependency Graph

```
Phase 1 (dir renames)
  └─→ Phase 2 (import paths)
        └─→ Phase 3 (tool name strings)
              └─→ Phase 4 (parity fixtures + docs map)
                    └─→ Phase 8 (verification)

Phase 5 (agent kit content)  ← independent, can start anytime
  └─→ Phase 6 (init enhancement)
        └─→ Phase 8 (verification)

Phase 7 (documentation) ← depends on Phase 3 (needs final names)
  └─→ Phase 8 (verification)
```

Phases 1–4 are strictly sequential (each breaks the build until the next completes).
Phase 5 is independent and can be developed in parallel with Phases 1–4.
Phase 6 depends on Phase 5.
Phase 7 depends on Phase 3 (needs final tool names).
Phase 8 depends on all other phases.

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Missed tool name in string literal | LOW | HIGH (broken tool) | Grep verification in Phase 8 + parity tests |
| Import path typo | LOW | HIGH (build failure) | `go vet` catches immediately |
| Agent kit content quality | MEDIUM | MEDIUM (poor agent UX) | Adapt from proven upstream templates |
| Parity fixture key mismatch | LOW | MEDIUM (false test failure) | Mechanical rename, verified by test run |
| `embed.FS` path issues | LOW | LOW (build failure) | Well-documented stdlib feature, tested in Phase 5 |
