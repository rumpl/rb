package styles

import (
	"strings"

	"github.com/charmbracelet/glamour/v2/ansi"
)

func toChroma(style ansi.StylePrimitive) string {
	var s []string

	if style.Color != nil {
		s = append(s, *style.Color)
	}
	if style.BackgroundColor != nil {
		s = append(s, "bg:"+*style.BackgroundColor)
	}
	if style.Italic != nil && *style.Italic {
		s = append(s, "italic")
	}
	if style.Bold != nil && *style.Bold {
		s = append(s, "bold")
	}
	if style.Underline != nil && *style.Underline {
		s = append(s, "underline")
	}

	return strings.Join(s, " ")
}

func uintPtr(u uint) *uint {
	return &u
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
