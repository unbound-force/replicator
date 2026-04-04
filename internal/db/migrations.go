package db

// Schema migrations -- compatible with cyborg-swarm's libSQL schema.
// Each migration is idempotent (IF NOT EXISTS).

const migrationEvents = `
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    payload TEXT NOT NULL DEFAULT '{}',
    project_key TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    agent_name TEXT GENERATED ALWAYS AS (json_extract(payload, '$.agent_name')) STORED
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_project ON events(project_key);
CREATE INDEX IF NOT EXISTS idx_events_created ON events(created_at);
`

const migrationAgents = `
CREATE TABLE IF NOT EXISTS agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    project_path TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'worker',
    status TEXT NOT NULL DEFAULT 'active',
    task_description TEXT,
    registered_at TEXT NOT NULL DEFAULT (datetime('now')),
    last_seen_at TEXT NOT NULL DEFAULT (datetime('now')),
    metadata TEXT NOT NULL DEFAULT '{}'
);
`

const migrationCells = `
CREATE TABLE IF NOT EXISTS beads (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL DEFAULT 'task',
    status TEXT NOT NULL DEFAULT 'open',
    priority INTEGER NOT NULL DEFAULT 1,
    parent_id TEXT,
    project_key TEXT NOT NULL DEFAULT '',
    assigned_to TEXT,
    labels TEXT NOT NULL DEFAULT '[]',
    blocked_by TEXT NOT NULL DEFAULT '[]',
    metadata TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    closed_at TEXT,
    close_reason TEXT,
    FOREIGN KEY (parent_id) REFERENCES beads(id)
);

CREATE INDEX IF NOT EXISTS idx_beads_status ON beads(status);
CREATE INDEX IF NOT EXISTS idx_beads_type ON beads(type);
CREATE INDEX IF NOT EXISTS idx_beads_parent ON beads(parent_id);
CREATE INDEX IF NOT EXISTS idx_beads_project ON beads(project_key);
`

const migrationCellEvents = `
CREATE TABLE IF NOT EXISTS cell_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cell_id TEXT NOT NULL,
    type TEXT NOT NULL,
    payload TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (cell_id) REFERENCES beads(id)
);

CREATE INDEX IF NOT EXISTS idx_cell_events_cell ON cell_events(cell_id);
`
