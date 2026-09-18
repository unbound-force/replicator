#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
PATCHER="$SCRIPT_DIR/patch-homebrew-cask.sh"
FIXTURE="$SCRIPT_DIR/testdata/replicator-v0.5.0.rb"
ARCHIVE_NAME="replicator_0.5.0_darwin_arm64.tar.gz"
LINUX_ARM64_SHA="d35cf51192f4bc3eb92d32c2a63304fdbc561243a2bb8e406d0a5c7f9d1a83f1"

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

if [ ! -x "$PATCHER" ]; then
  fail "patcher is not executable: $PATCHER"
fi

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

new_case() {
  CASE_DIR="$WORK/$1"
  mkdir -p "$CASE_DIR"
  cp "$FIXTURE" "$CASE_DIR/replicator.rb"
  printf 'signed darwin archive fixture\n' > "$CASE_DIR/$ARCHIVE_NAME"
  ARCHIVE_SHA=$(sha256sum "$CASE_DIR/$ARCHIVE_NAME" | awk '{ print $1 }')
  printf '%s  %s\n' "$ARCHIVE_SHA" "$ARCHIVE_NAME" > "$CASE_DIR/checksums.txt"
}

assert_success() {
  "$PATCHER" \
    "$CASE_DIR/$ARCHIVE_NAME" \
    "$CASE_DIR/checksums.txt" \
    "$CASE_DIR/replicator.rb" >/dev/null
}

assert_failure_preserves_cask() {
  cp "$CASE_DIR/replicator.rb" "$CASE_DIR/before.rb"
  if "$PATCHER" \
    "$CASE_DIR/$ARCHIVE_NAME" \
    "$CASE_DIR/checksums.txt" \
    "$CASE_DIR/replicator.rb" >"$CASE_DIR/stdout" 2>"$CASE_DIR/stderr"; then
    fail "$1: expected failure"
  fi
  grep -q '^::error::' "$CASE_DIR/stderr" || \
    fail "$1: failure did not emit an ::error:: annotation"
  cmp -s "$CASE_DIR/before.rb" "$CASE_DIR/replicator.rb" || \
    fail "$1: original cask changed on failure"
}

new_case happy
assert_success
sed "7c\\      sha256 \"$ARCHIVE_SHA\"" "$FIXTURE" > "$CASE_DIR/expected.rb"
cmp -s "$CASE_DIR/expected.rb" "$CASE_DIR/replicator.rb" || \
  fail "happy: patched cask differs from exact expected fixture"
grep -q '#{staged_path}' "$CASE_DIR/replicator.rb" || \
  fail "happy: cask uses {{staged_path}} instead of Ruby interpolation #{staged_path}"

new_case stray-comment
printf '\n# note: darwin_arm64 builds are notarized\n' >> "$CASE_DIR/replicator.rb"
assert_success
grep -q "sha256 \"$LINUX_ARM64_SHA\"" "$CASE_DIR/replicator.rb" || \
  fail "stray-comment: linux_arm64 SHA changed"

new_case trailing-comment
sed -i '/linux_arm64.tar.gz"$/s/$/ # darwin_arm64/' "$CASE_DIR/replicator.rb"
assert_success
grep -q "sha256 \"$LINUX_ARM64_SHA\"" "$CASE_DIR/replicator.rb" || \
  fail "trailing-comment: linux_arm64 SHA changed"

new_case missing-darwin
sed -i '/darwin_arm64/d' "$CASE_DIR/replicator.rb"
assert_failure_preserves_cask "missing darwin URL"

new_case duplicate-darwin
printf '  url "https://example.invalid/replicator_darwin_arm64.tar.gz"\n' >> \
  "$CASE_DIR/replicator.rb"
assert_failure_preserves_cask "duplicate darwin URL"

new_case reordered
awk '
  NR == 7 { sha = $0; next }
  NR == 8 { print; print sha; next }
  { print }
' "$CASE_DIR/replicator.rb" > "$CASE_DIR/reordered.rb"
mv "$CASE_DIR/reordered.rb" "$CASE_DIR/replicator.rb"
assert_failure_preserves_cask "URL before sha256"

new_case stale-candidate
cat > "$CASE_DIR/replicator.rb" <<'CASK'
cask "replicator" do
  sha256 "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  url "https://example.invalid/replicator_linux_arm64.tar.gz"
  url "https://example.invalid/replicator_darwin_arm64.tar.gz"
  sha256 "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
end
CASK
assert_failure_preserves_cask "stale checksum candidate"

new_case missing-manifest
: > "$CASE_DIR/checksums.txt"
assert_failure_preserves_cask "missing manifest entry"

new_case duplicate-manifest
cat "$CASE_DIR/checksums.txt" >> "$CASE_DIR/checksums.txt.copy"
cat "$CASE_DIR/checksums.txt" >> "$CASE_DIR/checksums.txt.copy"
mv "$CASE_DIR/checksums.txt.copy" "$CASE_DIR/checksums.txt"
assert_failure_preserves_cask "duplicate manifest entry"

new_case mismatched-manifest
printf '%064d  %s\n' 0 "$ARCHIVE_NAME" > "$CASE_DIR/checksums.txt"
assert_failure_preserves_cask "mismatched manifest entry"

echo "PASS: Homebrew cask integrity regression suite"
