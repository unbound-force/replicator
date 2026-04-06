## Context

Replicator has `setup` (per-machine: `~/.config/swarm-tools/` + SQLite) but no `init` (per-repo). The `uf` CLI needs `replicator init` to exist so it can delegate to it during `uf init`, following the same pattern as `dewey init`.

The existing `setup.go` pattern: cobra command → `runSetup()` function → filesystem ops → print status. Follow this exactly.

## Goals / Non-Goals

### Goals
- Add `replicator init` command that creates `.hive/cells.json` in the target directory
- Support `--path` flag for specifying a non-CWD target
- Idempotent (exit 0 whether creating or already exists)
- No database, no git, no network

### Non-Goals
- Creating `.hive/memories.jsonl` or other hive files (those are created by specific tools when needed)
- Running `git add` or `git commit` (that's `hive_sync`)
- Initializing the SQLite database (that's `replicator setup`)
- Adding `.hive/` to `.gitignore` (user decides tracking policy)

## Decisions

**D1: Follow the `setup.go` pattern.** New file `cmd/replicator/init.go` with `initCmd() *cobra.Command` and `runInit(targetDir string) error`. Register in `main.go` via `root.AddCommand(initCmd())`.

**D2: Seed file is `cells.json` with `[]`.** This matches the format `hive_sync` writes. An empty JSON array signals "initialized but no cells yet" rather than the ambiguity of a missing file.

**D3: `--path` defaults to `.` (current directory).** Same UX as `dewey init` which operates on CWD by default.

**D4: Testable via `t.TempDir()`.** The `runInit` function takes a directory path, making it directly testable without mocking.
