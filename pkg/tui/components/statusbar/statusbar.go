package statusbar

import (
	"charm.land/lipgloss/v2"

	"github.com/rumpl/rb/pkg/tui/core"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/version"
)

// StatusBar represents the status bar component that displays key bindings help
type StatusBar struct {
	width int
	help  core.KeyMapHelp
}

// New creates a new StatusBar instance
func New(help core.KeyMapHelp) StatusBar {
	return StatusBar{
		help: help,
	}
}

// SetWidth sets the width of the status bar
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// View renders the status bar
func (s *StatusBar) View() string {
	versionText := styles.MutedStyle.Render("rb " + version.Version)

	return styles.BaseStyle.
		Width(s.width).
		PaddingLeft(1).
		PaddingRight(1).
		Align(lipgloss.Right).
		Render(versionText)
}
