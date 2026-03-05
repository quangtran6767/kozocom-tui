package calendar

import (
	"encoding/json"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/services"
)

type Model struct {
	Year    int
	Month   time.Month
	Data    *AttendanceData
	Loading bool
	Spinner spinner.Model
	Token   string
}

func New(token string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	now := time.Now()

	return Model{
		Year:    now.Year(),
		Month:   now.Month(),
		Spinner: s,
		Token:   token,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.fetchDataCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.Month--
			if m.Month < time.January {
				m.Year--
				m.Month = time.December
			}
			m.Loading = true
			cmds = append(cmds, m.fetchDataCmd())
		case "l", "right":
			m.Month++
			if m.Month > time.December {
				m.Year++
				m.Month = time.January
			}
			m.Loading = true
			cmds = append(cmds, m.fetchDataCmd())
		case "r":
			m.Loading = true
			cmds = append(cmds, m.fetchDataCmd())
		}

	case messages.AttendanceLogsMsg:
		m.Loading = false
		if raw, ok := msg.Data.(json.RawMessage); ok {
			var data AttendanceData
			if err := json.Unmarshal(raw, &data); err == nil {
				m.Data = &data
			} else {
				config.DebugLog.Println("calendar Update: error unmarshaling attendance logs", err)
			}
		}

	case messages.AttendanceLogsFailMsg:
		m.Loading = false
		// Render error state inside RenderCalendar based on empty data maybe, or add an Error field.
	}

	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View(width, height int) string {
	if m.Loading {
		return "\n  " + m.Spinner.View() + " Loading attendance logs..."
	}

	return RenderCalendar(m.Year, int(m.Month), width, height, m.Data)
}

func (m Model) fetchDataCmd() tea.Cmd {
	return services.FetchAttendanceLogs(m.Token, m.Year, int(m.Month))
}

func (m Model) PanelBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev month"),
		),
		key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next month"),
		),
		key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}
