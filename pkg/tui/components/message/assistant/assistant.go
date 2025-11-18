package assistant

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/rumpl/rb/pkg/tui/components/markdown"
	"github.com/rumpl/rb/pkg/tui/components/spinner"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents an assistant message view
type Component struct {
	message      *types.Message
	width        int
	height       int
	spinner      spinner.Spinner
	themeManager *styles.Manager
}

// New creates a new assistant message component
func New(msg *types.Message, themeManager *styles.Manager) layout.Model {
	return &Component{
		message:      msg,
		width:        80,
		height:       1,
		spinner:      spinner.New(spinner.ModeBoth, themeManager),
		themeManager: themeManager,
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

	s := c.themeManager.GetTheme()
	// Account for border (1 char) + padding (2 chars left + 2 chars right) = 5 chars total
	// But we also need to account for the border itself, so available width is width - border - padding
	availableWidth := max(
		// 1 for border, 4 for padding (2 left + 2 right)
		c.width-1-4, 10)
	rendered, err := markdown.NewRenderer(availableWidth, c.themeManager).Render(c.message.Content)
	if err != nil {
		// Fallback: wrap content manually if markdown renderer fails
		wrapped := wrapText(c.senderPrefix(c.message.Sender)+c.message.Content, availableWidth)
		return s.AssistantMessageBorderStyle.Width(c.width).Render(wrapped)
	}

	content := c.senderPrefix(c.message.Sender) + strings.TrimRight(rendered, "\n\r\t ")
	// Constrain the content to available width to ensure it doesn't break layout
	content = constrainWidth(content, availableWidth)
	return s.AssistantMessageBorderStyle.Width(c.width).Render(content)
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
	var content string
	s := c.themeManager.GetTheme()
	if c.message.Content == "" {
		content = c.spinner.View()
	} else {
		availableWidth := max(width-1-4, 10)
		rendered, err := markdown.NewRenderer(availableWidth, c.themeManager).Render(c.message.Content)
		if err != nil {
			wrapped := wrapText(c.senderPrefix(c.message.Sender)+c.message.Content, availableWidth)
			content = s.AssistantMessageBorderStyle.Width(width).Render(wrapped)
		} else {
			result := c.senderPrefix(c.message.Sender) + strings.TrimRight(rendered, "\n\r\t ")
			result = constrainWidth(result, availableWidth)
			content = s.AssistantMessageBorderStyle.Width(width).Render(result)
		}
	}
	return strings.Count(content, "\n") + 1
}

func (c *Component) SetMessage(msg *types.Message) {
	c.message = msg
}

func (c *Component) senderPrefix(sender string) string {
	if sender == "" {
		return ""
	}
	theme := c.themeManager.GetTheme()
	return theme.AgentBadgeStyle.Render("["+sender+"]") + "\n"
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

// constrainWidth ensures that each line of content doesn't exceed the specified width
func constrainWidth(content string, maxWidth int) string {
	if maxWidth <= 0 {
		return content
	}

	var lines []string
	for _, line := range strings.Split(content, "\n") {
		// Use lipgloss to get the actual display width (accounting for ANSI codes)
		lineWidth := lipgloss.Width(line)
		if lineWidth > maxWidth {
			// Line is too wide, need to wrap it
			// Try to break at word boundaries first
			words := strings.Fields(line)
			if len(words) == 0 {
				// No words, just truncate
				lines = append(lines, truncateLine(line, maxWidth))
				continue
			}

			var currentLine strings.Builder
			currentWidth := 0

			for _, word := range words {
				wordWidth := lipgloss.Width(word)
				spaceWidth := 1 // space character width

				switch {
				case currentWidth == 0:
					// First word on line
					currentLine.WriteString(word)
					currentWidth = wordWidth
				case currentWidth+spaceWidth+wordWidth <= maxWidth:
					// Word fits on current line
					currentLine.WriteString(" " + word)
					currentWidth += spaceWidth + wordWidth
				default:
					// Word doesn't fit, start new line
					lines = append(lines, currentLine.String())
					currentLine.Reset()
					currentLine.WriteString(word)
					currentWidth = wordWidth
				}
			}

			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
			}
		} else {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

// truncateLine truncates a line to fit within maxWidth, preserving ANSI codes
func truncateLine(line string, maxWidth int) string {
	if lipgloss.Width(line) <= maxWidth {
		return line
	}

	// Simple truncation - in practice, we'd want to preserve ANSI codes
	// For now, just truncate and add ellipsis if needed
	truncated := ""
	for _, r := range line {
		if lipgloss.Width(truncated+string(r)) > maxWidth-3 {
			return truncated + "..."
		}
		truncated += string(r)
	}
	return truncated
}

func stripANSI(s string) string {
	// Basic ANSI escape sequence removal
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			i++
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}
