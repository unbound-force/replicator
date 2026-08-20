## ADDED Requirements

### Requirement: Dedicated cask publish job

The release pipeline MUST include a dedicated `publish-cask` job that is solely responsible for publishing the Homebrew cask to `unbound-force/homebrew-tap`.

#### Scenario: Successful release with signing
- **GIVEN** the `release` job completed successfully AND the `sign-macos` job completed successfully
- **WHEN** the `publish-cask` job runs
- **THEN** it MUST download the final darwin_arm64 archive from the GitHub Release, compute the SHA-256 from the downloaded artifact, patch the cask file, and push to the Homebrew tap

#### Scenario: Successful release without signing (no secrets)
- **GIVEN** the `release` job completed successfully AND the `sign-macos` job was skipped (no signing secrets)
- **WHEN** the `publish-cask` job runs
- **THEN** it MUST download the unsigned darwin_arm64 archive from the GitHub Release, compute the SHA-256 from the downloaded artifact, patch the cask file, and push to the Homebrew tap

#### Scenario: Release job failed
- **GIVEN** the `release` job failed
- **WHEN** the GitHub Actions job graph evaluates the `publish-cask` job
- **THEN** the `publish-cask` job MUST NOT run

#### Scenario: sign-macos failed (partial or complete failure)
- **GIVEN** the `release` job completed successfully AND the `sign-macos` job failed (regardless of whether it re-uploaded the signed archive before failing)
- **WHEN** the GitHub Actions job graph evaluates the `publish-cask` job
- **THEN** the `publish-cask` job MUST NOT run. A `sign-macos` failure indicates the release is in an indeterminate state and manual intervention is needed.

#### Scenario: Workflow cancelled during sign-macos
- **GIVEN** the `release` job completed successfully AND the workflow was cancelled during `sign-macos`
- **WHEN** the GitHub Actions job graph evaluates the `publish-cask` job
- **THEN** the `publish-cask` job MUST NOT run (cancelled is not success or skipped)

#### Scenario: HOMEBREW_TAP_TOKEN unavailable
- **GIVEN** the `release` and `sign-macos` jobs completed successfully AND the `HOMEBREW_TAP_TOKEN` secret is not configured or has expired
- **WHEN** the `publish-cask` job attempts to push to the Homebrew tap
- **THEN** the job MUST fail with a clear error AND the GitHub Release artifacts MUST remain intact (the job has read-only contents permission)

### Requirement: SHA format validation

The `publish-cask` job MUST validate that the computed SHA-256 is a 64-character hexadecimal string before using it for cask patching. If validation fails, the job MUST fail with an error annotation.

### Requirement: SHA verification before tap push

The `publish-cask` job MUST verify that the patched cask file contains the computed SHA in the correct context (associated with the `darwin_arm64` section) before pushing to the Homebrew tap.

#### Scenario: SHA patching succeeds
- **GIVEN** the cask file has been patched with the computed SHA
- **WHEN** the verification step runs
- **THEN** the computed SHA MUST appear in the cask file AND the push to the tap MUST proceed

#### Scenario: SHA patching fails
- **GIVEN** the cask file was patched but the computed SHA does not appear in the result
- **WHEN** the verification step runs
- **THEN** the job MUST fail with an error annotation AND the push to the tap MUST NOT proceed

## MODIFIED Requirements

### Requirement: sign-macos job scope

Previously: The `sign-macos` job was responsible for signing, notarization, archive re-upload, checksum update, AND Homebrew cask patching and tap publishing.

The `sign-macos` job MUST be limited to signing, notarization, archive re-upload, and checksum update. It MUST NOT publish to the Homebrew tap.

#### Scenario: sign-macos completes
- **GIVEN** the signing and notarization steps succeed
- **WHEN** the `sign-macos` job completes
- **THEN** the signed archives and updated checksums MUST be uploaded to the GitHub Release AND no Homebrew tap operations MUST occur

## REMOVED Requirements

None.
