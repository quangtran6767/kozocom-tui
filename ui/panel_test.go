package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// -------------------------------------------------------------------
// Dimensions — use lipgloss.Width/Height to verify output geometry
//
// RenderPanel passes `width` directly and `innerHeight = height - 2`
// to the lipgloss style, so:
//   - rendered width  == requested width  (border accounted inline)
//   - rendered height == requested height - 2  (top + bottom border rows)
// -------------------------------------------------------------------

func TestRenderPanel_Dimensions_MatchRequestedSize(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
	}{
		{"standard 80x24", 80, 24},
		{"small 20x10", 20, 10},
		{"wide 120x30", 120, 30},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := RenderPanel("Panel", "content", tc.width, tc.height, false)
			gotW := lipgloss.Width(out)
			gotH := lipgloss.Height(out)
			wantH := tc.height - 2 // innerHeight = height - 2 (border rows)
			if gotW != tc.width {
				t.Errorf("width: got %d, want %d", gotW, tc.width)
			}
			if gotH != wantH {
				t.Errorf("height: got %d, want %d (height-2)", gotH, wantH)
			}
		})
	}
}

func TestRenderPanel_Dimensions_NoTitle_MatchRequestedSize(t *testing.T) {
	out := RenderPanel("", "content", 60, 20, false)
	gotW := lipgloss.Width(out)
	gotH := lipgloss.Height(out)
	wantH := 20 - 2 // innerHeight = height - 2
	if gotW != 60 {
		t.Errorf("width: got %d, want %d", gotW, 60)
	}
	if gotH != wantH {
		t.Errorf("height: got %d, want %d (height-2)", gotH, wantH)
	}
}

// -------------------------------------------------------------------
// Guard clauses — negative/zero dims must not panic
// -------------------------------------------------------------------

func TestRenderPanel_DoesNotPanic_WithZeroSize(t *testing.T) {
	_ = RenderPanel("Title", "content", 0, 0, false)
}

func TestRenderPanel_DoesNotPanic_WithNegativeSize(t *testing.T) {
	// innerWidth/innerHeight are clamped to 0 by the guard clauses
	_ = RenderPanel("Title", "content", -10, -5, false)
}

// -------------------------------------------------------------------
// Content & title presence — strip ANSI with charmbracelet/x/ansi
// -------------------------------------------------------------------

func TestRenderPanel_OutputContainsContent(t *testing.T) {
	content := "hello world"
	out := RenderPanel("", content, 80, 24, false)
	plain := ansi.Strip(out)
	if !strings.Contains(plain, content) {
		t.Errorf("expected output to contain %q\ngot (plain):\n%s", content, plain)
	}
}

func TestRenderPanel_WithTitle_OutputContainsTitle(t *testing.T) {
	title := "My Panel"
	out := RenderPanel(title, "body", 80, 24, false)
	plain := ansi.Strip(out)
	if !strings.Contains(plain, title) {
		t.Errorf("expected output to contain title %q\ngot (plain):\n%s", title, plain)
	}
}

func TestRenderPanel_WithTitle_OutputContainsContent(t *testing.T) {
	content := "some body text"
	out := RenderPanel("Header", content, 80, 24, false)
	plain := ansi.Strip(out)
	if !strings.Contains(plain, content) {
		t.Errorf("expected output to contain content %q\ngot (plain):\n%s", content, plain)
	}
}

// -------------------------------------------------------------------
// Branch differences
// -------------------------------------------------------------------

func TestRenderPanel_EmptyTitle_DiffersFromWithTitle(t *testing.T) {
	withTitle := RenderPanel("Header", "body", 80, 24, false)
	withoutTitle := RenderPanel("", "body", 80, 24, false)
	if withTitle == withoutTitle {
		t.Error("output should differ when title is set vs empty")
	}
}

func TestRenderPanel_FocusedVsUnfocused_DifferentOutput(t *testing.T) {
	// Different border colors produce different ANSI codes
	focused := RenderPanel("Panel", "body", 80, 24, true)
	unfocused := RenderPanel("Panel", "body", 80, 24, false)
	if focused == unfocused {
		t.Error("focused and unfocused panels should produce different output (border color differs)")
	}
}
