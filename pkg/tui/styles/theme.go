package styles

import (
	"image/color"

	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/glamour/v2/ansi"
)

// ThemeName represents the name of a theme.
type ThemeName string

const (
	// ThemeDark is the dark theme (Tokyo Night-inspired with pink accents).
	ThemeDark ThemeName = "dark"
	// ThemeLight is the light theme.
	ThemeLight ThemeName = "light"
)

// Theme contains all style definitions for a theme.
type Theme struct {
	// Color definitions
	Colors ThemeColors

	// Markdown and syntax highlighting
	MarkdownStyle ansi.StyleConfig

	// Base Styles
	BaseStyle lipgloss.Style
	AppStyle  lipgloss.Style

	// Text Styles
	HighlightStyle     lipgloss.Style
	MutedStyle         lipgloss.Style
	SubtleStyle        lipgloss.Style
	SecondaryStyle     lipgloss.Style
	BoldStyle          lipgloss.Style
	ItalicStyle        lipgloss.Style
	ToolCallTitleStyle lipgloss.Style

	// Status Styles
	SuccessStyle    lipgloss.Style
	ErrorStyle      lipgloss.Style
	WarningStyle    lipgloss.Style
	InfoStyle       lipgloss.Style
	ActiveStyle     lipgloss.Style
	InProgressStyle lipgloss.Style
	PendingStyle    lipgloss.Style

	// Layout Styles
	CenterStyle lipgloss.Style

	// Border Styles
	BorderStyle                 lipgloss.Style
	BorderedBoxStyle            lipgloss.Style
	BorderedBoxFocusedStyle     lipgloss.Style
	UserMessageBorderStyle      lipgloss.Style
	AssistantMessageBorderStyle lipgloss.Style

	// Dialog Styles
	DialogStyle             lipgloss.Style
	DialogWarningStyle      lipgloss.Style
	DialogTitleStyle        lipgloss.Style
	DialogTitleWarningStyle lipgloss.Style
	DialogTitleInfoStyle    lipgloss.Style
	DialogContentStyle      lipgloss.Style
	DialogSeparatorStyle    lipgloss.Style
	DialogLabelStyle        lipgloss.Style
	DialogValueStyle        lipgloss.Style
	DialogQuestionStyle     lipgloss.Style
	DialogOptionsStyle      lipgloss.Style
	DialogHelpStyle         lipgloss.Style

	// Command Palette Styles
	PaletteSelectedStyle   lipgloss.Style
	PaletteUnselectedStyle lipgloss.Style
	PaletteCategoryStyle   lipgloss.Style
	PaletteDescStyle       lipgloss.Style

	// Diff Styles
	DiffAddStyle       lipgloss.Style
	DiffRemoveStyle    lipgloss.Style
	DiffUnchangedStyle lipgloss.Style
	DiffContextStyle   lipgloss.Style

	// Syntax highlighting UI element styles
	LineNumberStyle lipgloss.Style
	SeparatorStyle  lipgloss.Style

	// Tool Call Styles
	ToolCallArgs      lipgloss.Style
	ToolCallArgKey    lipgloss.Style
	ToolCallResult    lipgloss.Style
	ToolCallResultKey lipgloss.Style

	// Input Styles
	InputStyle           textarea.Styles
	EditorStyle          lipgloss.Style
	SuggestionGhostStyle lipgloss.Style

	// Scrollbar
	TrackStyle lipgloss.Style
	ThumbStyle lipgloss.Style

	// Notification Styles
	NotificationStyle        lipgloss.Style
	NotificationInfoStyle    lipgloss.Style
	NotificationWarningStyle lipgloss.Style
	NotificationErrorStyle   lipgloss.Style

	// Completion Styles
	CompletionBoxStyle       lipgloss.Style
	CompletionSelectedStyle  lipgloss.Style
	CompletionNormalStyle    lipgloss.Style
	CompletionDescStyle      lipgloss.Style
	CompletionNoResultsStyle lipgloss.Style

	// Badge Styles
	AgentBadgeStyle    lipgloss.Style
	TransferBadgeStyle lipgloss.Style

	// Deprecated styles (kept for backward compatibility)
	StatusStyle lipgloss.Style
	ActionStyle lipgloss.Style
	ChatStyle   lipgloss.Style

	// Selection Styles
	SelectionStyle lipgloss.Style

	// Spinner Styles
	SpinnerCharStyle          lipgloss.Style
	SpinnerTextBrightestStyle lipgloss.Style
	SpinnerTextBrightStyle    lipgloss.Style
	SpinnerTextDimStyle       lipgloss.Style
	SpinnerTextDimmestStyle   lipgloss.Style
}

// ThemeColors contains the color palette for a theme.
type ThemeColors struct {
	// Background colors
	Background    color.Color
	BackgroundAlt color.Color

	// Primary accent colors
	Accent    color.Color
	AccentDim color.Color

	// Status colors
	Success color.Color
	Error   color.Color
	Warning color.Color
	Info    color.Color

	// Text hierarchy
	TextPrimary   color.Color
	TextSecondary color.Color
	TextMuted     color.Color
	TextSubtle    color.Color

	// Border colors
	BorderPrimary   color.Color
	BorderSecondary color.Color
	BorderMuted     color.Color
	BorderWarning   color.Color
	BorderError     color.Color

	// Diff colors
	DiffAddBg    color.Color
	DiffRemoveBg color.Color
	DiffAddFg    color.Color
	DiffRemoveFg color.Color

	// UI element colors
	LineNumber color.Color
	Separator  color.Color

	// Interactive element colors
	Selected    color.Color
	SelectedFg  color.Color
	Hover       color.Color
	Placeholder color.Color

	// Badge colors
	AgentBadge    color.Color
	TransferBadge color.Color

	// Spinner colors
	SpinnerDim       color.Color
	SpinnerBright    color.Color
	SpinnerBrightest color.Color
}

// GetTheme returns a Theme for the given theme name.
func GetTheme(name ThemeName) Theme {
	switch name {
	case ThemeLight:
		return createLightTheme()
	case ThemeDark:
		return createDarkTheme()
	default:
		return createDarkTheme()
	}
}

// GetChromaStyle returns the appropriate Chroma style for the theme.
func (t Theme) GetChromaStyle() *chroma.Style {
	style, err := chroma.NewStyle("rb", t.getChromaTheme())
	if err != nil {
		panic(err)
	}
	return style
}

// getChromaTheme converts the markdown style to Chroma theme entries.
func (t Theme) getChromaTheme() chroma.StyleEntries {
	md := t.MarkdownStyle.CodeBlock
	return chroma.StyleEntries{
		chroma.Text:                toChroma(md.Chroma.Text),
		chroma.Error:               toChroma(md.Chroma.Error),
		chroma.Comment:             toChroma(md.Chroma.Comment),
		chroma.CommentPreproc:      toChroma(md.Chroma.CommentPreproc),
		chroma.Keyword:             toChroma(md.Chroma.Keyword),
		chroma.KeywordReserved:     toChroma(md.Chroma.KeywordReserved),
		chroma.KeywordNamespace:    toChroma(md.Chroma.KeywordNamespace),
		chroma.KeywordType:         toChroma(md.Chroma.KeywordType),
		chroma.Operator:            toChroma(md.Chroma.Operator),
		chroma.Punctuation:         toChroma(md.Chroma.Punctuation),
		chroma.Name:                toChroma(md.Chroma.Name),
		chroma.NameBuiltin:         toChroma(md.Chroma.NameBuiltin),
		chroma.NameTag:             toChroma(md.Chroma.NameTag),
		chroma.NameAttribute:       toChroma(md.Chroma.NameAttribute),
		chroma.NameClass:           toChroma(md.Chroma.NameClass),
		chroma.NameDecorator:       toChroma(md.Chroma.NameDecorator),
		chroma.NameFunction:        toChroma(md.Chroma.NameFunction),
		chroma.LiteralNumber:       toChroma(md.Chroma.LiteralNumber),
		chroma.LiteralString:       toChroma(md.Chroma.LiteralString),
		chroma.LiteralStringEscape: toChroma(md.Chroma.LiteralStringEscape),
		chroma.GenericDeleted:      toChroma(md.Chroma.GenericDeleted),
		chroma.GenericEmph:         toChroma(md.Chroma.GenericEmph),
		chroma.GenericInserted:     toChroma(md.Chroma.GenericInserted),
		chroma.GenericStrong:       toChroma(md.Chroma.GenericStrong),
		chroma.GenericSubheading:   toChroma(md.Chroma.GenericSubheading),
		chroma.Background:          toChroma(md.Chroma.Background),
	}
}
