package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/components/content"
	"github.com/quangtran6767/kozocom-tui/components/footer"
	"github.com/quangtran6767/kozocom-tui/components/sidebar"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type appModel struct {
	sidebar sidebar.Model
	content content.Model
	footer  footer.Model
	width   int
	height  int
	ready   bool
}

func newAppModel() appModel {
	return appModel{
		sidebar: sidebar.New(),
		content: content.New(),
		footer:  footer.New(),
	}
}

func (m appModel) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.RequestWindowSize()
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		dims := ui.CalculateLayout(m.width, m.height)
		m.sidebar.SetSize(dims.SidebarWidth, dims.SidebarHeight)
		m.content.SetSize(dims.ContentWidth, dims.TopHeight)
		m.footer.SetSize(dims.ContentWidth, dims.BottomHeight)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	dims := ui.CalculateLayout(m.width, m.height)

	sidebarPanel := ui.RenderPanel(
		"Sidebar",
		m.sidebar.View(),
		dims.SidebarWidth,
		dims.SidebarHeight,
	)

	contentPanel := ui.RenderPanel(
		"Content",
		m.content.View(),
		dims.ContentWidth,
		dims.TopHeight,
	)

	footerPanel := ui.RenderPanel(
		"Footer",
		m.footer.View(),
		dims.ContentWidth,
		dims.BottomHeight,
	)

	layout := ui.RenderLayout(sidebarPanel, contentPanel, footerPanel)

	v := tea.NewView(layout)
	v.AltScreen = true

	return v
}
