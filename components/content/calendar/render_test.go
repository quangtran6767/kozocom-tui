package calendar

import (
	"strings"
	"testing"
	"time"
)

func TestRenderHeader(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	t.Run("Current month", func(t *testing.T) {
		res := renderHeader(currentYear, currentMonth)
		if !strings.Contains(res, "●") {
			t.Errorf("Expected '●' for current month, got: %s", res)
		}
	})

	t.Run("Future month", func(t *testing.T) {
		futureMonth := currentMonth + 1
		futureYear := currentYear
		if futureMonth > 12 {
			futureMonth = 1
			futureYear++
		}
		res := renderHeader(futureYear, futureMonth)
		if !strings.Contains(res, "(+") {
			t.Errorf("Expected '(+' for future month, got: %s", res)
		}
	})

	t.Run("Past month", func(t *testing.T) {
		res := renderHeader(currentYear-1, currentMonth) // 1 year ago
		if !strings.Contains(res, "(-") {
			t.Errorf("Expected '(-' relative indicator for past month, got: %s", res)
		}
	})
}

func TestRenderSummary(t *testing.T) {
	t.Run("Nil data", func(t *testing.T) {
		res := renderSummary(nil)
		if !strings.Contains(res, "No data available") {
			t.Errorf("Expected 'No data available', got: %s", res)
		}
	})

	t.Run("With data", func(t *testing.T) {
		data := &AttendanceData{
			ActualWork:    10.5,
			DaysOffTaken:  1.0,
			UnpaidLeave:   2.5,
			HoursOverTime: 5.0,
		}
		res := renderSummary(data)
		if !strings.Contains(res, "Summary") {
			t.Errorf("Expected 'Summary', got: %s", res)
		}
		if !strings.Contains(res, "10.5") {
			t.Errorf("Expected to see ActualWork 10.5, got: %s", res)
		}
		if !strings.Contains(res, "1.0") {
			t.Errorf("Expected to see DaysOffTaken 1.0, got: %s", res)
		}
	})
}
