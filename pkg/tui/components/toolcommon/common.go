package toolcommon

import (
	"encoding/json"
	"strings"

	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

func Icon(status types.ToolStatus, themeManager *styles.Manager) string {
	theme := themeManager.GetTheme()
	switch status {
	case types.ToolStatusPending:
		return "⊙"
	case types.ToolStatusRunning:
		return "⚙"
	case types.ToolStatusCompleted:
		return theme.SuccessStyle.Render("✓")
	case types.ToolStatusError:
		return theme.ErrorStyle.Render("✗")
	case types.ToolStatusConfirmation:
		return theme.WarningStyle.Render("?")
	default:
		return theme.WarningStyle.Render("?")
	}
}

func FormatToolResult(content string, width int, themeManager *styles.Manager) string {
	theme := themeManager.GetTheme()
	var formattedContent string
	var m map[string]any
	if err := json.Unmarshal([]byte(content), &m); err != nil {
		formattedContent = content
	} else if buf, err := json.MarshalIndent(m, "", "  "); err != nil {
		formattedContent = content
	} else {
		formattedContent = string(buf)
	}

	padding := theme.ToolCallResult.Padding().GetHorizontalPadding()
	availableWidth := max(width-2-padding, 10) // Minimum readable width

	lines := wrapLines(formattedContent, availableWidth)

	header := "output"
	if len(lines) > 10 {
		lines = lines[:10]
		header = "output (truncated)"
		lines = append(lines, wrapLines("...", availableWidth)...)
	}

	trimmedContent := strings.Join(lines, "\n")
	if trimmedContent != "" {
		return "\n" + theme.ToolCallResult.Render(theme.ToolCallResultKey.Render("\n-> "+header+":")+"\n"+trimmedContent)
	}

	return ""
}

func wrapLines(text string, width int) []string {
	if width <= 0 {
		return strings.Split(text, "\n")
	}

	var lines []string

	for line := range strings.SplitSeq(text, "\n") {
		for len(line) > width {
			lines = append(lines, line[:width])
			line = line[width:]
		}

		lines = append(lines, line)
	}

	return lines
}

// RenderToolMessage wraps arbitrary tool output in the same container used for assistant
// messages so both share identical padding, border, and background treatment.
func RenderToolMessage(width int, content string, themeManager *styles.Manager) string {
	theme := themeManager.GetTheme()
	trimmed := strings.TrimRight(content, "\n")
	return theme.AssistantMessageBorderStyle.Width(width).Render(trimmed)
}

// ContentWidthFromContainer returns the usable width inside an assistant message
// once borders and padding are removed. Ensures a minimum width for readability.
func ContentWidthFromContainer(width int) int {
	inner := width - 5 // 1 border + 2 padding on each side
	if inner < 10 {
		return 10
	}
	return inner
}
