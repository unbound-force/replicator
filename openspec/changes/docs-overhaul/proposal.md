## Why

The README describes Phase 0 with 4 tools. The actual codebase has completed all 5 phases with 53 MCP tools, 8 CLI commands, full swarm orchestration, Dewey memory proxy, and parity testing. A developer or AI agent looking at the README gets a fundamentally wrong picture of what replicator can do.

Beyond the README, there is no MCP tool reference document -- the 53 tools and their argument schemas are only discoverable by reading source code or calling `tools/list` at runtime. There is no CONTRIBUTING.md, no CHANGELOG.md, and no architecture diagram.

GoDoc coverage is already excellent (100% of exported functions/types have doc comments), and AGENTS.md is comprehensive for AI agents. The gap is in user-facing documentation.

## What Changes

### New Capabilities
- `replicator docs` CLI command: auto-generates a markdown tool reference from the live tool registry (always in sync with code)
- `docs/tools.md`: complete MCP tool reference with schemas and hand-written usage examples for key tools
- Mermaid architecture diagram in README showing the MCP request flow

### Modified Capabilities
- `README.md`: rewritten from scratch with accurate status, full tool inventory, all CLI commands, environment variables, MCP client config examples, badges, and architecture tree
- `CONTRIBUTING.md`: new file covering dev setup, testing conventions, PR workflow, and spec-first requirements
- `CHANGELOG.md`: retroactive entries for v0.1.0 and v0.2.0

## Impact

- `cmd/replicator/docs.go` (new, ~80 lines)
- `cmd/replicator/docs_test.go` (new, ~40 lines)
- `cmd/replicator/main.go` (add `docsCmd()` registration)
- `docs/tools.md` (new, auto-generated + hand-curated examples)
- `README.md` (full rewrite)
- `CONTRIBUTING.md` (new)
- `CHANGELOG.md` (new)

## Constitution Alignment

### I. Autonomous Collaboration
**PASS**: The `replicator docs` command produces a self-describing artifact (markdown) that any consumer can read without consulting the producing tool.

### II. Composability First
**PASS**: Documentation does not introduce dependencies. The `docs` command uses only the existing tool registry -- no new imports.

### III. Observable Quality
**PASS**: The tool reference is auto-generated from the same registry that serves MCP requests, guaranteeing accuracy. The `docs` command's output can be diffed against `docs/tools.md` to detect drift.

### IV. Testability
**PASS**: The `docs` command is testable by verifying its output contains all registered tools. No external services needed.
