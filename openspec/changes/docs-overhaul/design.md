## Context

The README is stuck at Phase 0. The codebase has 53 tools, 8 CLI commands, 190+ tests, and a working Homebrew distribution. GoDoc coverage is 100% and AGENTS.md is comprehensive, but user-facing docs are stale or missing.

## Goals / Non-Goals

### Goals
- Rewrite README to accurately reflect the current codebase
- Create a `replicator docs` command that auto-generates tool reference markdown
- Create `docs/tools.md` combining auto-generated schemas with hand-written examples
- Add CONTRIBUTING.md and CHANGELOG.md
- Add mermaid architecture diagram to README

### Non-Goals
- Generating a docs website (static site, Hugo, etc.)
- Adding `Example_xxx` test functions for godoc (low priority polish)
- Creating a SECURITY.md (no public vulnerability reporting process yet)
- Updating spec artifacts (they're historical records)

## Decisions

**D1: `replicator docs` iterates the registry.** The command creates an in-memory store, registers all tools (same wiring as `serve.go`), then iterates `registry.List()` to emit markdown. This guarantees the output matches what `tools/list` returns at runtime.

**D2: Category grouping via naming convention.** Tools are grouped by prefix: `hive_*` → Hive, `swarmmail_*` → Swarm Mail, `swarm_*` → Swarm, `hivemind_*` → Memory. The `docs` command derives the category from the tool name prefix.

**D3: Output format is GitHub-flavored markdown.** Each tool gets: `### tool_name`, description paragraph, and `inputSchema` as a fenced JSON code block. The `--output` flag writes to a file; default is stdout.

**D4: `docs/tools.md` is partially generated, partially hand-written.** The auto-generated section (schema tables) is wrapped in `<!-- BEGIN AUTO-GENERATED -->` / `<!-- END AUTO-GENERATED -->` markers. Hand-written examples live outside the markers and are preserved across regeneration.

**D5: README mermaid diagram shows the MCP request flow.** A `flowchart LR` diagram: `Agent --> |stdin JSON-RPC| MCP Server --> Registry --> Tool Handler --> Domain Logic --> SQLite`. This renders natively on GitHub.

## Risks / Trade-offs

**Risk: None.** This is documentation. No production behavior changes. The only new code is the `docs` command which reads the registry (read-only) and writes markdown (stdout or file).

**Trade-off: Hand-written examples require maintenance.** When a tool's schema changes, the auto-generated schema section updates automatically but the hand-written examples may drift. Mitigation: keep examples minimal (3-4 key tools, not all 53).
