// Package config manages replicator configuration.
//
// The global database lives at ~/.config/swarm-tools/swarm.db
// (shared with the TypeScript cyborg-swarm for migration compatibility).
package config

import (
	"os"
	"path/filepath"
)

// Config holds runtime configuration.
type Config struct {
	// DatabasePath is the path to the global SQLite database.
	DatabasePath string

	// DeweyURL is the Dewey MCP server endpoint.
	DeweyURL string

	// ZenAPIKey is the OpenCode Zen API key for LLM calls.
	ZenAPIKey string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		DatabasePath: envOr("REPLICATOR_DB", defaultDatabasePath()),
		DeweyURL:     envOr("DEWEY_MCP_URL", "http://localhost:3333/mcp/"),
		ZenAPIKey:    os.Getenv("ZEN_API_KEY"),
	}
}

func defaultDatabasePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "swarm.db"
	}
	dir := filepath.Join(home, ".config", "swarm-tools")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "swarm.db")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
