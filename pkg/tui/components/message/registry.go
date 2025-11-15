package message

import (
	"sync"

	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/types"
)

// ComponentBuilder is a function that creates a message component.
type ComponentBuilder func(msg *types.Message) layout.Model

// Registry manages message component builders.
type Registry struct {
	mu       sync.RWMutex
	builders map[types.MessageType]ComponentBuilder
}

func NewRegistry() *Registry {
	return &Registry{
		builders: make(map[types.MessageType]ComponentBuilder),
	}
}

func (r *Registry) Register(messageType types.MessageType, builder ComponentBuilder) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.builders[messageType] = builder
}

func (r *Registry) Get(messageType types.MessageType) (ComponentBuilder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	builder, exists := r.builders[messageType]
	return builder, exists
}
