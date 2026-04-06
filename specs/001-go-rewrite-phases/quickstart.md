# Quickstart: Replicator

**Branch**: `001-go-rewrite-phases` | **Date**: 2026-04-04

## Prerequisites

- **Go 1.25+**: `go version`
- **Git 2.20+**: `git --version` (required for worktree operations in Phase 2+)
- **Dewey** (optional): Required for memory tools in Phase 3+. See [Dewey docs](https://github.com/unbound-force/dewey).

## Build

```bash
# Build the binary
make build

# Output: bin/replicator
```

## Test

```bash
# Run all tests
make test

# Run all tests with race detector
go test ./... -count=1 -race

# Run tests for a specific package
go test ./internal/hive/... -count=1 -race

# Run a specific test
go test ./internal/hive/... -run TestCreateCell -count=1

# Quick check (vet + test)
make check
```

## Run

### MCP Server (for AI agents)

```bash
# Start the MCP server on stdio
make serve

# Or directly:
./bin/replicator serve
```

### CLI Commands

```bash
# List work items
./bin/replicator cells

# Print version
./bin/replicator version

# Phase 4 commands (when implemented):
# ./bin/replicator doctor    # Check dependencies
# ./bin/replicator stats     # Activity metrics
# ./bin/replicator query     # SQL analytics
# ./bin/replicator setup     # Initialize environment
```

## Development Workflow

### Adding a New MCP Tool

1. **Domain logic**: Add function to `internal/{domain}/{file}.go`
2. **Tests**: Add tests to `internal/{domain}/{file}_test.go`
3. **Tool registration**: Add tool to `internal/tools/{domain}/tools.go`
4. **Wire up**: Call `{domain}.Register(reg, store)` in `cmd/replicator/serve.go`
5. **Verify**: `make check`

### Example: Adding a new hive tool

```go
// 1. internal/hive/cells.go
func StartCell(store *db.Store, id string) error {
    // ...
}

// 2. internal/hive/cells_test.go
func TestStartCell(t *testing.T) {
    store := testStore(t)
    // ...
}

// 3. internal/tools/hive/tools.go
func hiveStart(store *db.Store) *registry.Tool {
    return &registry.Tool{
        Name: "hive_start",
        // ...
    }
}

// 4. internal/tools/hive/tools.go Register()
func Register(reg *registry.Registry, store *db.Store) {
    // ... existing tools ...
    reg.Register(hiveStart(store))
}
```

### Adding a New CLI Command

```go
// 1. cmd/replicator/doctor.go
func runDoctor(cfg *config.Config, w io.Writer) error {
    // Domain logic via internal/doctor/
}

// 2. cmd/replicator/main.go
root.AddCommand(doctorCmd())
```

### Database Migrations

Add new `CREATE TABLE IF NOT EXISTS` statements to `internal/db/migrations.go`:

```go
const migrationMessages = `
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- ...
);
`

// Add to the migrations slice in db.go:
var migrations = []string{
    migrationEvents,
    migrationAgents,
    migrationCells,
    migrationCellEvents,
    migrationMessages,  // new
}
```

## Project Layout

```
cmd/replicator/       CLI entrypoint (cobra commands)
internal/
  config/             Configuration (env vars, paths)
  db/                 SQLite connection + migrations
  hive/               Cell domain logic (work items)
  mcp/                MCP JSON-RPC server (stdio)
  tools/
    registry/         Tool registration system
    hive/             Hive tool handlers
specs/                Feature specifications
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REPLICATOR_DB` | `~/.config/uf/replicator/replicator.db` | SQLite database path |
| `DEWEY_MCP_URL` | `http://localhost:3333/mcp/` | Dewey MCP endpoint |
| `ZEN_API_KEY` | (none) | OpenCode Zen API key |

## CI

The CI pipeline runs:

```bash
go vet ./...
go test ./... -count=1 -race
go build -o bin/replicator ./cmd/replicator
```

## Release

Releases are built via [GoReleaser](https://goreleaser.com/):

```bash
# Local build for all platforms
goreleaser build --snapshot --clean

# Targets: linux/darwin (amd64, arm64), windows (amd64)
```
