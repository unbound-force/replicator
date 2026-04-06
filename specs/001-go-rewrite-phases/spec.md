# Feature Specification: Go Rewrite -- Remaining Phases

**Feature Branch**: `001-go-rewrite-phases`  
**Created**: 2026-04-04  
**Status**: Ready  
**Input**: User description: "all the other phases to create a GoLang replacement for swarm-tools"  
**References**: [Replicator Phase 0](https://github.com/unbound-force/replicator), [cyborg-swarm](https://github.com/unbound-force/cyborg-swarm)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Complete Hive Tool Suite (Priority: P1)

An AI coding agent connects to the replicator MCP server and uses the full hive tool suite to manage work items: create cells, create epics with subtasks, query by status/type, update status, start work, find the next ready cell, and sync to git. The agent's workflow is identical to using the TypeScript version -- same tool names, same argument schemas, same response shapes.

**Why this priority**: Hive is the foundation for all work tracking. Phase 0 delivered 4 of 11 hive tools. The remaining 7 tools (`hive_create_epic`, `hive_query`, `hive_start`, `hive_ready`, `hive_sync`, `hive_session_start`, `hive_session_end`) complete the work management layer that every other feature depends on.

**Independent Test**: Create an epic with 3 subtasks via MCP, query for ready cells, start one, complete it, then sync. All operations succeed with correct response shapes matching the TypeScript version.

**Acceptance Scenarios**:

1. **Given** the MCP server is running, **When** an agent calls `hive_create_epic` with a title and subtask array, **Then** the epic and all subtasks are created with correct parent-child relationships.
2. **Given** cells exist with dependencies, **When** an agent calls `hive_ready`, **Then** only unblocked cells are returned, sorted by priority.
3. **Given** a cell is open, **When** an agent calls `hive_start`, **Then** the cell status changes to `in_progress` and the timestamp is recorded.
4. **Given** cells have been modified, **When** `hive_sync` is called, **Then** the hive state is serialized and committed to the project's git repository.
5. **Given** an agent begins a session, **When** `hive_session_start` is called, **Then** previous session handoff notes are returned if available.

---

### User Story 2 - Agent Messaging (Priority: P1)

An AI coding agent initializes a swarm mail session, sends messages to other agents, checks its inbox, reserves files for exclusive editing, and releases reservations. Multiple agents working on the same project can coordinate without conflicts.

**Why this priority**: Swarm mail is the communication backbone for multi-agent coordination. Without messaging and file reservations, parallel workers cannot coordinate safely.

**Independent Test**: Initialize two agent sessions, send a message from agent A to agent B, verify it appears in B's inbox, reserve a file path, verify the reservation is enforced, then release it.

**Acceptance Scenarios**:

1. **Given** a project path, **When** `swarmmail_init` is called with an agent name, **Then** the agent is registered and a session is created.
2. **Given** two registered agents, **When** agent A calls `swarmmail_send` with a message to agent B, **Then** the message is persisted and appears in B's `swarmmail_inbox`.
3. **Given** agent A has sent a message, **When** agent B calls `swarmmail_read_message` with the message ID, **Then** the full message body is returned.
4. **Given** a file path, **When** agent A calls `swarmmail_reserve` with exclusive mode, **Then** the file is locked and agent B's attempt to reserve the same file fails with a conflict error.
5. **Given** agent A holds a reservation, **When** agent A calls `swarmmail_release`, **Then** the file is unlocked and available for other agents.
6. **Given** a message requires acknowledgment, **When** agent B calls `swarmmail_ack`, **Then** the message is marked as acknowledged.

---

### User Story 3 - Swarm Orchestration (Priority: P1)

An AI coding agent decomposes a task into subtasks, spawns parallel workers in isolated git worktrees, tracks progress, merges results, and reviews completed work. The full swarm coordination lifecycle works end-to-end.

**Why this priority**: Orchestration is the core differentiator of the tool. Task decomposition, parallel execution, and worktree isolation are the primary reasons agents use this system.

**Independent Test**: Decompose a task description into subtasks, create worktrees for 2 parallel workers, report progress, complete the subtasks, merge worktrees back, and verify the merged result is clean.

**Acceptance Scenarios**:

1. **Given** a task description, **When** `swarm_decompose` is called, **Then** a structured decomposition prompt is returned with file assignments.
2. **Given** a project path and commit hash, **When** `swarm_worktree_create` is called, **Then** an isolated git worktree is created at the expected path.
3. **Given** an active worktree, **When** `swarm_spawn_subtask` is called, **Then** the subtask prompt includes the correct context, files, and worktree path.
4. **Given** a worker is executing, **When** `swarm_progress` is called with a percentage, **Then** the progress is recorded and visible via `swarm_status`.
5. **Given** a completed worktree, **When** `swarm_worktree_merge` is called, **Then** the worker's commits are cherry-picked onto the main branch.
6. **Given** all subtasks are complete, **When** `swarm_status` is called, **Then** it shows all subtasks as completed with their outcomes.

---

### User Story 4 - Memory and Context (Priority: P2)

An AI coding agent stores learnings, searches for relevant context, and retrieves past decisions through the memory tools. Memory operations are proxied through Dewey's semantic search service for cross-repo context.

**Why this priority**: Memory enables agents to learn from past sessions and avoid repeating mistakes. It's important but not blocking for the core coordination workflow.

**Independent Test**: Store a learning with tags, search for it by semantic query, retrieve it by ID, and validate the response matches the stored content.

**Acceptance Scenarios**:

1. **Given** Dewey is available, **When** `hivemind_store` is called with information and tags, **Then** the memory is persisted through Dewey and a deprecation warning is included in the response.
2. **Given** memories exist, **When** `hivemind_find` is called with a query, **Then** semantically relevant results are returned in the expected response format.
3. **Given** Dewey is unavailable, **When** any memory tool is called, **Then** a clear error message is returned explaining the backend is unreachable.
4. **Given** a secondary memory tool is called (`hivemind_get`, `hivemind_stats`, etc.), **Then** a deprecation message is returned with the Dewey replacement tool name.

---

### User Story 5 - CLI Operations (Priority: P2)

A developer uses the `replicator` CLI to set up their environment, check health, query swarm activity, and inspect work items. The CLI provides the same observability as the TypeScript `swarm` CLI but starts instantly.

**Why this priority**: CLI commands are the developer-facing interface. They're important for observability but don't block agent operations (which use MCP tools).

**Independent Test**: Run `replicator setup`, `replicator doctor`, `replicator cells`, `replicator stats`, and `replicator query` and verify correct output for each.

**Acceptance Scenarios**:

1. **Given** a fresh environment, **When** `replicator doctor` is run, **Then** it checks for required dependencies and reports their status.
2. **Given** an active swarm with completed work, **When** `replicator stats` is run, **Then** it displays health metrics and activity summary.
3. **Given** a database with events, **When** `replicator query` is run with a preset name, **Then** it executes the SQL analytics query and displays results.
4. **Given** a project with cells, **When** `replicator cells` is run with a status filter, **Then** matching cells are displayed as formatted output.

---

### User Story 6 - Parity Verification (Priority: P3)

A maintainer runs a parity test suite that compares the replicator's MCP tool responses against the TypeScript cyborg-swarm's responses for the same inputs. All tools produce identical response shapes and semantically equivalent results.

**Why this priority**: Parity verification is the final quality gate before the replicator can replace cyborg-swarm. It depends on all other stories being complete.

**Independent Test**: For each of the 70 unique MCP tools, send identical arguments to both the Go and TypeScript servers and compare response shapes. Zero shape mismatches.

**Acceptance Scenarios**:

1. **Given** both servers are running, **When** identical `hive_create` arguments are sent to each, **Then** both return the same JSON shape with matching field names and types.
2. **Given** the parity suite completes, **When** results are inspected, **Then** all 70 unique tools show "shape match" status.
3. **Given** any response shape mismatch, **When** the parity report is generated, **Then** it identifies the specific field, expected type, and actual type.

---

### Edge Cases

- What happens when the SQLite database file is locked by another process? The system should retry with exponential backoff and report a clear error if the lock cannot be acquired within a reasonable timeout.
- What happens when `swarm_worktree_create` is called but the git working directory has uncommitted changes? The worktree is created from the specified commit hash, not the working directory state, so uncommitted changes do not affect worktree creation.
- What happens when `swarm_worktree_merge` encounters a conflict? The merge fails with a clear error listing the conflicting files, and the worktree is NOT cleaned up (allowing manual resolution).
- What happens when Dewey is down during `hivemind_store`? The tool returns a structured error with code `DEWEY_UNAVAILABLE` and a hint to check the Dewey service.
- What happens when a CLI command is run but no database exists? The database is auto-created with the full schema on first access.
- How does the system handle the TypeScript cyborg-swarm's database during migration? The schema is wire-compatible -- the Go binary reads and writes the same SQLite database file at `~/.config/uf/replicator/replicator.db`.

## Requirements *(mandatory)*

### Functional Requirements

#### Hive (Phase 1)
- **FR-001**: All 11 hive MCP tools MUST be implemented with argument schemas and response shapes matching the TypeScript version.
- **FR-002**: The `hive_create_epic` tool MUST support atomic creation of an epic with N subtasks in a single call.
- **FR-003**: The `hive_ready` tool MUST return only unblocked cells (no open dependencies), sorted by priority.
- **FR-004**: The `hive_sync` tool MUST serialize cell state and commit to the project's `.uf/replicator/` directory.

#### Swarm Mail (Phase 1)
- **FR-005**: All 9 swarm mail MCP tools MUST be implemented with matching schemas and response shapes.
- **FR-006**: File reservations MUST enforce exclusive locking -- concurrent reserve attempts on the same path MUST fail with a conflict error.
- **FR-007**: Messages MUST persist across process restarts (stored in SQLite).

#### Orchestration (Phase 2)
- **FR-008**: All 16 swarm orchestration MCP tools MUST be implemented.
- **FR-009**: The `swarm_worktree_create` tool MUST create isolated git worktrees using `git worktree add`.
- **FR-010**: The `swarm_worktree_merge` tool MUST cherry-pick worker commits back to the main branch and detect conflicts.
- **FR-011**: The `swarm_decompose` tool MUST generate decomposition prompts that match the TypeScript version's structure.

#### Memory (Phase 3)
- **FR-012**: The `hivemind_store` and `hivemind_find` tools MUST proxy through Dewey's MCP endpoint.
- **FR-013**: The 6 secondary hivemind tools MUST return structured deprecation messages with replacement tool names.
- **FR-014**: When Dewey is unavailable, proxied tools MUST return a clear error rather than silently failing.

#### CLI (Phase 4)
- **FR-015**: The binary MUST start in under 50 milliseconds (cold start, no warm cache).
- **FR-016**: CLI commands MUST include at minimum: `serve`, `cells`, `doctor`, `stats`, `query`, `version`.
- **FR-017**: The `doctor` command MUST check for Dewey availability, database accessibility, and git configuration.

#### Parity (Phase 5)
- **FR-018**: A parity test suite MUST compare response shapes for all implemented tools against the TypeScript version.
- **FR-019**: The parity report MUST list every tool with pass/fail status and any shape differences.

#### Cross-Cutting
- **FR-020**: The database schema MUST be compatible with cyborg-swarm's SQLite database (same table names, column names, and types).
- **FR-021**: The MCP server MUST handle the `initialize`, `tools/list`, and `tools/call` methods per the MCP specification.
- **FR-022**: All database operations MUST use WAL mode for concurrent read access.
- **FR-023**: The binary MUST be distributable as a single file with zero runtime dependencies.

### Key Entities

- **Cell**: A work item with ID, title, description, type (task/bug/feature/epic/chore), status (open/in_progress/blocked/closed), priority, and optional parent-child relationships.
- **Agent**: A registered AI agent with name, project path, role, status, and metadata. Agents interact via swarm mail.
- **Message**: An inter-agent communication with sender, recipients, subject, body, importance level, and thread ID.
- **Reservation**: A file lock held by an agent for exclusive editing, with path, TTL, and exclusivity flag.
- **Worktree**: An isolated git working directory created for a specific subtask, with task ID, branch name, and status.
- **Memory/Learning**: A piece of knowledge stored via Dewey, with information text, tags, and semantic embedding.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 70 unique MCP tools are implemented and pass their individual test suites with zero failures.
- **SC-002**: The binary starts and responds to the first MCP request in under 50 milliseconds (measured from process launch to first JSON-RPC response).
- **SC-003**: The parity test suite shows 100% response shape match for all tools compared to the TypeScript version.
- **SC-004**: The binary size is under 20MB for all target platforms (macOS/Linux arm64 and amd64).
- **SC-005**: The test suite achieves at least 80% line coverage across all internal packages.
- **SC-006**: Two parallel agents can coordinate via swarm mail (send messages, reserve files, report progress) without data corruption or deadlocks.
- **SC-007**: The full suite of 70 tools completes implementation within 8 months of Phase 0 completion.
- **SC-008**: The CLI `doctor` command detects and reports all required dependencies within 2 seconds.

## Assumptions

- Phase 0 (scaffold) is complete: MCP server, SQLite, tool registry, and 4 hive tools are working.
- The database schema at `~/.config/uf/replicator/replicator.db` is the shared state between the Go and TypeScript versions during migration.
- Dewey is the canonical memory backend; Ollama embedding operations are handled by Dewey, not by replicator.
- The eval system (`swarm-evals`) remains in TypeScript and is NOT part of this rewrite.
- The dashboard web UI (`swarm-dashboard`) remains in TypeScript and is NOT part of this rewrite.
- Git operations (worktree create/merge/cleanup) shell out to the `git` binary rather than using a pure-Go git library.
- The `swarm-queue` package (BullMQ/Redis) is deferred -- it can be added later if Redis-based queuing is needed.

## Dependencies

- **Phase 0 (complete)**: MCP server, SQLite, tool registry, 4 hive tools.
- **Dewey**: Must be running for memory proxy tools to function. Graceful degradation when unavailable.
- **Git**: Must be installed for worktree operations. The `doctor` command checks for this.
- **cyborg-swarm**: The TypeScript version serves as the behavioral reference for parity testing.
