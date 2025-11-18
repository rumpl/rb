package editfile

import (
	"encoding/json"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour/v2"

	"github.com/rumpl/rb/pkg/tools/builtin"
	"github.com/rumpl/rb/pkg/tui/components/spinner"
	"github.com/rumpl/rb/pkg/tui/components/toolcommon"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

type ToggleDiffViewMsg struct{}

// Component is a specialized component for rendering edit_file tool calls.
type Component struct {
	message      *types.Message
	renderer     *glamour.TermRenderer
	spinner      spinner.Spinner
	width        int
	height       int
	sessionState *service.SessionState
	themeManager *styles.Manager
}

func New(
	msg *types.Message,
	renderer *glamour.TermRenderer,
	sessionState *service.SessionState,
	themeManager *styles.Manager,
) layout.Model {
	return &Component{
		message:      msg,
		renderer:     renderer,
		spinner:      spinner.New(spinner.ModeSpinnerOnly, themeManager),
		width:        80,
		height:       1,
		sessionState: sessionState,
		themeManager: themeManager,
	}
}

func (c *Component) SetSize(width, height int) tea.Cmd {
	c.width = width
	c.height = height
	return nil
}

func (c *Component) Init() tea.Cmd {
	if c.message.ToolStatus == types.ToolStatusPending || c.message.ToolStatus == types.ToolStatusRunning {
		return c.spinner.Init()
	}
	return nil
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	if c.message.ToolStatus == types.ToolStatusPending || c.message.ToolStatus == types.ToolStatusRunning {
		var cmd tea.Cmd
		var model layout.Model
		model, cmd = c.spinner.Update(msg)
		c.spinner = model.(spinner.Spinner)
		return c, cmd
	}

	return c, nil
}

func (c *Component) View() string {
	msg := c.message
	var args builtin.EditFileArgs
	if err := json.Unmarshal([]byte(msg.ToolCall.Function.Arguments), &args); err != nil {
		return ""
	}

	displayName := msg.ToolDefinition.DisplayName()
	theme := c.themeManager.GetTheme()
	content := fmt.Sprintf("%s %s %s", toolcommon.Icon(msg.ToolStatus, c.themeManager), theme.ToolCallTitleStyle.Render(displayName), theme.MutedStyle.Render(args.Path))

	if msg.ToolStatus == types.ToolStatusPending || msg.ToolStatus == types.ToolStatusRunning {
		content += " " + c.spinner.View()
	}

	// Account for border (1 char) + padding (2 left + 2 right) = 5 chars total
	availableWidth := c.width - 1 - 4 // 1 for border, 4 for padding (2 left + 2 right)
	if availableWidth < 10 {
		availableWidth = 10 // Minimum readable width
	}

	if msg.ToolCall.Function.Arguments != "" {
		content += "\n\n" + theme.ToolCallResult.Render(renderEditFile(msg.ToolCall, availableWidth, c.sessionState.SplitDiffView, msg.ToolStatus, c.themeManager))
	}

	var resultContent string
	if (msg.ToolStatus == types.ToolStatusCompleted || msg.ToolStatus == types.ToolStatusError) && msg.Content != "" {
		resultContent = toolcommon.FormatToolResult(msg.Content, availableWidth, c.themeManager)
	}

	return toolcommon.RenderToolMessage(c.width, content+resultContent, c.themeManager)
}
