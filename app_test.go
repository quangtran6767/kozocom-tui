package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
)

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

	// Since we selected the Attendance Log, active view in content should trigger initialization
	// But mostly we just want to ensure the update propagates and doesn't panic
	if newM.state != StateMain {
		t.Errorf("Expected state to remain StateMain, changed to %v", newM.state)
	}
}
