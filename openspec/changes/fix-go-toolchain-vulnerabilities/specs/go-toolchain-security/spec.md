## ADDED Requirements

### Requirement: Patched Go toolchain baseline

The project MUST declare Go 1.26.6 as its minimum Go version in `go.mod`. CI and release builds MUST select their Go toolchain from that declaration rather than maintain a separate version value.

#### Scenario: CI selects the patched toolchain

- **GIVEN** `go.mod` declares Go 1.26.6
- **WHEN** the `Build and Test` workflow configures Go from `go.mod`
- **THEN** the workflow uses Go 1.26.6 or a compatible later patch release

#### Scenario: Release selects the patched toolchain

- **GIVEN** `go.mod` declares Go 1.26.6
- **WHEN** the reusable GoReleaser workflow configures the release environment from `go.mod`
- **THEN** release binaries are built with Go 1.26.6 or a compatible later patch release

### Requirement: Current toolchain documentation

Current source-build and contributor documentation MUST identify Go 1.26.6 or later as the supported toolchain baseline. Historical specifications and changelog records MUST retain the versions that applied when those records were written.

#### Scenario: Developer checks the build prerequisite

- **GIVEN** a developer reads the current project or contributor documentation
- **WHEN** the developer checks the Go prerequisite
- **THEN** the documentation states Go 1.26.6 or later

#### Scenario: Historical records remain accurate

- **GIVEN** a historical specification or changelog entry records an earlier Go version
- **WHEN** current toolchain documentation is updated
- **THEN** the historical version record remains unchanged

### Requirement: Toolchain vulnerability verification

The patched toolchain MUST pass `govulncheck ./...` without any reachable standard-library vulnerability. The existing vet, race test, coverage ratchet, build, parity, and release checks MUST remain enabled. The implementation MUST pass those checks without weakening their configuration or thresholds.

#### Scenario: Reachable standard-library vulnerabilities are resolved

- **GIVEN** the project is built and analyzed with Go 1.26.6 or later
- **WHEN** `govulncheck ./...` scans all packages
- **THEN** the scan reports no reachable standard-library vulnerability

#### Scenario: Toolchain update preserves quality gates

- **GIVEN** the module and current documentation use the new toolchain baseline
- **WHEN** the CI-equivalent and release verification commands run
- **THEN** every existing gate passes with its original flags and thresholds

## MODIFIED Requirements

None. This change adds a toolchain security capability without changing an existing requirement.

## REMOVED Requirements

None. No existing requirements are removed by this change.
<!-- scaffolded by uf v0.17.0 -->
