## MODIFIED Requirements

### Requirement: Per-repo directory

`replicator init` MUST create `[repo]/.uf/replicator/cells.json` instead of `[repo]/.hive/cells.json`.

#### Scenario: Fresh init
- **GIVEN** a directory without `.uf/replicator/`
- **WHEN** `replicator init` is run
- **THEN** `.uf/replicator/cells.json` is created containing `[]`

#### Scenario: Hive sync
- **GIVEN** cells exist in the database
- **WHEN** `hive_sync` is called
- **THEN** cells are written to `.uf/replicator/cells.json` and `git add .uf/replicator/` is run

### Requirement: Per-repo log file

The MCP server MUST write logs to `[repo]/.uf/replicator/replicator.log` instead of `[repo]/.unbound-force/replicator.log`.

#### Scenario: Server startup
- **GIVEN** a project directory
- **WHEN** `replicator serve` starts
- **THEN** `.uf/replicator/replicator.log` is created (truncating any existing file)

### Requirement: Per-machine database

The database MUST be stored at `~/.config/uf/replicator/replicator.db` instead of `~/.config/swarm-tools/swarm.db`.

#### Scenario: Setup
- **GIVEN** a fresh machine
- **WHEN** `replicator setup` is run
- **THEN** `~/.config/uf/replicator/` is created with `replicator.db`

#### Scenario: Doctor check
- **GIVEN** the config directory exists
- **WHEN** `replicator doctor` is run
- **THEN** the config_dir check verifies `~/.config/uf/replicator/` exists
