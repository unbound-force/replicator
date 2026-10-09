## Context

The module currently declares Go 1.25.13, while issue #131 records four reachable vulnerabilities observed with Go 1.26.5. Go 1.26.6 fixes those advisories, but a current scan reports nine newer reachable standard-library vulnerabilities in Go 1.26.6. Go 1.26.9 fixes all thirteen findings. The primary CI workflow uses `actions/setup-go` with `go-version-file: go.mod`, and the reusable release workflow follows the same pattern. Current source-build documentation still advertises Go 1.25 or later.

The proposal found this change aligned with all four constitution principles. It preserves artifact-based collaboration and standalone operation, changes no MCP output, and relies on automated local verification.

## Goals / Non-Goals

### Goals

- Establish Go 1.26.9 as the minimum project toolchain.
- Ensure CI and release builds inherit the patched version from one canonical declaration.
- Keep current source-build and contributor documentation consistent with the module.
- Verify that the standard-library findings are gone and every existing quality gate still passes.

### Non-Goals

- Change application behavior, package APIs, dependencies, or database schemas.
- Restructure CI or release workflows that already read `go.mod`.
- Change action references, `govulncheck` pinning, coverage thresholds, test flags, or release policy.
- Rewrite historical specifications, scan records, or changelog entries.
- Add a permanent test that hardcodes one Go patch version.

## Decisions

### Use `go.mod` as the single source of truth

Set the `go` directive to `1.26.9`. Both the primary CI workflow and the reusable release workflow already read this file, so a workflow-specific version would duplicate configuration and create another drift path. This decision supports Composability First by keeping local and hosted builds on the standard Go module mechanism.

Do not add a separate `toolchain` directive. The `go` directive supplies the required minimum version and is already the input consumed by the build workflows.

### Update current documentation, not historical records

Change the Go badge and source-build prerequisite in `README.md`, the contributor prerequisite in `CONTRIBUTING.md`, and current toolchain references in `AGENTS.md`. Leave older plans, archived change records, scan evidence, and changelog entries untouched because they describe the state at a specific time.

### Verify through existing observable gates

Run `govulncheck ./...` with the patched toolchain. Then run the repository's CI-equivalent vet, race test, coverage ratchet, build, parity, and release checks. These commands provide direct evidence for Observable Quality and Testability without adding code solely to assert a configuration literal.

The implementation must derive the final command set from `.github/workflows/` and preserve all flags and thresholds. No gate value may be relaxed to make the upgrade pass.

### Track the user-facing prerequisite change

The minimum source-build version appears in public documentation. Before the implementation PR is merged, prepare the required `unbound-force/website` issue content and have a maintainer file it. This planning change does not perform GitHub mutations.

## Risks / Trade-offs

- Go 1.26.9 may expose compiler, vet, test, or coverage differences across the entire module. Running the full gate set catches these regressions before release.
- Local commands use the developer's ambient Go installation. The documented prerequisite and `go.mod` directive make the requirement explicit, but they do not install Go outside CI.
- A later compatible patch may be selected by toolchain resolution. The security requirement is a minimum of 1.26.9, not a ban on later patched releases.
- The reusable release workflow lives in another repository. Its current contract reads `go.mod`; verification should confirm that behavior remains true at the pinned workflow revision rather than changing the pin.
<!-- scaffolded by uf v0.17.0 -->
