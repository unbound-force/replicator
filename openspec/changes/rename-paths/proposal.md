## Why

Replicator uses three different directory names inherited from the upstream fork: `.hive/` (per-repo cells), `.unbound-force/` (per-repo logs), and `~/.config/swarm-tools/` (per-machine database). These names are inconsistent, carry baggage from the TypeScript origin, and don't align with the `uf` CLI ecosystem where `uf` is the common root namespace.

Consolidating under `.uf/replicator/` (per-repo) and `~/.config/uf/replicator/` (per-machine) gives every replicator artifact a predictable, discoverable location within the UF namespace.

## What Changes

### Modified Capabilities
- Per-repo hive state moves from `[repo]/.hive/cells.json` to `[repo]/.uf/replicator/cells.json`
- Per-repo MCP log moves from `[repo]/.unbound-force/replicator.log` to `[repo]/.uf/replicator/replicator.log`
- Per-machine database moves from `~/.config/swarm-tools/swarm.db` to `~/.config/uf/replicator/replicator.db`
- Per-machine config directory moves from `~/.config/swarm-tools/` to `~/.config/uf/replicator/`

### No Migration
Clean break. Users run `replicator setup` to create the new directories. Old paths remain until manually removed. No auto-detection, no symlinks, no data copying.

## Impact

- 6 Go source files (path constants)
- 3 Go test files (assertion strings)
- 5 active documentation files (README, AGENTS.md, etc.)
- 10 completed spec artifact files (historical path references)
- Zero behavioral changes beyond the path locations

## Constitution Alignment

### I. Autonomous Collaboration
**PASS**: Path naming is a local concern. No inter-hero protocol changes.

### II. Composability First
**PASS**: Replicator remains independently installable. The new paths don't depend on other tools being present. The `~/.config/uf/` parent directory may be shared with other `uf` tools but each tool manages its own subdirectory.

### III. Observable Quality
**PASS**: The `doctor` command checks for the config directory at the new path. The path is documented in README and AGENTS.md.

### IV. Testability
**PASS**: All tests use `t.TempDir()` for filesystem operations. Path changes are string constant swaps with no logic changes.
