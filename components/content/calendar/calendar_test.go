package calendar

import (
	"encoding/json"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
)

func keyMsg(s string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Text: s, Code: rune(s[0])}
}

func TestNavigateCalendar(t *testing.T) {
	m := New("dummy-token")
	m.Year = 2024
	m.Month = time.January

	t.Run("Navigate Left (h) wraps to previous year", func(t *testing.T) {
		model, _ := m.Update(keyMsg("h"))
		if model.Year != 2023 || model.Month != time.December {
			t.Errorf("Expected Dec 2023, got %s %d", model.Month, model.Year)
		}
		if !model.Loading {
			t.Errorf("Expected loading to be true after navigation")
		}
	})

	t.Run("Navigate Right (l) to next month", func(t *testing.T) {
		model, _ := m.Update(keyMsg("l"))
		if model.Year != 2024 || model.Month != time.February {
			t.Errorf("Expected Feb 2024, got %s %d", model.Month, model.Year)
		}
	})

	t.Run("Reset to today (t)", func(t *testing.T) {
		model, _ := m.Update(keyMsg("t"))
		now := time.Now()
		if model.Year != now.Year() || model.Month != now.Month() {
			t.Errorf("Expected today %s %d, got %s %d", now.Month(), now.Year(), model.Month, model.Year)
		}
	})
}

func TestReceiveAttendanceLogs(t *testing.T) {
	m := New("dummy-token")
	m.Loading = true

	msg := messages.AttendanceLogsMsg{
		Data: json.RawMessage(`{"id":"data-1","actual_work":10}`),
	}

	model, _ := m.Update(msg)

	if model.Loading {
		t.Errorf("Expected Loading to become false after receiving logs")
	}
	if model.Data == nil {
		t.Errorf("Expected Data to be unmarshaled, got nil")
	} else if model.Data.ActualWork != 10 {
		t.Errorf("Expected ActualWork 10, got %f", model.Data.ActualWork)
	}
}
