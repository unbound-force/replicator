# Replicator

Multi-agent coordination for AI coding agents. Single Go binary, zero runtime dependencies.

> Go rewrite of [cyborg-swarm](https://github.com/unbound-force/cyborg-swarm) (TypeScript). Same tools, same protocol, faster startup, simpler distribution.

## Status: Phase 0 (Scaffold)

Working:
- [x] SQLite database with hive schema (cells, events, agents)
- [x] MCP JSON-RPC server (stdio transport)
- [x] 4 tools: `hive_cells`, `hive_create`, `hive_close`, `hive_update`
- [x] CLI: `replicator serve`, `replicator cells`, `replicator version`

Planned:
- [ ] Phase 1: Remaining hive tools + swarm mail messaging
- [ ] Phase 2: Swarm orchestration (decompose, spawn, worktrees)
- [ ] Phase 3: Memory (Dewey proxy, Zen LLM client)
- [ ] Phase 4: Full CLI (setup, doctor, stats, query, dashboard)
- [ ] Phase 5: Parity testing against cyborg-swarm

## Install

```bash
# From source
go install github.com/unbound-force/replicator/cmd/replicator@latest

# Or download binary from releases
# https://github.com/unbound-force/replicator/releases
```

## Usage

```bash
# Initialize a project for swarm operations
replicator init

# Start MCP server (for AI agent connections)
replicator serve

# List hive cells
replicator cells

# Version
replicator version
```

## Development

```bash
make build    # Build binary to bin/replicator
make test     # Run all tests
make vet      # Go vet
make check    # vet + test
make serve    # Build and run MCP server
```

## Architecture

```
cmd/replicator/     CLI entrypoint (cobra)
internal/
  config/           Configuration (env vars, defaults)
  db/               SQLite connection + migrations
  hive/             Cell (work item) domain logic
  mcp/              MCP JSON-RPC server
  tools/
    registry/       Tool registration framework
    hive/           Hive MCP tool handlers
```

## Credits

Go rewrite of [cyborg-swarm](https://github.com/unbound-force/cyborg-swarm), originally forked from [swarm-tools](https://github.com/joelhooks/swarm-tools) by [Joel Hooks](https://github.com/joelhooks). See [LICENSE](LICENSE).
