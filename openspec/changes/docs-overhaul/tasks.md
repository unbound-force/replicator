## 1. `replicator docs` command

- [x] 1.1 Create `cmd/replicator/docs.go` with `docsCmd() *cobra.Command` and `runDocs(outputPath string) error`. Wire in-memory store + full tool registration (hive, swarmmail, swarm, memory). Iterate `registry.List()`, group by name prefix (`hive_` / `swarmmail_` / `swarm_` / `hivemind_`), emit markdown with `### tool_name`, description, and fenced JSON `inputSchema` per tool.
- [x] 1.2 Add `--output` flag (default stdout). When set, write to file instead of stdout.
- [x] 1.3 Register `docsCmd()` in `cmd/replicator/main.go`
- [x] 1.4 Create `cmd/replicator/docs_test.go`: verify output contains all 53 tool names, verify category headers exist (Hive, Swarm Mail, Swarm, Memory), verify output is valid markdown

## 2. Tool reference document

- [x] 2.1 Run `replicator docs` and write output to `docs/tools.md`
- [x] 2.2 Add hand-written usage examples (request/response JSON) for 4 key tools: `hive_create`, `swarmmail_send`, `swarm_decompose`, `hivemind_store`

## 3. README rewrite

- [x] 3.1 Rewrite `README.md` with: project description, badges (CI workflow, Go version, license), accurate status (all phases complete, 53 tools, 190+ tests, 15MB binary), install section (Homebrew + `go install` + binary download), full usage section (all 9 CLI commands including `docs`), environment variables table, MCP client config example (`opencode.json` snippet), mermaid architecture diagram (flowchart: Agent → MCP Server → Registry → Tool → Domain → SQLite), complete package tree (all 15+ packages), credits

## 4. Supporting docs

- [x] 4.1 Create `CONTRIBUTING.md`: prerequisites (Go 1.25+, git), dev setup (`git clone` + `make check`), testing conventions (stdlib only, `db.OpenMemory()`, `t.TempDir()`, `httptest.NewServer`), PR workflow (branch naming, conventional commits, review council), spec-first development requirement (link to AGENTS.md)
- [x] 4.2 Create `CHANGELOG.md` with retroactive entries: v0.1.0 (Phase 0 scaffold, 4 tools, MCP server, SQLite, CI), v0.2.0 (Phases 1-5 complete, 53 tools, macOS signing, Homebrew distribution, Dewey proxy, parity testing, init command)

## 5. Final updates

- [x] 5.1 Update AGENTS.md CLI commands table to include `replicator docs`
- [x] 5.2 Run `make check` -- all tests pass
- [x] 5.3 Run `./bin/replicator docs | head -20` -- verify tool reference output
