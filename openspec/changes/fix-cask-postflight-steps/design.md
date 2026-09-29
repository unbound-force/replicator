## Context

Homebrew 7.0 deprecated the `postflight do ... end` cask stanza in favor
of the declarative `postflight_steps` DSL. The deprecation will become a
hard error after 2027-12-11.

GoReleaser generates the Homebrew cask from `.goreleaser.yaml`. The
`hooks.post.install` key maps to the deprecated `postflight do` block.
GoReleaser does not yet support native `postflight_steps` generation
(tracking: goreleaser/goreleaser#6873), but the `custom_block` key injects
arbitrary Ruby verbatim into the generated cask.

The existing test infrastructure (`patch-homebrew-cask_test.sh`) validates
SHA patching logic and fixture integrity but does not exercise the
goreleaser configuration's cask stanza content. A reversion of the
`custom_block` to legacy `hooks.post.install` would pass all existing
tests.

## Goals / Non-Goals

### Goals
- Replace `hooks.post.install` with `custom_block` containing
  `postflight_steps` DSL
- Update the test fixture to match the new DSL
- Add a deterministic regression test that validates the `custom_block`
  content against three invariants: `postflight_steps do` present, legacy
  `postflight do` absent, `#{staged_path}` Ruby interpolation intact
- All tests pass without network access or goreleaser binary

### Non-Goals
- Native GoReleaser `postflight_steps` support (blocked on upstream)
- Changes to quarantine removal behavior
- Changes to Go source code
- Changes to the job graph, permissions, or secrets

## Decisions

**D1: Use `custom_block` as the migration vehicle.** GoReleaser's
`custom_block` injects content verbatim into the generated cask. This is
the documented escape hatch for DSL features that GoReleaser does not yet
model natively. A comment above the stanza tracks the upstream issue for
future migration when native support ships.

**D2: Validate the goreleaser config, not the rendered cask.** Since the
`custom_block` is injected verbatim, its content in `.goreleaser.yaml` is
byte-identical to what appears in the generated cask. Extracting and
validating it at test time is equivalent to validating the rendered output,
without requiring goreleaser or network access. This satisfies the
"deterministic, no network" constraint.

**D3: Use awk for `custom_block` extraction.** The `custom_block` value in
`.goreleaser.yaml` is a YAML literal block scalar (`|`). Its content starts
on the line after `custom_block: |` and continues while lines are indented
deeper than the `custom_block` key. Awk can reliably extract this without
a YAML parser, avoiding new dependencies. The extraction is validated by
the assertions that follow it.

**D4: Three-invariant assertion set.** The regression test checks:
1. `postflight_steps do` is present (new DSL adopted)
2. `postflight do` without `_steps` is absent (legacy form removed)
3. `#{staged_path}` is present (Ruby interpolation, not Go template
   `{{staged_path}}` or literal text)

These three checks cover the complete DSL transition. If any is violated,
the cask will either emit deprecation warnings, fail on Homebrew 7.0+,
or target a nonexistent path at install time.

**D5: Add to existing test file.** The rendered-cask test is added to
`patch-homebrew-cask_test.sh` as a separate section. It validates the
goreleaser config that feeds the patcher's fixture, keeping all cask
integrity tests in one file and one CI step.

## Risks / Trade-offs

**Risk: `custom_block` extraction relies on YAML indentation.** The awk
extractor assumes the literal block scalar follows standard YAML
indentation rules. If `.goreleaser.yaml` is reformatted with non-standard
indentation, the extraction may fail. This is mitigated by the subsequent
assertions: if extraction produces garbage, the assertions fail.

**Risk: GoReleaser changes `custom_block` semantics.** If a future
GoReleaser version wraps `custom_block` content in a method or modifies
whitespace, the cask output may diverge from the raw YAML content. The
tracking comment and upstream issue reference (goreleaser/goreleaser#6873)
flag this for future migration.

**Trade-off: awk over YAML parser.** A YAML parser would be more robust
but would require a new dependency (Python/PyYAML in CI or a Go YAML
library). The awk approach is dependency-free and sufficient for the
literal block scalar pattern used here.
