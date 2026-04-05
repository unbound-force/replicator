## 1. GoReleaser Configuration

- [x] 1.1 Delete `.goreleaser.yml` and create `.goreleaser.yaml` with GoReleaser v2 config: single build entry for `cmd/replicator`, `CGO_ENABLED=0`, darwin/linux x amd64/arm64, ldflags for `main.Version`, `main.commit`, `main.date`
- [x] 1.2 Add `archives` section with `tar.gz` default, `zip` for windows, name template `replicator_{{ .Version }}_{{ .Os }}_{{ .Arch }}`
- [x] 1.3 Add `checksum` section with `checksums.txt`
- [x] 1.4 Add `changelog` section with conventional commit grouping (Features, Bug Fixes, Documentation, Others) and exclude `chore:` commits
- [x] 1.5 Add `homebrew_casks` section: name `replicator`, description, homepage, `directory: Casks`, `skip_upload: true`, post-install quarantine removal hook (`xattr -dr com.apple.quarantine`), repository `unbound-force/homebrew-tap`
- [x] 1.6 Run `goreleaser check` to validate the config (if goreleaser is installed, otherwise validate YAML syntax)

## 2. Release Workflow

- [x] 2.1 Create `.github/workflows/release.yml` triggered on `push: tags: ['v*']` with `permissions: contents: write`
- [x] 2.2 Add steps: checkout (fetch-depth 0), setup-go (go-version-file: go.mod), run goreleaser-action v7 with `release --clean`
- [x] 2.3 Add step to upload generated cask file to the GitHub Release as an artifact: `gh release upload "${GITHUB_REF_NAME}" dist/homebrew/Casks/replicator.rb --clobber`
- [x] 2.4 Add `GITHUB_TOKEN` env var from secrets for both GoReleaser and cask upload steps

## 3. Makefile Updates

- [x] 3.1 Add `release` target: `goreleaser release --snapshot --clean` for local dry-run testing
- [x] 3.2 Add `install` target: `go install $(LDFLAGS) ./cmd/replicator` to install to GOPATH/bin

## 4. Version Ldflags

- [x] 4.1 Add `commit` and `date` variables to `cmd/replicator/main.go` (alongside existing `Version`)
- [x] 4.2 Update `Makefile` `LDFLAGS` to include `-X main.commit=$(COMMIT) -X main.date=$(DATE)` variables
- [x] 4.3 Update `versionCmd` to display version, commit, and date when available

## 5. Verify

- [x] 5.1 Run `make build` -- binary builds with version info
- [x] 5.2 Run `make check` -- all existing tests pass (no regressions)
- [x] 5.3 Run `./bin/replicator version` -- displays version, commit, date
