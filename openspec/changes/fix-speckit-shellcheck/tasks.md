<!--
  [P] marks tasks eligible for parallel execution.
  Add [P] when a task: (a) touches different files from
  other [P] tasks in the group, (b) has no dependency
  on prior tasks in the group, (c) can safely execute
  without ordering constraints.
  Do NOT add [P] when tasks modify the same file —
  parallel workers will cause merge conflicts.
-->

## 1. Establish the remediation baseline

- [ ] 1.1 Run ShellCheck or the available MegaLinter command against `.specify/scripts/bash/`; record the exact reported rule IDs and affected lines.
- [ ] 1.2 Confirm each planned source edit preserves the current scaffold script behavior and is suitable for upstream contribution.
- [ ] 1.3 If a source fix is proven unsafe, stop and obtain explicit authorization before proposing a documented exception limited to the affected file and ShellCheck rule; record the unsafe-correction evidence, authorization reference, and exact scope in this change.

## 2. Repair the vendored scaffold scripts

- [ ] 2.1 [P] Update `common.sh` to resolve its SC2120, SC2155, SC2221, and SC2222 findings without changing command outcomes or generated output.
- [ ] 2.2 [P] Update `check-prerequisites.sh` to make its `common.sh` source analyzable by ShellCheck while retaining the existing runtime path resolution.
- [ ] 2.3 [P] Update `create-new-feature.sh` to resolve its reported SC2155 finding while preserving command failure propagation.
- [ ] 2.4 [P] Update `setup-plan.sh` to resolve its SC1091 finding without changing scaffold behavior.
- [ ] 2.5 [P] Update `setup-tasks.sh` to resolve its SC1091 finding without changing scaffold behavior.
- [ ] 2.6 [P] Update `update-agent-context.sh` to resolve its SC1091, SC2155, and SC2034 findings without changing scaffold behavior.

## 3. Add regression verification

- [ ] 3.1 Add one focused shell smoke test for a representative changed command-substitution failure path; induce the substituted command to fail and assert the containing function or script returns a non-zero status.
- [ ] 3.2 Add a static policy assertion that `.mega-linter.yml` retains `BASH_SHELLCHECK` and introduces no global rule suppression or broad ShellCheck path exclusion.
- [ ] 3.3 Run the focused regression checks and the available ShellCheck or MegaLinter command; record the outcome and any environment limitation.

## 4. Validate governance and documentation

- [ ] 4.1 Verify the implementation remains aligned with the proposal's PASS assessments for Autonomous Collaboration, Composability First, Observable Quality, and Testability.
- [ ] 4.2 Check whether README, AGENTS.md, or generated scaffold documentation needs an update; document why no update is needed if the behavior remains unchanged.
- [ ] 4.3 Run the CI-equivalent checks required by `.github/workflows/`, including `make check`; report any unavailable hosted-linter verification separately.
