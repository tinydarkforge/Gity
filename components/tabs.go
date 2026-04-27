package components

import (
	"strings"

	"github.com/tinydarkforge/intake/app"
)

var tabLabels = []string{"Tasks", "Create", "Settings"}

// Tabs renders the three top-level navigation tabs.
func Tabs(active app.Screen) string {
	var parts []string
	for i, label := range tabLabels {
		s := app.Screen(i)
		if s == active {
			parts = append(parts, app.StyleTabActive.Render(label))
		} else {
			parts = append(parts, app.StyleTab.Render(label))
		}
	}
	return strings.Join(parts, app.StyleDim.Render("│"))
}
