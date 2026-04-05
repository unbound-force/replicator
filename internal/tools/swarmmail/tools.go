// Package swarmmail registers MCP tools for swarm mail operations.
package swarmmail

import (
	"encoding/json"
	"fmt"

	"github.com/unbound-force/replicator/internal/db"
	"github.com/unbound-force/replicator/internal/swarmmail"
	"github.com/unbound-force/replicator/internal/tools/registry"
)

// Register adds all swarm mail tools to the registry.
func Register(reg *registry.Registry, store *db.Store) {
	reg.Register(swarmmailInit(store))
	reg.Register(swarmmailSend(store))
	reg.Register(swarmmailInbox(store))
	reg.Register(swarmmailReadMessage(store))
	reg.Register(swarmmailReserve(store))
	reg.Register(swarmmailRelease(store))
	reg.Register(swarmmailReleaseAll(store))
	reg.Register(swarmmailReleaseAgent(store))
	reg.Register(swarmmailAck(store))
	reg.Register(swarmmailHealth(store))
}

func swarmmailInit(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_init",
		Description: "Initialize swarm mail session. Registers the agent.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["project_path"],
			"properties": {
				"agent_name":       {"type": "string"},
				"project_path":     {"type": "string"},
				"task_description": {"type": "string"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				AgentName       string `json:"agent_name"`
				ProjectPath     string `json:"project_path"`
				TaskDescription string `json:"task_description"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			agentName := input.AgentName
			if agentName == "" {
				agentName = "default"
			}
			agent, err := swarmmail.Init(store, agentName, input.ProjectPath, input.TaskDescription)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(agent, "", "  ")
			return string(out), err
		},
	}
}

func swarmmailSend(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_send",
		Description: "Send a message to other agents via swarm mail.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["to", "subject", "body"],
			"properties": {
				"to":           {"type": "array", "items": {"type": "string"}},
				"subject":      {"type": "string"},
				"body":         {"type": "string"},
				"importance":   {"type": "string", "enum": ["low", "normal", "high", "urgent"]},
				"thread_id":    {"type": "string"},
				"ack_required": {"type": "boolean"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				From string `json:"from"`
				swarmmail.SendInput
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			fromAgent := input.From
			if fromAgent == "" {
				fromAgent = "unknown"
			}
			if err := swarmmail.Send(store, fromAgent, input.SendInput); err != nil {
				return "", err
			}
			return `{"status": "sent"}`, nil
		},
	}
}

func swarmmailInbox(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_inbox",
		Description: "Fetch inbox (max 5 messages, bodies excluded).",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"limit":       {"type": "number", "maximum": 5},
				"urgent_only": {"type": "boolean"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				AgentName  string `json:"agent_name"`
				Limit      int    `json:"limit"`
				UrgentOnly bool   `json:"urgent_only"`
			}
			if len(args) > 0 {
				json.Unmarshal(args, &input)
			}
			agentName := input.AgentName
			if agentName == "" {
				agentName = "default"
			}
			summaries, err := swarmmail.Inbox(store, agentName, input.Limit, input.UrgentOnly)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(summaries, "", "  ")
			return string(out), err
		},
	}
}

func swarmmailReadMessage(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_read_message",
		Description: "Fetch one message body by ID.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["message_id"],
			"properties": {
				"message_id": {"type": "number"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				MessageID int `json:"message_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			msg, err := swarmmail.ReadMessage(store, input.MessageID)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(msg, "", "  ")
			return string(out), err
		},
	}
}

func swarmmailReserve(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_reserve",
		Description: "Reserve file paths for exclusive editing.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["paths"],
			"properties": {
				"paths":       {"type": "array", "items": {"type": "string"}},
				"exclusive":   {"type": "boolean"},
				"reason":      {"type": "string"},
				"ttl_seconds": {"type": "number"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				AgentName  string   `json:"agent_name"`
				Paths      []string `json:"paths"`
				Exclusive  bool     `json:"exclusive"`
				Reason     string   `json:"reason"`
				TTLSeconds int      `json:"ttl_seconds"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			agentName := input.AgentName
			if agentName == "" {
				agentName = "default"
			}
			reservations, err := swarmmail.Reserve(store, agentName, input.Paths, input.Exclusive, input.Reason, input.TTLSeconds)
			if err != nil {
				return "", err
			}
			out, err := json.MarshalIndent(reservations, "", "  ")
			return string(out), err
		},
	}
}

func swarmmailRelease(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_release",
		Description: "Release file reservations.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"paths":           {"type": "array", "items": {"type": "string"}},
				"reservation_ids": {"type": "array", "items": {"type": "number"}}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Paths          []string `json:"paths"`
				ReservationIDs []int    `json:"reservation_ids"`
			}
			if len(args) > 0 {
				json.Unmarshal(args, &input)
			}
			if err := swarmmail.Release(store, input.Paths, input.ReservationIDs); err != nil {
				return "", err
			}
			return `{"status": "released"}`, nil
		},
	}
}

func swarmmailReleaseAll(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_release_all",
		Description: "Release all file reservations in the project.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			if err := swarmmail.ReleaseAll(store, ""); err != nil {
				return "", err
			}
			return `{"status": "all_released"}`, nil
		},
	}
}

func swarmmailReleaseAgent(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_release_agent",
		Description: "Release all file reservations for a specific agent.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["agent_name"],
			"properties": {
				"agent_name": {"type": "string"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				AgentName string `json:"agent_name"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			if err := swarmmail.ReleaseAgent(store, input.AgentName); err != nil {
				return "", err
			}
			return fmt.Sprintf(`{"status": "released", "agent": %q}`, input.AgentName), nil
		},
	}
}

func swarmmailAck(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_ack",
		Description: "Acknowledge a message.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"required": ["message_id"],
			"properties": {
				"message_id": {"type": "number"}
			}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				MessageID int `json:"message_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", err
			}
			if err := swarmmail.Ack(store, input.MessageID); err != nil {
				return "", err
			}
			return `{"status": "acknowledged"}`, nil
		},
	}
}

func swarmmailHealth(store *db.Store) *registry.Tool {
	return &registry.Tool{
		Name:        "swarmmail_health",
		Description: "Check swarm mail database health.",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {}
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			// Verify the database is accessible by counting key tables.
			var agentCount, msgCount, resCount int
			store.DB.QueryRow("SELECT COUNT(*) FROM agents").Scan(&agentCount)
			store.DB.QueryRow("SELECT COUNT(*) FROM messages").Scan(&msgCount)
			store.DB.QueryRow("SELECT COUNT(*) FROM reservations").Scan(&resCount)

			result := map[string]any{
				"status":       "healthy",
				"agents":       agentCount,
				"messages":     msgCount,
				"reservations": resCount,
			}
			out, err := json.MarshalIndent(result, "", "  ")
			return string(out), err
		},
	}
}
