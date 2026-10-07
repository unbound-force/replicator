## Why

MegaLinter's ShellCheck audit reports SC2155, SC2221, SC2222, and SC1091 findings in the vendored Speckit scripts under `.specify/scripts/bash/`. The findings block the standardized linter check and release preflight, preventing otherwise releasable changes from progressing.

The affected files are upstream-managed scaffolding, so the remediation must retain their generated-script behavior and remain suitable for upstream adoption. Disabling ShellCheck globally or suppressing rules outside the affected path would weaken a protected quality gate.

## What Changes

- Classify and resolve the reported ShellCheck findings in the affected Speckit scaffold scripts with behavior-preserving, upstream-compatible changes.
- Separate declarations from command substitutions where required so command failures remain observable.
- Make sourced helper paths analyzable by ShellCheck, using localized directives only where static analysis cannot safely resolve an existing dynamic source.
- Add regression coverage that verifies the repaired script paths and lint configuration continue to detect unrelated ShellCheck findings.
- Permit a documented ShellCheck exception for only the demonstrated file and rule only if an upstream-compatible correction is demonstrated unsafe and the exception receives explicit authorization.

## Capabilities

### New Capabilities
- `speckit-scaffold-lint-compliance`: Vendored Speckit scaffold scripts satisfy the repository's ShellCheck policy without reducing lint coverage for other shell code.

### Modified Capabilities
- None.

### Removed Capabilities
- None.

## Impact

- Affected files: `.specify/scripts/bash/common.sh`, `check-prerequisites.sh`, `create-new-feature.sh`, `setup-plan.sh`, `setup-tasks.sh`, and `update-agent-context.sh`.
- The implementation will add one focused regression smoke test for changed command-substitution failure propagation and documentation only for an approved exception.
- `.mega-linter.yml` must retain `BASH_SHELLCHECK` and must not gain a broad exclusion or repository-wide rule suppression.
- The standardized CI lint job and release preflight should pass once the affected findings are resolved.

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: PASS

The change affects only repository-managed scaffold artifacts and CI validation. It introduces no runtime coupling or coordination dependency, and the scripts remain independently invocable.

### II. Composability First

**Assessment**: PASS

The standalone binary, database schema, and optional integrations are unchanged. The remediation does not add a dependency or require an external service.

### III. Observable Quality

**Assessment**: PASS

ShellCheck and MegaLinter provide reproducible machine-readable evidence for the fixed findings. The change explicitly preserves the standardized lint gate rather than masking its results.

### IV. Testability

**Assessment**: PASS

The implementation will add isolated regression checks for the affected script behavior and validate the lint outcome locally where tooling is available. No test requires a live external service.
