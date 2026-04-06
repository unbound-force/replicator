## Why

The `uf init` command delegates per-repo initialization to each sub-tool (`dewey init`, `gaze init`). Replicator lacks an `init` command, so `uf init` cannot prepare a repository for swarm operations. Without it, the `.hive/` directory only appears organically on first `hive_sync` call, leaving new projects in an ambiguous state where the directory structure doesn't signal that swarm coordination is available.

This blocks `unbound-force/unbound-force#82` (replace Swarm plugin with Replicator in the `uf` CLI).

## What Changes

### New Capabilities
- `replicator init` creates a `.hive/` directory with an empty `cells.json` in the current working directory
- `replicator init --path /some/dir` creates `.hive/` at a specified location
- Idempotent: safe to run multiple times (prints "already initialized" if `.hive/` exists)

## Impact

- `cmd/replicator/init.go` (new file, ~40 lines)
- `cmd/replicator/main.go` (add `initCmd()` registration)
- `README.md` (add `init` to usage section)
- `AGENTS.md` (add `init` to commands table)

## Constitution Alignment

### I. Autonomous Collaboration
**PASS**: `init` creates a well-known directory (`.hive/`) that other tools can discover without coordination.

### II. Composability First
**PASS**: No database or external service required. Works standalone before `replicator setup` has run.

### III. Observable Quality
**PASS**: Prints machine-parseable status (`initialized .hive/` or `already initialized`). Exit 0 on success, exit 1 on error.

### IV. Testability
**PASS**: Uses `t.TempDir()` for filesystem tests. No external dependencies.
