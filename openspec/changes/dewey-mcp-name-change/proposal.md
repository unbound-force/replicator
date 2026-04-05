## Why

Dewey PR #28 dropped the `dewey_` prefix from 4 MCP tool registration names to eliminate the `dewey_dewey_*` double-prefix that appeared in OpenCode agent UIs. Replicator's memory proxy calls Dewey directly via HTTP JSON-RPC and references the old tool names. These calls now fail because Dewey no longer registers tools under those names.

## What Changes

### Modified Capabilities
- `hivemind_store`: Update the Dewey proxy call from `dewey_dewey_store_learning` to `store_learning` (the new internal MCP tool name)
- `hivemind_find`: Update the Dewey proxy call from `dewey_dewey_semantic_search` to `semantic_search` (the new internal MCP tool name)
- Deprecation warning text: Update agent-facing tool names from `dewey_dewey_*` to `dewey_*` (the correct OpenCode-prefixed names)

## Impact

- `internal/memory/proxy.go` -- 2 tool name strings in `Call()` invocations, 2 deprecation warning strings, 2 comments
- `internal/memory/proxy_test.go` -- 2 test assertions for expected method names
- No behavioral changes -- same proxy logic, same error handling, same deprecation pattern
- No changes to deprecated tool mappings (they reference agent-facing names like `dewey_search` which are unaffected)

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: PASS

This change updates a string constant in a proxy client. Replicator continues to communicate with Dewey through well-defined JSON-RPC artifacts. No coupling changes.

### II. Composability First

**Assessment**: PASS

Replicator remains independently installable. When Dewey is unavailable, the proxy returns structured `DEWEY_UNAVAILABLE` errors. This change does not introduce new dependencies.

### III. Observable Quality

**Assessment**: PASS

The deprecation warnings in tool responses are updated to show the correct agent-facing tool names (`dewey_store_learning` instead of `dewey_dewey_store_learning`), improving the quality of guidance agents receive.

### IV. Testability

**Assessment**: PASS

Tests use `httptest.NewServer` to mock Dewey. The mock matches on the tool name string, so updating the test assertions maintains isolation without requiring a live Dewey instance.
