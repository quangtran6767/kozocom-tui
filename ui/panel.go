package ui

import "charm.land/lipgloss/v2"

func RenderPanel(title, content string, width, height int) string {
	innerWidth := width - 2
	innerHeight := height - 2

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	// Create border box
	boxStyle := lipgloss.NewStyle().
		Border(PanelBorder).
		BorderForeground(BorderColor).
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
