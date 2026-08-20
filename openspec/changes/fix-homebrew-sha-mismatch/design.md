## Context

The release pipeline in `.github/workflows/release.yml` has a sequencing bug where the `sign-macos` job both replaces the darwin_arm64 archive with a signed version (changing its SHA) and publishes the Homebrew cask in the same job. This creates a partial-failure window and a skip-failure mode that result in SHA mismatches for Homebrew users.

The original `homebrew-distribution` design (D4) anticipated that a future `sign-macos` job would patch checksums before pushing to the tap, but the implementation coupled signing and publishing into a single job without transactional guarantees.

## Goals / Non-Goals

### Goals
- `brew install unbound-force/tap/replicator` works for both signed and unsigned releases
- SHA in the Homebrew cask always matches the actual release asset
- No partial-failure window where the release archive has a different SHA than the cask
- Pipeline works correctly whether signing secrets are available or not (Composability First)

### Non-Goals
- Post-release `verify-installation` CI job (future enhancement)
- Changes to the org-infra reusable workflows
- Changes to `.goreleaser.yaml` (the `skip_upload: true` pattern is correct and stays)
- darwin_amd64 support (currently excluded in `.goreleaser.yaml` via `ignore:`)

## Decisions

**D1: Dedicated `publish-cask` job.** Extract all Homebrew tap logic from `sign-macos` into a new `publish-cask` job. This job runs after both `release` and `sign-macos` complete (or after `release` alone when signing is skipped). The key principle: compute the SHA from the actual final artifact as-downloaded from the GitHub Release, not from a local file that may differ from the uploaded asset.

**D2: Use `always()` with explicit result checks.** The `publish-cask` job uses `if: always() && needs.release.result == 'success' && (needs.sign-macos.result == 'success' || needs.sign-macos.result == 'skipped')`. The `always()` is required because GitHub Actions skips downstream jobs when an upstream job is skipped — without it, the `publish-cask` job would not run in the unsigned path (when `sign-macos` is skipped). The explicit result checks ensure `publish-cask` runs only when the release succeeded AND signing either succeeded or was intentionally skipped. If `sign-macos` **fails** or the workflow is **cancelled**, `publish-cask` does not run — a signing failure indicates the release is in an indeterminate state requiring manual intervention.

**D3: Download from release, not from prior job.** The `publish-cask` job downloads the archive directly from the GitHub Release rather than receiving it via job artifacts or outputs. This ensures the SHA is computed from the exact bytes a user would download, closing the TOCTOU gap. This aligns with Observable Quality (Principle III) -- the integrity check is against the actual distributed artifact.

**D4: Keep `sign-macos` focused on signing only.** Remove lines 193-237 (the "Update Homebrew cask" step) from `sign-macos`. This simplifies the job and eliminates the partial-failure window. The job becomes: import cert, prepare notary key, download archive, sign/notarize, re-upload signed archive, update checksums, cleanup.

**D5: Minimal permissions for `publish-cask`.** The new job needs `contents: read` (to download release assets via `gh`) plus the `HOMEBREW_TAP_TOKEN` secret (to push to the tap repo). It does not need `contents: write` on the replicator repo.

## Risks / Trade-offs

**Risk: `sign-macos` failure after archive re-upload.** If `sign-macos` fails after re-uploading the signed archive but before completing, `publish-cask` will NOT run (the D2 condition excludes `failure` results). This is intentional: a `sign-macos` failure indicates the release is in an indeterminate state. The release will have a signed archive but no cask update — manual intervention is needed. This is safer than the current state (where the same scenario produces a broken cask) because at least the release artifacts remain intact and the failure is visible in the workflow run. Recovery: re-run the workflow, or manually download the archive, compute SHA, and push to the tap.

**Risk: `HOMEBREW_TAP_TOKEN` secret availability.** The `publish-cask` job requires the `HOMEBREW_TAP_TOKEN` secret to push to the tap. If this secret is missing, the job will fail. This is the same requirement as the current `sign-macos` job, so no new risk is introduced. The failure will be visible in the workflow run.

**Trade-off: Extra job in pipeline.** Adding a job adds ~2-3 minutes to the release pipeline. This is acceptable given the pipeline already takes 20-30 minutes for notarization. The added reliability outweighs the time cost.
