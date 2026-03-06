package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/components/content"
	"github.com/quangtran6767/kozocom-tui/messages"
)

func TestQDoesNotQuitWhileDayOffFormIsOpen(t *testing.T) {
	app := newAppModel()
	app.state = StateMain
	app.token = "token"
	app.content.SetToken("token")
	app.switchPanel(PanelContent)
	app.content.ActivateView(content.ViewDayOff)

	model, _ := app.updateMain(tea.KeyPressMsg{Code: 'n', Text: "n"})
	updated, ok := model.(appModel)
	if !ok {
		t.Fatalf("expected app model after opening form, got %T", model)
	}
	if !updated.content.ShouldBlockGlobalQuit() {
		t.Fatal("expected day-off form to block the global q shortcut")
	}

	model, cmd := updated.updateMain(tea.KeyPressMsg{Code: 'q', Text: "q"})
	updated, ok = model.(appModel)
	if !ok {
		t.Fatalf("expected app model after pressing q, got %T", model)
	}
	if updated.state != StateMain {
		t.Fatalf("expected q to keep the app running, got state %v", updated.state)
	}

	if cmd != nil {
		if msg := cmd(); msg != nil {
			if _, ok := msg.(tea.QuitMsg); ok {
				t.Fatal("expected q inside the form to avoid quitting the app")
			}
		}
	}
}

func TestAppUpdateMain_PanelSwitching(t *testing.T) {
	m := newAppModel()
	m.state = StateMain // Set state to main directly for testing

	// Press '2' to switch to content
	msg := tea.KeyPressMsg{Text: "2", Code: '2'}
	updatedModel, _ := m.Update(msg)
	newM := updatedModel.(appModel)

	if newM.activePanel != PanelContent {
		t.Errorf("Expected activePanel to switch to PanelContent after pressing '2'")
	}
	if !newM.content.IsFocused() {
		t.Errorf("Expected content panel to be focused")
	}

	// Press '1' to switch to sidebar
	msg = tea.KeyPressMsg{Text: "1", Code: '1'}
	updatedModel, _ = newM.Update(msg)
	newM = updatedModel.(appModel)

	if newM.activePanel != PanelSidebar {
		t.Errorf("Expected activePanel to switch to PanelSidebar after pressing '1'")
	}
	if newM.content.IsFocused() {
		t.Errorf("Expected content panel to be blurred")
	}
}

func TestAppUpdateMain_SidebarItemSelected(t *testing.T) {
	m := newAppModel()
	m.state = StateMain

	// Simulate selecting AttendanceLog from sidebar
	msg := messages.SidebarItemSelectedMsg{Item: messages.MenuAttendanceLog}
	updatedModel, _ := m.Update(msg)
	newM := updatedModel.(appModel)

	if newM.state != StateMain {
		t.Errorf("Expected state to remain StateMain, changed to %v", newM.state)
	}
	if newM.activePanel != PanelContent {
		t.Errorf("Expected sidebar selection to focus content panel")
	}
	if !newM.content.IsFocused() {
		t.Errorf("Expected content panel to be focused after sidebar selection")
	}
}

func TestAppUpdateMain_SidebarItemSelectedDayOff(t *testing.T) {
	m := newAppModel()
	m.enterMainState("user@example.com", "7", "test-token")

	msg := messages.SidebarItemSelectedMsg{Item: messages.MenuDayOffRequest}
	updatedModel, _ := m.Update(msg)
	newM := updatedModel.(appModel)

	if newM.activePanel != PanelContent {
		t.Fatal("Expected day-off selection to focus content panel")
	}
	if !newM.content.IsFocused() {
		t.Fatal("Expected content panel to stay focused for day-off view")
	}
	if len(newM.content.PanelBindings()) == 0 {
		t.Fatal("Expected day-off panel bindings to be available after activation")
	}
}

func TestEnterMainState_DefaultFocusAndContent(t *testing.T) {
	m := newAppModel()

	cmd := m.enterMainState("user@example.com", "7", "test-token")

	if cmd == nil {
		t.Fatal("Expected enterMainState to return initialization commands")
	}
	if m.state != StateMain {
		t.Fatalf("Expected state to switch to main, got %v", m.state)
	}
	if m.token != "test-token" {
		t.Fatalf("Expected app token to be stored, got %q", m.token)
	}
	if m.activePanel != PanelSidebar {
		t.Fatalf("Expected sidebar to be active on startup, got %v", m.activePanel)
	}
	if !m.sidebar.IsFocused() {
		t.Fatal("Expected sidebar to be focused on startup")
	}
	if m.content.IsFocused() {
		t.Fatal("Expected content to remain unfocused on startup")
	}
	if len(m.content.PanelBindings()) == 0 {
		t.Fatal("Expected default attendance log content to be activated on startup")
	}
}
