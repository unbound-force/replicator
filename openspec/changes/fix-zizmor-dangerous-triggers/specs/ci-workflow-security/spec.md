## ADDED Requirements

### Requirement: Pull-request workflow trust model

`docs/ci-workflow-security.md` MUST document the dependency-review workflow and the generated Fullsend shim. For each workflow, it MUST record its event, effective permissions, secret availability or forwarding, checkout behavior, immutable action or reusable-workflow reference, and trust boundary. The documentation SHALL identify whether untrusted pull-request code can execute with a privileged token or repository secret.

#### Scenario: Documented dependency-review workflow
- **GIVEN** the dependency-review workflow receives a pull-request event
- **WHEN** a maintainer reviews its workflow policy
- **THEN** the event, read-only default permissions, Dependabot guard, secret availability, absence of checkout, write-token scope, and trust boundary are documented

#### Scenario: Documented generated workflow exception
- **GIVEN** a generated external workflow cannot be redesigned locally without breaking its supported integration
- **WHEN** an exception is required for a zizmor finding
- **THEN** the documentation identifies the generated shim, its pinned reusable workflow, forwarded secrets, and the specific exception rationale

### Requirement: Narrow security-gate exceptions

A zizmor suppression MUST apply only to the `dangerous-triggers` finding in `fullsend.yaml`. The suppression MUST state its security rationale, cite issue #127 as the explicit human authorization record, and MUST NOT disable zizmor or unrelated workflow checks.

#### Scenario: Generated workflow requires an exception
- **GIVEN** an upstream-supported redesign is unavailable and issue #127 authorizes the exception
- **WHEN** the repository adds a zizmor suppression for the generated workflow
- **THEN** the suppression is limited to `fullsend.yaml` and `dangerous-triggers`, and preserves all other linter rules

### Requirement: Dependency-review pull-request execution

The dependency-review workflow MUST use the `pull_request` event for pull requests targeting `main`. It MUST retain read-only default permissions, Dependabot-only guards for comment and approval operations, full-SHA action pins, and no checkout of pull-request code. It MUST NOT execute untrusted pull-request code with a write-capable token or repository secret.

#### Scenario: Untrusted pull request triggers dependency review
- **GIVEN** a pull request to `main` is opened from an untrusted source
- **WHEN** the dependency-review workflow runs
- **THEN** the workflow executes under `pull_request`, does not check out pull-request code, and does not expose privileged credentials to untrusted code

#### Scenario: Policy regression is detected
- **GIVEN** the dependency-review workflow is changed to `pull_request_target`, adds checkout, or broadens default permissions
- **WHEN** workflow policy tests run
- **THEN** the tests fail with the violated security policy identified

#### Scenario: Reduced pull-request token limits a write operation
- **GIVEN** a same-repository, fork, or Dependabot pull request has restricted token permissions
- **WHEN** the dependency-review workflow cannot complete a comment or approval write
- **THEN** the workflow does not restore `pull_request_target` or broaden permissions, and the implementation records the safe failure or skip behavior for maintainers
