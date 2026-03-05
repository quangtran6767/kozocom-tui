package content

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/components/content/calendar"
)

type Model struct {
	width               int
	height              int
	focused             bool
	calendar            calendar.Model
	calendarInitialized bool
}

func New() Model {
	return Model{}
}

func (m *Model) SetToken(token string) tea.Cmd {
	m.calendar = calendar.New(token)
	m.calendarInitialized = true
	return m.calendar.Init()
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
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

func (m Model) PanelBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll down"),
		),
		key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.calendarInitialized {
		return m, nil
	}

	var cmd tea.Cmd
	// Only pass key messages if content is focused
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focused {
			m.calendar, cmd = m.calendar.Update(msg)
		}
	default:
		// Always pass other messages (like API responses)
		m.calendar, cmd = m.calendar.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if !m.calendarInitialized {
		return "Initializing content..."
	}
	return m.calendar.View(m.width, m.height)
}
