## ADDED Requirements

### Requirement: Auto-generated tool reference

`replicator docs` MUST output a markdown document listing all registered MCP tools with their names, descriptions, and argument schemas.

#### Scenario: Generate tool reference to stdout
- **GIVEN** the tool registry has 53 tools registered
- **WHEN** `replicator docs` is run
- **THEN** stdout contains a markdown document with all 53 tools grouped by category

#### Scenario: Generate tool reference to file
- **GIVEN** a valid output path
- **WHEN** `replicator docs --output docs/tools.md` is run
- **THEN** the file is created with the tool reference content

#### Scenario: Tool count matches registry
- **GIVEN** a new tool is added to the registry
- **WHEN** `replicator docs` is run after rebuild
- **THEN** the new tool appears in the output without any manual doc updates

### Requirement: Accurate README

The README MUST accurately reflect the current state of the codebase: tool count, CLI commands, architecture, install methods, and status.

#### Scenario: Developer reads README
- **GIVEN** a developer visits the GitHub repo
- **WHEN** they read the README
- **THEN** the status section reflects all 5 completed phases, the tool count is 53, and all 8 CLI commands are documented

### Requirement: Contributing guide

A CONTRIBUTING.md MUST exist explaining development setup, testing conventions, and PR workflow for human contributors.

#### Scenario: New contributor setup
- **GIVEN** a developer clones the repo
- **WHEN** they read CONTRIBUTING.md
- **THEN** they can build, test, and submit a PR by following the documented steps

## MODIFIED Requirements

### Requirement: CLI command list

The CLI MUST include a `docs` command alongside the existing `init`, `setup`, `serve`, `cells`, `doctor`, `stats`, `query`, and `version` commands.
