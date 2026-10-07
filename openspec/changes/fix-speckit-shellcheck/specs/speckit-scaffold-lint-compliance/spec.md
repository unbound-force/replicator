## ADDED Requirements

### Requirement: Vendored Speckit scripts pass ShellCheck

The repository MUST resolve the reported SC2155, SC2221, SC2222, and SC1091 findings in `.specify/scripts/bash/` using behavior-preserving changes compatible with the upstream Speckit scaffold. Command substitutions SHALL be assigned separately from variable declarations when a combined declaration would hide a command failure.

#### Scenario: Command substitution preserves failure status
- **GIVEN** an affected scaffold script assigns the result of a command substitution
- **WHEN** the substituted command fails
- **THEN** the script MUST retain the command's failure status instead of masking it through a declaration assignment

#### Scenario: Helper sourcing is analyzable
- **GIVEN** an affected scaffold script sources `common.sh`
- **WHEN** ShellCheck evaluates the script
- **THEN** the source reference MUST be statically analyzable or have a localized, documented directive explaining why static resolution is unsafe

### Requirement: ShellCheck coverage remains protected

The repository MUST retain `BASH_SHELLCHECK` in MegaLinter and MUST NOT add a global ShellCheck disablement, repository-wide rule suppression, or exception broader than the demonstrated file and ShellCheck rule to resolve this issue.

#### Scenario: Unrelated shell code remains linted
- **GIVEN** a shell file outside `.specify/scripts/bash/`
- **WHEN** the standardized linter job runs
- **THEN** ShellCheck MUST continue to evaluate that file under the existing policy

### Requirement: Exceptions require explicit justification

If an upstream-compatible script correction is demonstrated unsafe, an exception MAY suppress only the demonstrated ShellCheck rule in the affected file. The exception MUST document the incompatibility, the exact file and rule scope, and the explicit authorization for the exception.

#### Scenario: Scoped exception is documented
- **GIVEN** the implementation cannot safely correct an upstream-managed finding locally
- **WHEN** a ShellCheck exclusion is proposed
- **THEN** the change MUST document why the correction is unsafe, limit the exception to the affected file and rule, and record explicit authorization before applying it

### Requirement: Regression verification covers the reported findings

The implementation MUST run ShellCheck over all six affected scripts, MUST provide one isolated executable smoke test for a representative changed command-substitution failure path, and MUST verify that `BASH_SHELLCHECK` remains enabled without a broad suppression. The implementation MUST demonstrate that the standardized lint check can complete without the reported findings.

#### Scenario: Repaired scripts and policy are verified
- **GIVEN** the remediation is applied
- **WHEN** the focused smoke test, the policy assertion, and the applicable linter command run
- **THEN** the smoke test MUST preserve non-zero command failure status, the policy assertion MUST confirm ShellCheck remains enabled, and the linter MUST pass without the reported findings without requiring a live external service

## MODIFIED Requirements

None.

## REMOVED Requirements

None.
