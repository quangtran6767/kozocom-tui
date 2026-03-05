package content

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/components/content/calendar"
)

type ContentView int

const (
	ViewNone ContentView = iota
	ViewCalendar
)

type Model struct {
	width               int
	height              int
	focused             bool
	activeView          ContentView
	calendar            calendar.Model
	calendarInitialized bool
	token               string
}

func New() Model {
	return Model{
		activeView: ViewNone,
	}
}

func (m *Model) SetToken(token string) tea.Cmd {
	m.token = token
	return nil
}

func (m *Model) ActivateView(view ContentView) tea.Cmd {
	m.activeView = view

	if view == ViewCalendar && !m.calendarInitialized && m.token != "" {
		m.calendar = calendar.New(m.token)
		m.calendarInitialized = true
		return m.calendar.Init()
	}

	return nil
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
	if m.activeView == ViewNone {
		return m, nil
	}

	var cmd tea.Cmd

	switch m.activeView {
	case ViewCalendar:
		if m.calendarInitialized {
			if _, isKeyMsg := msg.(tea.KeyMsg); isKeyMsg && !m.focused {
				return m, nil
			}
			m.calendar, cmd = m.calendar.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.activeView == ViewNone {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Welcome to Kozocom TUI\n\nPlease select an item from the sidebar.")
	}

	if m.activeView == ViewCalendar {
		if !m.calendarInitialized {
			return "Initializing calendar..."
		}
		return m.calendar.View(m.width, m.height)
	}

	return ""
}
