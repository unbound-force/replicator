#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

# Exercise the repository-resolution path regardless of the caller environment.
unset SPECIFY_FEATURE

get_repo_root() {
    return 1
}

if get_current_branch >/dev/null; then
    printf 'expected get_current_branch to propagate get_repo_root failure\n' >&2
    exit 1
fi

if ! grep -Fxq '  - BASH_SHELLCHECK' "$REPO_ROOT/.mega-linter.yml"; then
    printf 'expected BASH_SHELLCHECK to remain enabled\n' >&2
    exit 1
fi

if grep -Eq '^(BASH_SHELLCHECK_(DISABLE|FILTER_REGEX_EXCLUDE)|SHELLCHECK_DISABLE)' "$REPO_ROOT/.mega-linter.yml"; then
    printf 'unexpected broad ShellCheck suppression or exclusion\n' >&2
    exit 1
fi
