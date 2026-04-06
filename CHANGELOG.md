# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/) and
this project adheres to [Semantic Versioning](https://semver.org/).

## [0.2.0] - 2026-04-06

### Added
- 53 MCP tools across 4 categories: Hive (11), Swarm Mail (10),
  Swarm (24), Memory (8)
- Swarm orchestration with git worktree isolation
- Agent messaging with file reservations
- Dewey memory proxy with graceful degradation
- CLI commands: init, doctor, stats, query, setup
- Parity testing engine (100% shape match vs TypeScript)
- macOS code signing and notarization
- Homebrew distribution via `brew install unbound-force/tap/replicator`
- GoReleaser v2 release pipeline (darwin-arm64, linux-amd64, linux-arm64)
- Dewey MCP tool name update (dewey#28 prefix drop)
- `replicator init` command for per-repo setup
- Constitution and expanded AGENTS.md

### Changed
- Version command now displays commit hash and build date
- Makefile: added release, install targets

## [0.1.0] - 2026-04-04

### Added
- Initial release: Phase 0 scaffold
- MCP JSON-RPC server (stdio transport)
- SQLite database via `modernc.org/sqlite` (pure Go, no CGo)
- Tool registry framework
- 4 hive tools: `hive_cells`, `hive_create`, `hive_close`, `hive_update`
- CLI: `replicator serve`, `replicator cells`, `replicator version`
- 16 tests across 3 packages
- CI workflow (go vet + go test + go build)
- MIT LICENSE with Joel Hooks attribution

[0.2.0]: https://github.com/unbound-force/replicator/releases/tag/v0.2.0
[0.1.0]: https://github.com/unbound-force/replicator/releases/tag/v0.1.0
