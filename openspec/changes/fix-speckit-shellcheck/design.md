## Context

The standardized CI workflow enables `BASH_SHELLCHECK`. Six vendored Speckit scaffold scripts in `.specify/scripts/bash/` currently report SC2155, SC2221, SC2222, and SC1091, blocking both the lint job and release preflight. These scripts are upstream-managed but are versioned in this repository, so remediation must avoid changing their observable scaffold behavior.

The proposal assesses all four constitution principles as PASS. This design preserves independent script execution, adds no runtime dependency, retains machine-verifiable lint evidence, and uses isolated checks for regression verification.

## Goals / Non-Goals

### Goals
- Correct reported ShellCheck findings with minimal upstream-compatible script edits.
- Preserve command failure propagation when separating declaration and assignment.
- Make `common.sh` sourcing visible to ShellCheck without changing runtime path resolution.
- Retain ShellCheck for all repository shell code and demonstrate the affected findings are resolved.

### Non-Goals
- Updating, replacing, or otherwise re-vendoring Speckit.
- Disabling `BASH_SHELLCHECK`, suppressing a rule across the repository, or changing release quality gates.
- Altering the generated features or user-facing behavior of the scaffold scripts.
- Modifying unrelated MegaLinter, workflow, or release configuration.

## Decisions

### Prefer source fixes over exclusions

Implementation will first make behavior-preserving changes in the six affected scripts. Separating `local` declarations from command substitution retains the command status, satisfying ShellCheck while improving error observability. Sourcing diagnostics will be addressed with static hints or directives placed adjacent to the source statement only when the existing runtime calculation cannot be expressed statically.

This keeps Observable Quality intact: the existing lint gate remains authoritative, produces machine-readable results, and continues to inspect unrelated paths.

### Protect lint scope

`.mega-linter.yml` will continue to enable `BASH_SHELLCHECK`. A file- and rule-specific exception is a last resort, not a planned implementation. It requires documented proof that a local correction is unsafe and explicit human authorization before it can be added. This avoids weakening the project’s protected CI quality gate.

### Verify behavior and lint outcome separately

One focused shell smoke test will exercise a representative changed command-substitution failure path and assert its non-zero status. ShellCheck will cover all six affected scripts, and a small policy assertion will ensure `BASH_SHELLCHECK` remains enabled without a broad suppression. The implementation will also run the available ShellCheck or MegaLinter command. This follows the proposal’s Testability alignment by making the original failure and the repaired result reproducible in isolation.

## Risks / Trade-offs

- Upstream template updates could overwrite local changes. Keep edits narrow and use standard shell constructs so they are suitable for upstream contribution.
- ShellCheck directives can hide genuine future problems if applied broadly. Directives, if needed, must be rule- and location-specific with a reason.
- The full hosted MegaLinter environment may not be available locally. Focused local verification can establish the regression result, while the standardized CI job remains the final integration signal.
