package shell

import (
	"fmt"
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/rumpl/rb/pkg/tui/components/markdown"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents a shell output message view
type Component struct {
	message      *types.Message
	width        int
	height       int
	themeManager *styles.Manager
}

// New creates a new shell output message component
func New(msg *types.Message, themeManager *styles.Manager) layout.Model {
	return &Component{
		message:      msg,
		width:        80,
		height:       1,
		themeManager: themeManager,
	}
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	return c, nil
}

func (c *Component) View() string {
	availableWidth := max(c.width-1, 10)
	if rendered, err := markdown.NewRenderer(availableWidth, c.themeManager).Render(fmt.Sprintf("```console\n%s\n```", c.message.Content)); err == nil {
		return strings.TrimRight(rendered, "\n\r\t ")
	}
	wrapped := wrapText(c.message.Content, availableWidth)
	return wrapped
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

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// wrapText wraps text to the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var lines []string
	for _, line := range strings.Split(text, "\n") {
		// Strip ANSI codes to get actual text length
		cleanLine := stripANSI(line)
		for len(cleanLine) > width {
			// Find the last space before width to break at word boundary
			breakPoint := width
			if idx := strings.LastIndex(cleanLine[:width], " "); idx > width/2 {
				breakPoint = idx + 1
			}
			lines = append(lines, line[:breakPoint])
			line = line[breakPoint:]
			cleanLine = cleanLine[breakPoint:]
			// Remove leading spaces from continuation
			line = strings.TrimLeft(line, " ")
			cleanLine = strings.TrimLeft(cleanLine, " ")
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
