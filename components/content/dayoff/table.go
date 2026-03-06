package dayoff

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
)

type tableModel struct {
	table       table.Model
	items       []messages.DayOffRequestItem
	loaded      bool
	errMsg      string
	width       int
	height      int
	periodLabel string
	page        int
	totalPages  int
	totalItems  int
	perPage     int
}

func newTableModel() tableModel {
	t := table.New(
		table.WithColumns(dayOffColumns(72)),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Bold(false)
	s.Cell = s.Cell.Padding(0, 1)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return tableModel{
		table: t,
	}
}

func (m *tableModel) UpdateItems(items []messages.DayOffRequestItem) {
	m.items = items
	m.loaded = true
	m.errMsg = ""
	rows := make([]table.Row, len(items))
	for i, item := range items {
		tStart, _ := time.Parse("2006-01-02 15:04:05", item.StartTime)
		tEnd, _ := time.Parse("2006-01-02 15:04:05", item.EndTime)
		fStart := item.StartTime
		fEnd := item.EndTime
		if !tStart.IsZero() {
			fStart = tStart.Format("2006-01-02 15:04")
		}
		if !tEnd.IsZero() {
			fEnd = tEnd.Format("2006-01-02 15:04")
		}

		rows[i] = table.Row{
			fStart,
			fEnd,
			emptyFallback(item.LeaveType, "-"),
			emptyFallback(item.Status, "-"),
			fmt.Sprintf("%.1f", item.NumberLeavesDays),
			compactText(item.Reason),
		}
	}
	m.table.SetRows(rows)
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tableModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63"))
	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	header := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Day-off Requests"),
		metaStyle.Render(tableSummary(m)),
	)

	body := m.table.View()
	switch {
	case m.errMsg != "":
		body = errorStyle.Render(m.errMsg)
	case !m.loaded:
		body = metaStyle.Render("Loading day-off requests...")
	case len(m.items) == 0:
		body = metaStyle.Render("No day-off requests found.")
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, "", body)
}

func (m *tableModel) SetMeta(period time.Time, page, totalPages, totalItems, perPage int) {
	m.periodLabel = period.Format("January 2006")
	m.page = page
	m.totalPages = totalPages
	m.totalItems = totalItems
	m.perPage = perPage
}

func (m *tableModel) SetSize(w, h int) {
	width := maxInt(1, w)
	height := maxInt(1, h)
	m.width = width
	m.height = height
	m.table.SetColumns(dayOffColumns(width))
	m.table.SetWidth(width)
	m.table.SetHeight(maxInt(3, height-3))
}

func (m *tableModel) SetLoading() {
	m.loaded = false
	m.errMsg = ""
	m.items = nil
	m.table.SetRows(nil)
}

func (m *tableModel) SetError(err string) {
	m.loaded = true
	m.errMsg = err
	m.items = nil
	m.table.SetRows(nil)
}

func (m tableModel) PanelBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new day off request"),
		),
		key.NewBinding(
			key.WithKeys("up", "down", "j", "k"),
			key.WithHelp("↑/↓", "navigate"),
		),
		key.NewBinding(
			key.WithKeys("left", "right", "h", "l"),
			key.WithHelp("←/→", "month"),
		),
		key.NewBinding(
			key.WithKeys("H", "L"),
			key.WithHelp("H/L", "year"),
		),
		key.NewBinding(
			key.WithKeys("[", "]"),
			key.WithHelp("[/]", "page"),
		),
	}
}

func dayOffColumns(width int) []table.Column {
	widths := fitColumnWidths(
		width,
		[]int{12, 12, 10, 9, 4, 16},
		[]int{15, 15, 12, 10, 5, 30},
	)

	return []table.Column{
		{Title: "From", Width: widths[0]},
		{Title: "To", Width: widths[1]},
		{Title: "Type", Width: widths[2]},
		{Title: "Status", Width: widths[3]},
		{Title: "Days", Width: widths[4]},
		{Title: "Reason", Width: widths[5]},
	}
}

func fitColumnWidths(total int, minWidths, preferredWidths []int) []int {
	count := len(preferredWidths)
	if count == 0 {
		return nil
	}

	widths := make([]int, count)
	for i := range widths {
		widths[i] = 1
	}

	if total <= count {
		return widths
	}

	remaining := total - count
	growToTargets(widths, minWidths, &remaining)
	growToTargets(widths, preferredWidths, &remaining)

	growOrder := []int{5, 0, 1, 2, 3, 4}
	for remaining > 0 {
		for _, idx := range growOrder {
			widths[idx]++
			remaining--
			if remaining == 0 {
				break
			}
		}
	}

	return widths
}

func growToTargets(widths, targets []int, remaining *int) {
	for *remaining > 0 {
		progressed := false
		for i := range widths {
			target := maxInt(1, targets[i])
			if widths[i] >= target {
				continue
			}

			widths[i]++
			*remaining = *remaining - 1
			progressed = true
			if *remaining == 0 {
				return
			}
		}

		if !progressed {
			return
		}
	}
}

func tableSummary(m tableModel) string {
	if m.errMsg != "" {
		return m.periodLabel + " | Unable to load requests"
	}
	if !m.loaded {
		return m.periodLabel + " | Loading day-off requests..."
	}
	if m.totalItems == 0 {
		return m.periodLabel + " | 0 requests"
	}
	return fmt.Sprintf("%s | Page %d/%d | %d per page | %d requests",
		m.periodLabel,
		m.page,
		m.totalPages,
		m.perPage,
		m.totalItems,
	)
}

func compactText(value string) string {
	value = strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	return emptyFallback(value, "-")
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
