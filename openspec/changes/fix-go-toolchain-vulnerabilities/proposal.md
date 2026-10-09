## Why

`govulncheck ./...` reports four reachable vulnerabilities from issue #131 in the Go 1.26.5 standard library used for local builds. Although Go 1.26.6 fixes those advisories, the current vulnerability database reports nine additional reachable standard-library vulnerabilities in Go 1.26.6. Go 1.26.9 fixes all thirteen findings. The project needs one explicit minimum toolchain version so development, CI, and release builds produce binaries without the known vulnerable code.

This change addresses GitHub issue #131. It also prevents the documented source-build requirements from lagging behind the version selected by `go.mod`.

## What Changes

- Set the module's Go version to 1.26.9 so CI and the reusable release workflow select the patched toolchain.
- Update current developer and source-build documentation to require Go 1.26.9 or later.
- Keep the existing `govulncheck`, test, coverage, parity, build, and release gates intact and use them to verify the upgrade.
- Leave historical specifications, changelog entries, action pins, and protected quality thresholds unchanged.

## Capabilities

### New Capabilities
- `go-toolchain-security`: Define the minimum patched Go toolchain used by development, CI, and release builds, with vulnerability and regression verification requirements.

### Modified Capabilities
- None.

### Removed Capabilities
- None.

## Impact

The implementation affects `go.mod` and current toolchain references in `README.md`, `CONTRIBUTING.md`, and `AGENTS.md`. CI and release workflows already read their Go version from `go.mod`, so they should not need structural changes.

The patch changes the compiler and standard library used across all packages and release targets. Verification therefore includes the full local CI-equivalent checks, coverage ratchets, `govulncheck ./...`, and a GoReleaser dry run.

Changing the source-build prerequisite affects public documentation. Repository policy requires website documentation tracking before the implementation PR is merged.

## Constitution Alignment

Assessed against the Unbound Force org constitution.

### I. Autonomous Collaboration

**Assessment**: PASS

The change preserves the existing artifact-based workflows and MCP response contracts. It introduces no runtime coupling or new coordination path.

### II. Composability First

**Assessment**: PASS

Replicator remains a standalone Go binary. Raising the build toolchain to a patched Go release adds no runtime service or mandatory external integration.

### III. Observable Quality

**Assessment**: PASS

The existing machine-readable MCP outputs remain unchanged. Automated vulnerability scanning, tests, coverage ratchets, builds, and release checks provide reproducible evidence for the toolchain update.

### IV. Testability

**Assessment**: PASS

The change uses the existing isolated test suite and CI-equivalent checks. Verification requires no live external service, and no coverage or test gate is weakened.
<!-- scaffolded by uf v0.17.0 -->
