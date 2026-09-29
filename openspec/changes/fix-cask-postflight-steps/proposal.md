## Why

`brew install unbound-force/tap/replicator` emits repeated deprecation
warnings on Homebrew 7.0+:

```
Warning: Calling `postflight` is deprecated! Use `postflight_steps` instead.
```

The `postflight` stanza will become a hard error after 2027-12-11, at which
point the cask will fail to install entirely.

The root cause is `.goreleaser.yaml`: the `hooks.post.install` key emits a
`postflight do` block in the generated cask. Homebrew 7.0 replaced
`postflight` with the declarative `postflight_steps` DSL.

GoReleaser does not yet have native `postflight_steps` support (tracking:
goreleaser/goreleaser#6873). The `custom_block` escape hatch is available
and injects content verbatim into the generated cask.

Fixes: https://github.com/unbound-force/replicator/issues/111

## What Changes

Two changes to the release configuration and one to the test fixture:

1. **Replace deprecated hooks.** Remove `hooks.post.install` from
   `.goreleaser.yaml` and replace it with a `custom_block` containing the
   `postflight_steps` DSL (`on_macos do`, `run`, `#{staged_path}`).
2. **Update test fixture.** Update
   `.github/scripts/testdata/replicator-v0.5.0.rb` to use
   `postflight_steps do`, `on_macos do`, `run`, and `#{staged_path}`.
3. **Add rendered-cask regression test.** Add a deterministic test case
   that extracts the `custom_block` from `.goreleaser.yaml` and verifies
   `postflight_steps do` is present, legacy `postflight do` is absent,
   and `#{staged_path}` remains valid Ruby interpolation.

## Capabilities

### New Capabilities
- `postflight_steps cask stanza`: Homebrew cask uses the new declarative
  `postflight_steps` DSL instead of the deprecated `postflight` block.
- `rendered-cask regression test`: Deterministic test that validates the
  goreleaser configuration emits correct Homebrew syntax without network
  access.

### Modified Capabilities
- `quarantine removal`: Unchanged behavior (still removes
  `com.apple.quarantine` from the replicator binary on macOS); only the
  Homebrew DSL syntax changes.

### Removed Capabilities
- None

## Impact

- **Files**: `.goreleaser.yaml` (replace `hooks.post.install` with
  `custom_block`), `.github/scripts/testdata/replicator-v0.5.0.rb`
  (update fixture stanza syntax),
  `.github/scripts/patch-homebrew-cask_test.sh` (add rendered-cask
  regression test).
- **Users**: No functional change to quarantine removal behavior.
  Deprecation warnings disappear on Homebrew 7.0+. Future-proofs
  against 2027-12-11 hard error.
- **No Go source code changes.** Only release configuration and CI test
  infrastructure are affected.

## Constitution Alignment

### I. Autonomous Collaboration

**Assessment**: N/A

No MCP tools, tool output shapes, or inter-agent communication paths are
affected. This change modifies release configuration only.

### II. Composability First

**Assessment**: PASS

Replicator MUST be independently installable via Homebrew. A deprecated
stanza that will become a hard error blocks the distribution channel.
This change restores clean installation on Homebrew 7.0+ and prevents
future breakage.

### III. Observable Quality

**Assessment**: PASS

The regression test deterministically validates the goreleaser
configuration emits correct Homebrew syntax. If the `custom_block` is
reverted to legacy `postflight`, the test catches it in pull-request CI.

### IV. Testability

**Assessment**: PASS

The rendered-cask test extracts the `custom_block` from `.goreleaser.yaml`
at test time, requires no network access, no goreleaser binary, and no
Ruby interpreter. It verifies the DSL transition and Ruby interpolation
syntax deterministically.
