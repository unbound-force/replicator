package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewStyles_BufferHasNoColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	if s.HasColor {
		t.Error("expected HasColor=false for bytes.Buffer (non-TTY)")
	}
}

func TestNewStyles_RendererNotNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	if s.Renderer == nil {
		t.Fatal("expected non-nil Renderer")
	}
}

func TestIndicator_PlainText_Pass(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	got := s.Indicator("pass")
	if got != "[PASS]" {
		t.Errorf("Indicator(pass) = %q, want %q", got, "[PASS]")
	}
}

func TestIndicator_PlainText_Warn(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	got := s.Indicator("warn")
	if got != "[WARN]" {
		t.Errorf("Indicator(warn) = %q, want %q", got, "[WARN]")
	}
}

func TestIndicator_PlainText_Fail(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	got := s.Indicator("fail")
	if got != "[FAIL]" {
		t.Errorf("Indicator(fail) = %q, want %q", got, "[FAIL]")
	}
}

func TestIndicator_UnknownStatus(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	got := s.Indicator("unknown")
	if got != "unknown" {
		t.Errorf("Indicator(unknown) = %q, want %q", got, "unknown")
	}
}

func TestIndicator_NoANSI_InBuffer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyles(&buf)

	// All indicators should be plain text (no ANSI escape sequences)
	// when writing to a non-TTY buffer.
	for _, status := range []string{"pass", "warn", "fail"} {
		got := s.Indicator(status)
		if strings.Contains(got, "\x1b[") {
			t.Errorf("Indicator(%s) contains ANSI escape sequences in non-TTY mode: %q", status, got)
		}
	}
}
