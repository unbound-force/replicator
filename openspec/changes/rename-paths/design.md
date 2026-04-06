## Context

Three different directory names from the upstream fork need consolidation:

| Current | New | Purpose |
|---------|-----|---------|
| `[repo]/.hive/cells.json` | `[repo]/.uf/replicator/cells.json` | Per-repo cell state |
| `[repo]/.unbound-force/replicator.log` | `[repo]/.uf/replicator/replicator.log` | Per-repo MCP session log |
| `~/.config/swarm-tools/swarm.db` | `~/.config/uf/replicator/replicator.db` | Per-machine SQLite database |

## Goals / Non-Goals

### Goals
- Consolidate all replicator paths under `.uf/replicator/` (per-repo) and `~/.config/uf/replicator/` (per-machine)
- Rename database file from `swarm.db` to `replicator.db`
- Update all Go source, tests, documentation, and spec artifacts

### Non-Goals
- Auto-migration from old paths (clean break)
- Changing the `.unbound-force/config.yaml` ecosystem config (that's UF-level, not replicator-specific)
- Changing `.opencode/` or `.specify/` directory structures
- Modifying the LICENSE attribution text

## Decisions

**D1: Pure find-and-replace.** Every change is a string constant swap. No logic modifications, no new functions, no behavioral changes.

**D2: No migration path.** The old `~/.config/swarm-tools/swarm.db` is abandoned. Users run `replicator setup` to create the new directory. This is a clean break -- the cyborg-swarm TypeScript version can no longer share the same database.

**D3: Update spec artifacts.** Historical spec files will be updated to reflect the new paths for consistency, since the old paths will cause confusion if someone reads them.

**D4: Database renamed to `replicator.db`.** Complete naming break from `swarm.db`. Matches the tool name.

## Risks / Trade-offs

**Risk: Existing users lose data.** Anyone with cells/events in `~/.config/swarm-tools/swarm.db` will start fresh. Mitigation: this is a pre-release tool with no external users yet. The old database can be manually copied if needed.

**Trade-off: Breaks cyborg-swarm compatibility.** The Go and TypeScript versions can no longer share a database. This is intentional -- replicator is replacing cyborg-swarm, not coexisting with it long-term.
