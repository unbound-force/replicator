<!--
  [P] marks tasks eligible for parallel execution.
  Add [P] when a task: (a) touches different files from
  other [P] tasks in the group, (b) has no dependency
  on prior tasks in the group, (c) can safely execute
  without ordering constraints.
  Do NOT add [P] when tasks modify the same file --
  parallel workers will cause merge conflicts.
  Tasks without [P] run sequentially first, then [P]
  tasks run in parallel.
-->

## 1. Establish The Patched Toolchain

- [x] 1.1 Record the pre-change `go version` and `govulncheck ./...` results, confirming the reachable standard-library findings from issue #131 before changing the toolchain.
- [ ] 1.2 Update the `go` directive in `go.mod` from 1.25.13 to 1.26.9 without adding a separate `toolchain` directive or changing dependencies.
- [x] 1.3 Confirm `.github/workflows/ci.yml` still selects Go through `go-version-file: go.mod` and that the pinned reusable release workflow has the same contract; do not add a duplicate workflow version.

## 2. Update Current Documentation

- [ ] 2.1 [P] Update the Go badge and source-build prerequisite in `README.md` to Go 1.26.9 or later.
- [ ] 2.2 [P] Update the Go prerequisite in `CONTRIBUTING.md` to Go 1.26.9 or later.
- [ ] 2.3 [P] Update current language, toolchain, and active-technology references in `AGENTS.md` to Go 1.26.9 or later while preserving historical spec and changelog records.

## 3. Verify Security And Quality Gates

- [ ] 3.1 Install the CI-pinned `govulncheck` version and run `govulncheck ./...` with Go 1.26.9 or later; confirm it reports no reachable standard-library vulnerabilities.
- [ ] 3.2 Run `.github/scripts/patch-homebrew-cask_test.sh`, `make vet`, `make check-coverage`, and `make build` without changing test flags, coverage thresholds, or workflow configuration.
- [ ] 3.3 Run `go test -tags parity ./test/parity/... -count=1 -race` and confirm all response-shape fixtures still pass.
- [ ] 3.4 Run `make release` and confirm the GoReleaser snapshot builds every configured target successfully.

## 4. Complete Governance Checks

- [ ] 4.1 Review `README.md`, `CONTRIBUTING.md`, `AGENTS.md`, GoDoc, and the OpenSpec artifacts for documentation completeness; confirm no GoDoc change is needed because no exported API changes.
- [ ] 4.2 Prepare the title and body for the required `unbound-force/website` documentation issue, and verify a maintainer files it before the implementation PR is merged.
- [ ] 4.3 Verify the completed change still passes all four constitution assessments from `proposal.md`: artifact-based collaboration is unchanged, Replicator remains standalone, quality evidence is reproducible, and tests require no external service.
<!-- scaffolded by uf v0.17.0 -->
<!-- spec-review: passed -->
