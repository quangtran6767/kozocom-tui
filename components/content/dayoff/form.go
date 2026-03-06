package dayoff

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/quangtran6767/kozocom-tui/messages"
)

const (
	formFieldStartTime = iota
	formFieldEndTime
	formFieldLeaveBalance
	formFieldApprover
	formFieldInvolvingPersons
	formFieldReason
)

const (
	defaultStartClock = "08:00:00"
	defaultEndClock   = "17:00:00"
)

type pickerKind int

const (
	pickerNone pickerKind = iota
	pickerLeaveBalance
	pickerApprover
	pickerInvolvingPersons
)

type pickerFocus int

const (
	pickerFocusSearch pickerFocus = iota
	pickerFocusList
)

type pickerOption struct {
	id    int
	label string
	meta  string
}

type formModel struct {
	width                  int
	height                 int
	focusIndex             int
	inputs                 []textinput.Model
	leaveBalance           messages.IDayLeaveBalance
	approvers              []messages.EmployeeItem
	allEmployees           []messages.EmployeeItem
	leaveDays              float64
	errMsg                 string
	submitting             bool
	picker                 pickerKind
	pickerCursor           int
	pickerFocus            pickerFocus
	pickerSearch           textinput.Model
	selectedLeaveBalanceID int
	selectedApproverID     int
	selectedInvolvingIDs   map[int]struct{}
}

func newFormModel() formModel {
	inputs := make([]textinput.Model, 6)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 255
		inputs[i].Prompt = ""
	}

	inputs[formFieldStartTime].Placeholder = "YYYY-MM-DD HH:MM:SS"
	inputs[formFieldEndTime].Placeholder = "YYYY-MM-DD HH:MM:SS"
	inputs[formFieldLeaveBalance].Placeholder = "Press Enter to choose"
	inputs[formFieldApprover].Placeholder = "Press Enter to choose"
	inputs[formFieldInvolvingPersons].Placeholder = "Press Enter to choose"
	inputs[formFieldReason].Placeholder = "Reason"

	pickerSearch := textinput.New()
	pickerSearch.Prompt = ""
	pickerSearch.Placeholder = "Type to filter"
	pickerSearch.CharLimit = 255

	return formModel{
		focusIndex:           formFieldStartTime,
		inputs:               inputs,
		pickerSearch:         pickerSearch,
		selectedInvolvingIDs: map[int]struct{}{},
	}.withDefaultTimes(time.Now())
}

func (m *formModel) SetLeaveBalance(data messages.IDayLeaveBalance) {
	m.leaveBalance = data
	if len(data.LeaveBalance) == 0 {
		m.selectedLeaveBalanceID = 0
		return
	}

	if !m.hasLeaveBalanceID(m.selectedLeaveBalanceID) {
		m.selectedLeaveBalanceID = data.LeaveBalance[0].ID
	}
}

func (m *formModel) SetApprovers(data []messages.EmployeeItem) {
	m.approvers = data
	if len(data) == 0 {
		m.selectedApproverID = 0
		return
	}

	if !m.hasApproverID(m.selectedApproverID) {
		m.selectedApproverID = data[0].AccountID
	}
}

func (m *formModel) SetAllEmployees(data []messages.EmployeeItem) {
	m.allEmployees = data
	if len(m.selectedInvolvingIDs) == 0 {
		return
	}

	valid := make(map[int]struct{}, len(m.selectedInvolvingIDs))
	for _, item := range data {
		if _, ok := m.selectedInvolvingIDs[item.AccountID]; ok {
			valid[item.AccountID] = struct{}{}
		}
	}
	m.selectedInvolvingIDs = valid
}

func (m *formModel) SetSize(w, h int) {
	m.width = maxInt(0, w)
	m.height = maxInt(0, h)

	inputWidth := maxInt(8, minInt(56, m.width-18))
	if m.width < 48 {
		inputWidth = maxInt(8, m.width-2)
	}
	for i := range m.inputs {
		m.inputs[i].SetWidth(inputWidth)
	}
	m.pickerSearch.SetWidth(maxInt(12, minInt(40, m.width-16)))
}

func (m *formModel) SetLeaveDays(days float64) {
	m.leaveDays = days
}

func (m *formModel) SetError(err string) {
	m.errMsg = err
	m.submitting = false
}

func (m *formModel) SetSubmitting(submitting bool) {
	m.submitting = submitting
	if submitting {
		m.errMsg = ""
	}
}

func (m *formModel) ResetFields() {
	width := m.width
	height := m.height
	leaveBalance := m.leaveBalance
	approvers := m.approvers
	allEmployees := m.allEmployees

	*m = newFormModel()
	m.width = width
	m.height = height
	m.leaveBalance = leaveBalance
	m.approvers = approvers
	m.allEmployees = allEmployees
	m.SetLeaveBalance(leaveBalance)
	m.SetApprovers(approvers)
	m.SetAllEmployees(allEmployees)
	m.SetSize(width, height)
}

func (m formModel) withDefaultTimes(now time.Time) formModel {
	date := now.Format("2006-01-02")
	m.inputs[formFieldStartTime].SetValue(date + " " + defaultStartClock)
	m.inputs[formFieldEndTime].SetValue(date + " " + defaultEndClock)
	return m
}

func (m *formModel) FocusFirst() tea.Cmd {
	m.focusIndex = formFieldStartTime
	return m.focusCurrent()
}

func (m *formModel) Blur() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.pickerSearch.Blur()
}

func (m *formModel) BuildPayload() (map[string]interface{}, string) {
	startTime := strings.TrimSpace(m.inputs[formFieldStartTime].Value())
	endTime := strings.TrimSpace(m.inputs[formFieldEndTime].Value())
	reason := strings.TrimSpace(m.inputs[formFieldReason].Value())

	if startTime == "" || endTime == "" || reason == "" {
		return nil, "Start time, end time, leave type, approver, and reason are required."
	}

	selectedLeaveBalance, ok := m.selectedLeaveBalance()
	if !ok {
		return nil, "Leave type is required."
	}

	if m.selectedApproverID == 0 {
		return nil, "Approver is required."
	}

	payload := map[string]interface{}{
		"start_time":                startTime,
		"end_time":                  endTime,
		"employee_leave_balance_id": selectedLeaveBalance.ID,
		"leave_type_id":             selectedLeaveBalance.LeaveTypeID,
		"reason":                    reason,
		"approvers":                 []int{m.selectedApproverID},
		"involving_persons":         m.selectedInvolvingList(),
		"status":                    1,
	}

	return payload, ""
}

func (m *formModel) focusCurrent() tea.Cmd {
	var cmd tea.Cmd
	for i := range m.inputs {
		if i == m.focusIndex && isTextField(i) {
			cmd = m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}
	m.pickerSearch.Blur()

	return cmd
}

func (m *formModel) moveFocus(delta int) tea.Cmd {
	m.focusIndex = (m.focusIndex + delta + len(m.inputs)) % len(m.inputs)
	return m.focusCurrent()
}

func (m *formModel) setPickerFocus(focus pickerFocus) tea.Cmd {
	m.pickerFocus = focus
	if focus == pickerFocusSearch {
		return m.pickerSearch.Focus()
	}
	m.pickerSearch.Blur()
	return nil
}

func (m *formModel) togglePickerFocus() tea.Cmd {
	if m.pickerFocus == pickerFocusSearch {
		return m.setPickerFocus(pickerFocusList)
	}
	return m.setPickerFocus(pickerFocusSearch)
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	if m.picker != pickerNone {
		return m.updatePicker(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			return m, m.moveFocus(1)
		case "shift+tab", "up":
			return m, m.moveFocus(-1)
		case "enter":
			if picker := pickerForField(m.focusIndex); picker != pickerNone {
				return m, m.openPicker(picker)
			}
		}
	}

	m.errMsg = ""

	if !isTextField(m.focusIndex) {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m formModel) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63"))

	cardWidth := minInt(maxInt(36, m.width-4), 84)
	if m.width > 0 {
		cardWidth = minInt(cardWidth, m.width)
	}
	if cardWidth <= 0 {
		cardWidth = 36
	}
	contentWidth := maxInt(18, cardWidth-6)
	labelWidth := 16
	if contentWidth < 56 {
		labelWidth = maxInt(10, contentWidth)
	}
	labelStyle = labelStyle.Width(labelWidth)

	bodyWidth := maxInt(12, contentWidth-labelWidth-2)
	for i := range m.inputs {
		if isTextField(i) {
			m.inputs[i].SetWidth(bodyWidth)
		}
	}

	if m.picker != pickerNone {
		return m.pickerScreenView(cardWidth, contentWidth, sectionTitleStyle, infoStyle, errorStyle)
	}

	rows := []string{
		renderInputRow("Start Time", m.inputs[formFieldStartTime], labelStyle, contentWidth),
		renderInputRow("End Time", m.inputs[formFieldEndTime], labelStyle, contentWidth),
		renderSelectorRow("Leave Type", m.selectedLeaveBalanceLabel(), m.focusIndex == formFieldLeaveBalance, bodyWidth, contentWidth < 56, labelStyle),
		renderSelectorRow("Approver", m.selectedApproverLabel(), m.focusIndex == formFieldApprover, bodyWidth, contentWidth < 56, labelStyle),
		renderSelectorRow("Involving", m.selectedInvolvingSummary(), m.focusIndex == formFieldInvolvingPersons, bodyWidth, contentWidth < 56, labelStyle),
		renderInputRow("Reason", m.inputs[formFieldReason], labelStyle, contentWidth),
	}

	sections := []string{
		sectionTitleStyle.Render("New Day-off Request"),
		infoStyle.Render(fmt.Sprintf("Taken %.1f | Total %.1f | Unpaid %.1f | Requested %.1f",
			m.leaveBalance.DaysOffTaken,
			m.leaveBalance.LeaveBalanceTotal,
			m.leaveBalance.DaysOffTakenUnpaid,
			m.leaveDays,
		)),
		"",
		rows[0],
		rows[1],
		rows[2],
		rows[3],
		rows[4],
		rows[5],
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	if m.submitting {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", infoStyle.Render("Submitting request..."))
	}

	if m.errMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorStyle.Render(m.errMsg))
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		infoStyle.Render(m.helpText()),
	)

	card := lipgloss.NewStyle().
		Width(cardWidth).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}

func (m formModel) PanelBindings() []key.Binding {
	if m.picker != pickerNone {
		bindings := []key.Binding{
			key.NewBinding(
				key.WithKeys("up", "down", "j", "k"),
				key.WithHelp("↑/↓", "navigate"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "choose"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "close picker"),
			),
		}
		if m.picker == pickerInvolvingPersons {
			bindings = append(bindings, key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "search/results"),
			), key.NewBinding(
				key.WithKeys("space"),
				key.WithHelp("space", "toggle"),
			))
		}
		return bindings
	}

	return []key.Binding{
		key.NewBinding(
			key.WithKeys("tab", "shift+tab"),
			key.WithHelp("tab", "next field"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose option"),
		),
		key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "submit"),
		),
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close form"),
		),
	}
}

func (m formModel) updatePicker(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		options := m.filteredPickerOptions()
		if m.picker == pickerInvolvingPersons {
			switch msg.String() {
			case "tab", "shift+tab":
				return m, m.togglePickerFocus()
			}
		}
		switch msg.String() {
		case "esc":
			return m, m.closePicker()
		case "up", "k":
			if len(options) > 0 && m.pickerCursor > 0 {
				m.pickerCursor--
			}
			return m, nil
		case "down", "j":
			if len(options) > 0 && m.pickerCursor < len(options)-1 {
				m.pickerCursor++
			}
			return m, nil
		case "enter":
			if m.picker == pickerInvolvingPersons {
				return m, m.closePicker()
			}
			if len(options) == 0 {
				return m, nil
			}
			m.applyPickerSelection(options[m.pickerCursor])
			return m, m.closePicker()
		}
		if m.picker == pickerInvolvingPersons && m.pickerFocus == pickerFocusList && key.Matches(msg, key.NewBinding(key.WithKeys("space"))) {
			if len(options) > 0 {
				m.applyPickerSelection(options[m.pickerCursor])
			}
			return m, nil
		}
		if m.picker == pickerInvolvingPersons && m.pickerFocus == pickerFocusList {
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.pickerSearch, cmd = m.pickerSearch.Update(msg)
	options := m.filteredPickerOptions()
	if len(options) == 0 {
		m.pickerCursor = 0
	} else if m.pickerCursor >= len(options) {
		m.pickerCursor = len(options) - 1
	}
	return m, cmd
}

func (m *formModel) openPicker(kind pickerKind) tea.Cmd {
	m.picker = kind
	m.pickerCursor = m.currentPickerCursor(kind)
	m.errMsg = ""
	m.pickerSearch.SetValue("")
	m.pickerSearch.CursorEnd()
	return m.setPickerFocus(pickerFocusSearch)
}

func (m *formModel) closePicker() tea.Cmd {
	m.picker = pickerNone
	m.pickerCursor = 0
	m.pickerFocus = pickerFocusSearch
	m.pickerSearch.SetValue("")
	m.pickerSearch.Blur()
	return m.focusCurrent()
}

func (m *formModel) filteredPickerOptions() []pickerOption {
	options := m.pickerOptions()
	query := strings.TrimSpace(strings.ToLower(m.pickerSearch.Value()))
	if query == "" {
		return options
	}

	filtered := make([]pickerOption, 0, len(options))
	for _, option := range options {
		if strings.Contains(strings.ToLower(option.label), query) || strings.Contains(strings.ToLower(option.meta), query) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func (m *formModel) pickerOptions() []pickerOption {
	switch m.picker {
	case pickerLeaveBalance:
		options := make([]pickerOption, 0, len(m.leaveBalance.LeaveBalance))
		for _, item := range m.leaveBalance.LeaveBalance {
			options = append(options, pickerOption{
				id:    item.ID,
				label: leaveBalanceLabel(item),
				meta:  fmt.Sprintf("remaining %.1f", item.RemainingDays),
			})
		}
		return options
	case pickerApprover:
		options := make([]pickerOption, 0, len(m.approvers))
		for _, item := range m.approvers {
			options = append(options, pickerOption{
				id:    item.AccountID,
				label: employeeLabel(item),
				meta:  fmt.Sprintf("account %d", item.AccountID),
			})
		}
		return options
	case pickerInvolvingPersons:
		options := make([]pickerOption, 0, len(m.allEmployees))
		for _, item := range m.allEmployees {
			options = append(options, pickerOption{
				id:    item.AccountID,
				label: employeeLabel(item),
				meta:  fmt.Sprintf("account %d", item.AccountID),
			})
		}
		return options
	default:
		return nil
	}
}

func (m *formModel) applyPickerSelection(option pickerOption) {
	switch m.picker {
	case pickerLeaveBalance:
		m.selectedLeaveBalanceID = option.id
	case pickerApprover:
		m.selectedApproverID = option.id
	case pickerInvolvingPersons:
		if _, ok := m.selectedInvolvingIDs[option.id]; ok {
			delete(m.selectedInvolvingIDs, option.id)
			return
		}
		m.selectedInvolvingIDs[option.id] = struct{}{}
	}
}

func (m formModel) pickerView(width int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63"))
	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
	baseStyle := lipgloss.NewStyle().
		Width(maxInt(12, width-2))

	options := m.filteredPickerOptions()
	lines := []string{
		titleStyle.Render(m.pickerTitle()),
		metaStyle.Render("Search"),
		m.pickerSearch.View(),
		"",
	}

	if len(options) == 0 {
		lines = append(lines, metaStyle.Render(m.emptyPickerMessage()))
	} else {
		for i, option := range options {
			prefix := "  "
			if i == m.pickerCursor {
				prefix = "> "
			}
			marker := ""
			if m.picker == pickerInvolvingPersons {
				if _, ok := m.selectedInvolvingIDs[option.id]; ok {
					marker = "[x] "
				} else {
					marker = "[ ] "
				}
			}

			row := fmt.Sprintf("%s%s%s", prefix, marker, option.label)
			if option.meta != "" {
				row += " (" + option.meta + ")"
			}
			row = truncateLabel(row, maxInt(12, width-4))

			if i == m.pickerCursor {
				lines = append(lines, selectedStyle.Render(baseStyle.Render(row)))
				continue
			}
			lines = append(lines, baseStyle.Render(row))
		}
	}

	return lipgloss.NewStyle().
		Padding(1, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (m formModel) pickerScreenView(cardWidth, contentWidth int, sectionTitleStyle, infoStyle, errorStyle lipgloss.Style) string {
	lines := []string{
		sectionTitleStyle.Render(m.pickerTitle()),
		infoStyle.Render(m.pickerSummary()),
		"",
		m.pickerView(contentWidth),
		"",
		infoStyle.Render(m.helpText()),
	}

	if m.errMsg != "" {
		lines = append(lines, "", errorStyle.Render(m.errMsg))
	}

	card := lipgloss.NewStyle().
		Width(cardWidth).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}

func (m formModel) pickerTitle() string {
	switch m.picker {
	case pickerLeaveBalance:
		return "Leave Type"
	case pickerApprover:
		return "Approver"
	case pickerInvolvingPersons:
		return "Involving"
	default:
		return ""
	}
}

func (m formModel) helpText() string {
	if m.picker == pickerInvolvingPersons {
		return "Type to filter | Tab: search/results | Space: toggle | Enter: done | Esc: cancel"
	}
	if m.picker != pickerNone {
		return "Type to filter | Enter: choose | Esc: close"
	}
	return "Tab/Shift+Tab: move | Enter: choose option | Ctrl+S: submit | Esc: cancel"
}

func (m formModel) selectedLeaveBalance() (messages.ILeaveBalance, bool) {
	for _, item := range m.leaveBalance.LeaveBalance {
		if item.ID == m.selectedLeaveBalanceID {
			return item, true
		}
	}
	return messages.ILeaveBalance{}, false
}

func (m formModel) selectedLeaveBalanceLabel() string {
	item, ok := m.selectedLeaveBalance()
	if !ok {
		if len(m.leaveBalance.LeaveBalance) == 0 {
			return "No leave types available"
		}
		return "Press Enter to choose"
	}
	return fmt.Sprintf("%s (remaining %.1f)", leaveBalanceLabel(item), item.RemainingDays)
}

func (m formModel) selectedApproverLabel() string {
	for _, item := range m.approvers {
		if item.AccountID == m.selectedApproverID {
			return employeeLabel(item)
		}
	}
	if len(m.approvers) == 0 {
		return "No approvers available"
	}
	return "Press Enter to choose"
}

func (m formModel) selectedInvolvingSummary() string {
	ids := m.selectedInvolvingList()
	if len(ids) == 0 {
		return "Optional"
	}

	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if employee, ok := m.employeeByID(id); ok {
			names = append(names, employeeLabel(employee))
		}
	}
	if len(names) == 0 {
		return fmt.Sprintf("%d selected", len(ids))
	}
	if len(names) <= 2 {
		return strings.Join(names, ", ")
	}
	return fmt.Sprintf("%d selected: %s, %s...", len(ids), names[0], names[1])
}

func (m formModel) selectedInvolvingList() []int {
	if len(m.selectedInvolvingIDs) == 0 {
		return []int{}
	}

	ids := make([]int, 0, len(m.selectedInvolvingIDs))
	for id := range m.selectedInvolvingIDs {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

func (m formModel) employeeByID(id int) (messages.EmployeeItem, bool) {
	for _, item := range m.allEmployees {
		if item.AccountID == id {
			return item, true
		}
	}
	for _, item := range m.approvers {
		if item.AccountID == id {
			return item, true
		}
	}
	return messages.EmployeeItem{}, false
}

func (m formModel) hasLeaveBalanceID(id int) bool {
	for _, item := range m.leaveBalance.LeaveBalance {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (m formModel) hasApproverID(id int) bool {
	for _, item := range m.approvers {
		if item.AccountID == id {
			return true
		}
	}
	return false
}

func renderInputRow(label string, input textinput.Model, labelStyle lipgloss.Style, width int) string {
	if width < 56 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Render(label),
			input.View(),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render(label), input.View())
}

func renderSelectorRow(label, value string, focused bool, width int, stacked bool, labelStyle lipgloss.Style) string {
	valueStyle := lipgloss.NewStyle().
		Width(maxInt(12, width)).
		Padding(0, 1).
		Border(lipgloss.NormalBorder())
	if focused {
		valueStyle = valueStyle.BorderForeground(lipgloss.Color("63"))
	} else {
		valueStyle = valueStyle.BorderForeground(lipgloss.Color("240"))
	}

	if value == "" {
		value = "Press Enter to choose"
	}
	renderedValue := valueStyle.Render(truncateLabel(value, maxInt(8, width-4)))

	if stacked {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Render(label),
			renderedValue,
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render(label), renderedValue)
}

func employeeLabel(item messages.EmployeeItem) string {
	if item.FullName != "" {
		return strings.TrimSpace(item.FullName)
	}
	if item.AccountName != "" {
		return strings.TrimSpace(item.AccountName)
	}
	if fullName := strings.TrimSpace(item.FirstName + " " + item.LastName); fullName != "" {
		return fullName
	}
	if item.Nickname != "" {
		return strings.TrimSpace(item.Nickname)
	}
	if item.Email != "" {
		return strings.TrimSpace(item.Email)
	}
	return fmt.Sprintf("Account %d", item.AccountID)
}

func truncateLabel(value string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(value, width, "...")
}

func leaveBalanceLabel(item messages.ILeaveBalance) string {
	if item.LeaveType != "" {
		return item.LeaveType
	}
	if item.Text != "" {
		return item.Text
	}
	return fmt.Sprintf("Leave Type %d", item.LeaveTypeID)
}

func (m formModel) currentPickerCursor(kind pickerKind) int {
	options := m.optionsForKind(kind)
	for i, option := range options {
		switch kind {
		case pickerLeaveBalance:
			if option.id == m.selectedLeaveBalanceID {
				return i
			}
		case pickerApprover:
			if option.id == m.selectedApproverID {
				return i
			}
		case pickerInvolvingPersons:
			if _, ok := m.selectedInvolvingIDs[option.id]; ok {
				return i
			}
		}
	}
	return 0
}

func (m formModel) optionsForKind(kind pickerKind) []pickerOption {
	original := m.picker
	m.picker = kind
	options := m.pickerOptions()
	m.picker = original
	return options
}

func (m formModel) pickerSummary() string {
	switch m.picker {
	case pickerLeaveBalance:
		return "Choose the leave balance to use for this request."
	case pickerApprover:
		return "Choose one approver for this request."
	case pickerInvolvingPersons:
		return "Choose any teammates to notify, then press Enter when finished."
	default:
		return ""
	}
}

func (m formModel) emptyPickerMessage() string {
	switch m.picker {
	case pickerLeaveBalance:
		if strings.TrimSpace(m.pickerSearch.Value()) != "" {
			return "No matching leave types."
		}
		return "No leave types available."
	case pickerApprover:
		if strings.TrimSpace(m.pickerSearch.Value()) != "" {
			return "No matching approvers."
		}
		return "No approvers available."
	case pickerInvolvingPersons:
		if strings.TrimSpace(m.pickerSearch.Value()) != "" {
			return "No matching people."
		}
		return "No employees available."
	default:
		return "No matching results."
	}
}

func isTextField(index int) bool {
	return index == formFieldStartTime || index == formFieldEndTime || index == formFieldReason
}

func pickerForField(index int) pickerKind {
	switch index {
	case formFieldLeaveBalance:
		return pickerLeaveBalance
	case formFieldApprover:
		return pickerApprover
	case formFieldInvolvingPersons:
		return pickerInvolvingPersons
	default:
		return pickerNone
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
