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
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Factory creates message components using the registry.
// It looks up registered component builders and falls back to a default component
// if no specific builder is registered for a message type.
type Factory struct {
	registry *Registry
}

func NewFactory(registry *Registry) *Factory {
	return &Factory{
		registry: registry,
	}
}

func (f *Factory) Create(msg *types.Message) layout.Model {
	if builder, ok := f.registry.Get(msg.Type); ok {
		return builder(msg)
	}

	return defaultmsg.New(msg)
}

var (
	defaultRegistry = newDefaultRegistry()
	defaultFactory  = NewFactory(defaultRegistry)
)

func newDefaultRegistry() *Registry {
	registry := NewRegistry()

	registry.Register(types.MessageTypeUser, user.New)
	registry.Register(types.MessageTypeAssistant, assistant.New)
	registry.Register(types.MessageTypeAssistantReasoning, reasoning.New)
	registry.Register(types.MessageTypeSpinner, spinmsg.New)
	registry.Register(types.MessageTypeError, errormsg.New)
	registry.Register(types.MessageTypeShellOutput, shell.New)
	registry.Register(types.MessageTypeCancelled, cancelled.New)
	registry.Register(types.MessageTypeWelcome, welcome.New)

	return registry
}

// New creates a new message component using the default factory
func New(msg *types.Message) layout.Model {
	return defaultFactory.Create(msg)
}
