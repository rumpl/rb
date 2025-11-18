package spinmsg

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/rumpl/rb/pkg/tui/components/spinner"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Component represents a spinner message view
type Component struct {
	message *types.Message
	width   int
	height  int
	spinner spinner.Spinner
}

// New creates a new spinner message component
func New(msg *types.Message, themeManager *styles.Manager) layout.Model {
	return &Component{
		message: msg,
		width:   80,
		height:  1,
		spinner: spinner.New(spinner.ModeBoth, themeManager),
	}
}

func (c *Component) Init() tea.Cmd {
	return c.spinner.Tick()
}

func (c *Component) Update(msg tea.Msg) (layout.Model, tea.Cmd) {
	s, cmd := c.spinner.Update(msg)
	c.spinner = s.(spinner.Spinner)
	return c, cmd
}

func (c *Component) View() string {
	return c.spinner.View()
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
	content := c.spinner.View()
	return strings.Count(content, "\n") + 1
}

func (c *Component) SetMessage(msg *types.Message) {
	c.message = msg
}
