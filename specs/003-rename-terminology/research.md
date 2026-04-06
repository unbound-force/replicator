# Research: Terminology Rename + Agent Kit

**Branch**: `003-rename-terminology` | **Date**: 2026-04-06

## Rename Impact Analysis

### Package Inventory

Analyzed all Go source files to map the rename blast radius.

#### Domain Packages (3 renames)

| Current | New | Files | Test Files | Total |
|---------|-----|-------|------------|-------|
| `internal/hive/` | `internal/org/` | 6 | 5 | 11 |
| `internal/swarmmail/` | `internal/comms/` | 3 | 3 | 6 |
| `internal/swarm/` | `internal/forge/` | 8 | 6 | 14 |

#### Tool Handler Packages (3 renames)

| Current | New | Files |
|---------|-----|-------|
| `internal/tools/hive/` | `internal/tools/org/` | 1 (`tools.go`) |
| `internal/tools/swarmmail/` | `internal/tools/comms/` | 1 (`tools.go`) |
| `internal/tools/swarm/` | `internal/tools/forge/` | 1 (`tools.go`) |

#### Unchanged Packages

| Package | Reason |
|---------|--------|
| `internal/memory/` | Contains `hivemind_*` tools — intentionally unchanged per FR-004 |
| `internal/tools/memory/` | Same — deprecated Dewey proxy stubs |
| `internal/tools/registry/` | Framework code, no tool-specific names |
| `internal/db/` | Database layer, no tool-specific names |
| `internal/config/` | Configuration, no tool-specific names |
| `internal/mcp/` | JSON-RPC server, dispatches by name from registry |
| `internal/ui/` | Lipgloss styles, no tool-specific names |
| `internal/gitutil/` | Git operations, no tool-specific names |
| `internal/doctor/` | Health checks, no tool-specific names |
| `internal/stats/` | Statistics, may reference tool names in queries |
| `internal/query/` | Preset SQL queries, may reference tool names |

### Import Path Consumers

Files that import the renamed packages (found via grep):

```
cmd/replicator/serve.go       — imports hive, swarmmail, swarm tool packages
cmd/replicator/docs.go        — imports hive, swarmmail, swarm tool packages
cmd/replicator/cells.go       — may import hive domain package
test/parity/parity_test.go    — imports all 4 tool packages
internal/tools/org/tools.go   — imports hive domain package (self)
internal/tools/comms/tools.go — imports swarmmail domain package (self)
internal/tools/forge/tools.go — imports swarm domain package (self)
```

### Tool Name String Occurrences

Files containing tool name strings (beyond the tool handler packages):

| File | Old References | Context |
|------|---------------|---------|
| `cmd/replicator/docs.go` | `"hive_"`, `"swarmmail_"`, `"swarm_"` | Category prefix map |
| `test/parity/parity_test.go` | `"hive_create"`, `"hive_cells"`, `"hive_close"`, etc. | Fixture tool name lookups |
| `test/parity/fixtures/hive.json` | Tool name keys | Parity fixture data |
| `test/parity/fixtures/swarmmail.json` | Tool name keys | Parity fixture data |
| `test/parity/fixtures/swarm.json` | Tool name keys | Parity fixture data |
| `internal/swarm/spawn.go` | May reference tool names in prompts | Subtask prompt generation |
| `internal/swarm/review.go` | May reference tool names in prompts | Review prompt generation |
| `internal/swarm/init.go` | May reference tool names | Init response |
| `internal/swarm/progress.go` | May reference tool names | Progress tracking |
| `internal/swarm/insights.go` | May reference tool names | Insight queries |
| `internal/query/presets.go` | May reference tool names in SQL | Preset analytics queries |
| `internal/stats/stats_test.go` | May reference tool names | Test fixtures |
| `internal/mcp/server_test.go` | May reference tool names | Server test fixtures |

### Parity Test Fixture Structure

Each fixture file is a JSON object mapping tool names to `{request, typescript_response}`:

```json
{
  "hive_cells": {
    "request": {"name": "hive_cells", "arguments": {}},
    "typescript_response": {"content": [{"type": "text", "text": "[]"}]}
  }
}
```

The tool name appears in:
1. The top-level key (e.g., `"hive_cells"`)
2. The `request.name` field
3. Hardcoded references in `parity_test.go` (e.g., `reg.Get("hive_create")`)

All three must be updated consistently.

### Docs Command Category Map

Current `categories` in `cmd/replicator/docs.go`:

```go
var categories = []struct {
    prefix string
    name   string
}{
    {"hive_", "Hive"},
    {"swarmmail_", "Swarm Mail"},
    {"swarm_", "Swarm"},
    {"hivemind_", "Memory"},
}
```

**Important ordering note**: The `swarm_` prefix must come after `swarmmail_` in the current code because `swarmmail_` starts with `swarm`. After the rename, `comms_` and `forge_` have no prefix overlap, so ordering is no longer critical. However, `hivemind_` must remain to categorize the 8 deprecated memory tools.

New categories:
```go
{"org_", "Org"},
{"comms_", "Comms"},
{"forge_", "Forge"},
{"hivemind_", "Memory"},
```

### Init Command Current Behavior

`cmd/replicator/init.go` currently:
1. Creates `.uf/replicator/` directory
2. Writes `cells.json` with `[]\n`
3. Returns early if `.uf/replicator/` already exists (full idempotency)

The early return on "already initialized" means the current init is all-or-nothing. For the agent kit, we need finer-grained skip logic: `.uf/replicator/` may exist but `.opencode/` may not.

**Design decision**: Change the init flow to:
1. Create `.uf/replicator/cells.json` (skip if exists)
2. Call `agentkit.Scaffold(targetDir, force)` (per-file skip logic)
3. Report results for both steps

This means removing the early return and checking each artifact independently.

## Agent Kit Content Research

### Upstream Source Analysis

The upstream `joelhooks/swarm-tools` Claude Code plugin provides:

**Commands** (in `claude-plugin/`):
- `/swarm:swarm` — Full Socratic orchestrator (decompose → spawn → monitor → review → complete)
- `/swarm:status` — Check swarm status
- `/swarm:handoff` — End session with handoff notes
- Plus others specific to Claude Code

**Skills** (in `claude-plugin/skills/`):
- `always-on-guidance` — General coding principles
- `swarm-coordination` — How to coordinate parallel workers
- Plus others

**Agents** (in `claude-plugin/agents/`):
- `coordinator` — Orchestrates the swarm
- `worker` — Executes subtasks
- `background-worker` — Runs in background

### Adaptation Strategy

For the replicator agent kit, we adapt the upstream content with these changes:

1. **Tool name references**: All `hive_*` → `org_*`, `swarmmail_*` → `comms_*`, `swarm_*` → `forge_*`
2. **Platform**: Claude Code → OpenCode conventions (`.opencode/` directory structure)
3. **Orchestrator**: Lightweight version of `/swarm:swarm` — delegates to MCP tools without the full Socratic planning flow
4. **Excluded**: `ralph` command (Codex-specific), `cli-builder`/`queue`/`skill-creator`/`skill-generator` (deferred per spec assumptions)

### Agent Kit File Manifest

| # | Path (under `.opencode/`) | Source | Description |
|---|--------------------------|--------|-------------|
| 1 | `command/forge.md` | Adapted from `/swarm:swarm` | Lightweight forge orchestrator |
| 2 | `command/org.md` | New | Quick org cell management |
| 3 | `command/inbox.md` | New | Check comms inbox |
| 4 | `command/forge-status.md` | Adapted from `/swarm:status` | Check forge status |
| 5 | `command/handoff.md` | Adapted from `/swarm:handoff` | Session handoff |
| 6 | `skills/always-on-guidance/SKILL.md` | Adapted | General coding guidance |
| 7 | `skills/forge-coordination/SKILL.md` | Adapted from `swarm-coordination` | Forge coordination patterns |
| 8 | `skills/replicator-cli/SKILL.md` | New | Replicator CLI reference |
| 9 | `skills/testing-patterns/SKILL.md` | New | Go testing patterns for replicator |
| 10 | `skills/system-design/SKILL.md` | New | System design principles |
| 11 | `skills/learning-systems/SKILL.md` | New | Learning from outcomes |
| 12 | `skills/forge-global/SKILL.md` | Adapted | Global forge skill |
| 13 | `agents/coordinator.md` | Adapted | Coordinator agent definition |
| 14 | `agents/worker.md` | Adapted | Worker agent definition |
| 15 | `agents/background-worker.md` | Adapted | Background worker definition |

### embed.FS Design

```go
package agentkit

import "embed"

//go:embed content/*
var content embed.FS

// Scaffold writes agent kit files to targetDir/.opencode/.
// Returns a list of results (created/skipped/overwritten).
func Scaffold(targetDir string, force bool) ([]ScaffoldResult, error) {
    // Walk content FS
    // For each file, check if target exists
    // Write if !exists || force
    // Return results
}

type ScaffoldResult struct {
    Path   string // relative path under .opencode/
    Action string // "created", "skipped", "overwritten"
}
```

The `embed` directive `//go:embed content/*` captures all files under `internal/agentkit/content/`. The `fs.WalkDir` function traverses the embedded filesystem. Each file is written to the corresponding path under `.opencode/` in the target directory.

## Technical Decisions

### Decision 1: Clean Break vs. Aliases

**Chosen**: Clean break (no backward-compatible aliases).

**Rationale**: Aliases add complexity (dual registration, documentation confusion) for minimal benefit. The rename is a major version change. Agents update their prompts once.

**Rejected**: Registering both old and new names pointing to the same handler. This would inflate the tool count from 53 to 98 and confuse `tools/list` consumers.

### Decision 2: Package Names

**Chosen**: `org`, `comms`, `forge` (short, Go-idiomatic).

**Rationale**: Go convention favors short package names. `org` is clear for "organization of work items". `comms` is the natural abbreviation of "communications". `forge` evokes "forging" work through orchestration.

**Rejected**: `organization` (too long), `comm` (ambiguous), `orchestrator` (too long, conflicts with existing coordinator concept).

### Decision 3: Agent Kit Location

**Chosen**: `internal/agentkit/` with `embed.FS`.

**Rationale**: Follows the existing `internal/` package pattern. `embed.FS` is stdlib — no new dependencies. Files are compiled into the binary, ensuring version consistency.

**Rejected**: External file distribution (requires separate install step), `go generate` (adds build complexity), runtime HTTP fetch (violates Composability First principle).

### Decision 4: Init Behavior Change

**Chosen**: Per-file skip logic with `--force` flag.

**Rationale**: The current all-or-nothing early return doesn't work when init creates artifacts in two locations (`.uf/` and `.opencode/`). Per-file skip logic is more user-friendly and matches the spec's acceptance scenarios.

**Rejected**: Separate `replicator init-kit` command (fragments the UX), always-overwrite (destructive, violates user expectations).
