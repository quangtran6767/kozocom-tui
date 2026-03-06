package sidebar

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type MenuItem struct {
	ID    messages.MenuItemID
	Label string
}

type Model struct {
	width   int
	height  int
	focused bool
	items   []MenuItem
	cursor  int
}

func New() Model {
	return Model{
		items: []MenuItem{
			{ID: messages.MenuAttendanceLog, Label: "Attendance Log"},
			{ID: messages.MenuDayOffRequest, Label: "Day-off Request"},
		},
		cursor:  0,
		focused: false,
	}
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
			key.WithHelp("↑/k", "menu up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "menu down"),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.items) > 0 {
				selected := m.items[m.cursor].ID
				return m, func() tea.Msg {
					return messages.SidebarItemSelectedMsg{Item: selected}
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	for i, item := range m.items {
		isSelected := m.cursor == i

		var text string
		if isSelected {
			text = "> " + item.Label
		} else {
			text = "  " + item.Label
		}

		var rendered string
		if isSelected {
			style := ui.SidebarSelectedItemStyle
			if !m.focused {
				style = style.Faint(true)
			}
			rendered = style.Render(text)
		} else {
			if m.focused {
				rendered = ui.SidebarItemStyle.Render(text)
			} else {
				rendered = text
			}
		}

		sb.WriteString(rendered + "\n")
	}

	return sb.String()
}
