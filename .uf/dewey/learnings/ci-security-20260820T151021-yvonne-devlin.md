---
tag: ci-security
author: yvonne-devlin
category: gotcha
created_at: 2026-08-20T15:10:21Z
identity: ci-security-20260820T151021-yvonne-devlin
tier: draft
---

In GitHub Actions workflow YAML, using `${{ steps.X.outputs.Y }}` directly inside a `run:` shell block is a script injection vector (even if the value is currently safe like a SHA hex string). The safer pattern is to pass the value via an `env:` binding on the step: `env: MY_VAR: ${{ steps.X.outputs.Y }}` and then reference `$MY_VAR` in the shell. This prevents the value from being interpolated into the shell command template before execution. The replicator release workflow was updated to consistently use this pattern for all step outputs (`RELEASE_TAG`, `GH_TOKEN`, `ARM64_SHA`). Additionally, an empty-string guard (`if [ -z "$VAR" ]; then ... exit 1; fi`) should be added when consuming step outputs, because GitHub Actions output passing can silently produce empty values if the output name is misspelled or the write was truncated.
