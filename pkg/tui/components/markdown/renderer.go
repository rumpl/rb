package markdown

import (
	"github.com/charmbracelet/glamour/v2"

	"github.com/rumpl/rb/pkg/tui/styles"
)

func NewRenderer(width int, themeManager *styles.Manager) *glamour.TermRenderer {
	theme := themeManager.GetTheme()
	style := theme.MarkdownStyle

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStyles(style),
	)
	return r
}
