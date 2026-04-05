# Research: Go Rewrite — Library Decisions

**Branch**: `001-go-rewrite-phases` | **Date**: 2026-04-04
**Purpose**: Document technical decisions for Go library selection and architectural approach.

## Decision 1: MCP Server — Hand-Rolled vs mcp-go

### Context

The MCP (Model Context Protocol) server handles JSON-RPC 2.0 over stdio. Phase 0 implemented a hand-rolled server in `internal/mcp/server.go` (212 lines). The question is whether to adopt an MCP library for the remaining phases.

### Options Evaluated

| Option | Pros | Cons |
|--------|------|------|
| **Keep hand-rolled** (chosen) | Already working, 212 lines, zero dependencies, full control over response shapes, tested | Must handle any new MCP methods manually |
| `github.com/mark3labs/mcp-go` | Community standard, handles protocol details | Adds dependency, may constrain response shapes, version churn risk |
| `github.com/metoro-io/mcp-golang` | Alternative implementation | Less mature, similar trade-offs |

### Decision

**Keep the hand-rolled MCP server.** The existing implementation handles `initialize`, `tools/list`, and `tools/call` — the only three methods needed. It's 212 lines with full test coverage. Adding a library would increase binary size and introduce a dependency that could constrain response shape compatibility with cyborg-swarm.

### Validation

The hand-rolled server already passes integration tests with real MCP clients (OpenCode, Claude). No protocol compliance issues have been reported.

---

## Decision 2: SQLite Driver — modernc.org/sqlite

### Context

Phase 0 chose `modernc.org/sqlite` (pure Go, no CGo). This decision is confirmed and not revisited.

### Rationale

- **Zero CGo**: Single binary distribution without C compiler requirements
- **Schema compatibility**: Works with the same SQLite database as cyborg-swarm's `better-sqlite3`
- **WAL mode**: Supports WAL journal mode for concurrent reads
- **Performance**: Adequate for the workload (single-digit millisecond queries on small datasets)

### Confirmed Settings

```go
dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", path)
sqlDB.SetMaxOpenConns(1) // Single writer for SQLite
```

---

## Decision 3: Git Operations — os/exec vs go-git

### Context

Swarm orchestration requires git worktree management: create, list, merge (cherry-pick), and cleanup. Two approaches are available.

### Options Evaluated

| Option | Pros | Cons |
|--------|------|------|
| **os/exec to git binary** (chosen) | Full git feature support, worktrees work correctly, small binary impact, matches spec assumption | Requires git installed, subprocess overhead |
| `github.com/go-git/go-git` | Pure Go, no external dependency | No worktree support, adds ~15MB to binary, incomplete git feature coverage |

### Decision

**Shell out to the `git` binary via `os/exec`.** This is explicitly stated in the spec's Assumptions section: "Git operations (worktree create/merge/cleanup) shell out to the `git` binary rather than using a pure-Go git library." The `doctor` command will verify git is installed.

### Implementation Pattern

```go
// internal/gitutil/git.go
func Run(dir string, args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("git %s: %w\nstderr: %s", args[0], err, stderr.String())
    }
    return strings.TrimSpace(stdout.String()), nil
}
```

---

## Decision 4: HTTP Client for Dewey Proxy — stdlib vs third-party

### Context

Memory tools proxy requests to Dewey's MCP endpoint over HTTP. Need an HTTP client.

### Decision

**Use stdlib `net/http`.** The proxy makes simple JSON-RPC POST requests. No need for retry logic, connection pooling, or advanced features that would justify a third-party client. The Dewey endpoint is localhost.

### Implementation Pattern

```go
// internal/memory/proxy.go
type Client struct {
    url    string
    http   *http.Client
}

func NewClient(deweyURL string) *Client {
    return &Client{
        url:  deweyURL,
        http: &http.Client{Timeout: 10 * time.Second},
    }
}

func (c *Client) Call(method string, params any) (json.RawMessage, error) {
    // JSON-RPC 2.0 POST to c.url
}
```

---

## Decision 5: CLI Framework — cobra (confirmed)

### Context

Phase 0 chose `github.com/spf13/cobra` for CLI routing. This is confirmed per CS-009 (`[MUST]`).

### Phase 4 Commands

New commands follow the existing pattern from `cmd/replicator/main.go`:

```go
root.AddCommand(doctorCmd())
root.AddCommand(statsCmd())
root.AddCommand(queryCmd())
root.AddCommand(setupCmd())
```

Each command delegates to a `Run(opts)` function in the corresponding `internal/` package per AP-002.

---

## Decision 6: Logging — charmbracelet/log

### Context

CS-008 requires `github.com/charmbracelet/log` for all application logging. Phase 0 has no logging — it's a stdio server that communicates via JSON-RPC.

### Decision

**Add `charmbracelet/log` as a dependency when logging is first needed** (likely Phase 2 for worktree operations or Phase 4 for CLI output). Log to stderr only — stdout is reserved for MCP JSON-RPC communication.

### Deferred Until

Phase 2 or Phase 4, whichever introduces the first non-MCP output path. Note: The MCP server communicates via stdio JSON-RPC and currently has no application logging. If logging is added, it MUST go to stderr only (stdout is reserved for MCP). The dependency should be added to `go.mod` in the same task that introduces the first `log.Info()` call.

---

## Decision 7: Parity Testing Approach

### Context

Phase 5 requires comparing MCP tool response shapes between Go and TypeScript implementations.

### Options Evaluated

| Option | Pros | Cons |
|--------|------|------|
| **Fixture-based comparison** (chosen) | Deterministic, no TypeScript server needed at test time, fast | Fixtures can go stale |
| Live dual-server comparison | Always current | Requires TypeScript server running, flaky, slow |
| Schema-based validation | Formal, precise | Requires maintaining JSON schemas for all 70 tools |

### Decision

**Fixture-based comparison.** Capture TypeScript responses once, store as JSON fixtures, compare Go responses against the fixture shapes at test time. Fixtures are refreshed manually when the TypeScript version changes.

### Shape Comparison Algorithm

```
ShapeMatch(expected, actual):
  if types differ → FAIL (report path, expected type, actual type)
  if both objects → recurse on each key
  if both arrays → recurse on first element (shape of array items)
  if both primitives → PASS (values may differ, shapes match)
```

---

## Decision 8: Test Isolation for Database Tests

### Context

All domain packages need database access for testing. The established pattern uses `db.OpenMemory()`.

### Decision

**Continue using in-memory SQLite for all unit tests.** Each test function gets its own `*db.Store` via the `testStore(t)` helper. No shared state between tests (TC-010). The `t.Cleanup()` function closes the store after each test.

For integration tests that need git (Phase 2), use `t.TempDir()` to create isolated git repositories.

---

## Dependency Summary

### Current Dependencies (Phase 0)

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | v1.10.2 |
| `modernc.org/sqlite` | Pure Go SQLite | v1.48.1 |

### Planned Additions

| Package | Phase | Purpose |
|---------|-------|---------|
| `github.com/charmbracelet/log` | 2 or 4 | Structured logging to stderr |

### Explicitly NOT Adding

| Package | Reason |
|---------|--------|
| `mcp-go` | Hand-rolled server is simpler and sufficient |
| `go-git` | No worktree support, large binary impact |
| `testify` | Convention pack TC-001 prohibits external test libraries |
| `gorm` / `sqlx` | `database/sql` is sufficient for the query patterns |
| `gin` / `echo` | No HTTP server needed (MCP is stdio, Dewey proxy is client-only) |
