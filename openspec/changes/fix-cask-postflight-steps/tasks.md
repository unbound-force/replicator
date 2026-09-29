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

## 1. Replace deprecated Homebrew cask stanza

- [x] 1.1 Remove `hooks.post.install` from `.goreleaser.yaml` and replace
  with `custom_block` containing `postflight_steps` DSL using `on_macos do`,
  `run`, and `#{staged_path}` Ruby interpolation.

- [x] 1.2 Add maintenance comment above `custom_block` noting that native
  GoReleaser `postflight_steps` support is not yet available, with tracking
  reference to goreleaser/goreleaser#6873.

## 2. Update test fixture

- [x] 2.1 Update `.github/scripts/testdata/replicator-v0.5.0.rb`: replace
  `postflight do` with `postflight_steps do`, `if OS.mac?` with
  `on_macos do`, `system_command` with `run`, and verify `#{staged_path}`
  is preserved.

## 3. Add rendered-cask regression test

- [x] 3.1 Add deterministic test case to
  `.github/scripts/patch-homebrew-cask_test.sh` that extracts the
  `custom_block` from `.goreleaser.yaml` and verifies:
  (a) `postflight_steps do` is present,
  (b) legacy `postflight do` is absent,
  (c) `#{staged_path}` remains valid Ruby interpolation.

- [x] 3.2 Add assertion to existing happy-path test case verifying
  `#{staged_path}` survives SHA patching.

## 4. Verification

- [x] 4.1 Run the Homebrew cask integrity regression suite.

- [x] 4.2 Run `go vet ./...` and `go test ./... -count=1 -race`.

- [x] 4.3 Verify constitution alignment: Composability First (Homebrew
  install works without deprecation warnings), Observable Quality
  (regression test catches reversion), Testability (deterministic,
  no network, no goreleaser binary).
