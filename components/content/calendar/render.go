package calendar

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/ui"
)

func renderHeader(year, month int) string {
	m := time.Month(month)
	now := time.Now()
	currentYear, currentMonth := now.Year(), now.Month()

	var title string
	if year == currentYear && m == currentMonth {
		title = fmt.Sprintf("◀  %s %d  ●  ▶", m.String(), year)
	} else {
		monthsDiff := (year-currentYear)*12 + int(m) - int(currentMonth)
		if monthsDiff > 0 {
			title = fmt.Sprintf("◀  %s %d  (+%d)  ▶", m.String(), year, monthsDiff)
		} else {
			title = fmt.Sprintf("◀  %s %d  (%d)  ▶", m.String(), year, monthsDiff)
		}
	}
	return ui.CalendarHeaderStyle.Width(56).Render(title)
}

func renderDaysOfWeek() string {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var renderedDays []string
	for _, d := range days {
		renderedDays = append(renderedDays, ui.CalendarDayHeaderStyle.Render(d))
	}
	// Note: 6*7 = 42 width, plus borders = 56
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedDays...)
}

func renderGrid(year int, month int, data *AttendanceData) string {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	todayStr := time.Now().Format("2006-01-02")
	firstWeekday := firstOfMonth.Weekday()

	// Weekday starts with Sunday = 0
	// We want Monday = 0
	startOffset := int(firstWeekday) - 1
	if startOffset < 0 {
		startOffset = 6
	}

	var grid [][]string
	var currentRow []string

	// Padding before first day
	for i := 0; i < startOffset; i++ {
		currentRow = append(currentRow, ui.CalendarEmptyDayStyle.Render(""))
	}

	for d := 1; d <= lastOfMonth.Day(); d++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, d)
		isToday := dateStr == todayStr
		weekday := time.Weekday((int(firstWeekday) + d - 1) % 7)
		isWeekend := weekday == time.Sunday || weekday == time.Saturday

		cellStyle := ui.CalendarDayStyle
		if isToday {
			cellStyle = ui.CalendarTodayStyle
		} else if isWeekend {
			cellStyle = ui.CalendarWeekendStyle
		}

		if data != nil {
			if log, ok := data.AttendanceLogs[dateStr]; ok {
				if log.HoursAttendance > 0 {
					cellStyle = ui.CalendarPresentStyle
				} else if log.HoursDayOff > 0 {
					// Check if paid or unpaid (simplification for color)
					hasPaid := false
					for _, do := range log.DayOffs {
						if do.SalaryEntitlementRate > 0 {
							hasPaid = true
							break
						}
					}
					if hasPaid {
						cellStyle = ui.CalendarPaidLeaveStyle
					} else {
						cellStyle = ui.CalendarUnpaidLeaveStyle
					}
				}
			}
		}

		currentRow = append(currentRow, cellStyle.Render(strconv.Itoa(d)))

		// Wrap after Sunday
		if len(currentRow) == 7 {
			grid = append(grid, currentRow)
			currentRow = nil
		}
	}

	// Pad end of last row
	if len(currentRow) > 0 {
		for len(currentRow) < 7 {
			currentRow = append(currentRow, ui.CalendarEmptyDayStyle.Render(""))
		}
		grid = append(grid, currentRow)
	}

	var renderedRows []string
	for _, row := range grid {
		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}

func renderSummary(data *AttendanceData) string {
	if data == nil {
		return ui.CalendarStatsBoxStyle.Render("No data available.")
	}

	stats := fmt.Sprintf(
		"Summary\n\n"+
			"Working Days: %.1f\n"+
			"Paid Leave:   %.1f\n"+
			"Unpaid Leave: %.1f\n"+
			"OT Hours:     %.1f",
		data.ActualWork,
		data.DaysOffTaken,
		data.UnpaidLeave,
		data.HoursOverTime,
	)

	return ui.CalendarStatsBoxStyle.Render(stats)
}

func RenderCalendar(year int, month int, width int, height int, data *AttendanceData) string {
	header := renderHeader(year, month)
	daysHeader := renderDaysOfWeek()
	grid := renderGrid(year, month, data)
	summary := renderSummary(data)

	// Combine components
	calBlock := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		daysHeader,
		grid,
	)

	cal := lipgloss.JoinHorizontal(lipgloss.Top,
		calBlock,
		lipgloss.NewStyle().MarginTop(2).Render(summary),
	)

	// Center everything if width > our calendar width
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, cal)
}
