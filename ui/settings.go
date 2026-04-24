package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tinydarkforge/gity/app"
	"github.com/tinydarkforge/gity/components"
	"github.com/tinydarkforge/gity/types"
)

type settingsField int

const (
	fieldRepo settingsField = iota
	fieldModel
	fieldOllamaHost
	fieldTimeout
	fieldMaxTurns
	fieldSound
	fieldDebug
	fieldCount
)

var fieldLabels = []string{
	"GitHub Repo",
	"Ollama Model",
	"Ollama Host",
	"Timeout (sec)",
	"Max Turns",
	"Sound",
	"Debug",
}

type SettingsModel struct {
	cfg     types.Config
	inputs  [fieldCount]textinput.Model
	focus   settingsField
	status  string
	width   int
	height  int
}

func NewSettings(cfg types.Config) SettingsModel {
	m := SettingsModel{cfg: cfg}
	vals := []string{
		cfg.Repo,
		cfg.Model,
		cfg.OllamaHost,
		strconv.Itoa(cfg.TimeoutSec),
		strconv.Itoa(cfg.MaxTurns),
		boolStr(cfg.Sound),
		boolStr(cfg.Debug),
	}
	for i := 0; i < int(fieldCount); i++ {
		ti := textinput.New()
		ti.SetValue(vals[i])
		ti.Width = 40
		m.inputs[i] = ti
	}
	m.inputs[0].Focus()
	return m
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) InterceptsKeys() bool { return true }

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case app.ConfigSavedMsg:
		m.status = "settings saved"

	case app.ErrMsg:
		m.status = "save failed: " + msg.Err.Error()

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			m.inputs[m.focus].Blur()
			m.focus = (m.focus + 1) % fieldCount
			m.inputs[m.focus].Focus()
		case "shift+tab", "up", "k":
			m.inputs[m.focus].Blur()
			if m.focus == 0 {
				m.focus = fieldCount - 1
			} else {
				m.focus--
			}
			m.inputs[m.focus].Focus()
		case "ctrl+s", "enter":
			if msg.String() == "ctrl+s" {
				return m, m.saveCmd()
			}
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m *SettingsModel) applyInputs() {
	m.cfg.Repo = m.inputs[fieldRepo].Value()
	m.cfg.Model = m.inputs[fieldModel].Value()
	m.cfg.OllamaHost = m.inputs[fieldOllamaHost].Value()
	if n, err := strconv.Atoi(m.inputs[fieldTimeout].Value()); err == nil {
		m.cfg.TimeoutSec = n
	}
	if n, err := strconv.Atoi(m.inputs[fieldMaxTurns].Value()); err == nil {
		m.cfg.MaxTurns = n
	}
	m.cfg.Sound = strings.ToLower(m.inputs[fieldSound].Value()) == "true"
	m.cfg.Debug = strings.ToLower(m.inputs[fieldDebug].Value()) == "true"
}

func (m SettingsModel) saveCmd() tea.Cmd {
	m.applyInputs()
	cfg := m.cfg
	return func() tea.Msg {
		if err := types.SaveConfig(cfg); err != nil {
			return app.ErrMsg{Err: err}
		}
		return app.ConfigSavedMsg{}
	}
}

func (m SettingsModel) Config() types.Config {
	m.applyInputs()
	return m.cfg
}

func (m SettingsModel) View() string {
	var lines []string
	lines = append(lines, "",
		app.StyleBold.Render("  Settings"),
		"",
	)
	for i := 0; i < int(fieldCount); i++ {
		label := app.StyleDim.Render(fmt.Sprintf("  %-16s", fieldLabels[i]))
		var input string
		if settingsField(i) == m.focus {
			input = app.StyleAccent.Render("> ") + m.inputs[i].View()
		} else {
			input = "  " + m.inputs[i].View()
		}
		lines = append(lines, label+input)
	}

	if m.status != "" {
		lines = append(lines, "", app.StyleStatusOK.Render("  "+m.status))
	}

	footer := components.FooterHint(m.width,
		app.Global.Back,
		app.Global.Quit,
	)
	lines = append(lines, "", footer)
	return strings.Join(lines, "\n")
}
