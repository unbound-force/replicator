---
tag: ci-workflow-testing
author: yvonne-devlin
category: pattern
created_at: 2026-08-20T15:10:27Z
identity: ci-workflow-testing-20260820T151027-yvonne-devlin
tier: draft
---

For CI-only changes (workflow YAML modifications with no Go source code), the pre-flight skill's file-scope filter correctly identifies that Go-specific tools (go test, golangci-lint) have no in-scope files and skips them. The relevant validation for workflow changes is: (1) `goreleaser check` if the GoReleaser config is affected, (2) `actionlint` if available, (3) structural review of the YAML job graph, and (4) runtime self-tests embedded in the workflow steps themselves (e.g., SHA format validation regex, grep verification). The spec review council (6 agents) and code review council (5 agents) both operated in spec-review and code-review modes respectively, providing comprehensive coverage even without traditional unit tests. The `make check` command still serves as the CI parity gate since it's marked "always in scope" by the pre-flight skill.
