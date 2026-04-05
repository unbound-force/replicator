package swarm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/unbound-force/replicator/internal/gitutil"
)

// initRepo creates a temporary git repo with an initial commit for worktree tests.
func initRepo(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("requires git")
	}

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if _, err := gitutil.Run(dir, args...); err != nil {
			t.Fatalf("git %s: %v", strings.Join(args, " "), err)
		}
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")

	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")

	return dir
}

func TestWorktreeCreate(t *testing.T) {
	repo := initRepo(t)
	commit, _ := gitutil.CurrentCommit(repo)

	result, err := WorktreeCreate(repo, "task-1", commit)
	if err != nil {
		t.Fatalf("WorktreeCreate: %v", err)
	}
	if result["status"] != "created" {
		t.Errorf("status = %v, want %q", result["status"], "created")
	}

	wtPath := result["worktree_path"].(string)
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("worktree directory should exist")
	}
}

func TestWorktreeCreate_MissingArgs(t *testing.T) {
	tests := []struct {
		name        string
		projectPath string
		taskID      string
		startCommit string
	}{
		{"missing project_path", "", "task", "abc"},
		{"missing task_id", "/tmp", "", "abc"},
		{"missing start_commit", "/tmp", "task", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := WorktreeCreate(tt.projectPath, tt.taskID, tt.startCommit)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestWorktreeMerge(t *testing.T) {
	repo := initRepo(t)
	commit, _ := gitutil.CurrentCommit(repo)

	// Create worktree and make a commit.
	WorktreeCreate(repo, "merge-task", commit)
	wtPath := filepath.Join(repo, ".worktrees", "merge-task")

	f := filepath.Join(wtPath, "new.txt")
	os.WriteFile(f, []byte("merge content\n"), 0o644)
	gitutil.Run(wtPath, "add", ".")
	gitutil.Run(wtPath, "commit", "-m", "worktree commit")

	result, err := WorktreeMerge(repo, "merge-task", commit)
	if err != nil {
		t.Fatalf("WorktreeMerge: %v", err)
	}
	if result["status"] != "merged" {
		t.Errorf("status = %v, want %q", result["status"], "merged")
	}

	// Verify file exists in main repo.
	mainFile := filepath.Join(repo, "new.txt")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Error("merged file should exist in main repo")
	}
}

func TestWorktreeMerge_MissingArgs(t *testing.T) {
	_, err := WorktreeMerge("", "task", "abc")
	if err == nil {
		t.Error("expected error for missing project_path")
	}
	_, err = WorktreeMerge("/tmp", "", "abc")
	if err == nil {
		t.Error("expected error for missing task_id")
	}
}

func TestWorktreeCleanup(t *testing.T) {
	repo := initRepo(t)
	commit, _ := gitutil.CurrentCommit(repo)

	WorktreeCreate(repo, "cleanup-task", commit)

	result, err := WorktreeCleanup(repo, "cleanup-task", false)
	if err != nil {
		t.Fatalf("WorktreeCleanup: %v", err)
	}
	if result["status"] != "cleaned" {
		t.Errorf("status = %v, want %q", result["status"], "cleaned")
	}
}

func TestWorktreeCleanup_Idempotent(t *testing.T) {
	repo := initRepo(t)

	// Cleanup a non-existent worktree should not error.
	result, err := WorktreeCleanup(repo, "nonexistent", false)
	if err != nil {
		t.Fatalf("WorktreeCleanup: %v", err)
	}
	if result["status"] != "cleaned" {
		t.Errorf("status = %v, want %q", result["status"], "cleaned")
	}
}

func TestWorktreeCleanup_MissingProjectPath(t *testing.T) {
	_, err := WorktreeCleanup("", "", false)
	if err == nil {
		t.Error("expected error for missing project_path")
	}
}

func TestWorktreeCleanup_MissingTaskID(t *testing.T) {
	if testing.Short() {
		t.Skip("requires git")
	}
	_, err := WorktreeCleanup("/tmp", "", false)
	if err == nil {
		t.Error("expected error for missing task_id when cleanup_all is false")
	}
}

func TestWorktreeList(t *testing.T) {
	repo := initRepo(t)

	result, err := WorktreeList(repo)
	if err != nil {
		t.Fatalf("WorktreeList: %v", err)
	}
	count := result["count"].(int)
	if count < 1 {
		t.Errorf("expected at least 1 worktree, got %d", count)
	}
}

func TestWorktreeList_MissingProjectPath(t *testing.T) {
	_, err := WorktreeList("")
	if err == nil {
		t.Error("expected error for missing project_path")
	}
}
