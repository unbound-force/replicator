# Replicator Agent Guide

## Overview

Replicator is the Go rewrite of cyborg-swarm. It provides multi-agent
coordination tools via the MCP protocol and a CLI for observability.

## Language & Toolchain

- Go 1.25+
- SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- CLI via `cobra`
- Tests via `go test` (stdlib)

## Critical Rules

### TDD Everything

All code changes follow Red-Green-Refactor:
1. Write a failing test
2. Write minimal code to pass
3. Refactor while tests stay green

### Database

Single global database at `~/.config/swarm-tools/swarm.db`.
Schema is compatible with cyborg-swarm's libSQL database.
Use in-memory databases for tests (`db.OpenMemory()`).

### MCP Protocol

Tools are registered via the `registry` package and served
over stdio JSON-RPC. Each tool has:
- A name (e.g., `hive_cells`)
- A description
- A JSON schema for arguments
- An execute function

### Naming Convention: The Hive Metaphor

| Concept | Name |
|---------|------|
| Work items | **Hive** |
| Individual item | **Cell** |
| Agent coordination | **Swarm** |
| Messaging | **Swarm Mail** |
| Parallel workers | **Workers** |
| Task orchestrator | **Coordinator** |
| File locks | **Reservations** |

## Constitution (Highest Authority)

The project constitution at `.specify/memory/constitution.md`
extends the Unbound Force org constitution (v1.1.0) with four
core principles:

1. **I. Autonomous Collaboration**: Tools are callable
   independently via MCP. Outputs are self-describing JSON.
   Inter-agent communication uses swarm mail.
2. **II. Composability First**: The binary works standalone.
   Dewey integration degrades gracefully. Database schema is
   compatible with cyborg-swarm.
3. **III. Observable Quality**: All tool responses are JSON.
   Response shapes match the TypeScript version (parity tests).
   Doctor reports are machine-readable.
4. **IV. Testability**: Database tests use in-memory SQLite.
   Git tests use `t.TempDir()`. HTTP tests use `httptest`.
   No external services required.

Constitution violations are CRITICAL severity and
non-negotiable.

## Behavioral Constraints

- **Zero-Waste Mandate**: No orphaned code, unused
  dependencies, or dead functions. Every file must serve a
  purpose traceable to a spec or tool.
- **CI Parity Gate**: Before marking any task complete, run
  the CI-equivalent checks locally. Read
  `.github/workflows/` for the exact commands -- do not
  rely on memory. Any failure blocks the task.
- **Intent Drift Detection**: Implementation must faithfully
  capture the spec's intent. The parity test suite verifies
  response shapes match the TypeScript version.
- **Automated Governance**: Constitution alignment is
  verified via the Constitution Check gate at planning
  time, not ad-hoc review.

## Coding Conventions

- **Formatting**: `gofmt` and `goimports` (enforced by
  golangci-lint).
- **Naming**: Standard Go conventions. PascalCase exported,
  camelCase unexported.
- **Comments**: GoDoc-style on all exported functions and
  types. Package-level doc comments on every package.
- **Error handling**: Return `error`. Wrap with
  `fmt.Errorf("context: %w", err)`. Use `errors.Is` for
  sentinel errors (not string comparison).
- **Import grouping**: Standard library, then third-party,
  then internal packages (separated by blank lines).
- **No global state**: Prefer dependency injection and
  functional style.
- **JSON tags**: Required on all struct fields intended for
  serialization.
- **Constants**: Use string-typed constants for enumerations.

## Testing Conventions

- **Framework**: Standard library `testing` package only.
  No testify, gomega, or external assertion libraries.
- **Assertions**: Use `t.Errorf` / `t.Fatalf` directly.
- **Test naming**: `TestXxx_Description` (e.g.,
  `TestCreateCell_Defaults`, `TestReadyCell_PriorityOrder`).
- **Test isolation**: `db.OpenMemory()` for database tests.
  `t.TempDir()` for filesystem/git tests.
  `httptest.NewServer` for HTTP tests.
- **No shared state**: Each test creates its own store and
  fixtures. No test depends on another test's side effects.
- **Parity tests**: Build tag `//go:build parity`. Compare
  Go response shapes against TypeScript fixtures.
- **Git tests**: Guard with `if testing.Short() { t.Skip() }`
  for tests that shell out to git.

## Knowledge Retrieval

Agents SHOULD prefer Dewey MCP tools over grep/glob/read
for cross-repo context, design decisions, and architectural
patterns.

### Tool Selection Matrix

| Query Intent | Dewey Tool | When to Use |
|-------------|-----------|-------------|
| Conceptual understanding | `dewey_semantic_search` | "How does X work?" |
| Keyword lookup | `dewey_search` | Known terms, FR numbers |
| Read specific page | `dewey_get_page` | Known document path |
| Relationship discovery | `dewey_find_connections` | "How are X and Y related?" |
| Similar documents | `dewey_similar` | "Find specs like this one" |
| Filtered semantic | `dewey_semantic_search_filtered` | Search within source type |
| Graph navigation | `dewey_traverse` | Dependency chain walking |

### Graceful Degradation (3-Tier)

**Tier 3 (Full Dewey)**: `dewey_semantic_search`,
`dewey_search`, `dewey_traverse`, and
`dewey_semantic_search_filtered` for comprehensive
cross-repo context.

**Tier 2 (Graph-only, no embedding model)**: `dewey_search`
and `dewey_traverse` for keyword and structural queries.

**Tier 1 (No Dewey)**: Direct file operations -- Read tool,
Grep tool, convention packs at
`.opencode/unbound/packs/`.

## Specification Framework

This project uses a two-tier specification framework:

| Tier | Tool | When to Use | Location |
|------|------|-------------|----------|
| Strategic | Speckit | 3+ stories, architecture | `specs/NNN-*/` |
| Tactical | OpenSpec | <3 stories, bug fix | `openspec/changes/` |

### Speckit Pipeline

```text
constitution -> specify -> clarify -> plan -> tasks
  -> analyze -> checklist -> implement
```

### OpenSpec Workflow

```text
propose -> design -> specs -> tasks -> apply -> archive
```

### Ordering Constraints

1. Constitution must exist before specs.
2. Spec before plan. Plan before tasks.
3. Tasks before implementation.
4. All checklists must pass before implementation.

### Task Completion Bookkeeping

Mark `- [ ]` to `- [x]` immediately after each task. Do
not batch completions.

### Documentation Validation Gate

Before marking any task complete, check whether changes
require updates to:

- `README.md` -- commands, flags, architecture
- `AGENTS.md` -- conventions, packages, patterns
- GoDoc comments -- exported functions and types
- Spec artifacts under `specs/`

### Spec Commit Gate

All spec artifacts MUST be committed and pushed before
implementation begins.

## Git & Workflow

- **Commit format**: Conventional Commits --
  `type: description` (feat, fix, docs, chore, refactor).
- **Branching**: Feature branches required. Speckit:
  `NNN-<name>`. OpenSpec: `opsx/<name>`.
- **Code review**: Required before merge.
- **Semantic versioning**: For releases.

## Commands

```bash
make build    # Build binary
make test     # Run all tests
make vet      # Go vet
make check    # Vet + test
make serve    # Build and run MCP server
make release  # GoReleaser dry-run
make install  # Install to GOPATH/bin
```

## Project Structure

```
cmd/replicator/       CLI entrypoint (cobra)
internal/
  config/             Configuration
  db/                 SQLite + migrations (7 tables)
  hive/               Cell domain logic (CRUD, epics, sessions, sync)
  swarmmail/          Agent messaging + file reservations
  swarm/              Orchestration (decompose, spawn, worktree, review, insights)
  memory/             Dewey proxy + deprecated tool stubs
  gitutil/            Git worktree operations (os/exec)
  doctor/             Health check engine
  stats/              Database statistics
  query/              Preset SQL queries
  mcp/                MCP JSON-RPC server
  tools/
    registry/         Tool registration framework
    hive/             Hive tool handlers (11 tools)
    swarmmail/        Swarm mail tool handlers (10 tools)
    swarm/            Swarm tool handlers (24 tools)
    memory/           Memory tool handlers (8 tools)
test/parity/          Shape comparison engine + fixtures
```

## Credits

Go rewrite of [cyborg-swarm](https://github.com/unbound-force/cyborg-swarm),
originally by [Joel Hooks](https://github.com/joelhooks).

## Active Technologies
- Go 1.25+ + `cobra` (CLI), `modernc.org/sqlite` (pure Go SQLite), stdlib `encoding/json` (MCP JSON-RPC), stdlib `os/exec` (git operations) (001-go-rewrite-phases)
- SQLite at `~/.config/swarm-tools/swarm.db` (WAL mode, compatible with cyborg-swarm) (001-go-rewrite-phases)

## Recent Changes
- 001-go-rewrite-phases: Added Go 1.25+ + `cobra` (CLI), `modernc.org/sqlite` (pure Go SQLite), stdlib `encoding/json` (MCP JSON-RPC), stdlib `os/exec` (git operations)
