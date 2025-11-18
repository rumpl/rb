package todotool

import (
	"charm.land/lipgloss/v2"

	"github.com/rumpl/rb/pkg/tui/styles"
)

func renderTodoIcon(status string, theme *styles.Theme) (string, lipgloss.Style) {
	switch status {
	case "pending":
		return "◯", theme.PendingStyle
	case "in-progress":
		return "◕", theme.InProgressStyle
	case "completed":
		return "✓", theme.MutedStyle
	default:
		return "?", theme.BaseStyle
	}
}
