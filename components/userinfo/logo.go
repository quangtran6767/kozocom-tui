package userinfo

import (
	"image/color"
	"sync/atomic"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const logoArt = `
█╗  ██╗███████╗
██║ ██╔╝╚══███╔╝
█████╔╝   ███╔╝
██╔═██╗  ███╔╝
██║  ██╗███████╗
╚═╝  ╚═╝╚══════╝
`

var colorPalette = []color.Color{
	lipgloss.Color("#1a1a9e"), // navy peak
	lipgloss.Color("#3a3ab0"),
	lipgloss.Color("#5a5ac0"),
	lipgloss.Color("#7070b8"),
	lipgloss.Color("#8a8a8a"), // gray valley (com circles)
	lipgloss.Color("#7070b8"),
	lipgloss.Color("#5a5ac0"),
	lipgloss.Color("#3a3ab0"),
}

const logoFPS = time.Second / 3

var lastID int64

func nextLogoID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

type logoTickMsg struct {
	id  int // Use to filter correct instance
	tag int // Prevent fast-forward
}

type logoModel struct {
	colorIndex int
	id         int
	tag        int
}

func newLogo() logoModel {
	return logoModel{
		id: nextLogoID(),
	}
}

func (m logoModel) tick() tea.Cmd {
	// Capture id and tag before into closure before return
	// Do not use m.id directly in closure because it will be changed
	id := m.id
	tag := m.tag

	return tea.Tick(logoFPS, func(t time.Time) tea.Msg {
		return logoTickMsg{
			id:  id,
			tag: tag,
		}
	})
}

func (m logoModel) Init() tea.Cmd {
	return m.tick()
}

func (m logoModel) Update(msg tea.Msg) (logoModel, tea.Cmd) {
	tickMsg, ok := msg.(logoTickMsg)
	if !ok {
		return m, nil
	}

	// Filter: only process message for this instance
	if tickMsg.id > 0 && tickMsg.id != m.id {
		return m, nil
	}

	// Filter: prevent fast-forward
	if tickMsg.tag > 0 && tickMsg.tag != m.tag {
		return m, nil
	}

	m.colorIndex = (m.colorIndex + 1) % len(colorPalette)
	m.tag++
	return m, m.tick()
}

func (m logoModel) View() string {
	color := colorPalette[m.colorIndex]
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(logoArt)
}
