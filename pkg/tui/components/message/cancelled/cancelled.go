package cancelled

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents a cancelled message view
type Component struct {
	message *types.Message
	width   int
	height  int
}

// New creates a new cancelled message component
func New(msg *types.Message) layout.Model {
	return &Component{
		message: msg,
		width:   80,
		height:  1,
	}
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	return c, nil
}

func (c *Component) View() string {
	return styles.WarningStyle.Render("⚠ stream cancelled ⚠")
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
	content := styles.WarningStyle.Render("⚠ stream cancelled ⚠")
	return strings.Count(content, "\n") + 1
}

func (c *Component) SetMessage(msg *types.Message) {
	c.message = msg
}
