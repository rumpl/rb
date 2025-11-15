package welcome

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/rumpl/rb/pkg/tui/components/markdown"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents a welcome message view
type Component struct {
	message *types.Message
	width   int
	height  int
}

// New creates a new welcome message component
func New(msg *types.Message) layout.Model {
	return &Component{
		message: msg,
		width:   80,
		height:  1,
	}
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	return c, nil
}

func (c *Component) View() string {
	// Render welcome message with a distinct style
	availableWidth := max(c.width-1, 10)
	rendered, err := markdown.NewRenderer(availableWidth).Render(c.message.Content)
	if err != nil {
		wrapped := wrapText(c.message.Content, availableWidth)
		return styles.MutedStyle.Render(wrapped)
	}
	return styles.MutedStyle.Render(strings.TrimRight(rendered, "\n\r\t "))
}

func (c *Component) SetSize(width, height int) tea.Cmd {
	c.width = width
	c.height = height
	return nil
}

func (c *Component) GetSize() (width, height int) {
	return c.width, c.height
}

func (c *Component) Height(width int) int {
	content := c.View()
	return strings.Count(content, "\n") + 1
}

func (c *Component) SetMessage(msg *types.Message) {
	c.message = msg
}

// wrapText wraps text to the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var lines []string
	for _, line := range strings.Split(text, "\n") {
		for len(line) > width {
			// Find the last space before width to break at word boundary
			breakPoint := width
			if idx := strings.LastIndex(line[:width], " "); idx > width/2 {
				breakPoint = idx + 1
			}
			lines = append(lines, line[:breakPoint])
			line = line[breakPoint:]
			// Remove leading spaces from continuation
			line = strings.TrimLeft(line, " ")
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
