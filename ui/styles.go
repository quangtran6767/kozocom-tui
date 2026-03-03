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
)
