package dayoff

import (
	"reflect"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/quangtran6767/kozocom-tui/messages"
)

func TestDayOffColumnsFitRequestedWidth(t *testing.T) {
	width := 72
	columns := dayOffColumns(width)

	total := 0
	for _, column := range columns {
		total += column.Width
	}

	if total > width {
		t.Fatalf("columns overflow width: got %d, want <= %d", total, width)
	}
}

func TestNewFormModelSetsTodayDefaultTimes(t *testing.T) {
	form := newFormModel()
	today := time.Now().Format("2006-01-02")

	if got := form.inputs[formFieldStartTime].Value(); got != today+" 08:00:00" {
		t.Fatalf("unexpected default start time: %q", got)
	}
	if got := form.inputs[formFieldEndTime].Value(); got != today+" 17:00:00" {
		t.Fatalf("unexpected default end time: %q", got)
	}
}

func TestFormViewShowsRealFields(t *testing.T) {
	form := newFormModel()
	form.SetSize(76, 24)
	form.SetLeaveBalance(messages.IDayLeaveBalance{
		LeaveBalance: []messages.ILeaveBalance{
			{ID: 12, LeaveTypeID: 4, Text: "Paid Leave", RemainingDays: 3.5},
		},
	})
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 1001, Email: "alice@example.com"},
	})
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, AccountName: "Bob"},
	})

	view := ansi.Strip(form.View())

	if strings.Contains(view, "Day Off Form Placeholder") {
		t.Fatal("placeholder text should not be rendered")
	}
	if strings.Contains(view, "Format: YYYY-MM-DD HH:MM:SS") {
		t.Fatal("format helper text should be removed")
	}
	if strings.Contains(view, "Quick Reference") {
		t.Fatal("quick reference should be removed")
	}

	for _, expected := range []string{"New Day-off Request", "Start Time", "End Time", "Leave Type", "Approver", "Reason", "Paid Leave", "alice@example.com"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected form view to contain %q", expected)
		}
	}

	lines := strings.Split(view, "\n")
	firstNonEmpty := -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			firstNonEmpty = i
			break
		}
	}
	if firstNonEmpty <= 0 {
		t.Fatalf("expected centered form with top padding, got %q", view)
	}
}

func TestBuildPayloadUsesSelectedLeaveBalance(t *testing.T) {
	form := newFormModel()
	form.SetLeaveBalance(messages.IDayLeaveBalance{
		LeaveBalance: []messages.ILeaveBalance{
			{ID: 12, LeaveTypeID: 4, LeaveType: "Paid Leave", RemainingDays: 3.5},
		},
	})
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 1001, AccountName: "Alice"},
		{AccountID: 1002, AccountName: "Eve"},
	})
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, AccountName: "Bob"},
	})

	form.inputs[formFieldStartTime].SetValue("2026-03-10 08:00:00")
	form.inputs[formFieldEndTime].SetValue("2026-03-10 17:00:00")
	form.inputs[formFieldReason].SetValue("Medical appointment")
	form.selectedApproverID = 1002
	form.selectedInvolvingIDs[2001] = struct{}{}

	payload, errMsg := form.BuildPayload()
	if errMsg != "" {
		t.Fatalf("expected payload to build successfully, got error %q", errMsg)
	}

	if got := payload["employee_leave_balance_id"]; got != 12 {
		t.Fatalf("unexpected leave balance id: %#v", got)
	}

	if got := payload["leave_type_id"]; got != 4 {
		t.Fatalf("unexpected leave type id: %#v", got)
	}

	if got := payload["approvers"]; !reflect.DeepEqual(got, []int{1002}) {
		t.Fatalf("unexpected approvers payload: %#v", got)
	}

	if got := payload["involving_persons"]; !reflect.DeepEqual(got, []int{2001}) {
		t.Fatalf("unexpected involving_persons payload: %#v", got)
	}
}

func TestPickerViewSupportsSearchAndMultiSelect(t *testing.T) {
	form := newFormModel()
	form.SetSize(76, 24)
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, AccountName: "Bob"},
		{AccountID: 2002, AccountName: "Charlie"},
	})

	form.picker = pickerInvolvingPersons
	form.selectedInvolvingIDs[2001] = struct{}{}
	form.pickerSearch.SetValue("bo")

	view := ansi.Strip(form.View())
	for _, expected := range []string{"Involving", "Choose any teammates to notify", "Bob", "[x]"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected picker view to contain %q, got %q", expected, view)
		}
	}
	if strings.Contains(view, "Charlie") {
		t.Fatalf("expected picker search to filter out non-matching options, got %q", view)
	}
	if strings.Contains(view, "Start Time") {
		t.Fatalf("expected picker view to hide the main form, got %q", view)
	}
}

func TestEscClosesPickerBeforeClosingForm(t *testing.T) {
	model := New("token")
	model.state = StateForm
	model.focused = true
	model.formModel.picker = pickerApprover

	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	if updated.state != StateForm {
		t.Fatalf("expected form to stay open, got state %v", updated.state)
	}
	if updated.formModel.picker != pickerNone {
		t.Fatal("expected esc to close the picker first")
	}
}

func TestOpenPickerStartsAtCurrentSelection(t *testing.T) {
	form := newFormModel()
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 1, AccountName: "Alice"},
		{AccountID: 2, AccountName: "Bob"},
	})
	form.selectedApproverID = 2

	form.openPicker(pickerApprover)

	if form.pickerCursor != 1 {
		t.Fatalf("expected picker cursor to start at selected approver, got %d", form.pickerCursor)
	}
}

func TestEnterClosesInvolvingPickerWithoutTogglingSelection(t *testing.T) {
	form := newFormModel()
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, AccountName: "Bob"},
	})
	form.picker = pickerInvolvingPersons

	updated, _ := form.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if updated.picker != pickerNone {
		t.Fatalf("expected picker to close, got %v", updated.picker)
	}
	if len(updated.selectedInvolvingIDs) != 0 {
		t.Fatalf("expected enter to finish without toggling selection, got %#v", updated.selectedInvolvingIDs)
	}
}

func TestInvolvingPickerUsesSpaceForSearchThenSelection(t *testing.T) {
	form := newFormModel()
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, FullName: "Bob Ray"},
	})
	form.openPicker(pickerInvolvingPersons)

	updated, _ := form.Update(tea.KeyPressMsg{Code: 'b', Text: "b"})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})

	if got := updated.pickerSearch.Value(); got != "b " {
		t.Fatalf("expected space to stay in search while search is focused, got %q", got)
	}
	if len(updated.selectedInvolvingIDs) != 0 {
		t.Fatalf("expected search focus space to avoid toggling selection, got %#v", updated.selectedInvolvingIDs)
	}

	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if updated.pickerFocus != pickerFocusList {
		t.Fatalf("expected tab to move picker focus to the results list, got %v", updated.pickerFocus)
	}

	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	if _, ok := updated.selectedInvolvingIDs[2001]; !ok {
		t.Fatalf("expected space to toggle the highlighted involving person, got %#v", updated.selectedInvolvingIDs)
	}
}

func TestLabelsFallbackWhenPrimaryFieldsAreEmpty(t *testing.T) {
	form := newFormModel()
	form.SetLeaveBalance(messages.IDayLeaveBalance{
		LeaveBalance: []messages.ILeaveBalance{
			{ID: 12, LeaveTypeID: 4, Text: "Paid Leave", RemainingDays: 3.5},
		},
	})
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 1001, Email: "alice@example.com"},
	})

	if got := form.selectedLeaveBalanceLabel(); !strings.Contains(got, "Paid Leave") {
		t.Fatalf("expected leave balance label fallback to use text, got %q", got)
	}
	if got := form.selectedApproverLabel(); got != "alice@example.com" {
		t.Fatalf("expected approver label fallback to use email, got %q", got)
	}
}

func TestEmployeeLabelUsesFullNameBeforeAccountFallback(t *testing.T) {
	form := newFormModel()
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 5, FullName: "Pham Tien", AccountName: "Account 5"},
	})

	if got := form.selectedApproverLabel(); got != "Pham Tien" {
		t.Fatalf("expected approver label to use full name, got %q", got)
	}
}

func TestRenderSelectorRowTruncatesLongValue(t *testing.T) {
	row := ansi.Strip(renderSelectorRow(
		"Approver",
		"Nguyen Van A with a very very very long label that should not overflow the content area",
		false,
		20,
		false,
		lipgloss.NewStyle().Width(16),
	))

	if strings.Contains(row, "very very very long label that should not overflow the content area") {
		t.Fatalf("expected selector row to truncate long value, got %q", row)
	}
	if !strings.Contains(row, "...") {
		t.Fatalf("expected selector row to show ellipsis for truncated value, got %q", row)
	}
}

func TestModelSurfacesLeaveBalanceFetchErrors(t *testing.T) {
	model := New("token")

	updated, _ := model.Update(messages.LeaveBalanceFailMsg{Error: "leave balance failed"})

	if updated.formModel.errMsg != "leave balance failed" {
		t.Fatalf("expected leave balance error to be stored, got %q", updated.formModel.errMsg)
	}
}

func TestTableViewShowsEmptyState(t *testing.T) {
	table := newTableModel()
	table.SetSize(90, 16)
	table.SetMeta(time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), 1, 1, 0, 5)
	table.UpdateItems(nil)

	view := ansi.Strip(table.View())
	if !strings.Contains(view, "No day-off requests found.") {
		t.Fatalf("expected empty state in table view, got %q", view)
	}
}

func TestTableViewShowsMappedRequestData(t *testing.T) {
	table := newTableModel()
	table.SetSize(110, 16)
	table.SetMeta(time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), 1, 1, 1, 5)
	table.UpdateItems([]messages.DayOffRequestItem{
		{
			RequestCode:      "920bfcc1-12e2-45d2-bc57-98b2372d8cc4",
			StartTime:        "2026-03-06 08:00:00",
			EndTime:          "2026-03-06 17:00:00",
			LeaveType:        "Paid Leave",
			Status:           "Pending",
			NumberLeavesDays: 1,
			Reason:           "Medical appointment",
		},
	})

	view := ansi.Strip(table.View())
	for _, expected := range []string{"March 2026", "Page 1/1", "Paid Leave", "Pending", "1.0", "Medical appoint"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected table view to contain %q, got %q", expected, view)
		}
	}

	if strings.Contains(view, "920bfcc1...") || strings.Contains(view, "Code") {
		t.Fatalf("expected request code column to be removed, got %q", view)
	}
}

func TestSyncTableStatePaginatesItems(t *testing.T) {
	model := New("token")
	model.tableModel.SetSize(110, 16)
	model.tableModel.loaded = true
	model.items = make([]messages.DayOffRequestItem, 6)
	for i := range model.items {
		model.items[i] = messages.DayOffRequestItem{
			StartTime:        "2026-03-06 08:00:00",
			EndTime:          "2026-03-06 17:00:00",
			LeaveType:        "Paid Leave",
			Status:           "Pending",
			NumberLeavesDays: float64(i + 1),
			Reason:           "Reason " + string(rune('A'+i)),
		}
	}

	model.pager.Page = 1
	model.syncTableState()

	view := ansi.Strip(model.tableModel.View())
	if !strings.Contains(view, "Page 2/2") {
		t.Fatalf("expected second page summary, got %q", view)
	}

	if !strings.Contains(view, "Reason F") {
		t.Fatalf("expected second page item, got %q", view)
	}

	if strings.Contains(view, "Reason A") {
		t.Fatalf("expected first page item to be excluded, got %q", view)
	}
}
