package services

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildDayOffRequestsURL(t *testing.T) {
	urlStr := buildDayOffRequestsURL(DayOffRequestQuery{
		DateStart: "2026-03",
		DateEnd:   "2026-03",
	})

	parsed, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("expected valid url, got error %v", err)
	}

	if got := parsed.Query().Get("date_start"); got != "2026-03" {
		t.Fatalf("unexpected date_start: %q", got)
	}

	if got := parsed.Query().Get("date_end"); got != "2026-03" {
		t.Fatalf("unexpected date_end: %q", got)
	}
}

func TestDecodeDayOffRequestsResponse(t *testing.T) {
	body := strings.NewReader(`{
		"meta": {
			"code": 200,
			"error_message": null
		},
		"data": {
			"data": [
				{
					"request_code": "920bfcc1-12e2-45d2-bc57-98b2372d8cc4",
					"start_time": "2026-03-06 08:00",
					"end_time": "2026-03-06 17:00",
					"reason": "Medical appointment",
					"status": "Pending",
					"leave_type": "Paid Leave",
					"total_days": 1,
					"created_at": "2026-03-06 08:42"
				}
			]
		}
	}`)

	items, errMsg := decodeDayOffRequestsResponse(body)
	if errMsg != "" {
		t.Fatalf("expected decode success, got %q", errMsg)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].RequestCode == "" {
		t.Fatal("expected request code to be decoded")
	}

	if items[0].LeaveType != "Paid Leave" {
		t.Fatalf("unexpected leave type: %q", items[0].LeaveType)
	}

	if items[0].Status != "Pending" {
		t.Fatalf("unexpected status: %q", items[0].Status)
	}
}

func TestDecodeLeaveBalanceResponseSupportsNumericStrings(t *testing.T) {
	body := strings.NewReader(`{
		"meta": {
			"code": 200,
			"error_message": null
		},
		"data": {
			"days_off_taken_text": "1 (Approved)",
			"days_off_taken": 1,
			"leave_balance_total": 19.25,
			"days_off_taken_unpaid": 1,
			"approver": "211",
			"involving_persons": ["2", "3"],
			"leave_balance": [
				{
					"id": 16,
					"year": 2025,
					"annual_leave_days": "16.00",
					"expiry_date": "2026-03-30",
					"leave_type_id": 1,
					"leave_type": "Phép năm",
					"remaining_days": 11,
					"text": "Annual leave for the year 2025"
				}
			]
		}
	}`)

	data, errMsg := decodeLeaveBalanceResponse(body)
	if errMsg != "" {
		t.Fatalf("expected decode success, got %q", errMsg)
	}

	if len(data.LeaveBalance) != 1 {
		t.Fatalf("expected 1 leave balance item, got %d", len(data.LeaveBalance))
	}
	if data.LeaveBalance[0].AnnualLeaveDays != 16 {
		t.Fatalf("unexpected annual leave days: %v", data.LeaveBalance[0].AnnualLeaveDays)
	}
	if data.LeaveBalance[0].LeaveType != "Phép năm" {
		t.Fatalf("unexpected leave type: %q", data.LeaveBalance[0].LeaveType)
	}
}

func TestDecodeEmployeesResponseMapsFullName(t *testing.T) {
	body := strings.NewReader(`{
		"meta": {
			"code": 200,
			"error_message": null
		},
		"data": [
			{
				"account_id": 5,
				"fullName": "Pham Tien",
				"nickname": null
			}
		]
	}`)

	items, errMsg := decodeEmployeesResponse(body)
	if errMsg != "" {
		t.Fatalf("expected decode success, got %q", errMsg)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 employee item, got %d", len(items))
	}
	if items[0].FullName != "Pham Tien" {
		t.Fatalf("expected full name to be mapped, got %q", items[0].FullName)
	}
}

func TestDecodeEmployeesResponseSupportsSnakeCaseFullName(t *testing.T) {
	body := strings.NewReader(`{
		"meta": {
			"code": 200,
			"error_message": null
		},
		"data": [
			{
				"account_id": 5,
				"full_name": "Pham Tien"
			}
		]
	}`)

	items, errMsg := decodeEmployeesResponse(body)
	if errMsg != "" {
		t.Fatalf("expected decode success, got %q", errMsg)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 employee item, got %d", len(items))
	}
	if items[0].FullName != "Pham Tien" {
		t.Fatalf("expected snake_case full name to be mapped, got %q", items[0].FullName)
	}
}
