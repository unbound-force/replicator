package forge

import (
	"strings"
	"testing"
)

func TestSelectStrategy_RiskBased(t *testing.T) {
	tests := []string{"Fix security vulnerability", "Database migration", "Breaking API change"}
	for _, task := range tests {
		strategy, err := SelectStrategy(task, "")
		if err != nil {
			t.Fatalf("SelectStrategy(%q): %v", task, err)
		}
		if strategy != "risk-based" {
			t.Errorf("SelectStrategy(%q) = %q, want %q", task, strategy, "risk-based")
		}
	}
}

func TestSelectStrategy_FileBased(t *testing.T) {
	tests := []string{"Refactor the auth module", "Rename package", "Move files to new directory"}
	for _, task := range tests {
		strategy, err := SelectStrategy(task, "")
		if err != nil {
			t.Fatalf("SelectStrategy(%q): %v", task, err)
		}
		if strategy != "file-based" {
			t.Errorf("SelectStrategy(%q) = %q, want %q", task, strategy, "file-based")
		}
	}
}

func TestSelectStrategy_FeatureBased(t *testing.T) {
	strategy, err := SelectStrategy("Add user authentication", "")
	if err != nil {
		t.Fatalf("SelectStrategy: %v", err)
	}
	if strategy != "feature-based" {
		t.Errorf("strategy = %q, want %q", strategy, "feature-based")
	}
}

func TestSelectStrategy_Empty(t *testing.T) {
	_, err := SelectStrategy("", "")
	if err == nil {
		t.Error("expected error for empty task")
	}
}

func TestSelectStrategy_ContextInfluence(t *testing.T) {
	// Context can also trigger strategy selection.
	strategy, err := SelectStrategy("Update the API", "this is a critical security fix")
	if err != nil {
		t.Fatalf("SelectStrategy: %v", err)
	}
	if strategy != "risk-based" {
		t.Errorf("strategy = %q, want %q", strategy, "risk-based")
	}
}

func TestPlanPrompt(t *testing.T) {
	prompt := PlanPrompt("Build auth system", "file-based", "Go project", 3)

	if !strings.Contains(prompt, "file-based") {
		t.Error("prompt should mention strategy")
	}
	if !strings.Contains(prompt, "Build auth system") {
		t.Error("prompt should contain task")
	}
	if !strings.Contains(prompt, "Go project") {
		t.Error("prompt should contain context")
	}
	if !strings.Contains(prompt, "3") {
		t.Error("prompt should mention max subtasks")
	}
}

func TestPlanPrompt_Defaults(t *testing.T) {
	prompt := PlanPrompt("task", "", "", 0)
	if !strings.Contains(prompt, "feature-based") {
		t.Error("default strategy should be feature-based")
	}
	if !strings.Contains(prompt, "5") {
		t.Error("default max subtasks should be 5")
	}
}

func TestDecompose(t *testing.T) {
	prompt := Decompose("Build auth system", "Go project", 4)

	if !strings.Contains(prompt, "Build auth system") {
		t.Error("prompt should contain task")
	}
	if !strings.Contains(prompt, "4") {
		t.Error("prompt should mention max subtasks")
	}
	if !strings.Contains(prompt, "epic_title") {
		t.Error("prompt should include output format")
	}
}

func TestValidateDecomposition_Valid(t *testing.T) {
	input := `{
		"epic_title": "Auth System",
		"subtasks": [
			{"title": "Add login", "files": ["auth.go"], "priority": 1}
		]
	}`

	result, err := ValidateDecomposition(input)
	if err != nil {
		t.Fatalf("ValidateDecomposition: %v", err)
	}
	if result["epic_title"] != "Auth System" {
		t.Errorf("epic_title = %v, want %q", result["epic_title"], "Auth System")
	}
}

func TestValidateDecomposition_FromMarkdown(t *testing.T) {
	input := "Here is the plan:\n```json\n" + `{
		"epic_title": "Auth",
		"subtasks": [{"title": "Login"}]
	}` + "\n```\n"

	result, err := ValidateDecomposition(input)
	if err != nil {
		t.Fatalf("ValidateDecomposition: %v", err)
	}
	if result["epic_title"] != "Auth" {
		t.Errorf("epic_title = %v, want %q", result["epic_title"], "Auth")
	}
}

func TestValidateDecomposition_Empty(t *testing.T) {
	_, err := ValidateDecomposition("")
	if err == nil {
		t.Error("expected error for empty response")
	}
}

func TestValidateDecomposition_MissingEpicTitle(t *testing.T) {
	_, err := ValidateDecomposition(`{"subtasks": [{"title": "a"}]}`)
	if err == nil {
		t.Error("expected error for missing epic_title")
	}
}

func TestValidateDecomposition_MissingSubtasks(t *testing.T) {
	_, err := ValidateDecomposition(`{"epic_title": "test"}`)
	if err == nil {
		t.Error("expected error for missing subtasks")
	}
}

func TestValidateDecomposition_EmptySubtasks(t *testing.T) {
	_, err := ValidateDecomposition(`{"epic_title": "test", "subtasks": []}`)
	if err == nil {
		t.Error("expected error for empty subtasks")
	}
}

func TestValidateDecomposition_SubtaskMissingTitle(t *testing.T) {
	_, err := ValidateDecomposition(`{"epic_title": "test", "subtasks": [{"files": ["a.go"]}]}`)
	if err == nil {
		t.Error("expected error for subtask missing title")
	}
}

func TestValidateDecomposition_InvalidJSON(t *testing.T) {
	_, err := ValidateDecomposition("not json at all")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
