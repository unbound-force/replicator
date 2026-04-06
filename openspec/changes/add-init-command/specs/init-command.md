## ADDED Requirements

### Requirement: Per-repo initialization

`replicator init` MUST create a `.hive/` directory with an empty `cells.json` file in the target directory.

#### Scenario: Fresh directory
- **GIVEN** a directory without `.hive/`
- **WHEN** `replicator init` is run
- **THEN** `.hive/cells.json` is created containing `[]`
- **AND** stdout prints `initialized .hive/`
- **AND** exit code is 0

#### Scenario: Already initialized
- **GIVEN** a directory with existing `.hive/`
- **WHEN** `replicator init` is run
- **THEN** no files are modified
- **AND** stdout prints `already initialized`
- **AND** exit code is 0

#### Scenario: Custom path
- **GIVEN** a valid directory at `/some/path`
- **WHEN** `replicator init --path /some/path` is run
- **THEN** `.hive/cells.json` is created at `/some/path/.hive/cells.json`

#### Scenario: Invalid path
- **GIVEN** a path that does not exist and cannot be created
- **WHEN** `replicator init --path /nonexistent/path` is run
- **THEN** an error is printed
- **AND** exit code is 1

### Requirement: No external dependencies

`replicator init` MUST NOT require the SQLite database, git, network access, or any other external service to function.
