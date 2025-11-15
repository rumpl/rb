package message

import (
	"fmt"
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/rumpl/rb/pkg/tui/components/markdown"
	"github.com/rumpl/rb/pkg/tui/components/spinner"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Model represents a view that can render a message
type Model interface {
	layout.Model
	layout.Sizeable
	SetMessage(msg *types.Message)
}

// messageModel implements Model
type messageModel struct {
	message *types.Message
	width   int
	height  int
	focused bool
	spinner spinner.Spinner
}

// New creates a new message view
func New(msg *types.Message) *messageModel {
	return &messageModel{
		message: msg,
		width:   80, // Default width
		height:  1,  // Will be calculated
		focused: false,
		spinner: spinner.New(spinner.ModeBoth),
	}
}

// Bubble Tea Model methods

// Init initializes the message view
func (mv *messageModel) Init() tea.Cmd {
	if mv.message.Type == types.MessageTypeSpinner {
		return mv.spinner.Tick()
	}
	return nil
}

func (mv *messageModel) SetMessage(msg *types.Message) {
	mv.message = msg
}

// Update handles messages and updates the message view state
func (mv *messageModel) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	if mv.message.Type == types.MessageTypeSpinner {
		s, cmd := mv.spinner.Update(msg)
		mv.spinner = s.(spinner.Spinner)
		return mv, cmd
	}
	return mv, nil
}

// View renders the message view
func (mv *messageModel) View() string {
	return mv.Render(mv.width)
}

// Render renders the message view content
func (mv *messageModel) Render(width int) string {
	msg := mv.message
	switch msg.Type {
	case types.MessageTypeSpinner:
		return mv.spinner.View()
	case types.MessageTypeUser:
		return styles.UserMessageBorderStyle.Width(width).Render(msg.Content)
	case types.MessageTypeAssistant:
		if msg.Content == "" {
			return mv.spinner.View()
		}

		// Account for border (1 char) + padding (2 chars left + 2 chars right) = 5 chars total
		// But we also need to account for the border itself, so available width is width - border - padding
		availableWidth := width - 1 - 4 // 1 for border, 4 for padding (2 left + 2 right)
		if availableWidth < 10 {
			availableWidth = 10 // Minimum readable width
		}
		rendered, err := markdown.NewRenderer(availableWidth).Render(msg.Content)
		if err != nil {
			// Fallback: wrap content manually if markdown renderer fails
			wrapped := wrapText(senderPrefix(msg.Sender)+msg.Content, availableWidth)
			return styles.AssistantMessageBorderStyle.Width(width).Render(wrapped)
		}

		content := senderPrefix(msg.Sender) + strings.TrimRight(rendered, "\n\r\t ")
		// Constrain the content to available width to ensure it doesn't break layout
		content = constrainWidth(content, availableWidth)
		return styles.AssistantMessageBorderStyle.Width(width).Render(content)
	case types.MessageTypeAssistantReasoning:
		if msg.Content == "" {
			return mv.spinner.View()
		}
		// Render through the markdown renderer to ensure proper wrapping to width
		availableWidth := width - 1
		if availableWidth < 10 {
			availableWidth = 10
		}
		rendered, err := markdown.NewRenderer(availableWidth).Render(msg.Content)
		if err != nil {
			text := "Thinking: " + senderPrefix(msg.Sender) + msg.Content
			wrapped := wrapText(text, availableWidth)
			return styles.MutedStyle.Italic(true).Render(wrapped)
		}
		// Strip ANSI from inner rendering so muted style fully applies
		clean := stripANSI(strings.TrimRight(rendered, "\n\r\t "))
		thinkingText := "Thinking: " + senderPrefix(msg.Sender) + clean
		return styles.MutedStyle.Italic(true).Render(thinkingText)
	case types.MessageTypeShellOutput:
		availableWidth := width - 1
		if availableWidth < 10 {
			availableWidth = 10
		}
		if rendered, err := markdown.NewRenderer(availableWidth).Render(fmt.Sprintf("```console\n%s\n```", msg.Content)); err == nil {
			return strings.TrimRight(rendered, "\n\r\t ")
		}
		wrapped := wrapText(msg.Content, availableWidth)
		return wrapped
	case types.MessageTypeCancelled:
		return styles.WarningStyle.Render("⚠ stream cancelled ⚠")
	case types.MessageTypeWelcome:
		// Render welcome message with a distinct style
		availableWidth := width - 1
		if availableWidth < 10 {
			availableWidth = 10
		}
		rendered, err := markdown.NewRenderer(availableWidth).Render(msg.Content)
		if err != nil {
			wrapped := wrapText(msg.Content, availableWidth)
			return styles.MutedStyle.Render(wrapped)
		}
		return styles.MutedStyle.Render(strings.TrimRight(rendered, "\n\r\t "))
	case types.MessageTypeError:
		return styles.ErrorStyle.Render("│ " + msg.Content)
	default:
		return msg.Content
	}
}

func senderPrefix(sender string) string {
	if sender == "" {
		return ""
	}
	return styles.AgentBadgeStyle.Render("["+sender+"]") + "\n\n"
}

// Height calculates the height needed for this message view
func (mv *messageModel) Height(width int) int {
	content := mv.Render(width)
	return strings.Count(content, "\n") + 1
}

// Message returns the underlying message
func (mv *messageModel) Message() *types.Message {
	return mv.message
}

// Layout.Sizeable methods

// SetSize sets the dimensions of the message view
func (mv *messageModel) SetSize(width, height int) tea.Cmd {
	mv.width = width
	mv.height = height
	return nil
}

// GetSize returns the current dimensions
func (mv *messageModel) GetSize() (width, height int) {
	return mv.width, mv.height
}

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// wrapText wraps text to the specified width, handling ANSI codes
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
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

// constrainWidth ensures that each line of content doesn't exceed the specified width
// This is needed because markdown renderer might output lines wider than requested
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

				if currentWidth == 0 {
					// First word on line
					currentLine.WriteString(word)
					currentWidth = wordWidth
				} else if currentWidth+spaceWidth+wordWidth <= maxWidth {
					// Word fits on current line
					currentLine.WriteString(" " + word)
					currentWidth += spaceWidth + wordWidth
				} else {
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
