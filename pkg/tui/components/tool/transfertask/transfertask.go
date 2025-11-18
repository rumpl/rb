package transfertask

import (
	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour/v2"

	"github.com/rumpl/rb/pkg/tools/builtin"
	"github.com/rumpl/rb/pkg/tui/components/toolcommon"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component is a specialized component for rendering transfer_task tool calls.
type Component struct {
	message      *types.Message
	renderer     *glamour.TermRenderer
	width        int
	themeManager *styles.Manager
}

func New(
	msg *types.Message,
	renderer *glamour.TermRenderer,
	_ *service.SessionState,
	themeManager *styles.Manager,
) layout.Model {
	return &Component{
		message:      msg,
		renderer:     renderer,
		width:        80,
		themeManager: themeManager,
	}
}

func (c *Component) SetSize(width, height int) tea.Cmd {
	c.width = width
	return nil
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(tea.Msg) (layout.Model, tea.Cmd) {
	return c, nil
}

func (c *Component) View() string {
	var params builtin.TransferTaskArgs
	if err := json.Unmarshal([]byte(c.message.ToolCall.Function.Arguments), &params); err != nil {
		return "" // TODO: Partial tool call
	}

	theme := c.themeManager.GetTheme()
	badge := theme.TransferBadgeStyle.Render(c.message.Sender + " -> " + params.Agent + ":")
	content := theme.MutedStyle.Render(params.Task)

	body := badge + "\n" + content
	return toolcommon.RenderToolMessage(c.width, body, c.themeManager)
}
