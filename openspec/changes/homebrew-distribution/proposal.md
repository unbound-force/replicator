## Why

Replicator is a single Go binary with zero runtime dependencies, but currently the only install path is `go install` or downloading a release tarball manually. The Unbound Force ecosystem uses `brew install unbound-force/tap/<tool>` as the primary distribution channel -- `uf`, `dewey`, and `gaze` are all distributed this way. Replicator needs the same path so `uf setup` can install it via Homebrew instead of requiring Node.js/npm.

## What Changes

### New Capabilities
- **Homebrew distribution**: `brew install unbound-force/tap/replicator` installs the binary
- **Automated releases**: Pushing a `v*` tag triggers GoReleaser to build cross-platform binaries and publish to GitHub Releases
- **macOS quarantine removal**: Post-install hook removes quarantine attribute so the binary runs without Gatekeeper warnings

### Modified Capabilities
- **GoReleaser config**: Replace the existing minimal `.goreleaser.yml` with a full v2 config matching the `unbound-force/unbound-force` pattern (cross-compilation, Homebrew cask, changelog grouping)
- **Makefile**: Add `release` (local GoReleaser test) and `install` (build + copy to GOPATH/bin) targets

## Impact

- `.goreleaser.yml` → `.goreleaser.yaml` (renamed, v2 format)
- `.github/workflows/release.yml` (new)
- `Makefile` (modified -- 2 new targets)
- `unbound-force/homebrew-tap` (receives auto-generated `Casks/replicator.rb` on release)

## Constitution Alignment

### I. Autonomous Collaboration

**Assessment**: PASS

Replicator is distributed as a standalone binary. The Homebrew cask installs it independently -- no runtime coupling to other heroes. The release pipeline uses artifact-based communication (GitHub Releases + checksums).

### II. Composability First

**Assessment**: PASS

The binary is independently installable via Homebrew, `go install`, or direct download. No mandatory dependencies on other tools. The cask does not declare dependencies on `dewey` or `uf` -- they're optional runtime peers, not install-time requirements.

### III. Observable Quality

**Assessment**: PASS

GoReleaser produces machine-parseable artifacts: checksums.txt (SHA-256), structured changelogs, versioned archives with consistent naming (`replicator_<version>_<os>_<arch>.tar.gz`).

### IV. Testability

**Assessment**: PASS

The release pipeline is testable locally via `goreleaser check` (config validation) and `goreleaser release --snapshot --clean` (dry run without publishing). No external services required for verification.
