package ui

import "charm.land/lipgloss/v2"

func RenderPanel(title, content string, width, height int, focused bool) string {
	innerWidth := width - 2
	innerHeight := height - 2

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	// Decide border color based on which panel being focused on
	borderColor := BorderColor
	if focused {
		borderColor = ActiveBorderColor
	}

	// Create border box
	boxStyle := lipgloss.NewStyle().
		Border(PanelBorder).
		BorderForeground(borderColor).
		Width(width).
		Height(innerHeight)

	box := boxStyle.Render(content)

	if title == "" {
		return box
	}

	// Render legend text
	legend := LegendStyle.Render(title)

	// Use compositor to layer legend on top of box
	boxLayer := lipgloss.NewLayer(box)
	legendLayer := lipgloss.NewLayer(legend).X(2).Y(0)

	comp := lipgloss.NewCompositor(boxLayer, legendLayer)

	return comp.Render()
}
