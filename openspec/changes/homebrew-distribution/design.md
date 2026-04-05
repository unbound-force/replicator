## Context

The `unbound-force/unbound-force` repo has a proven GoReleaser + Homebrew cask pipeline that we follow exactly. The pattern: GoReleaser v2 config → GitHub Actions release workflow triggered by `v*` tags → cross-platform builds → Homebrew cask auto-published to `unbound-force/homebrew-tap`.

## Goals / Non-Goals

### Goals
- `brew install unbound-force/tap/replicator` installs the binary
- `replicator version` prints the correct semver after Homebrew install
- Cross-platform builds (darwin-arm64, darwin-amd64, linux-arm64, linux-amd64)
- macOS quarantine removal via post-install hook
- Automated release on `v*` tag push

### Non-Goals
- macOS code signing and notarization (deferred -- requires Apple Developer account + signing secrets)
- Windows Scoop/Chocolatey distribution
- Linux package managers (apt, yum)

## Decisions

**D1: Follow the unbound-force pattern exactly.** The `.goreleaser.yaml` structure, release workflow, cask template, and quarantine hook are copied from `unbound-force/unbound-force` with only the binary name and descriptions changed. This ensures consistency across the ecosystem and avoids inventing new patterns.

**D2: Rename `.goreleaser.yml` to `.goreleaser.yaml`.** GoReleaser v2 prefers `.yaml`. The existing file uses v2 syntax already but has the `.yml` extension.

**D3: No `dewey` cask dependency.** Unlike `unbound-force` (which declares `dewey` as a cask dependency), replicator works without Dewey -- it degrades gracefully with `DEWEY_UNAVAILABLE` errors. Dewey is an optional runtime peer.

**D4: `skip_upload: true` for the cask.** Same pattern as `unbound-force`: GoReleaser generates the cask file but does not push it directly. A future `sign-macos` job can patch darwin checksums with signed values before pushing to the tap. For now, a simple upload step pushes the generated cask directly.

**D5: Version via ldflags.** `main.Version` is set at build time via `-X main.Version={{.Tag}}`. This matches the existing `Makefile` pattern and the `cmd/replicator/main.go` `Version` variable.

## Risks / Trade-offs

**Risk: First-time cask tap setup.** The `unbound-force/homebrew-tap` repo must accept cask files at `Casks/replicator.rb`. If the tap structure doesn't support this, the release will fail on the cask push step. Mitigation: the tap already has `Casks/` directory from other tools.

**Trade-off: No code signing.** macOS users may see a Gatekeeper warning on first run if they don't use Homebrew (which applies its own quarantine removal). The post-install hook handles the Homebrew case; direct downloads require `xattr -dr com.apple.quarantine replicator` manually.
