## Context

Dewey PR #28 renamed 4 MCP tool registrations in `server.go`:

| Before | After | Agent sees (OpenCode prefix) |
|--------|-------|------------------------------|
| `dewey_semantic_search` | `semantic_search` | `dewey_semantic_search` |
| `dewey_similar` | `similar` | `dewey_similar` |
| `dewey_semantic_search_filtered` | `semantic_search_filtered` | `dewey_semantic_search_filtered` |
| `dewey_store_learning` | `store_learning` | `dewey_store_learning` |

Replicator's `internal/memory/proxy.go` calls Dewey directly via HTTP JSON-RPC, using the **internal** MCP tool names. These are now wrong.

There are two layers of naming:
1. **Internal names** (used in direct HTTP calls to Dewey's MCP endpoint): `store_learning`, `semantic_search`
2. **Agent-facing names** (what OpenCode agents see, with server prefix): `dewey_store_learning`, `dewey_semantic_search`

## Goals / Non-Goals

### Goals
- Update internal tool names in `proxy.go` to match Dewey's new registrations
- Update deprecation warning text to show correct agent-facing names
- Update test assertions to match new names

### Non-Goals
- Changing the MCP protocol usage pattern (the `Call` method structure stays as-is)
- Updating completed spec artifacts (historical records)
- Modifying deprecated tool mappings (`deprecated.go` -- those reference agent-facing names which are unchanged)

## Decisions

**D1: Simple string replacement.** All 8 changes are string literal updates -- no logic, no schema, no type changes. The proxy's behavior is unchanged; only the tool name strings are different.

**D2: Deprecation warnings use agent-facing names.** The `_warning` field in proxy responses tells agents what tool to use instead. These should show the OpenCode-prefixed names (`dewey_store_learning`, `dewey_semantic_search`) since that's what agents actually type.

**D3: Leave spec artifacts as-is.** `specs/001-go-rewrite-phases/plan.md` and `tasks.md` reference the old names in historical task descriptions. Updating them adds noise to git history with no functional benefit.

## Risks / Trade-offs

**Risk: None.** This is a mechanical find-and-replace with no behavioral impact. The proxy sends the correct tool name, Dewey receives it, behavior is identical.

**Trade-off: Not fixing the MCP protocol issue.** The `Call` method puts the tool name in the JSON-RPC `method` field rather than using standard `tools/call` with `params.name`. This works because Dewey accepts both patterns, but it's non-standard. Fixing this is out of scope for this change -- it would require restructuring the `Call` method and all tests.
