package dayoff

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/services"
	"time"
)

type ViewState int

const (
	StateList ViewState = iota
	StateForm
)

type Model struct {
	token    string
	width    int
	height   int
	state    ViewState
	focused  bool
	period   time.Time
	items    []messages.DayOffRequestItem
	pager    paginator.Model
	pageSize int

	tableModel tableModel
	formModel  formModel // Will define next
}

func New(token string) Model {
	now := time.Now()
	period := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	pager := paginator.New(paginator.WithPerPage(5))
	pager.KeyMap = paginator.KeyMap{
		PrevPage: key.NewBinding(
			key.WithKeys("[", "pgup"),
			key.WithHelp("[", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("]", "pgdown"),
			key.WithHelp("]", "next page"),
		),
	}

	m := Model{
		token:      token,
		state:      StateList,
		period:     period,
		pager:      pager,
		pageSize:   5,
		tableModel: newTableModel(),
		formModel:  newFormModel(),
	}
	m.tableModel.SetMeta(m.period, 1, 1, 0, m.pageSize)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchDayOffRequests(),
		services.FetchLeaveBalance(m.token),
		services.FetchApprovers(m.token),
		services.FetchAllEmployees(m.token),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}

		switch m.state {
		case StateList:
			switch msg.String() {
			case "n":
				m.formModel.ResetFields()
				m.state = StateForm
				return m, m.formModel.FocusFirst()
			case "left", "h":
				m.setPeriod(m.period.AddDate(0, -1, 0))
				return m, m.fetchDayOffRequests()
			case "right", "l":
				m.setPeriod(m.period.AddDate(0, 1, 0))
				return m, m.fetchDayOffRequests()
			case "H":
				m.setPeriod(m.period.AddDate(-1, 0, 0))
				return m, m.fetchDayOffRequests()
			case "L":
				m.setPeriod(m.period.AddDate(1, 0, 0))
				return m, m.fetchDayOffRequests()
			}

			prevPage := m.pager.Page
			m.pager, _ = m.pager.Update(msg)
			if m.pager.Page != prevPage {
				m.syncTableState()
				return m, nil
			}

		case StateForm:
			switch msg.String() {
			case "esc":
				m.state = StateList
				m.formModel.Blur()
				return m, nil
			case "ctrl+s":
				payload, errMsg := m.formModel.BuildPayload()
				if errMsg != "" {
					m.formModel.SetError(errMsg)
					return m, nil
				}

				m.formModel.SetSubmitting(true)
				return m, services.CreateDayOffRequest(m.token, payload)
			}
		}

	case messages.DayOffRequestsMsg:
		m.items = msg.Data
		m.pager.Page = 0
		m.syncTableState()
		return m, nil

	case messages.DayOffRequestsFailMsg:
		m.items = nil
		m.pager.Page = 0
		m.pager.SetTotalPages(0)
		m.tableModel.SetMeta(m.period, 1, 1, 0, m.pageSize)
		m.tableModel.SetError(msg.Error)
		return m, nil

	case messages.LeaveBalanceMsg:
		m.formModel.SetLeaveBalance(msg.Data)
		return m, nil

	case messages.ApproversMsg:
		m.formModel.SetApprovers(msg.Data)
		return m, nil

	case messages.AllEmployeesMsg:
		m.formModel.SetAllEmployees(msg.Data)
		return m, nil

	case messages.LeaveDaysCalcMsg:
		m.formModel.SetLeaveDays(msg.Result)
		return m, nil

	case messages.LeaveDaysCalcFailMsg:
		m.formModel.SetError(msg.Error)
		return m, nil

	case messages.CreateDayOffSuccessMsg:
		m.state = StateList
		m.formModel.ResetFields()
		m.formModel.Blur()
		return m, m.fetchDayOffRequests()

	case messages.CreateDayOffFailMsg:
		m.formModel.SetError(msg.Error)
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.state {
	case StateList:
		m.tableModel, cmd = m.tableModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateForm:
		m.formModel, cmd = m.formModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(maxWidth, maxHeight int) string {
	if maxWidth != m.width || maxHeight != m.height {
		m.width = maxWidth
		m.height = maxHeight
		childWidth := maxInt(1, maxWidth-2)
		childHeight := maxInt(1, maxHeight-2)
		m.tableModel.SetSize(childWidth, childHeight)
		m.formModel.SetSize(childWidth, childHeight)
	}

	if m.state == StateForm {
		return m.formModel.View()
	}

	return m.tableModel.View()
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) PanelBindings() []key.Binding {
	switch m.state {
	case StateList:
		return m.tableModel.PanelBindings()
	case StateForm:
		return m.formModel.PanelBindings()
	}
	return []key.Binding{}
}

func (m *Model) setPeriod(period time.Time) {
	m.period = time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, period.Location())
	m.pager.Page = 0
	m.items = nil
	m.syncTableState()
	m.tableModel.SetLoading()
}

func (m Model) fetchDayOffRequests() tea.Cmd {
	month := m.period.Format("2006-01")
	return services.FetchDayOffRequests(m.token, services.DayOffRequestQuery{
		DateStart: month,
		DateEnd:   month,
	})
}

func (m *Model) syncTableState() {
	totalItems := len(m.items)
	m.pager.PerPage = m.pageSize
	m.pager.SetTotalPages(totalItems)
	if m.pager.TotalPages > 0 && m.pager.Page >= m.pager.TotalPages {
		m.pager.Page = m.pager.TotalPages - 1
	}
	if m.pager.Page < 0 {
		m.pager.Page = 0
	}

	currentPage := 1
	totalPages := maxInt(1, m.pager.TotalPages)
	if totalItems > 0 {
		currentPage = m.pager.Page + 1
	}
	m.tableModel.SetMeta(m.period, currentPage, totalPages, totalItems, m.pageSize)

	start, end := m.pager.GetSliceBounds(totalItems)
	pageItems := m.items
	if totalItems > 0 {
		pageItems = m.items[start:end]
	}
	m.tableModel.UpdateItems(pageItems)
}
