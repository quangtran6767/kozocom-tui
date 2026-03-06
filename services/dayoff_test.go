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
