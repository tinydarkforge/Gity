package app

import "github.com/charmbracelet/bubbles/key"

// GlobalKeys are active on every screen.
type GlobalKeys struct {
	Create      key.Binding
	Settings    key.Binding
	Help        key.Binding
	Quit        key.Binding
	Back        key.Binding
	SwitchPane  key.Binding
}

var Global = GlobalKeys{
	Create:      key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "create")),
	Settings:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "settings")),
	Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Back:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	SwitchPane:  key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
}

// ListKeys are active in the issue list screen.
type ListKeys struct {
	Up      key.Binding
	Down    key.Binding
	Open    key.Binding
	Toggle  key.Binding
	Filter  key.Binding
	Refresh key.Binding
}

var List = ListKeys{
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Open:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
	Toggle:  key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "toggle state")),
	Filter:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
}

// CreateKeys are active in the create screen.
type CreateKeys struct {
	Next      key.Binding
	Prev      key.Binding
	Submit    key.Binding
	Cancel    key.Binding
	EditTitle key.Binding
	Confirm   key.Binding
	Regen     key.Binding
}

var Create = CreateKeys{
	Next:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
	Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
	Submit:  key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "submit")),
	Cancel:  key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "cancel")),
	Confirm: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm create")),
	Regen:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "regenerate")),
}

// DetailKeys are active in the detail screen.
type DetailKeys struct {
	Up      key.Binding
	Down    key.Binding
	Comment key.Binding
	Edit    key.Binding
	Assign  key.Binding
	Close   key.Binding
	Reopen  key.Binding
	Save    key.Binding
}

var Detail = DetailKeys{
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	Comment: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "comment")),
	Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Assign:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "assign me")),
	Close:   key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "close")),
	Reopen:  key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "reopen")),
	Save:    key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
}
