## Why

MegaLinter's zizmor audit blocks CI and release preflight because two workflows use `pull_request_target`. The dependency-review workflow can run safely on the unprivileged `pull_request` event. The Fullsend-generated workflow needs a documented trust model and, if no safe upstream redesign is available, the narrowest justified exception rather than a broad security-gate bypass.

## What Changes

- Replace `pull_request_target` with `pull_request` in the dependency-review workflow and preserve its Dependabot-only, no-checkout design.
- Update dependency-review workflow policy tests to enforce the unprivileged trigger and existing token, action-pinning, and fail-closed controls.
- Document Fullsend's generated workflow trust boundary, including its event guards, no-checkout constraint, pinned reusable workflow, and privileged dispatch requirements.
- Use an upstream-supported Fullsend hardening change when available; otherwise add only a finding-specific, documented zizmor suppression for the generated workflow.
- Verify zizmor, actionlint, workflow tests, and CI-equivalent checks without reducing unrelated security controls.

## Capabilities

### New Capabilities
- `ci-workflow-trust-models`: Documents the event, permissions, secrets, checkout behavior, and trust boundary for workflows that receive pull-request events.

### Modified Capabilities
- None.

### Removed Capabilities
- None.

## Impact

- Affected workflows: `.github/workflows/ci_dependencies.yml` and `.github/workflows/fullsend.yaml`.
- Affected tests: `internal/workflowtest/ci_dependencies_test.go`; retain its static assertions for the trigger, permissions, action pins, and no-checkout boundary.
- Affected documentation: `docs/ci-workflow-security.md`, containing a concise trust-boundary record for the two affected workflows.
- Affected configuration: `.github/zizmor.yml` for the specific, justified zizmor suppression.
- The dependency-review comment and approval behavior for forked or Dependabot pull requests must be checked under the reduced token permissions of `pull_request`.

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: PASS

The workflows remain independently triggered and their trust boundaries are recorded in repository artifacts. No runtime coupling between MCP tools or heroes is introduced.

### II. Composability First

**Assessment**: PASS

The change adds no runtime dependency to the Replicator binary. The optional Fullsend integration remains explicitly bounded and its generated workflow is handled through its supported upstream mechanism.

### III. Observable Quality

**Assessment**: PASS

zizmor, actionlint, and structural workflow tests provide reproducible evidence that the security controls remain active. The coverage strategy is limited to static assertions for trigger selection, permissions, action pins, and no-checkout behavior; any exception is limited to one finding with recorded rationale.

### IV. Testability

**Assessment**: PASS

The dependency-review workflow policy remains verified through isolated repository tests. Static workflow analysis and CI-equivalent checks verify observable configuration effects without requiring production secrets or external services.
