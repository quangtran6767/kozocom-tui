package ui

import lipgloss "charm.land/lipgloss/v2"

var (
	// Colors
	BorderColor = lipgloss.Color("#FF8899")
	LegendBg    = lipgloss.Color("#FF8899")
	LegendFg    = lipgloss.Color("#1A1A2E")
	BaseBg      = lipgloss.Color("#0F0F1A")

	// Border style: rounded corners
	PanelBorder = lipgloss.RoundedBorder()

	// Legend style: display panel like HTML <legend>
	LegendStyle = lipgloss.NewStyle().
			Background(LegendBg).
			Foreground(LegendFg).
			Bold(true).
			Padding(0, 1)
)
