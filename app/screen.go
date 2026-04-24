package app

// Screen identifies which top-level view is active.
type Screen int

const (
	ScreenNC       Screen = iota // two-pane Norton Commander mode (default)
	ScreenCreate                 // full-screen overlay: create issue
	ScreenDetail                 // full-screen overlay: issue detail
	ScreenSettings               // full-screen overlay: settings
)
