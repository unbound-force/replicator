package memory

import (
	"encoding/json"
	"testing"
)

func TestDeprecatedResponse_WithReplacement(t *testing.T) {
	tests := []struct {
		tool        string
		replacement string
	}{
		{"hivemind_get", "dewey_get_page"},
		{"hivemind_remove", "dewey_delete_page"},
		{"hivemind_stats", "dewey_health"},
		{"hivemind_index", "dewey_reload"},
		{"hivemind_sync", "dewey_reload"},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			resp := DeprecatedResponse(tt.tool)

			var parsed map[string]any
			if err := json.Unmarshal([]byte(resp), &parsed); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if parsed["deprecated"] != true {
				t.Error("deprecated should be true")
			}
			if parsed["tool"] != tt.tool {
				t.Errorf("tool = %v, want %q", parsed["tool"], tt.tool)
			}
			if parsed["replacement"] != tt.replacement {
				t.Errorf("replacement = %v, want %q", parsed["replacement"], tt.replacement)
			}
		})
	}
}

func TestDeprecatedResponse_NoReplacement(t *testing.T) {
	resp := DeprecatedResponse("hivemind_validate")

	var parsed map[string]any
	if err := json.Unmarshal([]byte(resp), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if parsed["deprecated"] != true {
		t.Error("deprecated should be true")
	}
	if parsed["replacement"] != "" {
		t.Errorf("replacement = %v, want empty string", parsed["replacement"])
	}

	msg, ok := parsed["message"].(string)
	if !ok {
		t.Fatal("message should be a string")
	}
	if msg == "" {
		t.Error("message should not be empty")
	}
}

func TestDeprecatedResponse_UnknownTool(t *testing.T) {
	resp := DeprecatedResponse("hivemind_unknown")

	var parsed map[string]any
	if err := json.Unmarshal([]byte(resp), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if parsed["deprecated"] != true {
		t.Error("deprecated should be true")
	}
	if parsed["tool"] != "hivemind_unknown" {
		t.Errorf("tool = %v, want %q", parsed["tool"], "hivemind_unknown")
	}
}

func TestDeprecatedResponse_ValidJSON(t *testing.T) {
	// Verify all known deprecated tools produce valid JSON.
	tools := []string{
		"hivemind_get", "hivemind_remove", "hivemind_validate",
		"hivemind_stats", "hivemind_index", "hivemind_sync",
	}

	for _, tool := range tools {
		resp := DeprecatedResponse(tool)
		if !json.Valid([]byte(resp)) {
			t.Errorf("DeprecatedResponse(%q) produced invalid JSON: %s", tool, resp)
		}
	}
}
