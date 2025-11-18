package statusbar

import (
	"charm.land/lipgloss/v2"

	"github.com/rumpl/rb/pkg/tui/core"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/version"
)

// StatusBar represents the status bar component that displays key bindings help
type StatusBar struct {
	width        int
	help         core.KeyMapHelp
	themeManager *styles.Manager
}

// New creates a new StatusBar instance
func New(help core.KeyMapHelp, themeManager *styles.Manager) StatusBar {
	return StatusBar{
		help:         help,
		themeManager: themeManager,
	}
}

// SetWidth sets the width of the status bar
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// View renders the status bar
func (s *StatusBar) View() string {
	theme := s.themeManager.GetTheme()
	versionText := theme.MutedStyle.Render("rb " + version.Version)

	return theme.BaseStyle.
		Width(s.width).
		PaddingLeft(1).
		PaddingRight(1).
		Align(lipgloss.Right).
		Render(versionText)
}
