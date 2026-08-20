---
tag: github-config
author: yvonne-devlin
category: gotcha
created_at: 2026-07-27T09:09:36Z
identity: github-config-20260727T090936-yvonne-devlin
tier: draft
---

The replicator repo's .github/settings.yml has require_code_owner_reviews: true on the main branch protection, but this setting was a complete no-op until a CODEOWNERS file was created. GitHub silently ignores the setting when no CODEOWNERS file exists -- it doesn't error or warn. When auditing repo settings, always verify that require_code_owner_reviews has a matching CODEOWNERS file, and that the team referenced in CODEOWNERS (@unbound-force/overlords) actually exists on GitHub.
