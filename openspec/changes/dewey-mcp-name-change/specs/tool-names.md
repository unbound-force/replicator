## MODIFIED Requirements

### Requirement: Dewey proxy store tool name

The `hivemind_store` proxy MUST call Dewey's MCP endpoint using the internal tool name `store_learning`.

Previously: The proxy called `dewey_dewey_store_learning`.

#### Scenario: Store a learning via proxy
- **GIVEN** Dewey is available at the configured MCP endpoint
- **WHEN** an agent calls `hivemind_store` with information and tags
- **THEN** the proxy sends a JSON-RPC request to Dewey with method `store_learning`

### Requirement: Dewey proxy find tool name

The `hivemind_find` proxy MUST call Dewey's MCP endpoint using the internal tool name `semantic_search`.

Previously: The proxy called `dewey_dewey_semantic_search`.

#### Scenario: Search memories via proxy
- **GIVEN** Dewey is available at the configured MCP endpoint
- **WHEN** an agent calls `hivemind_find` with a query
- **THEN** the proxy sends a JSON-RPC request to Dewey with method `semantic_search`

### Requirement: Deprecation warnings reference agent-facing names

Deprecation warnings in proxy responses MUST reference the correct agent-facing tool names: `dewey_store_learning` and `dewey_semantic_search` (with OpenCode's `dewey_` prefix, not the internal MCP name).

Previously: Warnings referenced `dewey_dewey_store_learning` and `dewey_dewey_semantic_search`.

#### Scenario: Store response includes correct deprecation hint
- **GIVEN** an agent calls `hivemind_store` successfully
- **WHEN** the response is returned
- **THEN** the `_warning` field contains `dewey_store_learning` (not `dewey_dewey_store_learning`)
