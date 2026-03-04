package userinfo

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type Model struct {
	spinner spinner.Model
	width   int
	height  int

	userEmail      string
	userID         string
	isCheckedIn    bool
	checkinLoading bool
}

func New() Model {
	return Model{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(ui.SpinnerStyle),
		),
	}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) SetUserInfo(email, userID string) {
	m.userEmail = email
	m.userID = userID
}

func (m *Model) SetCheckinStatus(checked bool) {
	m.isCheckedIn = checked
}

func (m *Model) SetCheckinLoading(loading bool) {
	m.checkinLoading = loading
}

func (m Model) IsCheckedIn() bool {
	return m.isCheckedIn
}

func (m Model) IsLoading() bool {
	return m.checkinLoading
}

func (m Model) PanelBindings() []key.Binding {
	return []key.Binding{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.checkinLoading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.userEmail == "" {
		return "Not logged in"
	}

	emailStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.LabelForeground)

	idStyle := lipgloss.NewStyle().
		Foreground(ui.LabelForeground)

	var statusText string
	var statusStyle lipgloss.Style

	if m.checkinLoading {
		statusText = m.spinner.View() + " Checking in..."
		statusStyle = lipgloss.NewStyle().Foreground(ui.LabelForeground)
	} else if m.isCheckedIn {
		statusText = "✅ Checked in"
		statusStyle = lipgloss.NewStyle().Foreground(ui.CheckinSuccessColor)
	} else {
		statusText = "⬜ Not checked in"
		statusStyle = lipgloss.NewStyle().Foreground(ui.LabelForeground)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		emailStyle.Render(m.userEmail),
		idStyle.Render(m.userID),
		statusStyle.Render(statusText),
	)
}
