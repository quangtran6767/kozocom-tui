package dayoff

import (
	"strings"
	"testing"
	"time"

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

func TestFormViewShowsRealFields(t *testing.T) {
	form := newFormModel()
	form.SetSize(76, 24)
	form.SetLeaveBalance(messages.IDayLeaveBalance{
		LeaveBalance: []messages.ILeaveBalance{
			{ID: 12, LeaveTypeID: 4, LeaveType: "Paid Leave", RemainingDays: 3.5},
		},
	})
	form.SetApprovers([]messages.EmployeeItem{
		{AccountID: 1001, AccountName: "Alice"},
	})
	form.SetAllEmployees([]messages.EmployeeItem{
		{AccountID: 2001, AccountName: "Bob"},
	})

	view := ansi.Strip(form.View())

	if strings.Contains(view, "Day Off Form Placeholder") {
		t.Fatal("placeholder text should not be rendered")
	}

	for _, expected := range []string{"Start Time", "End Time", "Leave Type", "Approvers", "Reason"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected form view to contain %q", expected)
		}
	}
}

func TestBuildPayloadUsesSelectedLeaveBalance(t *testing.T) {
	form := newFormModel()
	form.SetLeaveBalance(messages.IDayLeaveBalance{
		LeaveBalance: []messages.ILeaveBalance{
			{ID: 12, LeaveTypeID: 4, LeaveType: "Paid Leave", RemainingDays: 3.5},
		},
	})

	form.inputs[formFieldStartTime].SetValue("2026-03-10 08:00:00")
	form.inputs[formFieldEndTime].SetValue("2026-03-10 17:00:00")
	form.inputs[formFieldLeaveBalance].SetValue("12")
	form.inputs[formFieldApprovers].SetValue("1001,1002")
	form.inputs[formFieldInvolvingPersons].SetValue("2001")
	form.inputs[formFieldReason].SetValue("Medical appointment")

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
