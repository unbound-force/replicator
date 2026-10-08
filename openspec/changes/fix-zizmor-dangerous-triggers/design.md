## Context

Issue #127 reports zizmor `dangerous-triggers` findings for the dependency-review and Fullsend workflows. The dependency-review workflow currently uses `pull_request_target`, but it does not check out pull-request code and its privileged operations are limited to Dependabot paths. Fullsend is a generated shim that dispatches to a pinned external reusable workflow with privileged permissions and secrets; its remote workflow is the effective trust boundary.

The proposal's constitution alignment is PASS for all principles. This design preserves that alignment by recording independently reviewable workflow controls, retaining a standalone binary, producing machine-checked evidence, and relying on isolated static workflow tests.

## Goals / Non-Goals

### Goals
- Remove the avoidable privileged pull-request trigger from the repository-owned dependency-review workflow.
- Preserve existing least-privilege, no-checkout, pinned-action, Dependabot-only, and fail-closed protections.
- Establish a documented, auditable rationale for Fullsend's generated workflow and any narrow exception.
- Keep zizmor and actionlint active for all unrelated workflow findings.

### Non-Goals
- Redesign or fork Fullsend's upstream reusable workflow locally.
- Grant broader GitHub token permissions to preserve behavior on forked pull requests.
- Disable MegaLinter, zizmor, actionlint, or any other CI security gate.
- Change runtime behavior of the Replicator binary.

## Decisions

### Use `pull_request` for dependency review

The dependency-review workflow will use `pull_request` for `main` instead of `pull_request_target`. This removes the privileged base-repository execution context for a workflow driven by untrusted pull requests. The workflow will continue not to check out repository or pull-request code, and all static action references remain full-SHA pinned.

The implementation must verify whether GitHub's reduced token permissions affect Dependabot comments or approvals. If the operation cannot be safely performed for a token context, the workflow must fail closed or skip the write operation; it must not restore `pull_request_target` or broaden permissions solely to retain behavior.

### Treat Fullsend as an upstream-managed integration

`fullsend.yaml` is generated and delegates privileged dispatch to a pinned external reusable workflow. The local repository does not own the remote implementation, so the preferred resolution is an upstream-supported Fullsend hardening update. If one is not available, retain the generated workflow and use the smallest `dangerous-triggers` suppression scoped to `fullsend.yaml`. A suppression requires explicit human authorization before it changes `.mega-linter.yml`; issue #127 is the authorization record and must be cited with the rationale.

`docs/ci-workflow-security.md` will record each affected workflow's event, effective permissions, secret availability or forwarding, checkout behavior, immutable action or reusable-workflow reference, and trust boundary. This is maintainer documentation, not a new runtime contract or schema.

### Verify controls structurally

Extend the existing dependency-review workflow tests to reject `pull_request_target`, checkout actions, unpinned actions, broadened default permissions, and removal of Dependabot or fail-closed guards. This static test coverage, plus zizmor and actionlint, is the coverage strategy for the change. Validate the relevant same-repository, fork, and Dependabot token behavior without restoring privileged execution. Add local validation for any Fullsend suppression or trust-model document without attempting to simulate remote secrets or agent behavior.

This supports Observable Quality through reproducible lint and test output, and Testability through static local assertions. It avoids new runtime coupling, meeting Autonomous Collaboration and Composability First.

## Risks / Trade-offs

- `pull_request` may make write operations unavailable for forked or Dependabot pull requests. The implementation must explicitly validate these cases and prefer safe degradation over privilege escalation.
- An exception for Fullsend leaves an acknowledged external trust dependency. Its scope and rationale must be kept adjacent to the finding and revisited when Fullsend provides a supported hardening option.
- Static analysis cannot prove the behavior of the remote Fullsend workflow. The design limits this risk through SHA pinning, documented delegated trust, explicit permissions, and upstream ownership rather than claiming full local verification.
