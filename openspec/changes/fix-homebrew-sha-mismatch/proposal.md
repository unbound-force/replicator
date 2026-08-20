## Why

`brew install unbound-force/tap/replicator` fails on v0.4.0 with a SHA-256 mismatch. The Homebrew cask in `unbound-force/homebrew-tap` contains the pre-signing SHA for the darwin_arm64 archive, but the `sign-macos` job replaced the release asset with a signed binary that has a different SHA.

The root cause is a sequencing problem in `.github/workflows/release.yml`: the `sign-macos` job re-uploads the signed archive (changing its SHA) and then patches the cask in the same job. This creates two unguarded failure modes:

1. **Mid-job failure**: If `sign-macos` fails after re-uploading the signed archive (line 176) but before pushing the patched cask to the tap (line 234), the release has a new SHA but the cask retains the old one.
2. **Job skip**: When signing secrets are absent, `sign-macos` is skipped entirely. Because GoReleaser uses `skip_upload: true`, the cask is never pushed to the tap at all.

This blocks all macOS Homebrew users from installing replicator.

Fixes: https://github.com/unbound-force/replicator/issues/81

## What Changes

Extract Homebrew tap publishing from the `sign-macos` job into a dedicated `publish-cask` job in `.github/workflows/release.yml`. The new job computes the SHA from the actual final release artifact (signed or unsigned), patches the cask, and pushes to the Homebrew tap.

## Capabilities

### New Capabilities
- `publish-cask job`: Dedicated workflow job that publishes the Homebrew cask after all upstream jobs complete, computing SHA from the final downloaded artifact.

### Modified Capabilities
- `sign-macos job`: Simplified to only handle signing, notarization, archive re-upload, and checksum update. No longer responsible for Homebrew tap publishing.

### Removed Capabilities
- None

## Impact

- **File**: `.github/workflows/release.yml` -- restructure job graph to add `publish-cask` job, remove tap-publishing logic from `sign-macos`
- **Users**: Homebrew install will work correctly after the next release (v0.4.1 or v0.5.0)
- **Pipeline**: Adds one lightweight job (~2-3 min) to the release pipeline; removes equivalent logic from `sign-macos`
- **No Go source code changes**

## Constitution Alignment

Assessed against the Replicator project constitution (`.specify/memory/constitution.md`).

### I. Autonomous Collaboration

**Assessment**: N/A

This change modifies CI workflow structure only. No MCP tools, tool outputs, or inter-agent communication are affected. Artifact-based communication patterns remain unchanged.

### II. Composability First

**Assessment**: PASS

This fix directly restores compliance with Principle II. Replicator MUST be independently installable. The broken Homebrew install violates this principle. The fix ensures the Homebrew distribution channel works regardless of whether macOS signing is available -- the `publish-cask` job handles both the signed and unsigned paths.

### III. Observable Quality

**Assessment**: PASS

The `publish-cask` job adds a verification step that confirms the patched SHA appears in the cask file before pushing to the tap. This improves observability over the current approach where SHA patching failure is only detectable as a downstream `brew install` failure.

### IV. Testability

**Assessment**: PASS

The fix can be validated via `goreleaser check` (config validation) and by cutting a new release. The verification step in the `publish-cask` job provides a self-test that catches SHA patching failures at release time rather than at install time.
