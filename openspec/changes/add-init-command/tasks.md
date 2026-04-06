## 1. Implement init command

- [x] 1.1 Create `cmd/replicator/init.go` with `initCmd() *cobra.Command` (--path flag, default ".") and `runInit(targetDir string) error`
- [x] 1.2 In `runInit`: check if `.hive/` exists (print "already initialized", return nil), otherwise `os.MkdirAll` + `os.WriteFile("cells.json", "[]")`
- [x] 1.3 Register `initCmd()` in `cmd/replicator/main.go` via `root.AddCommand(initCmd())`

## 2. Tests

- [x] 2.1 Create `cmd/replicator/init_test.go` with tests: fresh dir creates `.hive/cells.json`, already-initialized is idempotent, custom --path works, cells.json content is `[]`

## 3. Documentation

- [x] 3.1 Add `replicator init` to the Usage section in `README.md`
- [x] 3.2 Add `make init` is NOT needed (clarify `init` is a CLI command, not a Makefile target) -- update AGENTS.md commands section to include `replicator init`

## 4. Verify

- [x] 4.1 Run `make check` -- all tests pass
- [x] 4.2 Run `./bin/replicator init` in a temp dir -- verify `.hive/cells.json` created
