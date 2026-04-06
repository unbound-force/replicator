package hive

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatCells_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatCells([]Cell{}, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No cells found") {
		t.Errorf("empty cells should show 'No cells found', got:\n%s", output)
	}
}

func TestFormatCells_RendersHeaders(t *testing.T) {
	var buf bytes.Buffer
	cells := []Cell{
		{ID: "cell-abc12345", Title: "Test task", Status: "open", Type: "task", Priority: 1},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	for _, header := range []string{"ID", "TITLE", "STATUS", "TYPE", "PRIORITY"} {
		if !strings.Contains(output, header) {
			t.Errorf("output missing header %q:\n%s", header, output)
		}
	}
}

func TestFormatCells_RendersData(t *testing.T) {
	var buf bytes.Buffer
	cells := []Cell{
		{ID: "cell-abc12345", Title: "Fix the bug", Status: "open", Type: "bug", Priority: 2},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "cell-abc") {
		t.Errorf("output missing truncated ID 'cell-abc':\n%s", output)
	}
	if !strings.Contains(output, "Fix the bug") {
		t.Errorf("output missing title:\n%s", output)
	}
	if !strings.Contains(output, "open") {
		t.Errorf("output missing status:\n%s", output)
	}
	if !strings.Contains(output, "bug") {
		t.Errorf("output missing type:\n%s", output)
	}
}

func TestFormatCells_TruncatesLongID(t *testing.T) {
	var buf bytes.Buffer
	cells := []Cell{
		{ID: "cell-abcdef1234567890", Title: "Task", Status: "open", Type: "task", Priority: 1},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	// ID should be truncated to 8 chars.
	if strings.Contains(output, "cell-abcdef1234567890") {
		t.Errorf("full ID should be truncated, got:\n%s", output)
	}
	if !strings.Contains(output, "cell-abc") {
		t.Errorf("output missing truncated ID:\n%s", output)
	}
}

func TestFormatCells_TruncatesLongTitle(t *testing.T) {
	var buf bytes.Buffer
	longTitle := "This is a very long title that should be truncated to keep the table compact"
	cells := []Cell{
		{ID: "cell-abc", Title: longTitle, Status: "open", Type: "task", Priority: 1},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, longTitle) {
		t.Errorf("long title should be truncated, got:\n%s", output)
	}
	if !strings.Contains(output, "...") {
		t.Errorf("truncated title should end with '...':\n%s", output)
	}
}

func TestFormatCells_MultipleStatuses(t *testing.T) {
	var buf bytes.Buffer
	cells := []Cell{
		{ID: "cell-001", Title: "Open", Status: "open", Type: "task", Priority: 1},
		{ID: "cell-002", Title: "WIP", Status: "in_progress", Type: "task", Priority: 1},
		{ID: "cell-003", Title: "Stuck", Status: "blocked", Type: "bug", Priority: 2},
		{ID: "cell-004", Title: "Done", Status: "closed", Type: "task", Priority: 0},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	for _, status := range []string{"open", "in_progress", "blocked", "closed"} {
		if !strings.Contains(output, status) {
			t.Errorf("output missing status %q:\n%s", status, output)
		}
	}
}

func TestFormatCells_NoANSI(t *testing.T) {
	var buf bytes.Buffer
	cells := []Cell{
		{ID: "cell-abc", Title: "Test", Status: "open", Type: "task", Priority: 1},
	}

	if err := FormatCells(cells, &buf); err != nil {
		t.Fatalf("FormatCells error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "\x1b[") {
		t.Errorf("output contains ANSI escape sequences in non-TTY mode:\n%s", output)
	}
}
