# Quickstart: Charm Bracelet CLI UX

**Branch**: `002-charm-ux` | **Date**: 2026-04-06

## Prerequisites

- Go 1.25+ installed
- Replicator repo cloned and on `002-charm-ux` branch
- `replicator setup` has been run (database exists)

## Verification Steps

### Step 1: Build

```bash
make build
```

Expected: Binary compiles with no errors. New dependencies (lipgloss,
log, termenv) resolve cleanly.

### Step 2: Doctor — Styled Output

```bash
./replicator doctor
```

Expected:
- Title line with pink bold "🩺 Replicator Doctor"
- Pass results in green with ✅ emoji
- Warnings in yellow with ⚠️ emoji
- Failures in red with ❌ emoji
- Rounded summary box with purple border showing counts
- Contextual message at bottom

### Step 3: Doctor — Pipe Fallback

```bash
./replicator doctor | cat
```

Expected:
- No ANSI escape codes in output
- Text indicators: `[PASS]`, `[WARN]`, `[FAIL]`
- No colored text, no box-drawing characters styled

Verify no escape codes:

```bash
./replicator doctor | grep -cP '\x1b\[' || echo "PASS: no ANSI codes"
```

### Step 4: Doctor — NO_COLOR

```bash
NO_COLOR=1 ./replicator doctor
```

Expected: Same plain-text output as piped, even though running in a TTY.

### Step 5: Cells — Styled Table

```bash
# Create test cells first (via MCP or direct DB)
./replicator cells
```

Expected:
- Bordered table with headers: ID, TITLE, STATUS, TYPE, PRIORITY
- Status column color-coded: green=open, yellow=in_progress, red=blocked, gray=closed
- Purple border lines

### Step 6: Cells — Empty State

```bash
# With no cells in database
./replicator cells
```

Expected: Styled "No cells found" message (not raw `[]`).

### Step 7: Cells — JSON Flag

```bash
./replicator cells --json
```

Expected: Raw JSON array output (backward compatible).

### Step 8: Cells — Pipe Fallback

```bash
./replicator cells | cat
```

Expected: Table structure preserved, no ANSI codes.

### Step 9: Version — Styled

```bash
./replicator version
```

Expected:
- Version number in bold
- Commit hash dimmed (gray)
- Build date dimmed (gray)

### Step 10: Setup — Styled

```bash
./replicator setup
```

Expected:
- Green ✓ indicators for success steps
- Red ✗ for failures
- Same color palette as doctor

### Step 11: Init — Styled

```bash
# In a new directory
./replicator init
```

Expected: Green "initialized .hive/"

```bash
# Run again
./replicator init
```

Expected: Dimmed "already initialized"

### Step 12: Stats — Styled

```bash
./replicator stats
```

Expected:
- Styled section headers (bold/pink)
- Consistent with doctor/cells palette

### Step 13: Query — Styled

```bash
./replicator query cells_by_status
```

Expected:
- Styled table with borders matching cells table
- Headers in bold purple

### Step 14: MCP Server Logging

```bash
# Start server (in a project directory)
./replicator serve 2>/dev/null &
SERVER_PID=$!

# Send a tool call
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./replicator serve 2>/dev/null

# Check log file
cat .uf/replicator/replicator.log
```

Expected:
- `.uf/replicator/` directory created
- `replicator.log` contains structured entries
- Each tool call logged with name + duration

### Step 15: Log Truncation

```bash
# Add a marker to the log
echo "MARKER_OLD_SESSION" >> .uf/replicator/replicator.log

# Restart serve
./replicator serve &
# ... send a request ...

# Verify marker is gone
grep -c "MARKER_OLD_SESSION" .uf/replicator/replicator.log
```

Expected: 0 matches — log was truncated on startup.

### Step 16: CLI Commands Don't Create Log

```bash
rm -f .uf/replicator/replicator.log
./replicator doctor
ls .uf/replicator/replicator.log 2>&1
```

Expected: "No such file or directory" — CLI commands don't create log files.

### Step 17: Full Test Suite

```bash
make check
```

Expected: All tests pass, `go vet` clean.

### Step 18: Docs Command Unchanged

```bash
./replicator docs
```

Expected: Markdown output, no lipgloss styling (docs is excluded per spec).

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Colors in piped output | Using `lipgloss.NewStyle()` instead of `renderer.NewStyle()` | Use `lipgloss.NewRenderer(w)` pattern |
| MCP responses contain ANSI | Logger writing to stdout | Logger must use `os.Stderr` + file, never stdout |
| `NO_COLOR` not respected | Old lipgloss version | Verify `go.sum` has lipgloss v1.1.0+ |
| Table too wide for terminal | No width constraint | Use `table.Width(termenv.Width())` |
| Log file not created | `.uf/replicator/` doesn't exist | `os.MkdirAll` on serve startup |
