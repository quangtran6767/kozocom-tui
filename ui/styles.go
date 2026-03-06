package ui

import lipgloss "charm.land/lipgloss/v2"

var (
	BaseBg = lipgloss.Color("#0F0F1A")

	// Legends
	LegendBg = lipgloss.Color("#FF8899")
	LegendFg = lipgloss.Color("#1A1A2E")

	// Title
	TitleForeground = lipgloss.Color("#FF8899")

	// Label
	LabelForeground = lipgloss.Color("#888888")

	// Checkin status
	CheckinSuccessColor = lipgloss.Color("#4ADE80")

	// Error
	ErrorForeground = lipgloss.Color("#FF4444")

	// Hint
	HintForeground = lipgloss.Color("#555555")

	// Box
	LoginForeground = lipgloss.Color("#3D3D5C")

	// Border style: rounded corners
	PanelBorder       = lipgloss.RoundedBorder()
	BorderColor       = lipgloss.Color("#3D3D5C")
	ActiveBorderColor = lipgloss.Color("#FF8899")

	// Legend style: display panel like HTML <legend>
	LegendStyle = lipgloss.NewStyle().
		// Background(LegendBg).
		// Foreground(LegendFg).
		Bold(true).
		Padding(0, 1)

	// Spinner style
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(LegendFg)

	// Calendar styles
	CalendarHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(TitleForeground).
				Align(lipgloss.Center)

	CalendarDayHeaderStyle = lipgloss.NewStyle().
				Width(6).
				Align(lipgloss.Center)

	CalendarDayStyle = lipgloss.NewStyle().
				Width(6).
				Height(3).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))

	CalendarTodayStyle = CalendarDayStyle.
				BorderForeground(lipgloss.Color("63")).
				Bold(true)

	CalendarPresentStyle = CalendarDayStyle.
				Background(lipgloss.Color("22")). // Green tint
				Foreground(lipgloss.Color("255"))

	CalendarPaidLeaveStyle = CalendarDayStyle.
				Background(lipgloss.Color("130")). // Yellow/Orange tint
				Foreground(lipgloss.Color("255"))

	CalendarUnpaidLeaveStyle = CalendarDayStyle.
					Background(lipgloss.Color("88")). // Red tint
					Foreground(lipgloss.Color("255"))

	CalendarWeekendStyle = CalendarDayStyle.
				Foreground(lipgloss.Color("240"))

	CalendarEmptyDayStyle = lipgloss.NewStyle().
				Width(6).
				Height(3)

	CalendarStatsBoxStyle = lipgloss.NewStyle().
				Padding(1, 4)

	// Sidebar styles
	SidebarSelectedItemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("205")). // Pinkish
					Bold(true)
	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")) // Gray

	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)
)
