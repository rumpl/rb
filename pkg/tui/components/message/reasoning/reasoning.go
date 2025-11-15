package reasoning

import (
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/rumpl/rb/pkg/tui/components/markdown"
	"github.com/rumpl/rb/pkg/tui/components/spinner"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents an assistant reasoning message view
type Component struct {
	message *types.Message
	width   int
	height  int
	spinner spinner.Spinner
}

// New creates a new assistant reasoning message component
func New(msg *types.Message) layout.Model {
	return &Component{
		message: msg,
		width:   80,
		height:  1,
		spinner: spinner.New(spinner.ModeBoth),
	}
}

func (c *Component) Init() tea.Cmd {
	if c.message.Content == "" {
		return c.spinner.Tick()
	}
	return nil
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	if c.message.Content == "" {
		s, cmd := c.spinner.Update(msg)
		c.spinner = s.(spinner.Spinner)
		return c, cmd
	}
	return c, nil
}

func (c *Component) View() string {
	if c.message.Content == "" {
		return c.spinner.View()
	}
	// Render through the markdown renderer to ensure proper wrapping to width
	availableWidth := c.width - 1
	if availableWidth < 10 {
		availableWidth = 10
	}
	rendered, err := markdown.NewRenderer(availableWidth).Render(c.message.Content)
	if err != nil {
		text := "Thinking: " + senderPrefix(c.message.Sender) + c.message.Content
		wrapped := wrapText(text, availableWidth)
		return styles.MutedStyle.Italic(true).Render(wrapped)
	}
	// Strip ANSI from inner rendering so muted style fully applies
	clean := stripANSI(strings.TrimRight(rendered, "\n\r\t "))
	thinkingText := "Thinking: " + senderPrefix(c.message.Sender) + clean
	return styles.MutedStyle.Italic(true).Render(thinkingText)
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

func senderPrefix(sender string) string {
	if sender == "" {
		return ""
	}
	return styles.AgentBadgeStyle.Render("["+sender+"]") + "\n"
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
