package sidebar

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type Model struct {
	width          int
	height         int
	menuHeight     int
	userInfoHeight int
	focused        bool

	userEmail string
	userID    string
}

func New() Model {
	return Model{}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) SetMenuSize(menuH, userInfoH int) {
	m.menuHeight = menuH
	m.userInfoHeight = userInfoH
}

func (m *Model) SetUserInfo(userEmail, userID string) {
	m.userEmail = userEmail
	m.userID = userID
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) IsFocused() bool {
	return m.focused
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) ViewMenu() string {
	return "Menu items here"
}

func (m Model) ViewUserInfo() string {
	if m.userEmail == "" {
		return "Not logged in"
	}

	emailStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.LabelForeground)

	idStyle := lipgloss.NewStyle().
		Foreground(ui.LabelForeground)

	return lipgloss.JoinVertical(lipgloss.Left,
		emailStyle.Render(m.userEmail),
		idStyle.Render(m.userID),
	)
}

func (m Model) View() string {
	return m.ViewMenu()
}
