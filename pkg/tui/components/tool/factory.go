package tool

import (
	"github.com/charmbracelet/glamour/v2"

	"github.com/rumpl/rb/pkg/tools/builtin"
	"github.com/rumpl/rb/pkg/tui/components/tool/defaulttool"
	"github.com/rumpl/rb/pkg/tui/components/tool/editfile"
	"github.com/rumpl/rb/pkg/tui/components/tool/readfile"
	"github.com/rumpl/rb/pkg/tui/components/tool/todotool"
	"github.com/rumpl/rb/pkg/tui/components/tool/transfertask"
	"github.com/rumpl/rb/pkg/tui/components/tool/writefile"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Create creates a tool component, using registered builders or falling back to default.
func (f *Factory) Create(
	msg *types.Message,
	renderer *glamour.TermRenderer,
	sessionState *service.SessionState,
) layout.Model {
	toolName := msg.ToolCall.Function.Name

	if builder, ok := f.GetBuilder(toolName); ok {
		return builder(msg, renderer, sessionState)
	}

	return defaulttool.New(msg, renderer, sessionState)
}

var (
	defaultRegistry = newDefaultRegistry()
	defaultFactory  = NewFactory(defaultRegistry)
)

func newDefaultRegistry() *Registry {
	reg := NewRegistry()

	reg.Register(builtin.ToolNameTransferTask, transfertask.New)
	reg.Register(builtin.ToolNameEditFile, editfile.New)
	reg.Register(builtin.ToolNameWriteFile, writefile.New)
	reg.Register(builtin.ToolNameReadFile, readfile.New)
	reg.Register(builtin.ToolNameCreateTodo, todotool.New)
	reg.Register(builtin.ToolNameCreateTodos, todotool.New)
	reg.Register(builtin.ToolNameUpdateTodo, todotool.New)
	reg.Register(builtin.ToolNameListTodos, todotool.New)

	return reg
}

func New(
	msg *types.Message,
	renderer *glamour.TermRenderer,
	sessionState *service.SessionState,
) layout.Model {
	return defaultFactory.Create(msg, renderer, sessionState)
}
