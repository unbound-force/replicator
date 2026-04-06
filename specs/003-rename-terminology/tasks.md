# Tasks: Terminology Rename + Agent Kit

**Input**: Design documents from `/specs/003-rename-terminology/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, quickstart.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Directory Renames [US1]

**Purpose**: Rename 6 Go package directories using `git mv` to preserve history. Foundation for all subsequent phases.

**⚠️ CRITICAL**: Build will break after this phase until Phase 2 completes. Phases 1+2 must be committed together.

- [x] T001 [P] [US1] Rename domain package `internal/hive/` → `internal/org/` via `git mv internal/hive internal/org`
- [x] T002 [P] [US1] Rename domain package `internal/swarmmail/` → `internal/comms/` via `git mv internal/swarmmail internal/comms`
- [x] T003 [P] [US1] Rename domain package `internal/swarm/` → `internal/forge/` via `git mv internal/swarm internal/forge`
- [x] T004 [P] [US1] Rename tool handler package `internal/tools/hive/` → `internal/tools/org/` via `git mv internal/tools/hive internal/tools/org`
- [x] T005 [P] [US1] Rename tool handler package `internal/tools/swarmmail/` → `internal/tools/comms/` via `git mv internal/tools/swarmmail internal/tools/comms`
- [x] T006 [P] [US1] Rename tool handler package `internal/tools/swarm/` → `internal/tools/forge/` via `git mv internal/tools/swarm internal/tools/forge`

**Checkpoint**: `ls internal/org/ internal/comms/ internal/forge/ internal/tools/org/ internal/tools/comms/ internal/tools/forge/` succeeds. Old directories do not exist.

---

## Phase 2: Package Declarations + Import Paths [US1]

**Purpose**: Update all `package` declarations in renamed files and fix all import paths across the codebase. Restores compilability.

**⚠️ CRITICAL**: Must complete before `go build` works. No task in this phase can be skipped.

### Package declarations (inside renamed directories)

- [x] T007 [P] [US1] Update `package hive` → `package org` in all 11 files under `internal/org/*.go` (cells.go, cells_test.go, epic.go, epic_test.go, format.go, format_test.go, session.go, session_test.go, start_ready_test.go, sync.go, sync_test.go)
- [x] T008 [P] [US1] Update `package swarmmail` → `package comms` in all 6 files under `internal/comms/*.go` (agent.go, agent_test.go, message.go, message_test.go, reservation.go, reservation_test.go)
- [x] T009 [P] [US1] Update `package swarm` → `package forge` in all 14 files under `internal/forge/*.go` (decompose.go, decompose_test.go, init.go, init_test.go, insights.go, insights_test.go, progress.go, progress_test.go, review.go, review_test.go, spawn.go, spawn_test.go, worktree.go, worktree_test.go). Also update the package-level GoDoc comment in `init.go`.
- [x] T010 [P] [US1] Update `package hive` → `package org` and package-level GoDoc comment in `internal/tools/org/tools.go`
- [x] T011 [P] [US1] Update `package swarmmail` → `package comms` and package-level GoDoc comment in `internal/tools/comms/tools.go`
- [x] T012 [P] [US1] Update `package swarm` → `package forge` and package-level GoDoc comment in `internal/tools/forge/tools.go`

### Import paths (consumers of renamed packages)

- [x] T013 [US1] Update import paths in `internal/tools/org/tools.go`: change `internal/hive` → `internal/org` and update all `hive.` qualifier references to `org.` (e.g., `hive.QueryCells` → `org.QueryCells`, `hive.CellQuery` → `org.CellQuery`)
- [x] T014 [US1] Update import paths in `internal/tools/comms/tools.go`: change `internal/swarmmail` → `internal/comms` and update all `swarmmail.` qualifier references to `comms.` (e.g., `swarmmail.Init` → `comms.Init`)
- [x] T015 [US1] Update import paths in `internal/tools/forge/tools.go`: change `internal/swarm` → `internal/forge` and update all `swarm.` qualifier references to `forge.` (e.g., `swarm.Init` → `forge.Init`, `swarm.SubtaskPrompt` → `forge.SubtaskPrompt`)
- [x] T016 [US1] Update import paths and aliases in `cmd/replicator/serve.go`: change `internal/tools/hive` → `internal/tools/org`, `internal/tools/swarmmail` → `internal/tools/comms` (alias `commstools`), `internal/tools/swarm` → `internal/tools/forge` (alias `forgetools`). Update `hive.Register` → `org.Register`, `swarmmailtools.Register` → `commstools.Register`, `swarmtools.Register` → `forgetools.Register`.
- [x] T017 [US1] Update import paths and aliases in `cmd/replicator/docs.go`: change `internal/tools/hive` → `internal/tools/org`, `internal/tools/swarmmail` → `internal/tools/comms` (alias `commstools`), `internal/tools/swarm` → `internal/tools/forge` (alias `forgetools`). Update `hive.Register` → `org.Register`, `swarmmailtools.Register` → `commstools.Register`, `swarmtools.Register` → `forgetools.Register`.
- [x] T018 [US1] Update import paths in `cmd/replicator/cells.go`: change `internal/hive` → `internal/org` and update all `hive.` qualifier references to `org.` (e.g., `hive.QueryCells` → `org.QueryCells`, `hive.FormatCells` → `org.FormatCells`, `hive.CellQuery` → `org.CellQuery`)
- [x] T019 [US1] Update import paths and aliases in `cmd/replicator/docs_test.go`: change `internal/tools/hive` → `internal/tools/org`, `internal/tools/swarmmail` → `internal/tools/comms` (alias `commstools`), `internal/tools/swarm` → `internal/tools/forge` (alias `forgetools`). Update all `hive.Register` → `org.Register`, `swarmmailtools.Register` → `commstools.Register`, `swarmtools.Register` → `forgetools.Register`.
- [x] T020 [US1] Update import paths in `internal/mcp/server_test.go`: change `internal/tools/hive` → `internal/tools/org` and update `hive.Register` → `org.Register`
- [x] T021 [US1] Update import paths and aliases in `test/parity/parity_test.go`: change `internal/tools/hive` → `internal/tools/org` (alias `orgtools`), `internal/tools/swarmmail` → `internal/tools/comms` (alias `commstools`), `internal/tools/swarm` → `internal/tools/forge` (alias `forgetools`). Update `hivetools.Register` → `orgtools.Register`, `swarmmailtools.Register` → `commstools.Register`, `swarmtools.Register` → `forgetools.Register`.

**Checkpoint**: `go vet ./...` passes. `go build ./cmd/replicator` succeeds.

---

## Phase 3: Tool Name Strings + Function Names [US1]

**Purpose**: Update the 45 MCP tool `Name:` string literals and Go function names from old prefixes to new prefixes. Also update event type strings in domain packages and prompt-generating code.

### Tool handler name strings (3 files, parallelizable)

- [x] T022 [P] [US1] In `internal/tools/org/tools.go`: rename all 11 tool `Name:` strings from `hive_*` → `org_*` (hive_cells→org_cells, hive_create→org_create, hive_close→org_close, hive_update→org_update, hive_create_epic→org_create_epic, hive_query→org_query, hive_start→org_start, hive_ready→org_ready, hive_sync→org_sync, hive_session_start→org_session_start, hive_session_end→org_session_end). Rename 11 Go function names from `hive*` → `org*` (e.g., `hiveCells` → `orgCells`, `hiveCreate` → `orgCreate`). Update the `Register` function's GoDoc comment. Update any description strings referencing old names.
- [x] T023 [P] [US1] In `internal/tools/comms/tools.go`: rename all 10 tool `Name:` strings from `swarmmail_*` → `comms_*` (swarmmail_init→comms_init, swarmmail_send→comms_send, swarmmail_inbox→comms_inbox, swarmmail_read_message→comms_read_message, swarmmail_reserve→comms_reserve, swarmmail_release→comms_release, swarmmail_release_all→comms_release_all, swarmmail_release_agent→comms_release_agent, swarmmail_ack→comms_ack, swarmmail_health→comms_health). Rename 10 Go function names from `swarmmail*` → `comms*` (e.g., `swarmmailInit` → `commsInit`). Update the `Register` function's GoDoc comment. Update any description strings referencing old names (e.g., "swarm mail" → "comms").
- [x] T024 [P] [US1] In `internal/tools/forge/tools.go`: rename all 24 tool `Name:` strings from `swarm_*` → `forge_*` (swarm_init→forge_init, swarm_select_strategy→forge_select_strategy, swarm_plan_prompt→forge_plan_prompt, swarm_decompose→forge_decompose, swarm_validate_decomposition→forge_validate_decomposition, swarm_subtask_prompt→forge_subtask_prompt, swarm_spawn_subtask→forge_spawn_subtask, swarm_complete_subtask→forge_complete_subtask, swarm_progress→forge_progress, swarm_complete→forge_complete, swarm_status→forge_status, swarm_record_outcome→forge_record_outcome, swarm_worktree_create→forge_worktree_create, swarm_worktree_merge→forge_worktree_merge, swarm_worktree_cleanup→forge_worktree_cleanup, swarm_worktree_list→forge_worktree_list, swarm_review→forge_review, swarm_review_feedback→forge_review_feedback, swarm_adversarial_review→forge_adversarial_review, swarm_evaluation_prompt→forge_evaluation_prompt, swarm_broadcast→forge_broadcast, swarm_get_strategy_insights→forge_get_strategy_insights, swarm_get_file_insights→forge_get_file_insights, swarm_get_pattern_insights→forge_get_pattern_insights). Rename 24 Go function names from `swarm*` → `forge*`. Update the `Register` function's GoDoc comment. Update any description strings referencing old names.

### Event type strings in domain packages

- [x] T025 [P] [US1] In `internal/forge/init.go`: update event type string `"swarm_init"` → `"forge_init"` in the DB insert and error message. Also in `internal/forge/init_test.go`: update `"swarm_init"` → `"forge_init"` in the DB query assertion (line 64) and error message (line 66)
- [x] T026 [P] [US1] In `internal/forge/progress.go`: update event type strings `"swarm_progress"` → `"forge_progress"`, `"swarm_complete"` → `"forge_complete"`, `"swarm_outcome"` → `"forge_outcome"` in all DB inserts and error messages. Also in `internal/forge/progress_test.go`: update `"swarm_progress"` → `"forge_progress"` (line 16), `"swarm_complete"` → `"forge_complete"` (line 49), `"swarm_outcome"` → `"forge_outcome"` (line 136) in DB query assertions
- [x] T027 [P] [US1] In `internal/forge/review.go`: update event type string `"swarm_broadcast"` → `"forge_broadcast"` in the DB insert. Also in `internal/forge/review_test.go`: update `"swarm_broadcast"` → `"forge_broadcast"` in DB query assertions (lines 139, 155)
- [x] T028 [P] [US1] In `internal/forge/insights.go`: update all `'swarm_outcome'` → `'forge_outcome'` in SQL queries (3 occurrences across GetStrategyInsights, GetFileInsights, GetPatternInsights)

### Prompt-generating code (tool name references in generated text)

- [x] T029 [P] [US1] In `internal/forge/spawn.go`: update tool name references in generated prompt text: `"swarmmail_reserve"` → `"comms_reserve"`, `"swarm_progress"` → `"forge_progress"`, `"swarm_complete"` → `"forge_complete"`
- [x] T030 [P] [US1] In `internal/forge/spawn_test.go`: update expected tool name strings in test assertions: `"swarmmail_reserve"` → `"comms_reserve"`, `"swarm_complete"` → `"forge_complete"`

**Checkpoint**: `go test ./internal/tools/... -count=1` passes. `go test ./internal/forge/... -count=1` passes. `grep -rn '"hive_\|"swarmmail_\|"swarm_' --include='*.go' internal/tools/ internal/forge/ | grep -v hivemind_` returns zero matches.

---

## Phase 4: Parity Fixtures + Docs + Query/Stats [US1, US4]

**Purpose**: Update parity test fixtures, docs category map, MCP server tests, and query/stats packages to use new tool names.

### Parity test fixtures (JSON files)

- [x] T031 [P] [US1] In `test/parity/fixtures/hive.json`: rename all top-level tool name keys from `hive_*` → `org_*` and update `request.name` fields inside each fixture entry to match
- [x] T032 [P] [US1] In `test/parity/fixtures/swarmmail.json`: rename all top-level tool name keys from `swarmmail_*` → `comms_*` and update `request.name` fields inside each fixture entry to match
- [x] T033 [P] [US1] In `test/parity/fixtures/swarm.json`: rename all top-level tool name keys from `swarm_*` → `forge_*` and update `request.name` fields inside each fixture entry to match

### Parity test code

- [x] T034 [US1] In `test/parity/parity_test.go`: update all hardcoded tool name strings — `"hive_cells"` → `"org_cells"`, `"hive_query"` → `"org_query"`, `"hive_ready"` → `"org_ready"`, `"hive_create"` → `"org_create"`, `"hive_close"` → `"org_close"`, `"hive_update"` → `"org_update"`, `"hive_start"` → `"org_start"`, `"hive_cells_with_data"` → `"org_cells_with_data"`, `"swarm_init"` → `"forge_init"`, `"swarm_worktree_list"` → `"forge_worktree_list"`. Update function names `testCellsWithData` references and all `reg.Get("hive_*")` calls to `reg.Get("org_*")`.

### Docs command

- [x] T035 [US1] In `cmd/replicator/docs.go`: update `categories` slice — `{"hive_", "Hive"}` → `{"org_", "Org"}`, `{"swarmmail_", "Swarm Mail"}` → `{"comms_", "Comms"}`, `{"swarm_", "Swarm"}` → `{"forge_", "Forge"}`. Update the `Long` description text from "Hive, Swarm Mail, Swarm, Memory" → "Org, Comms, Forge, Memory".
- [x] T036 [US1] In `cmd/replicator/docs_test.go`: update `TestWriteDocs_HasCategoryHeaders` expected headers from `"## Hive"`, `"## Swarm Mail"`, `"## Swarm"` → `"## Org"`, `"## Comms"`, `"## Forge"`.

### MCP server tests

- [x] T037 [US1] In `internal/mcp/server_test.go`: update all hardcoded tool name strings — `"hive_cells"` → `"org_cells"`, `"hive_create"` → `"org_create"`, `"hive_close"` → `"org_close"`, etc. in `TestToolsList`, `TestToolsCall_HiveCells_Empty`, `TestToolsCall_HiveCreate`, `TestToolsCall_CreateThenQuery`, `TestToolsCall_LogsToolName`, `TestToolsCall_LogsMultipleCalls`, `TestNewServer_NilLogger`. Update test function names if they reference old names (e.g., `TestToolsCall_HiveCells_Empty` → `TestToolsCall_OrgCells_Empty`).

### Query and stats packages (event type strings in SQL/tests)

- [x] T038 [P] [US1] In `internal/query/presets.go`: update SQL queries — `'swarm_%'` → `'forge_%'`, `'swarm_complete'` → `'forge_complete'`. Update display text "Swarm Completion Rate" → "Forge Completion Rate", "Total swarm events" → "Total forge events", "no swarm events" → "no forge events". Update constant `SwarmCompletionRate` → `ForgeCompletionRate` and its string value `"swarm_completion_rate"` → `"forge_completion_rate"`. Rename function `runSwarmCompletionRate` → `runForgeCompletionRate` and update its call site in the `Run` switch statement. Update the package-level GoDoc comment (line 5) from "swarm completion rates" → "forge completion rates".
- [x] T039 [P] [US1] In `internal/query/presets_test.go`: update all references to `SwarmCompletionRate` → `ForgeCompletionRate`, and event type strings `"swarm_init"` → `"forge_init"`, `"swarm_progress"` → `"forge_progress"`, `"swarm_complete"` → `"forge_complete"`. Update assertion strings `"Swarm Completion Rate"` → `"Forge Completion Rate"`, `"Total swarm events"` → `"Total forge events"`. Rename test functions `TestRun_SwarmCompletionRate_Empty` → `TestRun_ForgeCompletionRate_Empty` and `TestRun_SwarmCompletionRate_WithData` → `TestRun_ForgeCompletionRate_WithData`.
- [x] T040 [P] [US1] In `internal/stats/stats_test.go`: update event type strings `"swarm_init"` → `"forge_init"`, `"swarm_complete"` → `"forge_complete"`. Update assertion strings `"swarm_init"` → `"forge_init"`, `"swarm_complete"` → `"forge_complete"`.
- [x] T041 [P] [US1] In `internal/forge/insights_test.go`: update event type strings `"swarm_outcome"` → `"forge_outcome"` in all test fixture DB inserts

**Checkpoint**: `go vet ./...` passes. `go test ./... -count=1 -race` passes. `go test -tags parity ./test/parity/ -count=1` passes. `replicator docs` output shows Org/Comms/Forge/Memory categories.

---

## Phase 5: Agent Kit Content [US2, US3]

**Purpose**: Create the `internal/agentkit/` package with `embed.FS` scaffold logic and 15 template files. Independent of Phases 1–4.

### Scaffold package

- [x] T042 [US2] Create `internal/agentkit/agentkit.go`: define `embed.FS` with `//go:embed content/*` directive, `ScaffoldResult` struct (`Path string`, `Action string`), and `Scaffold(targetDir string, force bool) ([]ScaffoldResult, error)` function that walks the embedded FS and writes files to `targetDir/.opencode/`, with per-file skip logic (skip if exists and `force=false`)
- [x] T043 [US2] Create `internal/agentkit/agentkit_test.go`: test `Scaffold` in fresh `t.TempDir()` (creates 15 files), test skip-existing behavior (pre-create a file, verify "skipped" result), test `--force` overwrite (pre-create a file, verify "overwritten" result), test file count matches expected 15

### Command files (5 files)

- [x] T044 [P] [US2] Create `internal/agentkit/content/command/forge.md`: lightweight forge orchestrator command adapted from upstream `/swarm:swarm`. Must reference `forge_decompose`, `org_create_epic`, `comms_inbox`, `forge_status` tools. Provides workflow: decompose → create epic → spawn workers → monitor → review → complete.
- [x] T045 [P] [US2] Create `internal/agentkit/content/command/org.md`: org cell management command for quick cell CRUD. References `org_cells`, `org_create`, `org_update`, `org_close`, `org_start` tools.
- [x] T046 [P] [US2] Create `internal/agentkit/content/command/inbox.md`: comms inbox command for checking agent messages. References `comms_inbox`, `comms_read_message`, `comms_ack` tools.
- [x] T047 [P] [US2] Create `internal/agentkit/content/command/forge-status.md`: forge status command adapted from upstream `/swarm:status`. References `forge_status`, `org_cells`, `comms_inbox` tools.
- [x] T048 [P] [US2] Create `internal/agentkit/content/command/handoff.md`: session handoff command adapted from upstream `/swarm:handoff`. References `org_session_end`, `comms_release_all`, `org_sync` tools.

### Skill files (7 files)

- [x] T049 [P] [US2] Create `internal/agentkit/content/skills/always-on-guidance/SKILL.md`: general coding guidance skill adapted from upstream. Covers code quality, testing, error handling principles.
- [x] T050 [P] [US2] Create `internal/agentkit/content/skills/forge-coordination/SKILL.md`: forge coordination skill adapted from upstream `swarm-coordination`. Covers worker spawning, progress reporting, file reservation protocol using `comms_reserve`, `forge_progress`, `forge_complete` tools.
- [x] T051 [P] [US2] Create `internal/agentkit/content/skills/replicator-cli/SKILL.md`: replicator CLI reference skill. Documents all CLI commands (init, setup, serve, cells, doctor, stats, query, docs, version).
- [x] T052 [P] [US2] Create `internal/agentkit/content/skills/testing-patterns/SKILL.md`: Go testing patterns skill. Covers `db.OpenMemory()`, `t.TempDir()`, `httptest.NewServer`, parity tests, test naming conventions.
- [x] T053 [P] [US2] Create `internal/agentkit/content/skills/system-design/SKILL.md`: system design principles skill. Covers SOLID, DRY, dependency injection, interface abstractions.
- [x] T054 [P] [US2] Create `internal/agentkit/content/skills/learning-systems/SKILL.md`: learning from outcomes skill. Covers `forge_record_outcome`, `forge_get_strategy_insights`, `forge_get_file_insights`, `forge_get_pattern_insights` tools.
- [x] T055 [P] [US2] Create `internal/agentkit/content/skills/forge-global/SKILL.md`: global forge skill for cross-project coordination patterns.

### Agent files (3 files)

- [x] T056 [P] [US2] Create `internal/agentkit/content/agents/coordinator.md`: coordinator agent definition adapted from upstream. Defines role, available tools, workflow for orchestrating forge sessions.
- [x] T057 [P] [US2] Create `internal/agentkit/content/agents/worker.md`: worker agent definition adapted from upstream. Defines role, file reservation protocol, progress reporting, completion workflow.
- [x] T058 [P] [US2] Create `internal/agentkit/content/agents/background-worker.md`: background worker agent definition adapted from upstream. Defines role for non-interactive background tasks.

**Checkpoint**: `go test ./internal/agentkit/ -count=1` passes. `go build ./cmd/replicator` succeeds.

---

## Phase 6: Enhance `replicator init` [US2]

**Purpose**: Integrate agent kit scaffold into the `init` command. Add `--force` flag.

- [x] T059 [US2] In `cmd/replicator/init.go`: import `internal/agentkit`, remove the early-return on "already initialized" (replace with per-artifact skip logic), add `--force` bool flag, call `agentkit.Scaffold(targetDir, force)` after cells.json creation, render each `ScaffoldResult` with styled output (green for created, dim for skipped, yellow for overwritten). When overwriting with `--force`, print a warning that customizations will be lost. Update the `Short` description from "swarm operations" to "project operations". Update the `Long` description to mention agent kit scaffolding.
- [x] T060 [US2] In `cmd/replicator/init_test.go`: update `TestRunInit_FreshDirectory` to verify 16 files created (1 cells.json + 15 agent kit files). Add `TestRunInit_AgentKitSkipsExisting` — pre-create `.opencode/command/forge.md`, run init, verify file is not overwritten. Add `TestRunInit_ForceOverwrites` — pre-create a file, run init with force=true, verify file is overwritten. Update `TestRunInit_AlreadyInitialized` to verify agent kit files are still created even when `.uf/replicator/` already exists.

**Checkpoint**: `go test ./cmd/replicator/ -count=1 -run TestRunInit` passes. Manual smoke test: `go run ./cmd/replicator init --path /tmp/test-init` creates 16 files.

---

## Phase 7: Documentation Updates [US4]

**Purpose**: Update all documentation to use new terminology consistently.

- [x] T061 [P] [US4] In `AGENTS.md`: update the naming convention table (Hive→Org, Swarm Mail→Comms, Swarm→Forge), update the project structure section (hive/→org/, swarmmail/→comms/, swarm/→forge/, tools/hive/→tools/org/, tools/swarmmail/→tools/comms/, tools/swarm/→tools/forge/), add `internal/agentkit/` to the project structure, update MCP protocol section references, update any "hive" or "swarm" references in behavioral constraints and coding conventions. Update Active Technologies section to reference new package names.
- [x] T062 [P] [US4] In `README.md`: update the MCP Tools table categories (Hive→Org, Swarm Mail→Comms, Swarm→Forge), update the Architecture mermaid diagram domain label from "hive, swarm, mail" → "org, forge, comms", update the Package Layout section (hive/→org/, swarmmail/→comms/, swarm/→forge/, tools/hive/→tools/org/, tools/swarmmail/→tools/comms/, tools/swarm/→tools/forge/), add `agentkit/` to the package layout, update the Status section phase descriptions to use new names.
- [x] T063 [P] [US4] In `.specify/memory/constitution.md`: update Principle I "swarm mail messaging system" → "comms messaging system". Update Principle II "hive, swarm mail, orchestration" → "org, comms, orchestration" (or "org, comms, forge").
- [x] T064 [P] [US4] Verify `cmd/replicator/init.go` `Short` and `Long` descriptions were updated by T059 in Phase 6. If not, update `Short` from "swarm operations" → "project operations" and `Long` to mention agent kit scaffolding.

**Checkpoint**: `grep -rn 'hive_\|swarmmail_\|swarm_' AGENTS.md README.md .specify/memory/constitution.md | grep -v hivemind_ | grep -v historical` returns zero matches (excluding historical attribution and `hivemind_`).

---

## Phase 8: Verification [ALL]

**Purpose**: Full CI-equivalent verification pass. No code changes — validation only.

- [x] T065 [US1] Run `go vet ./...` — must pass with zero errors
- [x] T066 [US1] Run `go test ./... -count=1 -race` — must pass with zero failures
- [x] T067 [US1] Run `go build -o bin/replicator ./cmd/replicator` — must produce binary
- [x] T068 [US1] Run `go test -tags parity ./test/parity/ -count=1 -v` — must pass with 100% shape match
- [x] T069 [US4] Run grep verification: `grep -rn '"hive_\|"swarmmail_\|"swarm_' --include='*.go' | grep -v hivemind_` — must return zero matches (SC-005)
- [x] T070 [US2] Run init smoke test: build binary, run `./bin/replicator init --path $(mktemp -d)`, verify 16 files created (SC-002)
- [x] T071 [US1] Verify tool count: start MCP server, call `tools/list`, confirm 53 tools — 11 `org_*` + 10 `comms_*` + 24 `forge_*` + 8 `hivemind_*` (SC-001)

**Checkpoint**: All verification commands pass with zero failures. Feature is ready for review.

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (dir renames)
  └─→ Phase 2 (package decls + import paths)
        └─→ Phase 3 (tool name strings + event types)
              └─→ Phase 4 (parity fixtures + docs + query/stats)
                    └─→ Phase 8 (verification)

Phase 5 (agent kit content)  ← independent, can start anytime
  └─→ Phase 6 (init enhancement)
        └─→ Phase 8 (verification)

Phase 7 (documentation) ← depends on Phase 3 (needs final names)
  └─→ Phase 8 (verification)
```

- **Phase 1**: No dependencies — can start immediately
- **Phase 2**: Depends on Phase 1 — must complete before build works
- **Phase 3**: Depends on Phase 2 — needs compilable code
- **Phase 4**: Depends on Phase 3 — needs final tool names
- **Phase 5**: Independent — can run in parallel with Phases 1–4
- **Phase 6**: Depends on Phase 5 — needs agentkit package
- **Phase 7**: Depends on Phase 3 — needs final tool names decided
- **Phase 8**: Depends on ALL other phases — final validation

### Parallel Opportunities

- All Phase 1 tasks (T001–T006) can run in parallel
- All Phase 2 package declaration tasks (T007–T012) can run in parallel
- Phase 2 import path tasks (T013–T021) are independent per file
- All Phase 3 tasks can run in parallel (different files)
- Phase 4 fixture tasks (T031–T033) can run in parallel
- Phase 4 query/stats tasks (T038–T041) can run in parallel
- Phase 5 content files (T044–T058) can all run in parallel
- Phase 7 documentation tasks (T061–T064) can run in parallel
- **Phase 5 can run entirely in parallel with Phases 1–4**

### Commit Strategy

- **Commit 1**: Phases 1+2 together (dir renames + import fixes — build must work)
- **Commit 2**: Phase 3 (tool name strings)
- **Commit 3**: Phase 4 (parity fixtures + docs + query/stats)
- **Commit 4**: Phase 5 (agent kit content — can be committed independently)
- **Commit 5**: Phase 6 (init enhancement)
- **Commit 6**: Phase 7 (documentation)
- **Commit 7**: Phase 8 verification passes (no code changes, just confirmation)

<!-- spec-review: passed -->
