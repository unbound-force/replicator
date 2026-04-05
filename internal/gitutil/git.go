// Package gitutil provides thin wrappers around git commands for worktree
// management and commit operations.
//
// All functions shell out to the `git` binary. Tests use t.TempDir() with
// `git init` to create isolated repositories.
package gitutil

import (
	"fmt"
	"os/exec"
	"strings"
)

// WorktreeInfo describes a single git worktree entry.
type WorktreeInfo struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

// Run executes a git command in the given directory, returning stdout.
// If the command fails, the error includes stderr for diagnostics.
func Run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w\nstderr: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

// WorktreeAdd creates a new git worktree at worktreePath on the given branch,
// starting from startCommit.
func WorktreeAdd(projectPath, worktreePath, branch, startCommit string) error {
	_, err := Run(projectPath, "worktree", "add", "-b", branch, worktreePath, startCommit)
	if err != nil {
		return fmt.Errorf("worktree add: %w", err)
	}
	return nil
}

// WorktreeRemove removes a git worktree at the given path.
// projectPath is the main repository directory from which the worktree was created.
func WorktreeRemove(projectPath, worktreePath string) error {
	_, err := Run(projectPath, "worktree", "remove", worktreePath, "--force")
	if err != nil {
		return fmt.Errorf("worktree remove %s: %w", worktreePath, err)
	}
	return nil
}

// WorktreeList parses `git worktree list --porcelain` output into structured info.
func WorktreeList(projectPath string) ([]WorktreeInfo, error) {
	out, err := Run(projectPath, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("worktree list: %w", err)
	}

	var worktrees []WorktreeInfo
	var current WorktreeInfo

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "worktree "):
			current.Path = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "HEAD "):
			current.Commit = strings.TrimPrefix(line, "HEAD ")
		case strings.HasPrefix(line, "branch "):
			current.Branch = strings.TrimPrefix(line, "branch ")
		case line == "" || line == "bare":
			// End of a worktree block -- save if we have a path.
			if current.Path != "" {
				worktrees = append(worktrees, current)
			}
			current = WorktreeInfo{}
		}
	}
	// Capture the last entry if the output doesn't end with a blank line.
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	if worktrees == nil {
		worktrees = []WorktreeInfo{}
	}
	return worktrees, nil
}

// CherryPick identifies commits on worktreeBranch since startCommit and
// cherry-picks each into the current branch of projectPath.
func CherryPick(projectPath, worktreeBranch, startCommit string) error {
	// List commits from startCommit (exclusive) to the tip of worktreeBranch.
	out, err := Run(projectPath, "log", "--reverse", "--format=%H", startCommit+".."+worktreeBranch)
	if err != nil {
		return fmt.Errorf("list commits: %w", err)
	}

	if out == "" {
		return nil // No commits to cherry-pick.
	}

	for _, commit := range strings.Split(out, "\n") {
		commit = strings.TrimSpace(commit)
		if commit == "" {
			continue
		}
		if _, err := Run(projectPath, "cherry-pick", commit); err != nil {
			return fmt.Errorf("cherry-pick %s: %w", commit[:8], err)
		}
	}
	return nil
}

// CurrentCommit returns the HEAD commit SHA for the given project path.
func CurrentCommit(projectPath string) (string, error) {
	out, err := Run(projectPath, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("current commit: %w", err)
	}
	return out, nil
}
