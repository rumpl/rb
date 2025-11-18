package styles

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour/v2/ansi"
)

const (
	defaultListIndent = 2
	defaultMargin     = 2
)

// Dark theme color hex values
const (
	// Primary colors
	ColorAccentBlue      = "#FF69B4" // Hot pink
	ColorMutedBlue       = "#C76B98" // Muted pink
	ColorBackgroundAlt   = "#24283B" // Slightly lighter background (kept blue)
	ColorBorderSecondary = "#8B5A7D" // Dark mauve-pink
	ColorTextPrimary     = "#FFB3D9" // Light pink
	ColorTextSecondary   = "#FF8FC7" // Medium pink
	ColorSuccessGreen    = "#9ECE6A" // Soft green (kept for semantic meaning)
	ColorErrorRed        = "#F7768E" // Soft red (kept for semantic meaning)
	ColorWarningYellow   = "#E0AF68" // Soft yellow (kept for semantic meaning)

	// Spinner glow colors (transition from base pink towards white)
	ColorSpinnerDim       = "#FF85C0" // Light pink
	ColorSpinnerBright    = "#FFB3D9" // Much lighter pink
	ColorSpinnerBrightest = "#FFE0F0" // Very light pink, near white

	// Background colors
	ColorBackground = "#1A1B26" // Dark blue-black (kept blue)

	// Status colors
	ColorInfoCyan = "#FF91D7" // Bright pink

	// Badge colors
	ColorAgentBadge    = "#FF6BB5" // Vibrant pink
	ColorTransferBadge = "#FFAAE0" // Light pink

	// Diff colors
	ColorDiffAddBg    = "#20303B" // Dark blue-green
	ColorDiffRemoveBg = "#3C2A2A" // Dark red-brown

	// Line number and UI element colors
	ColorLineNumber = "#565F89" // Muted blue-grey
	ColorSeparator  = "#414868" // Dark blue-grey

	// Word-level diff highlight colors (visible but not harsh)
	ColorDiffWordAddBg    = "#2D4F3F" // Medium dark teal with green tint
	ColorDiffWordRemoveBg = "#4F2D3A" // Medium dark burgundy with red tint

	// Interactive element colors
	ColorSelected = "#5A3A5F" // Dark mauve for selected items
	ColorHover    = "#4A2D4F" // Slightly lighter than selected
)

// Chroma syntax highlighting colors (Monokai theme)
const (
	ChromaErrorFgColor             = "#F1F1F1"
	ChromaSuccessColor             = "#00D787"
	ChromaErrorBgColor             = "#F05B5B"
	ChromaCommentColor             = "#676767"
	ChromaCommentPreprocColor      = "#FF875F"
	ChromaKeywordColor             = "#00AAFF"
	ChromaKeywordReservedColor     = "#FF5FD2"
	ChromaKeywordNamespaceColor    = "#FF5F87"
	ChromaKeywordTypeColor         = "#6E6ED8"
	ChromaOperatorColor            = "#EF8080"
	ChromaPunctuationColor         = "#E8E8A8"
	ChromaNameBuiltinColor         = "#FF8EC7"
	ChromaNameTagColor             = "#B083EA"
	ChromaNameAttributeColor       = "#7A7AE6"
	ChromaNameDecoratorColor       = "#FFFF87"
	ChromaLiteralNumberColor       = "#6EEFC0"
	ChromaLiteralStringColor       = "#C69669"
	ChromaLiteralStringEscapeColor = "#AFFFD7"
	ChromaGenericDeletedColor      = "#FD5B5B"
	ChromaGenericSubheadingColor   = "#777777"
	ChromaBackgroundColor          = "#373737"
)

// ANSI color codes (8-bit color codes)
const (
	ANSIColor252 = "252"
	ANSIColor39  = "39"
	ANSIColor63  = "63"
	ANSIColor35  = "35"
	ANSIColor212 = "212"
	ANSIColor243 = "243"
	ANSIColor244 = "244"
)

func darkThemeColors() ThemeColors {
	return ThemeColors{
		// Background colors
		Background:    lipgloss.Color(ColorBackground),
		BackgroundAlt: lipgloss.Color(ColorBackgroundAlt),

		// Primary accent colors
		Accent:    lipgloss.Color(ColorAccentBlue),
		AccentDim: lipgloss.Color(ColorMutedBlue),

		// Status colors
		Success: lipgloss.Color(ColorSuccessGreen),
		Error:   lipgloss.Color(ColorErrorRed),
		Warning: lipgloss.Color(ColorWarningYellow),
		Info:    lipgloss.Color(ColorInfoCyan),

		// Text hierarchy
		TextPrimary:   lipgloss.Color(ColorTextPrimary),
		TextSecondary: lipgloss.Color(ColorTextSecondary),
		TextMuted:     lipgloss.Color(ColorMutedBlue),
		TextSubtle:    lipgloss.Color(ColorBorderSecondary),

		// Border colors
		BorderPrimary:   lipgloss.Color(ColorAccentBlue),
		BorderSecondary: lipgloss.Color(ColorBorderSecondary),
		BorderMuted:     lipgloss.Color(ColorBackgroundAlt),
		BorderWarning:   lipgloss.Color(ColorWarningYellow),
		BorderError:     lipgloss.Color(ColorErrorRed),

		// Diff colors
		DiffAddBg:    lipgloss.Color(ColorDiffAddBg),
		DiffRemoveBg: lipgloss.Color(ColorDiffRemoveBg),
		DiffAddFg:    lipgloss.Color(ColorSuccessGreen),
		DiffRemoveFg: lipgloss.Color(ColorErrorRed),

		// UI element colors
		LineNumber: lipgloss.Color(ColorLineNumber),
		Separator:  lipgloss.Color(ColorSeparator),

		// Interactive element colors
		Selected:    lipgloss.Color(ColorSelected),
		SelectedFg:  lipgloss.Color(ColorTextPrimary),
		Hover:       lipgloss.Color(ColorHover),
		Placeholder: lipgloss.Color(ColorMutedBlue),

		// Badge colors
		AgentBadge:    lipgloss.Color(ColorAgentBadge),
		TransferBadge: lipgloss.Color(ColorTransferBadge),

		// Spinner colors
		SpinnerDim:       lipgloss.Color(ColorSpinnerDim),
		SpinnerBright:    lipgloss.Color(ColorSpinnerBright),
		SpinnerBrightest: lipgloss.Color(ColorSpinnerBrightest),
	}
}

func createDarkTheme() Theme {
	colors := darkThemeColors()

	// Base Styles
	baseStyle := lipgloss.NewStyle().Foreground(colors.TextPrimary)
	appStyle := baseStyle.Padding(0, 1, 0, 1)

	theme := Theme{
		Colors: colors,

		// Markdown and syntax highlighting
		MarkdownStyle: getDarkMarkdownStyle(),

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

func getDarkMarkdownStyle() ansi.StyleConfig {
	h1Color := ColorAccentBlue
	h2Color := ColorAccentBlue
	h3Color := ColorTextSecondary
	h4Color := ColorTextSecondary
	h5Color := ColorTextSecondary
	h6Color := ColorMutedBlue
	linkColor := ColorAccentBlue
	strongColor := ColorTextPrimary
	codeColor := ColorTextPrimary
	blockquoteColor := ColorTextSecondary
	listColor := ColorTextPrimary
	hrColor := ColorBorderSecondary

	customDarkStyle := ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix:     "",
				BlockSuffix:     "",
				Color:           stringPtr(ANSIColor252),
				BackgroundColor: stringPtr(ColorBackgroundAlt),
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
				Color:       stringPtr(ANSIColor39),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           &h1Color,
				BackgroundColor: stringPtr(ANSIColor63),
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
			Color: stringPtr(ANSIColor35),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     stringPtr(ANSIColor212),
			Underline: boolPtr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  stringPtr(ANSIColor243),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           &codeColor,
				BackgroundColor: stringPtr(ColorBackgroundAlt),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color:           stringPtr(ANSIColor244),
					BackgroundColor: stringPtr(ColorBackgroundAlt),
				},
				Margin: uintPtr(defaultMargin),
			},
			Theme: "monokai",
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: stringPtr(ColorTextPrimary),
				},
				Error: ansi.StylePrimitive{
					Color:           stringPtr(ChromaErrorFgColor),
					BackgroundColor: stringPtr(ChromaErrorBgColor),
				},
				Comment: ansi.StylePrimitive{
					Color: stringPtr(ChromaCommentColor),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: stringPtr(ChromaCommentPreprocColor),
				},
				Keyword: ansi.StylePrimitive{
					Color: stringPtr(ChromaKeywordColor),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: stringPtr(ChromaKeywordReservedColor),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: stringPtr(ChromaKeywordNamespaceColor),
				},
				KeywordType: ansi.StylePrimitive{
					Color: stringPtr(ChromaKeywordTypeColor),
				},
				Operator: ansi.StylePrimitive{
					Color: stringPtr(ChromaOperatorColor),
				},
				Punctuation: ansi.StylePrimitive{
					Color: stringPtr(ChromaPunctuationColor),
				},
				Name: ansi.StylePrimitive{
					Color: stringPtr(ColorTextPrimary),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: stringPtr(ChromaNameBuiltinColor),
				},
				NameTag: ansi.StylePrimitive{
					Color: stringPtr(ChromaNameTagColor),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: stringPtr(ChromaNameAttributeColor),
				},
				NameClass: ansi.StylePrimitive{
					Color:     stringPtr(ChromaErrorFgColor),
					Underline: boolPtr(true),
					Bold:      boolPtr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: stringPtr(ChromaNameDecoratorColor),
				},
				NameFunction: ansi.StylePrimitive{
					Color: stringPtr(ChromaSuccessColor),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: stringPtr(ChromaLiteralNumberColor),
				},
				LiteralString: ansi.StylePrimitive{
					Color: stringPtr(ChromaLiteralStringColor),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: stringPtr(ChromaLiteralStringEscapeColor),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: stringPtr(ChromaGenericDeletedColor),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: stringPtr(ChromaSuccessColor),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: stringPtr(ChromaGenericSubheadingColor),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: stringPtr(ColorBackgroundAlt),
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

	customDarkStyle.List.Color = &listColor
	bg := ColorBackgroundAlt
	customDarkStyle.CodeBlock.BackgroundColor = &bg
	customDarkStyle.Code.BackgroundColor = &bg
	customDarkStyle.CodeBlock.Chroma.Background.BackgroundColor = &bg

	return customDarkStyle
}
