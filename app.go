package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/quangtran6767/kozocom-tui/components/auth"
	"github.com/quangtran6767/kozocom-tui/components/content"
	"github.com/quangtran6767/kozocom-tui/components/footer"
	"github.com/quangtran6767/kozocom-tui/components/sidebar"
	"github.com/quangtran6767/kozocom-tui/components/userinfo"
	"github.com/quangtran6767/kozocom-tui/keybindings"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/services"
	"github.com/quangtran6767/kozocom-tui/ui"
)

type PanelID int

const (
	PanelSidebar PanelID = iota
	PanelContent
	PanelFooter
)

type AppState int

const (
	StateAuth AppState = iota
	StateMain
)

type appModel struct {
	state       AppState
	auth        auth.Model
	activePanel PanelID
	userinfo    userinfo.Model
	sidebar     sidebar.Model
	content     content.Model
	footer      footer.Model
	help        help.Model
	width       int
	height      int
	ready       bool
}

func newAppModel() appModel {
	return appModel{
		state:    StateAuth,
		auth:     auth.New(),
		userinfo: userinfo.New(),
		sidebar:  sidebar.New(),
		content:  content.New(),
		footer:   footer.New(),
		help:     help.New(),
	}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		m.auth.Init(),
		tea.RequestWindowSize,
	)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.auth.SetSize(msg.Width, msg.Height)

		dims := ui.CalculateLayout(m.width, m.height)
		m.userinfo.SetSize(dims.SidebarWidth, dims.SidebarUserInfoHeight)
		m.sidebar.SetSize(dims.SidebarWidth, dims.SidebarHeight)
		m.content.SetSize(dims.ContentWidth, dims.TopHeight)
		m.footer.SetSize(dims.ContentWidth, dims.BottomHeight)
	}

	switch m.state {
	case StateAuth:
		return m.updateAuth(msg)
	case StateMain:
		return m.updateMain(msg)
	}

	return m, nil
}

func (m appModel) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	switch m.state {
	case StateAuth:
		v := tea.NewView(m.auth.View())
		v.AltScreen = true
		return v
	case StateMain:
		dims := ui.CalculateLayout(m.width, m.height)

		userInfoPanel := ui.RenderPanel(
			"User Info",
			m.userinfo.View(),
			dims.SidebarWidth,
			dims.SidebarUserInfoHeight,
			false, // User info panel doesn't need focus
		)

		sidebarPanel := ui.RenderPanel(
			"[1] Sidebar",
			m.sidebar.View(),
			dims.SidebarWidth,
			dims.SidebarHeight,
			m.activePanel == PanelSidebar,
		)

		contentPanel := ui.RenderPanel(
			"[2] Content",
			m.content.View(),
			dims.ContentWidth,
			dims.TopHeight,
			m.activePanel == PanelContent,
		)

		footerPanel := ui.RenderPanel(
			"[3] Footer",
			m.footer.View(),
			dims.ContentWidth,
			dims.BottomHeight,
			m.activePanel == PanelFooter,
		)

		layout := ui.RenderLayout(sidebarPanel, userInfoPanel, contentPanel, footerPanel)

		var panelBindings []key.Binding
		switch m.activePanel {
		case PanelSidebar:
			panelBindings = m.sidebar.PanelBindings()
		case PanelContent:
			panelBindings = m.content.PanelBindings()
		case PanelFooter:
			panelBindings = m.footer.PanelBindings()
		}

		km := keybindings.NewDynamicKeyMap(panelBindings)
		helpBar := m.help.View(km)

		final := lipgloss.JoinVertical(lipgloss.Left, layout, helpBar)

		v := tea.NewView(final)
		v.AltScreen = true
		return v
	}

	return tea.NewView("")
}

func (m *appModel) switchPanel(p PanelID) {
	m.sidebar.Blur()
	m.content.Blur()
	m.footer.Blur()

	m.activePanel = p

	switch p {
	case PanelSidebar:
		m.sidebar.Focus()
	case PanelContent:
		m.content.Focus()
	case PanelFooter:
		m.footer.Focus()
	}
}

func (m appModel) updateAuth(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.auth, cmd = m.auth.Update(msg)

	if m.auth.IsDone() {
		m.state = StateMain
		m.userinfo.SetUserInfo(m.auth.Email(), m.auth.UserID())
		return m, tea.Batch(
			m.userinfo.Init(),
			m.content.SetToken(m.auth.Token()),
			services.CheckinTodayStatus(m.auth.Token()),
		)
	}
	return m, cmd
}

func (m appModel) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.switchPanel(PanelSidebar)
		case "2":
			m.switchPanel(PanelContent)
		case "3":
			m.switchPanel(PanelFooter)
		case "C":
			if !m.userinfo.IsCheckedIn() && !m.userinfo.IsLoading() {
				m.userinfo.SetCheckinLoading(true)
				return m, tea.Batch(
					m.userinfo.Init(), // start spinner
					services.Checkin(m.auth.Token()),
				)
			}
		}
	case messages.CheckinStatusMsg:
		m.userinfo.SetCheckinStatus(msg.IsCheckedIn)
		return m, nil
	case messages.CheckinStatusFailMsg:
		return m, nil
	case messages.CheckinSuccessMsg:
		m.userinfo.SetCheckinLoading(false)
		m.userinfo.SetCheckinStatus(true)
		return m, nil
	case messages.CheckinFailMsg:
		m.userinfo.SetCheckinLoading(false)
		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.userinfo, cmd = m.userinfo.Update(msg)
	cmds = append(cmds, cmd)

	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
