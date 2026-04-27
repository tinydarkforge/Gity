package components

import (
	"github.com/tinydarkforge/intake/app"
)

// StatusBar renders a transient status line (success, error, or empty).
func StatusBar(width int, text string, isError bool) string {
	if text == "" {
		return ""
	}
	style := app.StyleStatusOK
	prefix := "✓ "
	if isError {
		style = app.StyleStatusErr
		prefix = "✗ "
	}
	msg := style.Render(prefix + text)
	return app.StyleFooter.Width(width).Render(msg)
}
