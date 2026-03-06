package dayoff

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
)

const (
	formFieldStartTime = iota
	formFieldEndTime
	formFieldLeaveBalance
	formFieldApprovers
	formFieldInvolvingPersons
	formFieldReason
)

type formModel struct {
	width        int
	height       int
	focusIndex   int
	inputs       []textinput.Model
	leaveBalance messages.IDayLeaveBalance
	approvers    []messages.EmployeeItem
	allEmployees []messages.EmployeeItem
	leaveDays    float64
	errMsg       string
	submitting   bool
}

func newFormModel() formModel {
	inputs := make([]textinput.Model, 6)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 255
	}

	inputs[formFieldStartTime].Placeholder = "YYYY-MM-DD HH:MM:SS"
	inputs[formFieldEndTime].Placeholder = "YYYY-MM-DD HH:MM:SS"
	inputs[formFieldLeaveBalance].Placeholder = "Employee leave balance ID"
	inputs[formFieldApprovers].Placeholder = "Approver account IDs, comma separated"
	inputs[formFieldInvolvingPersons].Placeholder = "Optional account IDs, comma separated"
	inputs[formFieldReason].Placeholder = "Reason"

	return formModel{
		focusIndex: formFieldStartTime,
		inputs:     inputs,
	}
}

func (m *formModel) SetLeaveBalance(data messages.IDayLeaveBalance) {
	m.leaveBalance = data
	if m.inputs[formFieldLeaveBalance].Value() == "" && len(data.LeaveBalance) > 0 {
		m.inputs[formFieldLeaveBalance].SetValue(strconv.Itoa(data.LeaveBalance[0].ID))
	}
}

func (m *formModel) SetApprovers(data []messages.EmployeeItem) {
	m.approvers = data
	if m.inputs[formFieldApprovers].Value() == "" && len(data) > 0 {
		m.inputs[formFieldApprovers].SetValue(strconv.Itoa(data[0].AccountID))
	}
}

func (m *formModel) SetAllEmployees(data []messages.EmployeeItem) {
	m.allEmployees = data
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

func (m *formModel) FocusFirst() tea.Cmd {
	m.focusIndex = formFieldStartTime
	return m.focusCurrent()
}

func (m *formModel) Blur() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func (m *formModel) BuildPayload() (map[string]interface{}, string) {
	startTime := strings.TrimSpace(m.inputs[formFieldStartTime].Value())
	endTime := strings.TrimSpace(m.inputs[formFieldEndTime].Value())
	reason := strings.TrimSpace(m.inputs[formFieldReason].Value())
	leaveBalanceIDRaw := strings.TrimSpace(m.inputs[formFieldLeaveBalance].Value())
	approverRaw := strings.TrimSpace(m.inputs[formFieldApprovers].Value())
	involvingRaw := strings.TrimSpace(m.inputs[formFieldInvolvingPersons].Value())

	if startTime == "" || endTime == "" || reason == "" || leaveBalanceIDRaw == "" || approverRaw == "" {
		return nil, "Start time, end time, leave type, approver, and reason are required."
	}

	leaveBalanceID, err := strconv.Atoi(leaveBalanceIDRaw)
	if err != nil {
		return nil, "Leave type must be a numeric employee leave balance ID."
	}

	var leaveTypeID int
	for _, item := range m.leaveBalance.LeaveBalance {
		if item.ID == leaveBalanceID {
			leaveTypeID = item.LeaveTypeID
			break
		}
	}
	if leaveTypeID == 0 {
		return nil, "Selected leave balance ID is not available."
	}

	approvers, err := parseIDList(approverRaw)
	if err != nil || len(approvers) == 0 {
		return nil, "Approvers must be a comma-separated list of numeric account IDs."
	}

	involvingPersons, err := parseIDList(involvingRaw)
	if err != nil {
		return nil, "Involving persons must be a comma-separated list of numeric account IDs."
	}

	payload := map[string]interface{}{
		"start_time":                startTime,
		"end_time":                  endTime,
		"employee_leave_balance_id": leaveBalanceID,
		"leave_type_id":             leaveTypeID,
		"reason":                    reason,
		"approvers":                 approvers,
		"involving_persons":         involvingPersons,
		"status":                    1,
	}

	return payload, ""
}

func (m *formModel) focusCurrent() tea.Cmd {
	var cmd tea.Cmd
	for i := range m.inputs {
		if i == m.focusIndex {
			cmd = m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}

	return cmd
}

func (m *formModel) moveFocus(delta int) tea.Cmd {
	m.focusIndex = (m.focusIndex + delta + len(m.inputs)) % len(m.inputs)
	return m.focusCurrent()
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			return m, m.moveFocus(1)
		case "shift+tab", "up":
			return m, m.moveFocus(-1)
		}
	}

	m.errMsg = ""

	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m formModel) View() string {
	contentWidth := maxInt(0, m.width)
	labelWidth := 14
	if contentWidth < 56 {
		labelWidth = maxInt(8, contentWidth)
	}

	labelStyle := lipgloss.NewStyle().
		Width(labelWidth).
		Foreground(lipgloss.Color("245"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63"))

	bodyWidth := maxInt(8, contentWidth-labelWidth-2)
	if contentWidth < 56 {
		bodyWidth = maxInt(8, contentWidth-2)
	}
	for i := range m.inputs {
		m.inputs[i].SetWidth(bodyWidth)
	}

	rows := []string{
		renderFieldRow("Start Time", m.inputs[formFieldStartTime], labelStyle, contentWidth),
		renderFieldRow("End Time", m.inputs[formFieldEndTime], labelStyle, contentWidth),
		renderFieldRow("Leave Type", m.inputs[formFieldLeaveBalance], labelStyle, contentWidth),
		renderFieldRow("Approvers", m.inputs[formFieldApprovers], labelStyle, contentWidth),
		renderFieldRow("Involving", m.inputs[formFieldInvolvingPersons], labelStyle, contentWidth),
		renderFieldRow("Reason", m.inputs[formFieldReason], labelStyle, contentWidth),
	}

	summary := []string{
		sectionTitleStyle.Render("New Day-off Request"),
		infoStyle.Render(fmt.Sprintf("Taken %.1f | Total %.1f | Unpaid %.1f | Requested %.1f",
			m.leaveBalance.DaysOffTaken,
			m.leaveBalance.LeaveBalanceTotal,
			m.leaveBalance.DaysOffTakenUnpaid,
			m.leaveDays,
		)),
		infoStyle.Render("Format: YYYY-MM-DD HH:MM:SS"),
	}

	referenceLines := []string{
		sectionTitleStyle.Render("Quick Reference"),
		infoStyle.Render("Leave types: " + compactLines(formatLeaveTypes(m.leaveBalance.LeaveBalance))),
		infoStyle.Render("Approvers: " + compactLines(formatEmployees(m.approvers))),
		infoStyle.Render("Employees: " + compactLines(formatEmployees(m.allEmployees))),
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		append(summary,
			"",
			rows[0],
			rows[1],
			rows[2],
			rows[3],
			rows[4],
			rows[5],
			"",
			referenceLines[0],
			referenceLines[1],
			referenceLines[2],
			referenceLines[3],
		)...,
	)

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
		infoStyle.Render("Tab/Shift+Tab: move | Ctrl+S: submit | Esc: cancel"),
	)

	return lipgloss.NewStyle().Width(contentWidth).Render(content)
}

func (m formModel) PanelBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("tab", "shift+tab"),
			key.WithHelp("tab", "next field"),
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

func formatLeaveTypes(items []messages.ILeaveBalance) []string {
	if len(items) == 0 {
		return []string{"No leave balances loaded yet."}
	}

	lines := make([]string, 0, minInt(len(items), 5)+1)
	for i, item := range items {
		if i == 5 {
			lines = append(lines, fmt.Sprintf("...and %d more", len(items)-i))
			break
		}

		lines = append(lines, fmt.Sprintf(
			"%d: %s (remaining %.1f)",
			item.ID,
			item.LeaveType,
			item.RemainingDays,
		))
	}

	return lines
}

func formatEmployees(items []messages.EmployeeItem) []string {
	if len(items) == 0 {
		return []string{"No employees loaded yet."}
	}

	lines := make([]string, 0, minInt(len(items), 5)+1)
	for i, item := range items {
		if i == 5 {
			lines = append(lines, fmt.Sprintf("...and %d more", len(items)-i))
			break
		}

		label := item.AccountName
		if label == "" {
			label = strings.TrimSpace(item.FirstName + " " + item.LastName)
		}

		lines = append(lines, fmt.Sprintf("%d: %s", item.AccountID, label))
	}

	return lines
}

func compactLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, " | ")
}

func renderFieldRow(label string, input textinput.Model, labelStyle lipgloss.Style, width int) string {
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

func parseIDList(raw string) ([]int, error) {
	if strings.TrimSpace(raw) == "" {
		return []int{}, nil
	}

	parts := strings.Split(raw, ",")
	ids := make([]int, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}

		id, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
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
