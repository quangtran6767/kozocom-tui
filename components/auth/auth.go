package auth

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/services"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type Phase int

const (
	PhaseCheckingAuth Phase = iota
	PhaseLoginForm
	PhaseLoggingIn
	PhaseDone
)

const (
	focusEmail    = 0
	focusPassword = 1
)

type Model struct {
	phase      Phase
	spinner    spinner.Model
	emailInput textinput.Model
	passInput  textinput.Model
	focusIndex int
	errMsg     string
	token      string
	userID     int
	width      int
	height     int
}

func New() Model {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(ui.SpinnerStyle),
	)

	// Email input
	ei := textinput.New()
	ei.Placeholder = "Email"
	ei.CharLimit = 64
	ei.SetWidth(30)

	// Password input
	pi := textinput.New()
	pi.Placeholder = "Password"
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '•'
	pi.CharLimit = 64
	pi.SetWidth(30)

	return Model{
		phase:      PhaseCheckingAuth,
		spinner:    s,
		emailInput: ei,
		passInput:  pi,
		focusIndex: focusEmail,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m Model) IsDone() bool {
	return m.phase == PhaseDone
}

func (m Model) Token() string {
	return m.token
}

func (m Model) UserID() int {
	return m.userID
}

// Init intialize the auth flow.
// Read token from file -> if has token call /me, if not display form
func (m Model) Init() tea.Cmd {
	token, err := config.LoadToken()
	if err != nil || token == "" {
		return func() tea.Msg { return messages.NoTokenMsg{} }
	}

	return tea.Batch(
		m.spinner.Tick,
		services.CheckAuth(token),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.phase {
	case PhaseCheckingAuth:
		return m.updateCheckingAuth(msg)
	case PhaseLoginForm:
		return m.updateLoginForm(msg)
	case PhaseLoggingIn:
		return m.updateLoggingIn(msg)
	}
	return m, nil
}

func (m Model) updateCheckingAuth(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.NoTokenMsg:
		m.phase = PhaseLoginForm
		return m, m.emailInput.Focus()
	case messages.AuthCheckSuccessMsg:
		m.userID = msg.UserID
		m.phase = PhaseDone
		return m, nil
	case messages.AuthCheckFailMsg:
		m.token = ""
		m.phase = PhaseLoginForm
		return m, m.emailInput.Focus()
	}

	// Forward spinner tick
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) updateLoginForm(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if m.focusIndex == focusEmail {
				m.focusIndex = focusPassword
				m.emailInput.Blur()
				return m, m.passInput.Focus()
			}
			m.focusIndex = focusEmail
			m.passInput.Blur()
			return m, m.emailInput.Focus()
		case "enter":
			email := m.emailInput.Value()
			pass := m.passInput.Value()

			if email == "" || pass == "" {
				m.errMsg = "Please fill in all fields"
				return m, nil
			}

			m.errMsg = ""
			m.phase = PhaseLoggingIn
			m.emailInput.Blur()
			m.passInput.Blur()

			return m, tea.Batch(
				m.spinner.Tick,
				services.Login(email, pass),
			)
		}
	}

	// Forward key events for input currently being focused on
	var cmd tea.Cmd
	if m.focusIndex == focusEmail {
		m.emailInput, cmd = m.emailInput.Update(msg)
	} else {
		m.passInput, cmd = m.passInput.Update(msg)
	}

	return m, cmd
}

func (m Model) updateLoggingIn(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.LoginSuccessMsg:
		m.token = msg.Token
		m.userID = msg.UserID
		m.phase = PhaseDone
		_ = config.SaveToken(msg.Token)
		return m, nil
	case messages.LoginFailMsg:
		m.errMsg = msg.Error
		m.phase = PhaseLoginForm
		return m, m.emailInput.Focus()
	}

	// Forward spinner tick
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.phase {
	case PhaseCheckingAuth:
		return m.viewCentered(
			fmt.Sprintf("%s Checking auth session...", m.spinner.View()),
		)
	case PhaseLoggingIn:
		return m.viewCentered(
			fmt.Sprintf("%s Logging In...", m.spinner.View()),
		)
	case PhaseLoginForm:
		return m.viewLoginForm()
	}
	return ""
}

func (m Model) viewCentered(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) viewLoginForm() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.TitleForeground).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Width(10).
		Foreground(ui.LabelForeground)

	errStyle := lipgloss.NewStyle().
		Foreground(ui.ErrorForeground).
		MarginTop(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(ui.HintForeground).
		MarginTop(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.LoginForeground).
		Padding(1, 2).
		Width(44)

	title := titleStyle.Render("Login")

	emailRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Email:"),
		m.emailInput.View(),
	)

	passRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Password:"),
		m.passInput.View(),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		emailRow,
		passRow,
	)

	if m.errMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			errStyle.Render("[Err]"+m.errMsg),
		)
	}

	content = lipgloss.JoinVertical(lipgloss.Left,
		content,
		hintStyle.Render("Tab: move between field • Enter: login"),
	)

	popup := boxStyle.Render(content)

	return m.viewCentered(popup)
}
