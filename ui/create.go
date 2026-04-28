package ui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tinydarkforge/intake/app"
	"github.com/tinydarkforge/intake/components"
	"github.com/tinydarkforge/intake/services"
	"github.com/tinydarkforge/intake/types"
)

type createStep int

const (
	stepChooseTemplate createStep = iota
	stepTitle
	stepBrief
	stepAgentTurns
	stepPreview
	stepCreating
	stepDone
)

type CreateModel struct {
	step      createStep
	templates []types.Template
	tmplIdx   int

	titleInput textinput.Model
	briefInput textarea.Model
	answerIn   textinput.Model

	agent     *services.Agent
	turn      int
	questions []string
	draft     types.Draft

	stream   components.StreamPane
	spinner  spinner.Model
	prog     *tea.Program
	cancelFn context.CancelFunc

	gh     *services.GitHub
	ollama *services.OllamaClient
	cfg    types.Config

	statusText string
	statusErr  bool
	createdURL string

	width  int
	height int
}

func NewCreate(gh *services.GitHub, ollama *services.OllamaClient, tmpls []types.Template, cfg types.Config) CreateModel {
	ti := textinput.New()
	ti.Placeholder = "Title (optional — agent will infer from context if blank)"
	ti.CharLimit = 120
	ti.Focus()

	bi := textarea.New()
	bi.Placeholder = "Paste anything: Slack message, error log, PR description, meeting notes…\nThe agent will extract structure from whatever you give it.\n\nctrl+s to start  ·  paste freely  ·  more context = better draft"
	bi.SetHeight(10)

	ai := textinput.New()
	ai.Placeholder = "Answer · enter to confirm · type 'skip' to finalize now…"

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = app.StyleAccent

	return CreateModel{
		templates:  tmpls,
		titleInput: ti,
		briefInput: bi,
		answerIn:   ai,
		spinner:    sp,
		stream:     components.NewStreamPane(80, 20),
		gh:         gh,
		ollama:     ollama,
		cfg:        cfg,
	}
}

// InterceptsKeys returns true whenever a text input is focused so the root
// model skips global shortcut handling and passes all keys here directly.
func (m CreateModel) InterceptsKeys() bool {
	switch m.step {
	case stepTitle, stepBrief, stepAgentTurns:
		return true
	}
	return false
}

func (m CreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m CreateModel) SetProgram(p *tea.Program) CreateModel {
	m.prog = p
	return m
}

func (m CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.stream.SetSize(m.width-4, m.height-12)
		m.briefInput.SetWidth(m.width - 6)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case services.TokenMsg:
		m.stream.Append(msg.Chunk)

	case services.StreamDoneMsg:
		m.stream.Append("")
		// stream preview done; advance to preview confirm step
		m.draft.Body = m.stream.Content()
		m.step = stepPreview

	case app.AgentResponseMsg:
		resp := msg.Resp
		if resp.IsReady() {
			m.draft.Title = resp.Title
			m.draft.Body = resp.Body
			m.step = stepPreview
		} else {
			m.questions = resp.Questions
			m.step = stepAgentTurns
			m.answerIn.Focus()
		}

	case app.AgentFinalMsg:
		m.draft.Title = msg.Resp.Title
		m.draft.Body = msg.Resp.Body
		m.step = stepPreview

	case app.IssueCreatedMsg:
		m.createdURL = msg.URL
		m.step = stepDone

	case app.ErrMsg:
		m.statusText = msg.Err.Error()
		m.statusErr = true
		m.step = stepTitle // go back to editable state

	case tea.KeyMsg:
		switch m.step {
		case stepChooseTemplate:
			cmds = append(cmds, m.handleTemplateKey(msg)...)
		case stepTitle:
			cmds = append(cmds, m.handleTitleKey(msg)...)
		case stepBrief:
			cmds = append(cmds, m.handleBriefKey(msg)...)
		case stepAgentTurns:
			cmds = append(cmds, m.handleAnswerKey(msg)...)
		case stepPreview:
			cmds = append(cmds, m.handlePreviewKey(msg)...)
		}
	}

	// bubble down to inputs
	switch m.step {
	case stepTitle:
		var cmd tea.Cmd
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepBrief:
		var cmd tea.Cmd
		m.briefInput, cmd = m.briefInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepAgentTurns:
		var cmd tea.Cmd
		m.answerIn, cmd = m.answerIn.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *CreateModel) handleTemplateKey(msg tea.KeyMsg) []tea.Cmd {
	switch msg.String() {
	case "j", "down":
		if m.tmplIdx < len(m.templates)-1 {
			m.tmplIdx++
		}
	case "k", "up":
		if m.tmplIdx > 0 {
			m.tmplIdx--
		}
	case "enter":
		m.step = stepTitle
		m.titleInput.Focus()
	}
	return nil
}

func (m *CreateModel) handleTitleKey(msg tea.KeyMsg) []tea.Cmd {
	switch msg.String() {
	case "enter", "tab":
		m.step = stepBrief
		m.briefInput.Focus()
	}
	return nil
}

func (m *CreateModel) handleBriefKey(msg tea.KeyMsg) []tea.Cmd {
	if msg.String() == "ctrl+s" {
		brief := strings.TrimSpace(m.briefInput.Value())
		if brief == "" {
			m.statusText = "paste some context first"
			m.statusErr = true
			return nil
		}
		title := strings.TrimSpace(m.titleInput.Value())
		// title is optional — agent will infer it
		tmpl := m.templates[m.tmplIdx]
		m.draft.Labels = tmpl.Labels
		m.draft.Template = tmpl.Filename

		// rich paste (>150 chars) → skip Q&A, go straight to finalize
		maxTurns := m.cfg.MaxTurns
		autoFinalize := len(brief) > 150
		m.agent = &services.Agent{
			O:        m.ollama,
			Template: tmpl,
			Title:    title,
			Brief:    brief,
			MaxTurns: maxTurns,
			Debug:    m.cfg.Debug,
		}
		m.turn = 1
		m.step = stepAgentTurns
		if autoFinalize {
			// rich paste — skip Q&A, draft immediately
			return []tea.Cmd{m.agentFinalCmd()}
		}
		return []tea.Cmd{m.agentTurnCmd()}
	}
	return nil
}

func (m *CreateModel) handleAnswerKey(msg tea.KeyMsg) []tea.Cmd {
	if msg.String() != "enter" {
		return nil
	}
	ans := strings.TrimSpace(m.answerIn.Value())
	m.answerIn.Reset()

	if strings.ToLower(ans) == "skip" || strings.ToLower(ans) == "do it" {
		return []tea.Cmd{m.agentFinalCmd()}
	}

	if m.agent != nil && len(m.questions) > 0 {
		qIdx := len(m.agent.History)
		q := ""
		if qIdx < len(m.questions) {
			q = m.questions[qIdx]
		}
		m.agent.History = append(m.agent.History, types.Turn{Question: q, Answer: ans})
	}

	// if we've answered all questions, do next turn
	if m.agent != nil && len(m.agent.History) >= len(m.questions) {
		m.turn++
		if m.turn > m.cfg.MaxTurns {
			return []tea.Cmd{m.agentFinalCmd()}
		}
		return []tea.Cmd{m.agentTurnCmd()}
	}
	return nil
}

func (m *CreateModel) handlePreviewKey(msg tea.KeyMsg) []tea.Cmd {
	switch msg.String() {
	case "y":
		m.step = stepCreating
		return []tea.Cmd{m.createIssueCmd()}
	case "r":
		m.step = stepAgentTurns
		return []tea.Cmd{m.agentFinalCmd()}
	case "esc":
		m.Reset()
	}
	return nil
}

func (m *CreateModel) agentTurnCmd() tea.Cmd {
	agent := m.agent
	turn := m.turn
	return func() tea.Msg {
		resp, err := agent.NextTurn(context.Background(), turn)
		if err != nil {
			return app.ErrMsg{Err: err}
		}
		return app.AgentResponseMsg{Resp: resp}
	}
}

func (m *CreateModel) agentFinalCmd() tea.Cmd {
	agent := m.agent
	return func() tea.Msg {
		resp, err := agent.Finalize(context.Background())
		if err != nil {
			return app.ErrMsg{Err: err}
		}
		return app.AgentFinalMsg{Resp: resp}
	}
}

func (m *CreateModel) createIssueCmd() tea.Cmd {
	gh := m.gh
	draft := m.draft
	return func() tea.Msg {
		url, err := gh.Create(context.Background(), draft)
		if err != nil {
			return app.ErrMsg{Err: err}
		}
		return app.IssueCreatedMsg{URL: url, Draft: draft}
	}
}

// Reset returns the model to the template-selection step.
func (m *CreateModel) Reset() {
	m.step = stepChooseTemplate
	m.agent = nil
	m.questions = nil
	m.draft = types.Draft{}
	m.stream.Reset()
	m.titleInput.Reset()
	m.briefInput.Reset()
	m.answerIn.Reset()
	m.statusText = ""
	m.statusErr = false
}

