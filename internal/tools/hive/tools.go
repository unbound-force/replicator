// Package hive registers MCP tools for hive cell operations.
package hive

import (
	"encoding/json"

	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/hive"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

// Register adds all hive tools to the registry.
func Register(reg *registry.Registry, store *db.Store) {
	reg.Register(hiveCells(store))
	reg.Register(hiveCreate(store))
	reg.Register(hiveClose(store))
	reg.Register(hiveUpdate(store))
}

func hiveCells(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "hive_cells",
		Description: "Query cells from the hive database with flexible filtering. Use to list work items, find by status/type, or look up a cell by partial ID.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"id":     {"type": "string", "description": "Partial cell ID to match"},
				"status": {"type": "string", "enum": ["open", "in_progress", "blocked", "closed"]},
				"type":   {"type": "string", "enum": ["task", "bug", "feature", "epic", "chore"]},
				"ready":  {"type": "boolean", "description": "If true, return only unblocked cells"},
				"limit":  {"type": "number", "description": "Max results (default 50)"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var q hive.CellQuery
			if len(args) > 0 {
				json.Unmarshal(args, &q)
			}
			cells, err := hive.QueryCells(store, q)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(cells, "", "  ")
			return string(out), err
		},
	}
}

func hiveCreate(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "hive_create",
		Description: "Create a new cell (work item) in the hive.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["title"],
			"properties": {
				"title":       {"type": "string"},
				"description": {"type": "string"},
				"type":        {"type": "string", "enum": ["task", "bug", "feature", "epic", "chore"]},
				"priority":    {"type": "number", "minimum": 0, "maximum": 3},
				"parent_id":   {"type": "string"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input hive.CreateCellInput
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			cell, err := hive.CreateCell(store, input)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(cell, "", "  ")
			return string(out), err
		},
	}
}

func hiveClose(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "hive_close",
		Description: "Close a cell with a reason.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["id", "reason"],
			"properties": {
				"id":     {"type": "string"},
				"reason": {"type": "string"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			if err := hive.CloseCell(store, input.ID, input.Reason); err != nil {
				return "", err
			}
			return `{"status": "closed"}`, nil
		},
	}
}

func hiveUpdate(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "hive_update",
		Description: "Update a cell's status, description, or priority.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["id"],
			"properties": {
				"id":          {"type": "string"},
				"status":      {"type": "string", "enum": ["open", "in_progress", "blocked", "closed"]},
				"description": {"type": "string"},
				"priority":    {"type": "number", "minimum": 0, "maximum": 3}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				ID          string  `json:"id"`
				Status      *string `json:"status,omitempty"`
				Description *string `json:"description,omitempty"`
				Priority    *int    `json:"priority,omitempty"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			if err := hive.UpdateCell(store, input.ID, input.Status, input.Description, input.Priority); err != nil {
				return "", err
			}
			return `{"status": "updated"}`, nil
		},
	}
}
