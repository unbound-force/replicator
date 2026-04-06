package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewTable_RendersHeaders(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	tbl := NewTable(s, []string{"ID", "NAME"}, [][]string{
		{"1", "alpha"},
		{"2", "beta"},
	})

	output := tbl.String()

	if !strings.Contains(output, "ID") {
		t.Errorf("table output missing header 'ID':\n%s", output)
	}
	if !strings.Contains(output, "NAME") {
		t.Errorf("table output missing header 'NAME':\n%s", output)
	}
}

func TestNewTable_RendersRows(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	tbl := NewTable(s, []string{"COL"}, [][]string{
		{"hello"},
		{"world"},
	})

	output := tbl.String()

	if !strings.Contains(output, "hello") {
		t.Errorf("table output missing row 'hello':\n%s", output)
	}
	if !strings.Contains(output, "world") {
		t.Errorf("table output missing row 'world':\n%s", output)
	}
}

func TestNewTable_EmptyRows(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	tbl := NewTable(s, []string{"A", "B"}, [][]string{})

	output := tbl.String()

	// Should still render headers even with no data rows.
	if !strings.Contains(output, "A") {
		t.Errorf("empty table missing header 'A':\n%s", output)
	}
}

func TestNewTable_HasBorderCharacters(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	tbl := NewTable(s, []string{"X"}, [][]string{{"y"}})

	output := tbl.String()

	// NormalBorder uses ASCII-range box-drawing characters like ─, │, etc.
	// At minimum, the output should contain pipe or box-drawing chars.
	if !strings.ContainsAny(output, "│|─-+") {
		t.Errorf("table output missing border characters:\n%s", output)
	}
}
