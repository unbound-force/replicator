## ADDED Requirements

### Requirement: Homebrew postflight steps

The GoReleaser Homebrew cask configuration MUST emit the `postflight_steps`
DSL. The generated cask MUST use `on_macos do` to remove the quarantine
attribute from the staged `replicator` binary, and MUST preserve Ruby
`#{staged_path}` interpolation.

#### Scenario: Generated cask uses the current Homebrew DSL
- **GIVEN** the release configuration contains the Homebrew cask definition
- **WHEN** GoReleaser renders a snapshot cask
- **THEN** the cask MUST contain `postflight_steps do` and `on_macos do`
- **AND** it MUST contain `#{staged_path}/replicator`

#### Scenario: Generated cask reintroduces the legacy stanza
- **GIVEN** a GoReleaser render produces a cask containing `postflight do`
- **WHEN** the cask regression suite runs
- **THEN** the suite MUST fail before the change can merge

### Requirement: Rendered Homebrew cask validation

The `Build and Test` job MUST render the Homebrew cask with the pinned
GoReleaser version before publishing a release. The cask regression suite MUST
validate the rendered cask contains `postflight_steps do`, does not contain the
legacy `postflight do` stanza, and preserves `#{staged_path}` Ruby
interpolation.

The GoReleaser setup action MUST receive arguments that allow it to complete
before the separate snapshot-rendering command runs.

#### Scenario: Rendered cask uses the current Homebrew DSL
- **GIVEN** the release configuration contains the Homebrew cask definition
- **WHEN** the `Build and Test` job runs the GoReleaser snapshot render
- **THEN** the regression suite MUST validate the generated cask contains
  `postflight_steps do`
- **AND** it MUST validate the generated cask retains `#{staged_path}`

#### Scenario: Missing rendered cask
- **GIVEN** the GoReleaser snapshot command does not produce the expected cask
- **WHEN** the cask regression suite runs
- **THEN** the suite MUST fail rather than validating only the source
  configuration

#### Scenario: GoReleaser setup action completes
- **GIVEN** the `Build and Test` job installs its pinned GoReleaser version
- **WHEN** the setup action runs
- **THEN** it MUST receive a non-rendering argument and complete before the
  snapshot-rendering command runs

### Requirement: Configuration-level cask validation

The repository MUST retain dependency-free checks for the GoReleaser
`custom_block` source configuration so accidental DSL regressions fail before
the snapshot-rendering step.

#### Scenario: custom_block removes the current DSL
- **GIVEN** the GoReleaser `custom_block` is empty or removes
  `postflight_steps do`
- **WHEN** the cask regression suite runs
- **THEN** the suite MUST fail
