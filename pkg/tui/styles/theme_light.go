package styles

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour/v2/ansi"
)

// Light theme color hex values
const (
	// Primary colors (light theme)
	LightColorAccent           = "#D6336C" // Deep rose/pink
	LightColorAccentDim        = "#E06B9F" // Lighter rose
	LightColorBackground       = "#FFFFFF" // Pure white
	LightColorBackgroundAlt    = "#F5F5F5" // Light gray
	LightColorBorderSecondary  = "#CCCCCC" // Medium gray
	LightColorTextPrimary      = "#1A1B26" // Dark blue-black
	LightColorTextSecondary    = "#414868" // Dark gray-blue
	LightColorSuccessGreen     = "#2E7D32" // Forest green
	LightColorErrorRed         = "#C62828" // Dark red
	LightColorWarningYellow    = "#F57C00" // Dark orange
	LightColorInfoBlue         = "#1976D2" // Blue
	LightColorSpinnerDim       = "#E06B9F" // Light rose
	LightColorSpinnerBright    = "#D6336C" // Deep rose
	LightColorSpinnerBrightest = "#B8004D" // Darker rose
	LightColorAgentBadge       = "#D6336C" // Deep rose
	LightColorTransferBadge    = "#E06B9F" // Light rose
	LightColorDiffAddBg        = "#C8E6C9" // Light green
	LightColorDiffRemoveBg     = "#FFCDD2" // Light red
	LightColorLineNumber       = "#9E9E9E" // Medium gray
	LightColorSeparator        = "#E0E0E0" // Light gray
	LightColorSelected         = "#E3F2FD" // Very light blue
	LightColorSelectedFg       = "#1A1B26" // Dark text
	LightColorHover            = "#F5F5F5" // Light gray
	LightColorPlaceholder      = "#9E9E9E" // Medium gray
)

// Light theme Chroma syntax highlighting colors
const (
	LightChromaErrorFgColor             = "#C62828"
	LightChromaSuccessColor             = "#2E7D32"
	LightChromaErrorBgColor             = "#FFEBEE"
	LightChromaCommentColor             = "#757575"
	LightChromaCommentPreprocColor      = "#F57C00"
	LightChromaKeywordColor             = "#1976D2"
	LightChromaKeywordReservedColor     = "#7B1FA2"
	LightChromaKeywordNamespaceColor    = "#C2185B"
	LightChromaKeywordTypeColor         = "#5E35B1"
	LightChromaOperatorColor            = "#C62828"
	LightChromaPunctuationColor         = "#424242"
	LightChromaNameBuiltinColor         = "#6A1B9A"
	LightChromaNameTagColor             = "#7B1FA2"
	LightChromaNameAttributeColor       = "#5E35B1"
	LightChromaNameDecoratorColor       = "#F57C00"
	LightChromaLiteralNumberColor       = "#00897B"
	LightChromaLiteralStringColor       = "#6D4C41"
	LightChromaLiteralStringEscapeColor = "#00897B"
	LightChromaGenericDeletedColor      = "#C62828"
	LightChromaGenericSubheadingColor   = "#616161"
	LightChromaBackgroundColor          = "#FAFAFA"
)

// ANSI color codes for light theme
const (
	LightANSIColor16  = "16"  // Black
	LightANSIColor18  = "18"  // Dark blue
	LightANSIColor125 = "125" // Purple
	LightANSIColor161 = "161" // Red
	LightANSIColor241 = "241" // Gray
	LightANSIColor242 = "242" // Gray
	LightANSIColor255 = "255" // White
)

func lightThemeColors() ThemeColors {
	return ThemeColors{
		// Background colors
		Background:    lipgloss.Color(LightColorBackground),
		BackgroundAlt: lipgloss.Color(LightColorBackgroundAlt),

		// Primary accent colors
		Accent:    lipgloss.Color(LightColorAccent),
		AccentDim: lipgloss.Color(LightColorAccentDim),

		// Status colors
		Success: lipgloss.Color(LightColorSuccessGreen),
		Error:   lipgloss.Color(LightColorErrorRed),
		Warning: lipgloss.Color(LightColorWarningYellow),
		Info:    lipgloss.Color(LightColorInfoBlue),

		// Text hierarchy
		TextPrimary:   lipgloss.Color(LightColorTextPrimary),
		TextSecondary: lipgloss.Color(LightColorTextSecondary),
		TextMuted:     lipgloss.Color(LightColorLineNumber),
		TextSubtle:    lipgloss.Color(LightColorSeparator),

		// Border colors
		BorderPrimary:   lipgloss.Color(LightColorAccent),
		BorderSecondary: lipgloss.Color(LightColorBorderSecondary),
		BorderMuted:     lipgloss.Color(LightColorSeparator),
		BorderWarning:   lipgloss.Color(LightColorWarningYellow),
		BorderError:     lipgloss.Color(LightColorErrorRed),

		// Diff colors
		DiffAddBg:    lipgloss.Color(LightColorDiffAddBg),
		DiffRemoveBg: lipgloss.Color(LightColorDiffRemoveBg),
		DiffAddFg:    lipgloss.Color(LightColorSuccessGreen),
		DiffRemoveFg: lipgloss.Color(LightColorErrorRed),

		// UI element colors
		LineNumber: lipgloss.Color(LightColorLineNumber),
		Separator:  lipgloss.Color(LightColorSeparator),

		// Interactive element colors
		Selected:    lipgloss.Color(LightColorSelected),
		SelectedFg:  lipgloss.Color(LightColorSelectedFg),
		Hover:       lipgloss.Color(LightColorHover),
		Placeholder: lipgloss.Color(LightColorPlaceholder),

		// Badge colors
		AgentBadge:    lipgloss.Color(LightColorAgentBadge),
		TransferBadge: lipgloss.Color(LightColorTransferBadge),

		// Spinner colors
		SpinnerDim:       lipgloss.Color(LightColorSpinnerDim),
		SpinnerBright:    lipgloss.Color(LightColorSpinnerBright),
		SpinnerBrightest: lipgloss.Color(LightColorSpinnerBrightest),
	}
}

func createLightTheme() Theme {
	colors := lightThemeColors()

	// Base Styles
	baseStyle := lipgloss.NewStyle().Foreground(colors.TextPrimary)
	appStyle := baseStyle.Padding(0, 1, 0, 1)

	theme := Theme{
		Colors: colors,

		// Markdown and syntax highlighting
		MarkdownStyle: getLightMarkdownStyle(),

		// Base Styles
		BaseStyle: baseStyle,
		AppStyle:  appStyle,

		// Text Styles
		HighlightStyle: baseStyle.Foreground(colors.Accent),
		MutedStyle:     baseStyle.Foreground(colors.TextMuted),
		SubtleStyle:    baseStyle.Foreground(colors.TextSubtle),
		SecondaryStyle: baseStyle.Foreground(colors.TextSecondary),
		BoldStyle:      baseStyle.Bold(true),
		ItalicStyle:    baseStyle.Italic(true),
		ToolCallTitleStyle: baseStyle.
			Foreground(colors.Accent).
			Bold(true).
			Background(colors.BackgroundAlt),

		// Status Styles
		SuccessStyle:    baseStyle.Foreground(colors.Success),
		ErrorStyle:      baseStyle.Foreground(colors.Error),
		WarningStyle:    baseStyle.Foreground(colors.Warning),
		InfoStyle:       baseStyle.Foreground(colors.Info),
		ActiveStyle:     baseStyle.Foreground(colors.Success),
		InProgressStyle: baseStyle.Foreground(colors.Warning),
		PendingStyle:    baseStyle.Foreground(colors.TextSecondary),

		// Layout Styles
		CenterStyle: baseStyle.Align(lipgloss.Center, lipgloss.Center),

		// Border Styles
		BorderStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderPrimary),

		BorderedBoxStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderSecondary).
			Padding(0, 1),

		BorderedBoxFocusedStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderPrimary).
			Padding(0, 1),

		UserMessageBorderStyle: baseStyle.
			Padding(1, 2).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(colors.BorderPrimary).
			Bold(true).
			Background(colors.BackgroundAlt),

		AssistantMessageBorderStyle: baseStyle.
			Padding(1, 2).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(colors.BorderSecondary).
			Background(colors.BackgroundAlt),

		// Dialog Styles
		DialogStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderSecondary).
			Foreground(colors.TextPrimary).
			Padding(1, 2).
			Align(lipgloss.Left),

		DialogWarningStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderWarning).
			Foreground(colors.TextPrimary).
			Padding(1, 2).
			Align(lipgloss.Left),

		DialogTitleStyle: baseStyle.
			Bold(true).
			Foreground(colors.TextSecondary).
			Align(lipgloss.Center),

		DialogTitleWarningStyle: baseStyle.
			Bold(true).
			Foreground(colors.Warning).
			Align(lipgloss.Center),

		DialogTitleInfoStyle: baseStyle.
			Bold(true).
			Foreground(colors.Info).
			Align(lipgloss.Center),

		DialogContentStyle: baseStyle.
			Foreground(colors.TextPrimary),

		DialogSeparatorStyle: baseStyle.
			Foreground(colors.BorderMuted),

		DialogLabelStyle: baseStyle.
			Bold(true).
			Foreground(colors.TextMuted),

		DialogValueStyle: baseStyle.
			Bold(true).
			Foreground(colors.TextSecondary),

		DialogQuestionStyle: baseStyle.
			Bold(true).
			Foreground(colors.TextPrimary).
			Align(lipgloss.Center),

		DialogOptionsStyle: baseStyle.
			Foreground(colors.TextMuted).
			Align(lipgloss.Center),

		DialogHelpStyle: baseStyle.
			Foreground(colors.TextMuted).
			Italic(true),

		// Command Palette Styles
		PaletteSelectedStyle: baseStyle.
			Background(colors.Selected).
			Foreground(colors.SelectedFg).
			Padding(0, 1),

		PaletteUnselectedStyle: baseStyle.
			Foreground(colors.TextPrimary).
			Padding(0, 1),

		PaletteCategoryStyle: baseStyle.
			Bold(true).
			Foreground(colors.TextMuted).
			MarginTop(1),

		PaletteDescStyle: baseStyle.
			Foreground(colors.TextMuted),

		// Diff Styles
		DiffAddStyle: baseStyle.
			Background(colors.DiffAddBg).
			Foreground(colors.DiffAddFg),

		DiffRemoveStyle: baseStyle.
			Background(colors.DiffRemoveBg).
			Foreground(colors.DiffRemoveFg),

		DiffUnchangedStyle: baseStyle.Background(colors.BackgroundAlt),

		DiffContextStyle: baseStyle,

		// Syntax highlighting UI element styles
		LineNumberStyle: baseStyle.Foreground(colors.LineNumber).Background(colors.BackgroundAlt),
		SeparatorStyle:  baseStyle.Foreground(colors.Separator).Background(colors.BackgroundAlt),

		// Tool Call Styles
		ToolCallArgs: baseStyle.
			BorderLeft(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderSecondary).
			Background(colors.BackgroundAlt),

		ToolCallArgKey: baseStyle.Bold(true).Foreground(colors.TextSecondary),

		ToolCallResult: baseStyle.
			BorderLeft(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderSecondary).
			Background(colors.BackgroundAlt),

		ToolCallResultKey: baseStyle.Bold(true).Foreground(colors.TextSecondary),

		// Input Styles
		InputStyle: textarea.Styles{
			Focused: textarea.StyleState{
				Placeholder: baseStyle.Foreground(colors.Placeholder),
			},
			Blurred: textarea.StyleState{
				Placeholder: baseStyle.Foreground(colors.Placeholder),
			},
			Cursor: textarea.CursorStyle{
				Color: colors.Accent,
			},
		},
		EditorStyle:          baseStyle.Background(colors.BackgroundAlt),
		SuggestionGhostStyle: baseStyle.Foreground(colors.TextMuted),

		// Scrollbar
		TrackStyle: lipgloss.NewStyle().Foreground(colors.BorderSecondary),
		ThumbStyle: lipgloss.NewStyle().Foreground(colors.Accent),

		// Notification Styles
		NotificationStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Success).
			Padding(0, 1),

		NotificationInfoStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Info).
			Padding(0, 1),

		NotificationWarningStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Warning).
			Padding(0, 1),

		NotificationErrorStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Error).
			Padding(0, 1),

		// Completion Styles
		CompletionBoxStyle: baseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderSecondary).
			Padding(0, 1),

		CompletionSelectedStyle: baseStyle.
			Foreground(colors.TextPrimary).
			Bold(true),

		CompletionNormalStyle: baseStyle.
			Foreground(colors.TextPrimary),

		CompletionDescStyle: baseStyle.
			Foreground(colors.TextSecondary).
			Italic(true),

		CompletionNoResultsStyle: baseStyle.
			Foreground(colors.TextMuted).
			Italic(true).
			Align(lipgloss.Center),

		// Badge Styles
		AgentBadgeStyle: baseStyle.
			Foreground(colors.AgentBadge).
			Bold(true).
			Background(colors.BackgroundAlt),

		TransferBadgeStyle: baseStyle.
			Foreground(colors.TransferBadge).
			Bold(true).
			Background(colors.BackgroundAlt),

		// Deprecated styles (kept for backward compatibility)
		StatusStyle: baseStyle.Foreground(colors.TextMuted),
		ActionStyle: baseStyle.Foreground(colors.TextSecondary),
		ChatStyle:   baseStyle,

		// Selection Styles
		SelectionStyle: baseStyle.
			Background(colors.Selected).
			Foreground(colors.SelectedFg),

		// Spinner Styles
		SpinnerCharStyle:          baseStyle.Foreground(colors.Accent),
		SpinnerTextBrightestStyle: baseStyle.Foreground(colors.SpinnerBrightest),
		SpinnerTextBrightStyle:    baseStyle.Foreground(colors.SpinnerBright),
		SpinnerTextDimStyle:       baseStyle.Foreground(colors.SpinnerDim),
		SpinnerTextDimmestStyle:   baseStyle.Foreground(colors.Accent),
	}

	return theme
}

func getLightMarkdownStyle() ansi.StyleConfig {
	h1Color := LightColorAccent
	h2Color := LightColorAccent
	h3Color := LightColorTextSecondary
	h4Color := LightColorTextSecondary
	h5Color := LightColorTextSecondary
	h6Color := LightColorLineNumber
	linkColor := LightColorInfoBlue
	strongColor := LightColorTextPrimary
	codeColor := LightColorTextPrimary
	blockquoteColor := LightColorTextSecondary
	listColor := LightColorTextPrimary
	hrColor := LightColorBorderSecondary

	customLightStyle := ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix:     "",
				BlockSuffix:     "",
				Color:           stringPtr(LightANSIColor16),
				BackgroundColor: stringPtr(LightColorBackgroundAlt),
			},
			Margin: uintPtr(0),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: &blockquoteColor,
			},
			Indent:      uintPtr(1),
			IndentToken: nil,
		},
		List: ansi.StyleList{
			LevelIndent: defaultListIndent,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr(LightANSIColor18),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           &h1Color,
				BackgroundColor: stringPtr(LightANSIColor125),
				Bold:            boolPtr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Color:  &h2Color,
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
				Color:  &h3Color,
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
				Color:  &h4Color,
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
				Color:  &h5Color,
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  &h6Color,
				Bold:   boolPtr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: boolPtr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Color: &strongColor,
			Bold:  boolPtr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  &hrColor,
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     &linkColor,
			Underline: boolPtr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: stringPtr(LightANSIColor18),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     stringPtr(LightANSIColor161),
			Underline: boolPtr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  stringPtr(LightANSIColor241),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           &codeColor,
				BackgroundColor: stringPtr(LightColorBackgroundAlt),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color:           stringPtr(LightANSIColor242),
					BackgroundColor: stringPtr(LightColorBackgroundAlt),
				},
				Margin: uintPtr(defaultMargin),
			},
			Theme: "monokailight",
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: stringPtr(LightColorTextPrimary),
				},
				Error: ansi.StylePrimitive{
					Color:           stringPtr(LightChromaErrorFgColor),
					BackgroundColor: stringPtr(LightChromaErrorBgColor),
				},
				Comment: ansi.StylePrimitive{
					Color: stringPtr(LightChromaCommentColor),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: stringPtr(LightChromaCommentPreprocColor),
				},
				Keyword: ansi.StylePrimitive{
					Color: stringPtr(LightChromaKeywordColor),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: stringPtr(LightChromaKeywordReservedColor),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: stringPtr(LightChromaKeywordNamespaceColor),
				},
				KeywordType: ansi.StylePrimitive{
					Color: stringPtr(LightChromaKeywordTypeColor),
				},
				Operator: ansi.StylePrimitive{
					Color: stringPtr(LightChromaOperatorColor),
				},
				Punctuation: ansi.StylePrimitive{
					Color: stringPtr(LightChromaPunctuationColor),
				},
				Name: ansi.StylePrimitive{
					Color: stringPtr(LightColorTextPrimary),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: stringPtr(LightChromaNameBuiltinColor),
				},
				NameTag: ansi.StylePrimitive{
					Color: stringPtr(LightChromaNameTagColor),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: stringPtr(LightChromaNameAttributeColor),
				},
				NameClass: ansi.StylePrimitive{
					Color:     stringPtr(LightChromaErrorFgColor),
					Underline: boolPtr(true),
					Bold:      boolPtr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: stringPtr(LightChromaNameDecoratorColor),
				},
				NameFunction: ansi.StylePrimitive{
					Color: stringPtr(LightChromaSuccessColor),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: stringPtr(LightChromaLiteralNumberColor),
				},
				LiteralString: ansi.StylePrimitive{
					Color: stringPtr(LightChromaLiteralStringColor),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: stringPtr(LightChromaLiteralStringEscapeColor),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: stringPtr(LightChromaGenericDeletedColor),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: stringPtr(LightChromaSuccessColor),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: stringPtr(LightChromaGenericSubheadingColor),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: stringPtr(LightColorBackgroundAlt),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	customLightStyle.List.Color = &listColor
	bg := LightColorBackgroundAlt
	customLightStyle.CodeBlock.BackgroundColor = &bg
	customLightStyle.Code.BackgroundColor = &bg
	customLightStyle.CodeBlock.Chroma.Background.BackgroundColor = &bg

	return customLightStyle
}
