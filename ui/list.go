package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tinydarkforge/intake/app"
	"github.com/tinydarkforge/intake/components"
	"github.com/tinydarkforge/intake/services"
	"github.com/tinydarkforge/intake/types"
)

type ListModel struct {
	gh      *services.GitHub
	issues  []types.Issue
	cursor  int
	state   string // "open" | "closed"
	filter  string
	loading bool
	errMsg  string
	width   int
	height  int
}

func NewList(gh *services.GitHub) ListModel {
	return ListModel{gh: gh, state: "open", loading: true}
}

func (m ListModel) Init() tea.Cmd {
	return m.fetchCmd()
}

func (m ListModel) fetchCmd() tea.Cmd {
	return func() tea.Msg {
		issues, err := m.gh.List(context.Background(), m.state, 30)
		if err != nil {
			return app.ErrMsg{Err: err}
		}
		return app.IssuesLoadedMsg{Issues: issues, State: m.state}
	}
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case app.IssuesLoadedMsg:
		m.issues = msg.Issues
		m.state = msg.State
		m.loading = false
		m.errMsg = ""
		m.cursor = 0

	case app.ErrMsg:
		m.loading = false
		m.errMsg = msg.Err.Error()

	case app.CursorToMsg:
		if msg.Row < len(m.visible()) {
			m.cursor = msg.Row
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, app.List.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, app.List.Down):
			if m.cursor < len(m.visible())-1 {
				m.cursor++
			}
		case key.Matches(msg, app.List.Toggle):
			if m.state == "open" {
				m.state = "closed"
			} else {
				m.state = "open"
			}
			m.loading = true
			return m, m.fetchCmd()
		case key.Matches(msg, app.List.Refresh):
			m.loading = true
			return m, m.fetchCmd()
		case key.Matches(msg, app.List.Open):
			vis := m.visible()
			if m.cursor < len(vis) {
				num := vis[m.cursor].Number
				return m, func() tea.Msg { return app.OpenIssueMsg{Number: num} }
			}
		}
	}
	return m, nil
}

func (m ListModel) SelectedIssue() *types.Issue {
	vis := m.visible()
	if m.cursor < len(vis) {
		return &vis[m.cursor]
	}
	return nil
}

// State returns the current filter state ("open" or "closed").
func (m ListModel) State() string { return m.state }

// Count returns the number of visible issues.
func (m ListModel) Count() int { return len(m.visible()) }

// SelectedIssueText returns a short status-bar friendly string for the
// currently highlighted issue, e.g. "#4782 QA: Lifecycle modal".
func (m ListModel) SelectedIssueText() string {
	vis := m.visible()
	if m.cursor < len(vis) {
		iss := vis[m.cursor]
		return fmt.Sprintf("#%d %s", iss.Number, iss.Title)
	}
	return ""
}

func (m ListModel) visible() []types.Issue {
	if m.filter == "" {
		return m.issues
	}
	q := strings.ToLower(m.filter)
	var out []types.Issue
	for _, i := range m.issues {
		if strings.Contains(strings.ToLower(i.Title), q) {
			out = append(out, i)
		}
	}
	return out
}

func (m ListModel) View() string {
	if m.loading {
		return app.StyleDim.Render("  loading issues…")
	}

	vis := m.visible()

	header := app.StyleBold.Render(fmt.Sprintf("  Issues (%s) — %d", m.state, len(vis)))
	var rows []string
	rows = append(rows, header, "")

	for i, issue := range vis {
		numStr := fmt.Sprintf("#%-5d", issue.Number)
		num := app.StyleDim.Render(numStr)

		// Render labels first so we know their width
		var labelParts []string
		for _, l := range issue.Labels {
			labelParts = append(labelParts, app.StyleLabel.Render(l.Name))
		}
		labelStr := strings.Join(labelParts, " ")
		labelW := lipgloss.Width(labelStr)

		// Fixed columns: 2 (left pad) + 7 (num) + 2 (gap) + 2 (gap before labels) = 13
		const fixedCols = 13
		labelSep := 0
		if labelW > 0 {
			labelSep = 2
		}
		available := m.width - fixedCols - labelW - labelSep
		if available < 10 {
			available = 10
		}

		// Rune-aware truncation
		title := issue.Title
		runes := []rune(title)
		if len(runes) > available {
			runes = runes[:available-1]
			title = string(runes) + "…"
		}
		// Pad title to fill available space so labels align at a fixed column
		title = fmt.Sprintf("%-*s", available, title)

		var row string
		if labelW > 0 {
			row = fmt.Sprintf("  %s  %s  %s", num, title, labelStr)
		} else {
			row = fmt.Sprintf("  %s  %s", num, title)
		}

		if i == m.cursor {
			row = app.StyleListItemSelected.Width(m.width).Render(row)
		}
		rows = append(rows, row)
	}

	if len(vis) == 0 {
		rows = append(rows, app.StyleDim.Render("  no issues found"))
	}

	content := strings.Join(rows, "\n")
	footer := components.FooterHint(m.width,
		app.List.Up, app.List.Down, app.List.Open,
		app.List.Toggle, app.List.Refresh,
		app.Global.Create, app.Global.Quit,
	)

	_ = lipgloss.NewStyle() // ensure import used
	return lipgloss.JoinVertical(lipgloss.Left, content, footer)
}

// NCLines returns exactly h content lines each w visible characters wide for
// use inside the Norton Commander frame renderer.
func (m ListModel) NCLines(w, h int) []string {
	vis := m.visible()

	// Column widths: 1 space + 6 num + 2 gap = 9 fixed; labels get 10; title fills rest.
	const numW = 6    // "#4783 "
	const labW = 10   // truncated label column
	const gutL = 2    // left gutter spaces
	const gutM = 1    // gap between num and title
	const gutR = 1    // gap before labels
	fixedW := gutL + numW + gutM + labW + gutR
	titleW := w - fixedW
	if titleW < 8 {
		titleW = 8
	}

	var lines []string

	if m.loading {
		lines = append(lines, styleListLine(truncPad("  loading issues…", w), false))
	} else if m.errMsg != "" {
		lines = append(lines, styleListLine(truncPad("  error: "+m.errMsg, w), false))
	}

	for i, issue := range vis {
		if len(lines) >= h {
			break
		}

		// cursor indicator
		indicator := " "
		if i == m.cursor {
			indicator = "▶"
		}

		numStr := truncPad(fmt.Sprintf("#%d", issue.Number), numW)

		titleStr := truncPad(issue.Title, titleW)

		// Collapse labels to a single truncated string.
		var labelNames []string
		for _, l := range issue.Labels {
			labelNames = append(labelNames, l.Name)
		}
		labStr := truncPad(strings.Join(labelNames, ","), labW)

		raw := indicator + " " + numStr + " " + titleStr + " " + labStr
		lines = append(lines, styleListLine(raw, i == m.cursor))
	}

	if len(vis) == 0 && !m.loading {
		msg := truncPad("  no issues", w)
		lines = append(lines, styleListLine(msg, false))
	}

	for len(lines) < h {
		lines = append(lines, styleListLine(strings.Repeat(" ", w), false))
	}
	return lines
}

func styleListLine(line string, selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Background(app.NCSelected).
			Foreground(app.NCSelFg).
			Render(line)
	}
	return lipgloss.NewStyle().
		Background(app.NCBg).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(line)
}
