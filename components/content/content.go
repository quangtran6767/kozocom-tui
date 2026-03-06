package content

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/components/content/calendar"
	"github.com/quangtran6767/kozocom-tui/components/content/dayoff"
)

type ContentView int

const (
	ViewNone ContentView = iota
	ViewCalendar
	ViewDayOff
)

type Model struct {
	width               int
	height              int
	focused             bool
	activeView          ContentView
	calendar            calendar.Model
	calendarInitialized bool
	dayoff              dayoff.Model
	dayoffInitialized   bool
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

	if view == ViewDayOff && !m.dayoffInitialized && m.token != "" {
		m.dayoff = dayoff.New(m.token)
		m.dayoffInitialized = true
		return m.dayoff.Init()
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
	if m.activeView == ViewCalendar && m.calendarInitialized {
		return m.calendar.PanelBindings()
	}
	if m.activeView == ViewDayOff && m.dayoffInitialized {
		return m.dayoff.PanelBindings()
	}
	return []key.Binding{}
}

func (m Model) ShouldBlockGlobalQuit() bool {
	if m.activeView == ViewDayOff && m.dayoffInitialized {
		return m.dayoff.ShouldBlockGlobalQuit()
	}
	return false
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
	case ViewDayOff:
		if m.dayoffInitialized {
			if _, isKeyMsg := msg.(tea.KeyMsg); isKeyMsg && !m.focused {
				return m, nil
			}
			if m.focused {
				m.dayoff.Focus()
			} else {
				m.dayoff.Blur()
			}
			m.dayoff, cmd = m.dayoff.Update(msg)
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

	if m.activeView == ViewDayOff {
		if !m.dayoffInitialized {
			return "Initializing day-off requests..."
		}
		return m.dayoff.View(m.width, m.height)
	}

	return ""
}
