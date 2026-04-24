package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"

	"github.com/tinydarkforge/gity/app"
)

// FooterHint renders a row of key hint pairs: "key  desc  key  desc …"
func FooterHint(width int, bindings ...key.Binding) string {
	var parts []string
	for _, b := range bindings {
		if !b.Enabled() {
			continue
		}
		k := app.StyleAccent.Render(b.Help().Key)
		d := app.StyleDim.Render(b.Help().Desc)
		parts = append(parts, k+" "+d)
	}
	hint := strings.Join(parts, app.StyleDim.Render("  "))
	return app.StyleFooter.Width(width).Render(hint)
}
