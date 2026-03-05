package keybindings

import "charm.land/bubbles/v2/key"

var (
	SwitchPanel1 = key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "sidebar"),
	)
	SwitchPanel2 = key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "content"),
	)
	SwitchPanel3 = key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "footer"),
	)
	Quit = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
	Checkin = key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "checkin"),
	)
)

func GlobalBindings() []key.Binding {
	return []key.Binding{SwitchPanel1, SwitchPanel2, SwitchPanel3, Checkin, Quit}
}

// PanelKeyProvider is an interface each panel must implement
// To provide unique keybindings for the help bar
// When activePanel changes, the app calls the PanelBindings() of that panel
// TO rebuild DynamicKeyMap
type PanelKeyProvider interface {
	PanelBindings() []key.Binding
}

type DynamicKeyMap struct {
	global []key.Binding
	panel  []key.Binding
}

func NewDynamicKeyMap(panel []key.Binding) *DynamicKeyMap {
	return &DynamicKeyMap{
		global: GlobalBindings(),
		panel:  panel,
	}
}

// ShortHelp render 1 line at the help bar
func (k DynamicKeyMap) ShortHelp() []key.Binding {
	all := make([]key.Binding, 0, len(k.panel)+len(k.global))
	all = append(all, k.global...)
	all = append(all, k.panel...)
	return all
}

func (k DynamicKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.panel,  // row1: panel-specifc
		k.global, // row2: global
	}
}
