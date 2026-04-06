// Package doctor runs health checks for the replicator environment.
//
// Checks verify that required dependencies (git, database, Dewey, config dir)
// are available and functional. Results include pass/fail/warn status and
// timing for each check.
package doctor

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/unbound-force/replicator/internal/config"
	"github.com/unbound-force/replicator/internal/db"
)

// CheckResult holds the outcome of a single health check.
type CheckResult struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"` // "pass", "fail", "warn"
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

// Run executes all health checks and returns the results.
// Individual check failures do not stop subsequent checks.
func Run(store *db.Store, cfg *config.Config) ([]CheckResult, error) {
	var results []CheckResult

	results = append(results, checkGit())
	results = append(results, checkDatabase(store))
	results = append(results, checkDewey(cfg.DeweyURL))
	results = append(results, checkConfigDir())

	return results, nil
}

// checkGit verifies that git is installed and returns its version.
func checkGit() CheckResult {
	start := time.Now()

	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	elapsed := time.Since(start)

	if err != nil {
		return CheckResult{
			Name:     "git",
			Status:   "fail",
			Message:  fmt.Sprintf("git not found: %v", err),
			Duration: elapsed,
		}
	}

	version := strings.TrimSpace(string(out))
	return CheckResult{
		Name:     "git",
		Status:   "pass",
		Message:  version,
		Duration: elapsed,
	}
}

// checkDatabase verifies the SQLite database is accessible.
func checkDatabase(store *db.Store) CheckResult {
	start := time.Now()

	err := store.DB.Ping()
	elapsed := time.Since(start)

	if err != nil {
		return CheckResult{
			Name:     "database",
			Status:   "fail",
			Message:  fmt.Sprintf("database ping failed: %v", err),
			Duration: elapsed,
		}
	}

	return CheckResult{
		Name:     "database",
		Status:   "pass",
		Message:  "SQLite database is accessible",
		Duration: elapsed,
	}
}

// checkDewey verifies the Dewey semantic search endpoint is reachable.
// Uses a simple HTTP GET -- a non-200 response is a warning, not a failure,
// because Dewey is optional for core operations.
func checkDewey(deweyURL string) CheckResult {
	start := time.Now()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(deweyURL)
	elapsed := time.Since(start)

	if err != nil {
		return CheckResult{
			Name:     "dewey",
			Status:   "warn",
			Message:  fmt.Sprintf("Dewey not reachable at %s: %v", deweyURL, err),
			Duration: elapsed,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CheckResult{
			Name:     "dewey",
			Status:   "warn",
			Message:  fmt.Sprintf("Dewey returned HTTP %d at %s", resp.StatusCode, deweyURL),
			Duration: elapsed,
		}
	}

	return CheckResult{
		Name:     "dewey",
		Status:   "pass",
		Message:  fmt.Sprintf("Dewey is reachable at %s", deweyURL),
		Duration: elapsed,
	}
}

// checkConfigDir verifies the config directory exists.
func checkConfigDir() CheckResult {
	start := time.Now()

	home, err := os.UserHomeDir()
	if err != nil {
		elapsed := time.Since(start)
		return CheckResult{
			Name:     "config_dir",
			Status:   "fail",
			Message:  fmt.Sprintf("cannot determine home directory: %v", err),
			Duration: elapsed,
		}
	}

	configDir := home + "/.config/uf/replicator"
	info, err := os.Stat(configDir)
	elapsed := time.Since(start)

	if os.IsNotExist(err) {
		return CheckResult{
			Name:     "config_dir",
			Status:   "fail",
			Message:  fmt.Sprintf("config directory does not exist: %s", configDir),
			Duration: elapsed,
		}
	}
	if err != nil {
		return CheckResult{
			Name:     "config_dir",
			Status:   "fail",
			Message:  fmt.Sprintf("cannot access config directory: %v", err),
			Duration: elapsed,
		}
	}
	if !info.IsDir() {
		return CheckResult{
			Name:     "config_dir",
			Status:   "fail",
			Message:  fmt.Sprintf("%s exists but is not a directory", configDir),
			Duration: elapsed,
		}
	}

	return CheckResult{
		Name:     "config_dir",
		Status:   "pass",
		Message:  fmt.Sprintf("config directory exists: %s", configDir),
		Duration: elapsed,
	}
}
