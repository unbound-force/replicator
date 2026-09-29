## Coverage Strategy

### What is being tested

The Homebrew cask postflight migration from `postflight do` (deprecated)
to `postflight_steps do` (Homebrew 7.0+ declarative DSL).

### Test surface

| Invariant | What it catches | How it is tested |
|-----------|----------------|-----------------|
| `postflight_steps do` present in `custom_block` | Missing DSL adoption | Grep extracted `custom_block` for literal `postflight_steps do` |
| `postflight do` (without `_steps`) absent from `custom_block` | Legacy stanza reversion | Negative grep: `postflight do` must NOT appear unless preceded by `_steps` |
| `#{staged_path}` present in `custom_block` | Ruby interpolation syntax intact | Grep for literal `#{staged_path}` — ensures neither Go template `{{staged_path}}` nor plain text was substituted |

### Test approach

**Deterministic rendered-cask test:**
- Extracts the `custom_block` literal block scalar from `.goreleaser.yaml`
  using awk (no external YAML parser, no network access, no goreleaser
  binary)
- Asserts all three invariants against the extracted content
- Fails if `custom_block` is not found or is empty (catches accidental
  removal)

**Existing fixture test (unchanged):**
- The happy-path case in `patch-homebrew-cask_test.sh` already asserts
  that `#{staged_path}` survives SHA patching (added in iteration 1)
- The fixture `replicator-v0.5.0.rb` uses `postflight_steps do` and
  serves as the expected output shape

### What is NOT tested

- Actual GoReleaser rendering (requires goreleaser binary + network)
- Homebrew cask installation (requires macOS + Homebrew)
- Notarization/signing interaction (tested by existing `sign-macos` job)

### Regression scenarios

| Scenario | Expected result |
|----------|----------------|
| `custom_block` reverted to `hooks.post.install` | Test fails: `custom_block` not found |
| `postflight_steps do` replaced with `postflight do` | Test fails: legacy stanza detected |
| `#{staged_path}` replaced with `{{staged_path}}` | Test fails: Ruby interpolation missing |
| `custom_block` removed entirely | Test fails: no content extracted |
