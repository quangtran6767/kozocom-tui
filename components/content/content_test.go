package content

import (
	"testing"
)

func TestContentActivateCalendarView(t *testing.T) {
	m := New()

	// Before setting token, calendar should not initialize
	m.ActivateView(ViewCalendar)
	if m.calendarInitialized {
		t.Errorf("Expected calendar NOT to initialize without a token")
	}

	// Set token and activate
	m.SetToken("test-token")
	m.ActivateView(ViewCalendar)

	if !m.calendarInitialized {
		t.Errorf("Expected calendar to initialize after token is set")
	}
	if m.activeView != ViewCalendar {
		t.Errorf("Expected active view to be ViewCalendar, got %v", m.activeView)
	}
}

func TestContentFocusBlur(t *testing.T) {
	m := New()

	if m.IsFocused() {
		t.Errorf("Expected not focused initially")
	}

	m.Focus()
	if !m.IsFocused() {
		t.Errorf("Expected focused after Focus()")
	}

	m.Blur()
	if m.IsFocused() {
		t.Errorf("Expected not focused after Blur()")
	}
}

func TestContentPanelBindings(t *testing.T) {
	m := New()

	bindings := m.PanelBindings()
	if len(bindings) != 0 {
		t.Errorf("Expected no bindings when no view active, got %d", len(bindings))
	}

	m.SetToken("test-token")
	m.ActivateView(ViewCalendar)

	bindings = m.PanelBindings()
	if len(bindings) == 0 {
		t.Errorf("Expected calendar bindings when calendar is active and initialized")
	}
}
