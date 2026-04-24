package app

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var noColor = os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"

// ---- NC palette ----------------------------------------------------------------

var (
	NCBg          = lipgloss.Color("#000087") // dark blue panel background
	NCSelected    = lipgloss.Color("#00AAAA") // cyan selected row background
	NCSelFg       = lipgloss.Color("#000000") // selected row foreground
	NCBorder      = lipgloss.Color("#FFFFFF") // active panel border
	NCBorderDim   = lipgloss.Color("#444466") // inactive panel border
	NCTitleActive = lipgloss.Color("#008888") // active panel title bar background
	NCTitleDim    = lipgloss.Color("#000066") // inactive panel title bar background
	NCFnBg        = lipgloss.Color("#00AAAA") // F-key button background
	NCFnFg        = lipgloss.Color("#000000") // F-key number foreground
	NCFnDescFg    = lipgloss.Color("#FFFFFF") // F-key description foreground
	NCStatus      = lipgloss.Color("#000087") // status bar background
	NCStatusFg    = lipgloss.Color("#00FFFF") // status bar foreground
	NCColHeader   = lipgloss.Color("#00AAAA") // column header background
)

// ---- Legacy palette (kept for overlay screens) ---------------------------------

var (
	colorAccent  = lipgloss.AdaptiveColor{Light: "#7B2D8B", Dark: "#C678DD"}
	colorFg      = lipgloss.AdaptiveColor{Light: "#282828", Dark: "#ABB2BF"}
	colorFgDim   = lipgloss.AdaptiveColor{Light: "#7C7C7C", Dark: "#5C6370"}
	colorBg      = lipgloss.AdaptiveColor{Light: "#F9F9F9", Dark: "#282C34"}
	colorBorder  = lipgloss.AdaptiveColor{Light: "#D4D4D4", Dark: "#3E4451"}
	colorFocus   = lipgloss.AdaptiveColor{Light: "#7B2D8B", Dark: "#C678DD"}
	colorSuccess = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#98C379"}
	colorWarning = lipgloss.AdaptiveColor{Light: "#E65100", Dark: "#E5C07B"}
	colorError   = lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#E06C75"}
	colorBlue    = lipgloss.AdaptiveColor{Light: "#1565C0", Dark: "#61AFEF"}
)

// ---- Style helpers -------------------------------------------------------------

var (
	StyleBase = lipgloss.NewStyle().
			Foreground(colorFg)

	StyleDim = lipgloss.NewStyle().
			Foreground(colorFgDim)

	StyleBold = lipgloss.NewStyle().
			Bold(true)

	StyleAccent = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess)

	StyleError = lipgloss.NewStyle().
			Foreground(colorError)

	StyleWarning = lipgloss.NewStyle().
			Foreground(colorWarning)

	StyleBlue = lipgloss.NewStyle().
			Foreground(colorBlue)

	// Panel borders
	BorderNormal = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	BorderFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorFocus)

	// Header
	StyleHeader = lipgloss.NewStyle().
			Background(colorAccent).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)

	// Footer / status
	StyleFooter = lipgloss.NewStyle().
			Foreground(colorFgDim).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder)

	StyleStatusOK  = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	StyleStatusErr = lipgloss.NewStyle().Foreground(colorError).Bold(true)

	// Tabs
	StyleTab = lipgloss.NewStyle().
			Foreground(colorFgDim).
			Padding(0, 2)

	StyleTabActive = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			Padding(0, 2).
			Underline(true)

	// List items
	StyleListItem = lipgloss.NewStyle().
			Foreground(colorFg).
			Padding(0, 1)

	StyleListItemSelected = lipgloss.NewStyle().
				Foreground(colorBg).
				Background(colorFocus).
				Padding(0, 1)

	// Label badge
	StyleLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorBlue).
			Padding(0, 1)
)

func init() {
	if noColor {
		disableColors()
	}
}

func disableColors() {
	colorAccent = lipgloss.AdaptiveColor{}
	colorFg = lipgloss.AdaptiveColor{}
	colorFgDim = lipgloss.AdaptiveColor{}
	colorBg = lipgloss.AdaptiveColor{}
	colorBorder = lipgloss.AdaptiveColor{}
	colorFocus = lipgloss.AdaptiveColor{}
	colorSuccess = lipgloss.AdaptiveColor{}
	colorWarning = lipgloss.AdaptiveColor{}
	colorError = lipgloss.AdaptiveColor{}
	colorBlue = lipgloss.AdaptiveColor{}
}
