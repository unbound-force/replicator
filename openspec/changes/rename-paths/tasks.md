## 1. Go source -- per-machine paths

- [x] 1.1 In `internal/config/config.go`: change `defaultDatabasePath()` from `~/.config/swarm-tools` to `~/.config/uf/replicator` and `swarm.db` to `replicator.db`. Update package comment.
- [x] 1.2 In `internal/doctor/checks.go`: change `checkConfigDir()` path from `/.config/swarm-tools` to `/.config/uf/replicator`
- [x] 1.3 In `cmd/replicator/setup.go`: change `runSetup()` path from `/.config/swarm-tools` to `/.config/uf/replicator`, update printed message

## 2. Go source -- per-repo paths

- [x] 2.1 In `cmd/replicator/init.go`: replace all `.hive` with `.uf/replicator` (dir path, messages, flag description)
- [x] 2.2 In `internal/hive/sync.go`: replace `.hive` with `.uf/replicator` in `Sync()` function and `git add` command
- [x] 2.3 In `cmd/replicator/serve.go`: replace `.unbound-force` with `.uf/replicator` in `setupLogger()` (log dir path and comment)

## 3. Go tests

- [x] 3.1 In `cmd/replicator/init_test.go`: replace all `.hive` with `.uf/replicator` in assertions
- [x] 3.2 In `internal/hive/sync_test.go`: replace all `.hive` with `.uf/replicator` in assertions
- [x] 3.3 In `cmd/replicator/serve_test.go`: replace all `.unbound-force` with `.uf/replicator` in assertions and comments

## 4. Active documentation

- [x] 4.1 In `README.md`: replace `~/.config/swarm-tools` with `~/.config/uf/replicator`, `swarm.db` with `replicator.db`
- [x] 4.2 In `AGENTS.md`: replace all `~/.config/swarm-tools` with `~/.config/uf/replicator`, `swarm.db` with `replicator.db`, `.hive/` with `.uf/replicator/` where referencing cell state
- [x] 4.3 In `.specify/memory/constitution.md`: replace `swarm.db` with `replicator.db` if referenced

## 5. Spec artifacts (historical)

- [x] 5.1 In `specs/001-go-rewrite-phases/`: replace `~/.config/swarm-tools/swarm.db` with `~/.config/uf/replicator/replicator.db` across spec.md, plan.md, tasks.md, quickstart.md
- [x] 5.2 In `specs/002-charm-ux/`: replace `.unbound-force/replicator.log` with `.uf/replicator/replicator.log` across spec.md, tasks.md, plan.md, quickstart.md, research.md
- [x] 5.3 In `openspec/changes/add-init-command/`: replace `~/.config/swarm-tools` with `~/.config/uf/replicator` in design.md

## 6. Verify

- [x] 6.1 Run `make check` -- all tests pass
- [x] 6.2 Grep verify: zero remaining `swarm-tools` in Go source/tests/active docs (excluding LICENSE, README credit link)
- [x] 6.3 Grep verify: zero remaining `.hive/` in Go source/tests (old per-repo path)
- [x] 6.4 Grep verify: zero remaining `.unbound-force/replicator` in Go source/tests (old log path)
