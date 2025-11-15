package defaulttool

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rumpl/rb/pkg/tools"
	"github.com/rumpl/rb/pkg/tui/styles"
)

func renderToolArgs(toolCall tools.ToolCall, width int) string {
	decoder := json.NewDecoder(strings.NewReader(toolCall.Function.Arguments))

	tok, err := decoder.Token()
	if err != nil {
		return ""
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		return ""
	}

	type kv struct {
		Key   string
		Value any
	}
	var kvs []kv

	for decoder.More() {
		tok, err := decoder.Token()
		if err != nil {
			return ""
		}
		key, ok := tok.(string)
		if !ok {
			return ""
		}

		var val any
		if err := decoder.Decode(&val); err != nil {
			return ""
		}

		kvs = append(kvs, kv{Key: key, Value: val})
	}
	_, _ = decoder.Token()

	style := styles.ToolCallArgs.Width(width)

	var md strings.Builder
	for i, kv := range kvs {
		if i > 0 {
			md.WriteString("\n")
		}

		var content string
		if v, ok := kv.Value.(string); ok {
			content = v
		} else {
			buf, err := json.MarshalIndent(kv.Value, "", "  ")
			if err != nil {
				content = fmt.Sprintf("%v", kv.Value)
			} else {
				content = string(buf)
			}
		}

		// Wrap content to fit within available width
		// Account for the key label and colon
		keyLabel := kv.Key + ":"
		keyWidth := len(keyLabel)
		contentWidth := width - keyWidth - 2 // 2 for spacing
		if contentWidth < 10 {
			contentWidth = 10
		}

		// Wrap long lines in content
		wrappedContent := wrapContentLines(content, contentWidth)

		fmt.Fprintf(&md, "%s:\n%s", styles.ToolCallArgKey.Render(kv.Key), wrappedContent)
		if !strings.HasSuffix(wrappedContent, "\n") {
			md.WriteString("\n")
		}
	}

	return "\n" + style.Render(strings.TrimSuffix(md.String(), "\n"))
}

// wrapContentLines wraps long lines in content to fit within the specified width
func wrapContentLines(content string, width int) string {
	if width <= 0 {
		return content
	}

	var lines []string
	for _, line := range strings.Split(content, "\n") {
		for len(line) > width {
			// Find the last space before width to break at word boundary
			breakPoint := width
			if idx := strings.LastIndex(line[:min(width, len(line))], " "); idx > width/2 {
				breakPoint = idx + 1
			}
			lines = append(lines, line[:breakPoint])
			line = strings.TrimLeft(line[breakPoint:], " ")
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
