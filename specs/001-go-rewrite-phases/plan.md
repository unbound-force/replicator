# Implementation Plan: Go Rewrite — Remaining Phases

**Branch**: `001-go-rewrite-phases` | **Date**: 2026-04-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-go-rewrite-phases/spec.md`

## Summary

Complete the Go rewrite of cyborg-swarm across 5 phases: finish the hive tool suite and add swarm mail (Phase 1), implement swarm orchestration with git worktree isolation (Phase 2), add memory/context tools proxied through Dewey (Phase 3), build CLI commands for observability (Phase 4), and verify parity against the TypeScript version (Phase 5). The existing Phase 0 scaffold provides the MCP server, SQLite database, tool registry, and 4 hive tools — all remaining work extends these established patterns.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: `cobra` (CLI), `modernc.org/sqlite` (pure Go SQLite), stdlib `encoding/json` (MCP JSON-RPC), stdlib `os/exec` (git operations)
**Storage**: SQLite at `~/.config/uf/replicator/replicator.db` (WAL mode)
**Testing**: `go test` (stdlib only, no testify — per TC-001)
**Target Platform**: macOS/Linux arm64+amd64, Windows amd64 (via goreleaser)
**Project Type**: CLI + MCP server (single binary)
**Performance Goals**: <50ms cold start to first MCP response (SC-002), <20MB binary (SC-004)
**Constraints**: Zero CGo, zero runtime dependencies, single-file distribution
**Scale/Scope**: 70 unique MCP tools, 6 CLI commands, ~80% line coverage target

## Constitution Check

*No `.specify/memory/constitution.md` found in this repo. Applying universal principles from AGENTS.md:*

- **TDD Everything**: All code follows Red-Green-Refactor. Tests use `db.OpenMemory()` for isolation. ✅
- **Single Binary**: No CGo, pure Go SQLite via `modernc.org/sqlite`. ✅
- **Schema Compatibility**: Database schema matches cyborg-swarm's libSQL tables. ✅
- **Convention Pack**: Go pack loaded from `.opencode/unbound/packs/go.md`. All `[MUST]` rules apply. ✅

## Project Structure

### Documentation (this feature)

```text
specs/001-go-rewrite-phases/
├── plan.md              # This file
├── research.md          # Library decisions and technical research
├── quickstart.md        # Build, test, and run commands
├── checklists/
│   └── requirements.md  # Spec quality checklist (complete)
└── tasks.md             # Task breakdown (created by /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/replicator/              # CLI entrypoint (cobra commands)
├── main.go                  # Root command, subcommand registration
├── serve.go                 # MCP server startup
├── cells.go                 # cells command (exists)
├── doctor.go                # Phase 4: health checks
├── stats.go                 # Phase 4: activity metrics
├── query.go                 # Phase 4: SQL analytics
└── setup.go                 # Phase 4: environment setup

internal/
├── config/                  # Configuration (exists)
│   └── config.go
├── db/                      # SQLite + migrations (exists)
│   ├── db.go
│   ├── db_test.go
│   └── migrations.go        # Extended with swarmmail/session tables
├── hive/                    # Cell domain logic (exists, extended)
│   ├── cells.go             # Existing: Create, Query, Close, Update
│   ├── cells_test.go
│   ├── epic.go              # Phase 1: CreateEpic (atomic epic+subtasks)
│   ├── epic_test.go
│   ├── session.go           # Phase 1: SessionStart, SessionEnd
│   ├── session_test.go
│   ├── sync.go              # Phase 1: Sync to git
│   └── sync_test.go
├── mcp/                     # MCP JSON-RPC server (exists)
│   ├── server.go
│   └── server_test.go
├── swarmmail/               # Phase 1: Agent messaging
│   ├── agent.go             # Agent registration, init
│   ├── agent_test.go
│   ├── message.go           # Send, inbox, read, ack
│   ├── message_test.go
│   ├── reservation.go       # File reservations (lock/unlock)
│   └── reservation_test.go
├── swarm/                   # Phase 2: Orchestration
│   ├── decompose.go         # Task decomposition prompts
│   ├── decompose_test.go
│   ├── worktree.go          # Git worktree create/merge/cleanup/list
│   ├── worktree_test.go
│   ├── progress.go          # Progress tracking, status, completion
│   ├── progress_test.go
│   ├── review.go            # Review prompts, feedback
│   ├── review_test.go
│   ├── spawn.go             # Subtask prompt generation, spawning
│   └── spawn_test.go
├── memory/                  # Phase 3: Dewey proxy
│   ├── proxy.go             # HTTP client to Dewey MCP endpoint
│   ├── proxy_test.go
│   ├── deprecated.go        # Deprecation stubs for secondary tools
│   └── deprecated_test.go
├── doctor/                  # Phase 4: Health checks
│   ├── checks.go            # Individual check functions
│   └── checks_test.go
├── stats/                   # Phase 4: Activity metrics
│   ├── stats.go
│   └── stats_test.go
├── query/                   # Phase 4: SQL analytics presets
│   ├── presets.go
│   └── presets_test.go
├── gitutil/                 # Shared git helpers (exec wrappers)
│   ├── git.go
│   └── git_test.go
└── tools/                   # MCP tool registrations (exists)
    ├── registry/
    │   └── registry.go
    ├── hive/
    │   └── tools.go         # Existing 4 tools + 7 new hive tools
    ├── swarmmail/            # Phase 1: 9 swarmmail tools
    │   └── tools.go
    ├── swarm/                # Phase 2: 16 swarm tools
    │   └── tools.go
    └── memory/               # Phase 3: 8 memory tools
        └── tools.go

test/                        # Phase 5: Parity testing
├── parity/
│   ├── parity_test.go       # Shape comparison harness
│   └── fixtures/            # Captured TypeScript responses
└── README.md
```

**Structure Decision**: Follows the existing Go project layout established in Phase 0. Domain logic lives in `internal/{domain}/` with corresponding tool registrations in `internal/tools/{domain}/`. Each domain package owns its types, business logic, and tests. The `cmd/replicator/` layer delegates to domain packages via `Run()` functions per AP-002. No new top-level directories are introduced except `test/` for cross-cutting parity tests in Phase 5.

---

## Phase 1: Hive Completion + Swarm Mail

**User Stories**: US1 (Complete Hive Tool Suite), US2 (Agent Messaging)
**Requirements**: FR-001 through FR-007, FR-020, FR-022
**Estimated Effort**: 2-3 weeks

### 1.1 Remaining Hive Tools (7 tools)

Complete the hive tool suite by adding the 7 remaining tools to match the TypeScript version.

| Tool | Domain Function | Package | Notes |
|------|----------------|---------|-------|
| `hive_create_epic` | `hive.CreateEpic()` | `internal/hive/epic.go` | Atomic: create epic + N subtasks in one transaction |
| `hive_query` | `hive.QueryCells()` | existing | Already implemented as `hive_cells` — verify schema parity, may need alias |
| `hive_start` | `hive.StartCell()` | `internal/hive/cells.go` | Set status=in_progress, record timestamp |
| `hive_ready` | `hive.ReadyCell()` | `internal/hive/cells.go` | Return highest-priority unblocked cell |
| `hive_sync` | `hive.Sync()` | `internal/hive/sync.go` | Serialize cells to `.uf/replicator/`, git add+commit |
| `hive_session_start` | `hive.SessionStart()` | `internal/hive/session.go` | Create session, return previous handoff notes |
| `hive_session_end` | `hive.SessionEnd()` | `internal/hive/session.go` | Save handoff notes for next session |

**Implementation approach**:
- `hive_create_epic`: Use `sql.Tx` for atomicity — insert epic row, then insert all subtask rows with `parent_id` pointing to the epic. Roll back on any failure.
- `hive_start`: Simple status update with timestamp. Reuse `UpdateCell` pattern.
- `hive_ready`: Query cells where `status='open'` and no parent cell is still open (unblocked). Order by `priority DESC`. Return first match.
- `hive_sync`: Shell out to `git` via `internal/gitutil/` to add and commit `.uf/replicator/` directory contents.
- `hive_session_start/end`: New `sessions` table in SQLite. Store handoff notes as JSON. Return previous session's notes on start.

**Database changes**:
- Add `sessions` table migration (session_id, agent_name, started_at, ended_at, handoff_notes, active_cell_id)

**Testing strategy**:
- Unit tests for each domain function using `db.OpenMemory()`
- `TestCreateEpic_AtomicRollback`: Verify that if subtask 3 of 5 fails, no rows are committed
- `TestReadyCell_UnblockedOnly`: Create parent+child cells, verify only unblocked cells are returned
- `TestSync_CommitsToGit`: Use `t.TempDir()` with a git repo, verify commit is created
- `TestSessionStart_ReturnsHandoff`: Create a session with handoff notes, start a new session, verify notes are returned

### 1.2 Swarm Mail (9 tools)

Implement the agent messaging and file reservation system.

| Tool | Domain Function | Package | Notes |
|------|----------------|---------|-------|
| `swarmmail_init` | `swarmmail.Init()` | `internal/swarmmail/agent.go` | Register agent, create session |
| `swarmmail_send` | `swarmmail.Send()` | `internal/swarmmail/message.go` | Persist message to SQLite |
| `swarmmail_inbox` | `swarmmail.Inbox()` | `internal/swarmmail/message.go` | Fetch messages (bodies excluded, max 5) |
| `swarmmail_read_message` | `swarmmail.ReadMessage()` | `internal/swarmmail/message.go` | Fetch one message body by ID |
| `swarmmail_ack` | `swarmmail.Ack()` | `internal/swarmmail/message.go` | Mark message as acknowledged |
| `swarmmail_reserve` | `swarmmail.Reserve()` | `internal/swarmmail/reservation.go` | Lock file paths for exclusive editing |
| `swarmmail_release` | `swarmmail.Release()` | `internal/swarmmail/reservation.go` | Release file reservations |
| `swarmmail_release_all` | `swarmmail.ReleaseAll()` | `internal/swarmmail/reservation.go` | Coordinator override: release all |
| `swarmmail_release_agent` | `swarmmail.ReleaseAgent()` | `internal/swarmmail/reservation.go` | Release all reservations for an agent |

**Database changes**:
- Add `messages` table migration (id, from_agent, to_agents JSON, subject, body, importance, thread_id, ack_required, acknowledged, created_at)
- Add `reservations` table migration (id, agent_name, path, exclusive, reason, ttl_seconds, created_at, expires_at)

**Implementation approach**:
- Messages: Standard CRUD with SQLite. `swarmmail_inbox` returns headers only (no body) for context efficiency. `swarmmail_read_message` fetches the full body.
- Reservations: `swarmmail_reserve` checks for existing exclusive reservations on the same path before inserting. Uses a transaction to prevent race conditions. Expired reservations (past TTL) are treated as released. Expiry is checked lazily — expired reservations are filtered out during `Reserve` conflict checks, not cleaned up proactively.
- Agent init: Upsert into `agents` table (existing schema in `internal/db/migrations.go` — columns `name`, `project_path`, `task_description`, `last_seen_at` map directly to `swarmmail_init` parameters), update `last_seen_at`.

**Testing strategy**:
- `TestReserve_ExclusiveConflict`: Agent A reserves a file, Agent B's reserve attempt fails with conflict error
- `TestReserve_TTLExpiry`: Reserve with short TTL, wait, verify re-reservation succeeds
- `TestInbox_ExcludesBodies`: Send a message, verify inbox response has no body field
- `TestSend_PersistsAcrossRestart`: Send message, close store, reopen, verify message is still there

### 1.3 Phase 1 Checkpoint

**Gate**: All 20 tools (4 existing + 7 hive + 9 swarmmail) pass tests. `make check` succeeds.

**Verification**:
```bash
go test ./... -count=1 -race
go vet ./...
```

**Expected test count**: ~50 tests (16 existing + ~34 new)

---

## Phase 2: Swarm Orchestration

**User Story**: US3 (Swarm Orchestration)
**Requirements**: FR-008 through FR-011
**Estimated Effort**: 3-4 weeks

### 2.1 Git Utilities

Before implementing orchestration tools, establish a shared `internal/gitutil/` package for git operations.

| Function | Description |
|----------|-------------|
| `gitutil.Run(dir string, args ...string) (string, error)` | Execute git command in directory, return stdout |
| `gitutil.WorktreeAdd(projectPath, worktreePath, branch, startCommit string) error` | `git worktree add` |
| `gitutil.WorktreeRemove(worktreePath string) error` | `git worktree remove` |
| `gitutil.WorktreeList(projectPath string) ([]Worktree, error)` | `git worktree list --porcelain` |
| `gitutil.CherryPick(projectPath, startCommit string) error` | Cherry-pick commits from worktree branch |
| `gitutil.CurrentCommit(projectPath string) (string, error)` | `git rev-parse HEAD` |

**Implementation approach**:
- Shell out to `git` binary via `os/exec`. This is the spec's stated approach (see Assumptions).
- All functions accept a directory path — no global state.
- Error messages include the git stderr output for debugging.

**Testing strategy**:
- Tests create real git repos in `t.TempDir()` using `git init`.
- `TestWorktreeAdd_CreatesDirectory`: Create a worktree, verify the directory exists and has the correct branch.
- `TestCherryPick_AppliesCommits`: Create a worktree, make commits, cherry-pick back, verify commits appear on main.
- `TestCherryPick_DetectsConflict`: Create conflicting changes, verify cherry-pick returns an error listing conflicting files.
- Guard with `testing.Short()` per TC-011 since these spawn subprocesses.

### 2.2 Orchestration Tools (16 tools)

| Tool | Domain Function | Package | Notes |
|------|----------------|---------|-------|
| `swarm_init` | `swarm.Init()` | `internal/swarm/` | Initialize swarm session |
| `swarm_select_strategy` | `swarm.SelectStrategy()` | `internal/swarm/decompose.go` | Recommend decomposition strategy |
| `swarm_plan_prompt` | `swarm.PlanPrompt()` | `internal/swarm/decompose.go` | Generate strategy-specific prompt |
| `swarm_decompose` | `swarm.Decompose()` | `internal/swarm/decompose.go` | Generate decomposition prompt |
| `swarm_validate_decomposition` | `swarm.ValidateDecomposition()` | `internal/swarm/decompose.go` | Validate response against schema |
| `swarm_subtask_prompt` | `swarm.SubtaskPrompt()` | `internal/swarm/spawn.go` | Generate prompt for spawned agent |
| `swarm_spawn_subtask` | `swarm.SpawnSubtask()` | `internal/swarm/spawn.go` | Prepare subtask for Task tool |
| `swarm_complete_subtask` | `swarm.CompleteSubtask()` | `internal/swarm/spawn.go` | Handle subtask completion |
| `swarm_progress` | `swarm.Progress()` | `internal/swarm/progress.go` | Report progress on subtask |
| `swarm_complete` | `swarm.Complete()` | `internal/swarm/progress.go` | Mark subtask complete with verification |
| `swarm_status` | `swarm.Status()` | `internal/swarm/progress.go` | Get swarm status by epic ID |
| `swarm_record_outcome` | `swarm.RecordOutcome()` | `internal/swarm/progress.go` | Record outcome for feedback scoring |
| `swarm_worktree_create` | `swarm.WorktreeCreate()` | `internal/swarm/worktree.go` | Create isolated git worktree |
| `swarm_worktree_merge` | `swarm.WorktreeMerge()` | `internal/swarm/worktree.go` | Cherry-pick commits back to main |
| `swarm_worktree_cleanup` | `swarm.WorktreeCleanup()` | `internal/swarm/worktree.go` | Remove worktree (idempotent) |
| `swarm_worktree_list` | `swarm.WorktreeList()` | `internal/swarm/worktree.go` | List active worktrees |

**Implementation approach**:
- Decomposition tools (`swarm_decompose`, `swarm_plan_prompt`, etc.) are primarily prompt generators — they construct structured text from inputs. No LLM calls; the agent calling the tool provides the LLM.
- Progress tracking uses the `events` table (already exists) to record progress events with `type='swarm_progress'` and a JSON payload.
- Worktree tools delegate to `internal/gitutil/`.
- `swarm_complete` runs verification (typecheck + tests) before allowing completion. The verification commands come from the project's build configuration.
- Review tools (`swarm_review`, `swarm_review_feedback`, `swarm_adversarial_review`, `swarm_evaluation_prompt`, `swarm_broadcast`, `swarm_get_strategy_insights`, `swarm_get_file_insights`, `swarm_get_pattern_insights`) are deferred to a follow-up — they depend on the events/outcomes data being populated first. The 16 core tools above are the priority.

**Database changes**:
- No new tables needed — progress events use the existing `events` table.
- Worktree state is tracked via the filesystem (git worktree list), not in SQLite.

**Testing strategy**:
- Decomposition tests verify prompt structure (contains expected sections, file lists, etc.)
- Worktree tests use real git repos in `t.TempDir()`
- Progress tests verify event recording and status aggregation
- `TestWorktreeCreate_UncommittedChanges`: Verify worktree creation succeeds even with dirty working directory (per edge case in spec)
- `TestWorktreeMerge_Conflict`: Verify merge failure returns conflicting file list and does NOT clean up the worktree (per edge case in spec)

### 2.3 Phase 2 Checkpoint

**Gate**: All 36 tools (20 from Phase 1 + 16 orchestration) pass tests. Git worktree integration tests pass.

**Expected test count**: ~85 tests

---

## Phase 3: Memory and Context

**User Story**: US4 (Memory and Context)
**Requirements**: FR-012 through FR-014
**Estimated Effort**: 1-2 weeks

### 3.1 Dewey Proxy Client

Create an HTTP client that proxies memory operations to Dewey's MCP endpoint.

| Component | Description |
|-----------|-------------|
| `memory.Client` | HTTP client struct with Dewey URL from config |
| `memory.Client.Call(method string, params any) (json.RawMessage, error)` | Send JSON-RPC request to Dewey |
| `memory.Client.Health() error` | Check if Dewey is reachable |

**Implementation approach**:
- Use stdlib `net/http` — no external HTTP client library needed.
- Dewey URL comes from `config.DeweyURL` (default: `http://localhost:3333/mcp/`).
- All calls are JSON-RPC 2.0 over HTTP POST.
- When Dewey is unreachable, return a structured error with code `DEWEY_UNAVAILABLE`.

### 3.2 Memory Tools (8 tools)

| Tool | Behavior | Notes |
|------|----------|-------|
| `hivemind_store` | Proxy to Dewey `dewey_dewey_store_learning` | Include deprecation warning in response |
| `hivemind_find` | Proxy to Dewey `dewey_dewey_semantic_search` | Include deprecation warning in response |
| `hivemind_get` | Return deprecation message | Point to `dewey_get_page` |
| `hivemind_remove` | Return deprecation message | Point to `dewey_delete_page` |
| `hivemind_validate` | Return deprecation message | No direct Dewey equivalent |
| `hivemind_stats` | Return deprecation message | Point to `dewey_health` |
| `hivemind_index` | Return deprecation message | Point to `dewey_reload` |
| `hivemind_sync` | Return deprecation message | Point to `dewey_reload` |

**Implementation approach**:
- `hivemind_store` and `hivemind_find` are live proxies — they forward the request to Dewey and return the response with an added deprecation notice.
- The remaining 6 tools return a structured deprecation message: `{"deprecated": true, "message": "...", "replacement": "dewey_xxx"}`.
- All tools handle Dewey unavailability gracefully per FR-014.

**Testing strategy**:
- Use `httptest.NewServer` to mock Dewey responses.
- `TestStore_ProxiesToDewey`: Verify the correct JSON-RPC method and params are forwarded.
- `TestStore_DeweyUnavailable`: Verify structured error when Dewey is down.
- `TestDeprecated_ReturnsMessage`: Verify each deprecated tool returns the correct replacement name.

### 3.3 Phase 3 Checkpoint

**Gate**: All 44 tools (36 + 8 memory) pass tests. Dewey proxy handles both available and unavailable states.

**Expected test count**: ~100 tests

---

## Phase 4: CLI Commands

**User Story**: US5 (CLI Operations)
**Requirements**: FR-015 through FR-017
**Estimated Effort**: 1-2 weeks

### 4.1 CLI Commands

| Command | Description | Domain Package |
|---------|-------------|----------------|
| `replicator doctor` | Check dependencies and health | `internal/doctor/` |
| `replicator stats` | Display activity metrics | queries `events` table |
| `replicator query <preset>` | Run SQL analytics queries | preset SQL in `internal/query/` |
| `replicator setup` | Initialize environment | create config dir, verify git, etc. |

**Existing commands** (no changes needed):
- `replicator serve` — MCP server (exists)
- `replicator cells` — List cells (exists)
- `replicator version` — Print version (exists)

### 4.2 Doctor Command

The `doctor` command checks for required dependencies and reports their status.

| Check | Description | Pass Criteria |
|-------|-------------|---------------|
| Git | `git --version` | Exit code 0, version ≥ 2.20 |
| Database | Open and ping SQLite | No error |
| Dewey | HTTP GET to health endpoint | 200 OK |
| Config dir | `~/.config/uf/replicator/` exists | Directory exists |

**Implementation approach** (per AP-002, AP-003):
- `doctor.Run(opts Options) (*Result, error)` — core logic, testable
- `Options` includes `io.Writer` for output, `Config` for paths
- Each check returns a `CheckResult{Name, Status, Message, Duration}`
- Output formatted as a table to stdout
- Must complete within 2 seconds (SC-008)

**Testing strategy**:
- Use `Options` struct with `bytes.Buffer` as writer per AP-003
- Mock external checks (git, Dewey) via interface injection
- `TestDoctor_AllPass`: All checks succeed, output contains ✓ for each
- `TestDoctor_DeweyDown`: Dewey check fails, output shows ✗ with message

### 4.3 Stats and Query Commands

- `stats`: Aggregate events table — count by type, recent activity, active agents. Pure SQL queries.
- `query`: Predefined SQL analytics queries (e.g., "agent activity last 24h", "cells by status", "swarm completion rate"). Queries are embedded as Go constants.

### 4.4 Phase 4 Checkpoint

**Gate**: All CLI commands execute correctly. `replicator doctor` completes in <2s. Binary starts in <50ms.

**Verification**:
```bash
time ./bin/replicator version  # Verify <50ms startup
./bin/replicator doctor        # Verify all checks run
```

**Expected test count**: ~115 tests

---

## Phase 5: Parity Testing

**User Story**: US6 (Parity Verification)
**Requirements**: FR-018, FR-019
**Estimated Effort**: 2-3 weeks

### 5.1 Parity Test Harness

Build a test harness that compares MCP tool response shapes between the Go replicator and the TypeScript cyborg-swarm.

**Approach**:
1. **Fixture capture**: Run each tool against the TypeScript server, capture response JSON to `test/parity/fixtures/`.
2. **Shape comparison**: For each tool, send the same arguments to the Go server and compare the response shape (field names, types, nesting) — not exact values.
3. **Report generation**: Output a table of all tools with pass/fail status and any shape differences.

**Implementation**:
- `test/parity/parity_test.go` — Go test file that iterates over fixtures
- `ShapeMatch(expected, actual json.RawMessage) (bool, []Difference)` — recursive JSON shape comparator
- Differences report: field name, expected type, actual type, path (e.g., `$.content[0].text`)

**Testing strategy**:
- Each tool gets a fixture file: `test/parity/fixtures/hive_create.json` containing `{request, typescript_response}`
- The test sends the same request to the Go server and compares shapes
- Tests are tagged with `//go:build parity` so they don't run in normal CI (they require the TypeScript server)

### 5.2 Phase 5 Checkpoint

**Gate**: Parity report shows 100% shape match for all implemented tools (SC-003).

**Expected final test count**: ~130 tests (115 unit/integration + ~15 parity)

---

## Cross-Cutting Concerns

### Database Migration Strategy

All new tables use `CREATE TABLE IF NOT EXISTS` for idempotency, matching the existing pattern in `internal/db/migrations.go`. New migrations are added to the `migrations` slice in order. No versioned migration system is needed — the IF NOT EXISTS pattern handles re-runs. When the Go binary opens an existing cyborg-swarm database, new tables (`sessions`, `messages`, `reservations`) are created automatically via IF NOT EXISTS migrations. Existing data is preserved.

### Error Handling Pattern

Follow the established pattern from `internal/hive/cells.go`:
```go
result, err := store.DB.Exec(query, args...)
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

All errors are wrapped with context per CS-006. Tool execute functions return errors that the MCP server wraps in a content block.

### Tool Registration Pattern

Follow the established pattern from `internal/tools/hive/tools.go`:
1. Domain logic in `internal/{domain}/` — pure functions that accept `*db.Store`
2. Tool registration in `internal/tools/{domain}/tools.go` — thin wrappers that unmarshal JSON args, call domain functions, marshal results
3. Registration function: `Register(reg *registry.Registry, store *db.Store)`
4. Called from `cmd/replicator/serve.go`

### Testing Pattern

Follow the established pattern from `internal/hive/cells_test.go`:
```go
func testStore(t *testing.T) *db.Store {
    t.Helper()
    store, err := db.OpenMemory()
    if err != nil {
        t.Fatalf("OpenMemory: %v", err)
    }
    t.Cleanup(func() { store.Close() })
    return store
}
```

- In-memory SQLite for unit tests
- `t.TempDir()` for filesystem tests
- `t.Helper()` on all test helpers
- `t.Cleanup()` for teardown
- No testify — stdlib `testing` only (TC-001)

## Complexity Tracking

No constitution violations identified. The project structure follows standard Go conventions with a single binary, no external services (except optional Dewey), and a single SQLite database.

| Concern | Decision | Rationale |
|---------|----------|-----------|
| Git operations | Shell out to `git` binary | Spec assumption. Pure-Go git libs (go-git) add ~15MB to binary and don't support worktrees well. |
| Dewey proxy | HTTP client, not embedded | Dewey is a separate service. Embedding would violate single-responsibility. |
| No mcp-go library | Hand-rolled JSON-RPC | Existing server.go is 212 lines and handles all needed methods. Adding a library would increase binary size and coupling for no benefit. See research.md. |
