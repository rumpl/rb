package message

import (
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/types"
)

// Model represents a view that can render a message
type Model interface {
	layout.Model
	layout.Sizeable
	SetMessage(msg *types.Message)
	Height(width int) int
}
