// Package memory registers MCP tools for semantic memory operations.
//
// Two primary tools (hivemind_store, hivemind_find) proxy to Dewey's semantic
// search. Six secondary tools return deprecation messages pointing to native
// Dewey equivalents.
package memory

import (
	"encoding/json"
	"errors"

	"github.com/unbound-force/replicator/internal/memory"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

// Register adds all memory tools to the registry.
func Register(reg *registry.Registry, client *memory.Client) {
	// Primary tools -- proxy to Dewey.
	reg.Register(hivemindStore(client))
	reg.Register(hivemindFind(client))

	// Deprecated tools -- return deprecation messages.
	reg.Register(hivemindDeprecated("hivemind_get", "Get specific memory by ID"))
	reg.Register(hivemindDeprecated("hivemind_remove", "Delete outdated/incorrect memory"))
	reg.Register(hivemindDeprecated("hivemind_validate", "Confirm memory is still accurate"))
	reg.Register(hivemindDeprecated("hivemind_stats", "Memory statistics and health check"))
	reg.Register(hivemindDeprecated("hivemind_index", "Index AI session directories"))
	reg.Register(hivemindDeprecated("hivemind_sync", "Sync learnings to git-backed team sharing"))
}

func hivemindStore(client *memory.Client) *registry.Tool {
	return &registry.Tool{
		Name:        "hivemind_store",
		Description: "Store a memory (learnings, decisions, patterns) with metadata and tags. Proxies to Dewey semantic search.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["information"],
			"properties": {
				"information": {"type": "string"},
				"tags":        {"type": "string"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Information string `json:"information"`
				Tags        string `json:"tags"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}

			result, err := client.Store(input.Information, input.Tags)
			if err != nil {
				var unavail *memory.UnavailableError
				if errors.As(err, &unavail) {
					return memory.UnavailableResponse(err), nil
				}
				return "", err
			}

			out, err := json.MarshalIndent(result, "", "  ")
			return string(out), err
		},
	}
}

func hivemindFind(client *memory.Client) *registry.Tool {
	return &registry.Tool{
		Name:        "hivemind_find",
		Description: "Search all memories by semantic similarity. Proxies to Dewey semantic search.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query":      {"type": "string"},
				"collection": {"type": "string"},
				"limit":      {"type": "number"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Query      string `json:"query"`
				Collection string `json:"collection"`
				Limit      int    `json:"limit"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}

			result, err := client.Find(input.Query, input.Collection, input.Limit)
			if err != nil {
				var unavail *memory.UnavailableError
				if errors.As(err, &unavail) {
					return memory.UnavailableResponse(err), nil
				}
				return "", err
			}

			out, err := json.MarshalIndent(result, "", "  ")
			return string(out), err
		},
	}
}

// hivemindDeprecated creates a tool that returns a deprecation message.
func hivemindDeprecated(name, description string) *registry.Tool {
	return &registry.Tool{
		Name:        name,
		Description: description + " (DEPRECATED)",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			return memory.DeprecatedResponse(name), nil
		},
	}
}
