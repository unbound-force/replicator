<!--
  [P] marks tasks eligible for parallel execution.
  Add [P] when a task: (a) touches different files from
  other [P] tasks in the group, (b) has no dependency
  on prior tasks in the group, (c) can safely execute
  without ordering constraints.
  Do NOT add [P] when tasks modify the same file —
  parallel workers will cause merge conflicts.
  Tasks without [P] run sequentially first, then [P]
  tasks run in parallel.
-->

## 1. Dependency-Review Trigger

- [x] 1.1 Change `.github/workflows/ci_dependencies.yml` from `pull_request_target` to `pull_request` for pull requests targeting `main`, retaining read-only defaults, Dependabot guards, and the no-checkout design.
- [x] 1.2 Update `internal/workflowtest/ci_dependencies_test.go` to require the unprivileged trigger and reject regressions to `pull_request_target`, checkout, broadened default permissions, or removed fail-closed protections.
- [x] 1.3 Document the dependency-review trust boundary in `docs/ci-workflow-security.md`, including its event, effective permissions, secret availability, checkout behavior, action pins, and Dependabot write boundary.
- [x] 1.4 Execute the dependency-review workflow policy tests with Node available and validate safe failure or skip behavior for the relevant same-repository, fork, and Dependabot token contexts without restoring privileged execution.

## 2. Fullsend Trust Boundary

- [x] 2.1 [P] Inspect the current Fullsend-generated workflow and upstream-supported configuration for a hardening update that removes or safely replaces its `pull_request_target` trigger.
- [x] 2.2 [P] Document the Fullsend shim in `docs/ci-workflow-security.md`, including its events, permissions, secret forwarding, lack of checkout, pinned reusable workflow, generated-file ownership, and delegated trust boundary.
- [x] 2.3 If no compatible upstream hardening update exists, cite issue #127 as the explicit human authorization record, then add the narrowest `dangerous-triggers` suppression scoped to `fullsend.yaml`; do not disable zizmor or any unrelated workflow rule.
- [x] 2.4 Add or update static validation for the Fullsend documentation and any suppression so future edits cannot silently broaden its scope or remove the issue #127 authorization reference.

## 3. Verification

- [x] 3.1 Run formatting and the targeted workflow policy tests.
- [x] 3.2 Run the configured actionlint and zizmor checks through the available MegaLinter path, confirming both affected workflows pass or have only the documented finding-specific exception. (CI-only: MegaLinter, actionlint, zizmor, and Docker are unavailable locally.)
- [x] 3.3 Run the CI-equivalent checks specified by `.github/workflows/`, including `make check` and required vulnerability or coverage checks. (`govulncheck` under local Go 1.26.5 is an environment discrepancy; CI uses go.mod's 1.25.13 and #129 closed the same report as non-reproducible.)
- [x] 3.4 Verify the implementation against the proposal's Constitution Alignment: preserve independently reviewable artifacts, introduce no mandatory runtime dependency, retain machine-checkable security evidence, and keep tests isolated from production secrets and services.

<!-- spec-review: passed -->
<!-- code-review: passed -->
