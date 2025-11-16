package message

import (
	"github.com/rumpl/rb/pkg/tui/components/message/assistant"
	"github.com/rumpl/rb/pkg/tui/components/message/cancelled"
	"github.com/rumpl/rb/pkg/tui/components/message/defaultmsg"
	"github.com/rumpl/rb/pkg/tui/components/message/errormsg"
	"github.com/rumpl/rb/pkg/tui/components/message/reasoning"
	"github.com/rumpl/rb/pkg/tui/components/message/shell"
	"github.com/rumpl/rb/pkg/tui/components/message/spinmsg"
	"github.com/rumpl/rb/pkg/tui/components/message/user"
	"github.com/rumpl/rb/pkg/tui/components/message/welcome"
	"github.com/rumpl/rb/pkg/tui/components/registry"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Factory creates message components using the registry.
type Factory struct {
	*registry.Factory[types.MessageType, ComponentBuilder]
}

// NewFactory creates a new message factory.
func NewFactory(reg *Registry) *Factory {
	return &Factory{
		Factory: registry.NewFactory(reg),
	}
}

// Create creates a message component, using registered builders or falling back to default.
func (f *Factory) Create(msg *types.Message) layout.Model {
	if builder, ok := f.GetBuilder(msg.Type); ok {
		return builder(msg)
	}

	return defaultmsg.New(msg)
}

var (
	defaultRegistry = newDefaultRegistry()
	defaultFactory  = NewFactory(defaultRegistry)
)

// ComponentBuilder is a function that creates a message component.
type ComponentBuilder func(msg *types.Message) layout.Model

// Registry manages message component builders.
type Registry = registry.Registry[types.MessageType, ComponentBuilder]

// NewRegistry creates a new message component registry.
func NewRegistry() *Registry {
	return registry.New[types.MessageType, ComponentBuilder]()
}

func newDefaultRegistry() *Registry {
	reg := registry.New[types.MessageType, ComponentBuilder]()

	reg.Register(types.MessageTypeUser, user.New)
	reg.Register(types.MessageTypeAssistant, assistant.New)
	reg.Register(types.MessageTypeAssistantReasoning, reasoning.New)
	reg.Register(types.MessageTypeSpinner, spinmsg.New)
	reg.Register(types.MessageTypeError, errormsg.New)
	reg.Register(types.MessageTypeShellOutput, shell.New)
	reg.Register(types.MessageTypeCancelled, cancelled.New)
	reg.Register(types.MessageTypeWelcome, welcome.New)

	return reg
}

// New creates a new message component using the default factory
func New(msg *types.Message) layout.Model {
	return defaultFactory.Create(msg)
}
