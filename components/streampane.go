package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/tinydarkforge/intake/app"
)

// StreamPane is a scrollable viewport that accumulates streamed tokens.
type StreamPane struct {
	vp      viewport.Model
	buf     strings.Builder
	focused bool
}

func NewStreamPane(width, height int) StreamPane {
	vp := viewport.New(width, height)
	vp.Style = app.BorderNormal
	return StreamPane{vp: vp}
}

func (s *StreamPane) SetSize(width, height int) {
	s.vp.Width = width
	s.vp.Height = height
}

func (s *StreamPane) SetFocused(f bool) {
	s.focused = f
	if f {
		s.vp.Style = app.BorderFocused
	} else {
		s.vp.Style = app.BorderNormal
	}
}

func (s *StreamPane) Append(chunk string) {
	s.buf.WriteString(chunk)
	s.vp.SetContent(s.buf.String())
	s.vp.GotoBottom()
}

func (s *StreamPane) SetContent(text string) {
	s.buf.Reset()
	s.buf.WriteString(text)
	s.vp.SetContent(text)
}

func (s *StreamPane) Reset() {
	s.buf.Reset()
	s.vp.SetContent("")
}

func (s *StreamPane) Content() string { return s.buf.String() }

func (s *StreamPane) View() string {
	return lipgloss.NewStyle().Render(s.vp.View())
}

// Viewport exposes the inner model so the parent can call vp.Update(msg).
func (s *StreamPane) Viewport() *viewport.Model { return &s.vp }
