package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/tinydarkforge/intake/app"
)

// Header renders a one-line banner: "intake  ·  repo/name  ·  model  [status]"
func Header(width int, repo, model, status string) string {
	left := app.StyleHeader.Render(fmt.Sprintf(" intake  ·  %s  ·  %s ", repo, model))
	right := ""
	if status != "" {
		right = app.StyleHeader.Render(fmt.Sprintf(" %s ", status))
	}
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}
	fill := app.StyleHeader.Render(fmt.Sprintf("%*s", gap, ""))
	return lipgloss.JoinHorizontal(lipgloss.Top, left, fill, right)
}
