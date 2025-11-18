package tool

import (
	"github.com/charmbracelet/glamour/v2"

	"github.com/rumpl/rb/pkg/tui/components/registry"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// ComponentBuilder is a function that creates a tool component.
type ComponentBuilder func(
	msg *types.Message,
	renderer *glamour.TermRenderer,
	sessionState *service.SessionState,
	themeManager *styles.Manager,
) layout.Model

// Registry manages tool component builders.
type Registry = registry.Registry[string, ComponentBuilder]

// NewRegistry creates a new tool component registry.
func NewRegistry() *Registry {
	return registry.New[string, ComponentBuilder]()
}

// Factory creates tool components using the registry.
type Factory struct {
	*registry.Factory[string, ComponentBuilder]
}

// NewFactory creates a new tool factory.
func NewFactory(reg *Registry) *Factory {
	return &Factory{
		Factory: registry.NewFactory[string, ComponentBuilder](reg),
	}
}
