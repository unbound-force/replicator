## ADDED Requirements

### Requirement: Homebrew installation

The replicator binary MUST be installable via `brew install unbound-force/tap/replicator`.

#### Scenario: Fresh install via Homebrew
- **GIVEN** the user has Homebrew installed and the `unbound-force/tap` tapped
- **WHEN** the user runs `brew install unbound-force/tap/replicator`
- **THEN** the `replicator` binary is installed and available in PATH

#### Scenario: Version after Homebrew install
- **GIVEN** replicator is installed via Homebrew
- **WHEN** the user runs `replicator version`
- **THEN** the output shows the correct semver tag (e.g., `replicator v0.1.0`)

### Requirement: Automated release pipeline

Pushing a `v*` tag to the repository MUST trigger an automated release that produces cross-platform binaries and a Homebrew cask.

#### Scenario: Tag-triggered release
- **GIVEN** all CI tests pass on main
- **WHEN** a `v*` tag is pushed (e.g., `v0.1.0`)
- **THEN** GoReleaser builds binaries for darwin-arm64, darwin-amd64, linux-arm64, linux-amd64

#### Scenario: Release artifacts
- **GIVEN** GoReleaser completes successfully
- **WHEN** the GitHub Release page is inspected
- **THEN** it contains archives, checksums.txt, and a generated cask file

### Requirement: macOS quarantine removal

The Homebrew cask MUST remove the macOS quarantine attribute after installation.

#### Scenario: No Gatekeeper warning after Homebrew install
- **GIVEN** replicator is installed via Homebrew on macOS
- **WHEN** the user runs `replicator version`
- **THEN** no Gatekeeper warning dialog appears

### Requirement: Local release testing

The Makefile MUST support local release testing without publishing.

#### Scenario: Dry run release
- **GIVEN** the developer has GoReleaser installed
- **WHEN** the developer runs `make release`
- **THEN** GoReleaser runs in snapshot mode and produces local artifacts in `dist/`
