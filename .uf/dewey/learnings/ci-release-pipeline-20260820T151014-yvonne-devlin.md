---
tag: ci-release-pipeline
author: yvonne-devlin
category: pattern
created_at: 2026-08-20T15:10:14Z
identity: ci-release-pipeline-20260820T151014-yvonne-devlin
tier: draft
---

When GitHub Actions workflows couple a signing/notarization job with Homebrew cask publishing in the same job, a TOCTOU race condition can emerge: the signing step re-uploads the archive (changing its SHA) and then the cask-patching step uses the new SHA, but if the job fails between those two operations, or if the job is skipped entirely, the Homebrew tap ends up with a stale or missing cask. The fix is to extract cask publishing into a dedicated downstream job that downloads the final artifact from the GitHub Release, computes the SHA from the actual downloaded artifact (not from a prior job's shell variable), and patches the cask. The job uses `always() && needs.release.result == 'success' && (needs.sign-macos.result == 'success' || needs.sign-macos.result == 'skipped')` to handle both signed and unsigned release paths. The `always()` is essential because GitHub Actions auto-skips downstream jobs when an upstream job is skipped — without `always()`, the cask would never publish for unsigned releases.
