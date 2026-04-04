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

## Commands

```bash
make build    # Build
make test     # Test
make vet      # Vet
make check    # Vet + test
make serve    # Build and run MCP server
```

## Project Structure

```
cmd/replicator/       CLI entrypoint
internal/
  config/             Configuration
  db/                 SQLite + migrations
  hive/               Cell domain logic
  mcp/                MCP JSON-RPC server
  tools/
    registry/         Tool registration
    hive/             Hive tool handlers
```

## Credits

Go rewrite of [cyborg-swarm](https://github.com/unbound-force/cyborg-swarm),
originally by [Joel Hooks](https://github.com/joelhooks).
