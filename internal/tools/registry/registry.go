// Package registry provides a tool registration system for MCP tools.
//
// Tools register themselves with a name, description, JSON schema for
// arguments, and an execute function. The MCP server dispatches
// tools/call requests to the registered handler.
package registry

import "encoding/json"

// Tool defines an MCP tool that can be called by AI agents.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
	Execute     ExecuteFunc     `json:"-"`
}

// ExecuteFunc is the handler for a tool invocation.
// It receives the raw JSON arguments and returns a result string.
type ExecuteFunc func(args json.RawMessage) (string, error)

// Registry holds registered MCP tools.
type Registry struct {
	tools map[string]*Tool
	order []string // preserve registration order for tools/list
}

// New creates an empty tool registry.
func New() *Registry {
	return &Registry{
		tools: make(map[string]*Tool),
	}
}

// Register adds a tool to the registry.
func (r *Registry) Register(t *Tool) {
	r.tools[t.Name] = t
	r.order = append(r.order, t.Name)
}

// Get returns a tool by name, or nil if not found.
func (r *Registry) Get(name string) *Tool {
	return r.tools[name]
}

// List returns all registered tools in registration order.
func (r *Registry) List() []*Tool {
	result := make([]*Tool, 0, len(r.order))
	for _, name := range r.order {
		if t, ok := r.tools[name]; ok {
			result = append(result, t)
		}
	}
	return result
}

// Count returns the number of registered tools.
func (r *Registry) Count() int {
	return len(r.tools)
}
