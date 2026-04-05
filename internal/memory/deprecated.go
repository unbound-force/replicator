package memory

import (
	"encoding/json"
	"fmt"
)

// deprecationMap maps deprecated hivemind tool names to their Dewey replacements.
var deprecationMap = map[string]string{
	"hivemind_get":      "dewey_get_page",
	"hivemind_remove":   "dewey_delete_page",
	"hivemind_validate": "", // no direct equivalent
	"hivemind_stats":    "dewey_health",
	"hivemind_index":    "dewey_reload",
	"hivemind_sync":     "dewey_reload",
}

// DeprecatedResponse returns a JSON string indicating the tool is deprecated.
// The response includes the tool name, a human-readable message, and the
// replacement tool (if one exists).
func DeprecatedResponse(toolName string) string {
	replacement, ok := deprecationMap[toolName]
	if !ok {
		// Unknown tool -- still return a deprecation message.
		replacement = ""
	}

	msg := fmt.Sprintf("%s is deprecated.", toolName)
	if replacement != "" {
		msg += fmt.Sprintf(" Use %s instead.", replacement)
	} else {
		msg += " No direct replacement is available."
	}

	resp := map[string]any{
		"deprecated":  true,
		"tool":        toolName,
		"message":     msg,
		"replacement": replacement,
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	return string(out)
}
