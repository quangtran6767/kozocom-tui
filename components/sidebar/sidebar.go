package sidebar

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Model manages the sidebar menu navigation.
// User info and checkin status are handled by the userinfo component.
type Model struct {
	width   int
	height  int
	focused bool
}

// New creates a new sidebar model.
func New() Model {
	return Model{}
}

// SetSize sets the width and height of the sidebar panel.
//
// @param w int - width
// @param h int - height
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Focus activates the sidebar panel.
func (m *Model) Focus() {
	m.focused = true
}

// Blur deactivates the sidebar panel.
func (m *Model) Blur() {
	m.focused = false
}

// IsFocused returns whether the sidebar is currently focused.
//
// @return bool - true if focused
func (m Model) IsFocused() bool {
	return m.focused
}

// PanelBindings returns sidebar-specific keybindings for the help bar.
//
// @return []key.Binding - list of keybindings
func (m Model) PanelBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "menu up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "menu down"),
		),
	}
}

// Init returns nil — no initialization commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the sidebar.
//
// @param msg tea.Msg - incoming message
// @return Model - updated model
// @return tea.Cmd - command to execute
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

// View renders the sidebar menu content.
//
// @return string - rendered output
func (m Model) View() string {
	return "Menu items here"
}
