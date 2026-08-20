<!--
  [P] marks tasks eligible for parallel execution.
  Add [P] when a task: (a) touches different files from
  other [P] tasks in the group, (b) has no dependency
  on prior tasks in the group, (c) can safely execute
  without ordering constraints.
  Do NOT add [P] when tasks modify the same file —
  parallel workers will cause merge conflicts.
  Tasks without [P] run sequentially first, then [P]
  tasks run in parallel.
-->

## 1. Restructure Release Workflow

All tasks in this group modify `.github/workflows/release.yml`.

- [x] 1.1 Remove Homebrew cask logic from `sign-macos` job: delete the "Update Homebrew cask" step (lines 193-237) including the `HOMEBREW_TAP_GITHUB_TOKEN` env var. The `sign-macos` job should end after "Replace release assets and update checksums" and "Cleanup signing materials".

- [x] 1.2 Add `publish-cask` job after `sign-macos` with `needs: [preflight, release, sign-macos]` and conditional `if: always() && needs.release.result == 'success' && (needs.sign-macos.result == 'success' || needs.sign-macos.result == 'skipped')`. Set `runs-on: ubuntu-latest`, `permissions: contents: read`, `timeout-minutes: 10`, `env: RELEASE_TAG: ${{ needs.preflight.outputs.tag }}`.

- [x] 1.3 Add step to `publish-cask`: "Download final darwin archive and compute SHA". Download `replicator_*_darwin_arm64.tar.gz` from the GitHub Release using `gh release download`. Compute SHA-256 with `shasum -a 256`. Validate the computed SHA is a 64-character hexadecimal string (`[[ "$SHA" =~ ^[0-9a-f]{64}$ ]]`); fail with an error annotation if validation fails. Output the validated SHA via `$GITHUB_OUTPUT`. Set `env: GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}` for `gh` CLI authentication.

- [x] 1.4 Add step to `publish-cask`: "Download cask template from release". Download `replicator.rb` from the GitHub Release using `gh release download`. Set `env: GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}` for `gh` CLI authentication.

- [x] 1.5 Add step to `publish-cask`: "Patch cask with actual SHA". Use `awk` to replace the darwin_arm64 SHA in the cask file (the awk pattern assumes the `sha256` directive follows the line containing `darwin_arm64` — this structural assumption is validated by the grep step). Verify the patched SHA appears in the correct context within the `darwin_arm64` section using `grep -q`; fail with `::error::` annotation if not found.

- [x] 1.6 Add step to `publish-cask`: "Push to Homebrew tap". Clone `unbound-force/homebrew-tap` using `HOMEBREW_TAP_GITHUB_TOKEN`, copy the patched cask, commit with message `feat: update replicator cask to v${VERSION}`, and push. Handle no-op case where cask is unchanged (`git diff --cached --quiet`).

## 2. Documentation

- [x] 2.1 [P] Update the workflow header comment (lines 1-16 of `release.yml`) to mention the `publish-cask` job in the pipeline description.

- [x] 2.2 [P] Update `CHANGELOG.md`: add an entry under Unreleased > Fixed noting the Homebrew cask SHA mismatch fix.

## 3. Verification

- [x] 3.1 Review the final workflow file to verify: (a) `sign-macos` has no Homebrew tap logic, (b) `publish-cask` job dependency chain is correct, (c) `publish-cask` handles both signed and skipped paths, (d) SHA is computed from the downloaded artifact (not a prior job variable), (e) verification step exists before tap push.

- [x] 3.2 Verify constitution alignment: Composability First (cask publishes regardless of signing availability), Observable Quality (SHA verification step catches failures at release time).

- [x] 3.3 Run `goreleaser check` to validate GoReleaser config is still valid after workflow changes.

<!-- spec-review: passed -->
<!-- code-review: passed -->
