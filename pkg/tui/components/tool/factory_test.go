package tool_test

import (
	"testing"

	"github.com/charmbracelet/glamour/v2"
	"github.com/stretchr/testify/assert"

	"github.com/rumpl/rb/pkg/tools"
	"github.com/rumpl/rb/pkg/tui/components/tool"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/types"
)

func TestToolFactory(t *testing.T) {
	t.Run("Create with registered builder", func(t *testing.T) {
		registry := tool.NewRegistry()
		factory := tool.NewFactory(registry)

		called := false
		customBuilder := func(
			msg *types.Message,
			renderer *glamour.TermRenderer,
			sessionState *service.SessionState,
		) layout.Model {
			called = true
			return nil
		}

		registry.Register("custom_tool", customBuilder)

		msg := &types.Message{
			ToolCall: tools.ToolCall{
				Function: tools.FunctionCall{
					Name: "custom_tool",
				},
			},
		}
		factory.Create(msg, nil, nil)

		assert.True(t, called, "Custom builder should be called")
	})

	t.Run("Create with unregistered tool falls back to default", func(t *testing.T) {
		registry := tool.NewRegistry()
		factory := tool.NewFactory(registry)

		msg := &types.Message{
			ToolCall: tools.ToolCall{
				Function: tools.FunctionCall{
					Name: "unknown_tool",
				},
			},
		}
		component := factory.Create(msg, nil, nil)

		assert.NotNil(t, component, "Should create default component")
	})

	t.Run("Default factory creates components", func(t *testing.T) {
		msg := &types.Message{
			ToolCall: tools.ToolCall{
				Function: tools.FunctionCall{
					Name: "unknown_tool",
				},
			},
		}
		component := tool.New(msg, nil, nil)

		assert.NotNil(t, component, "Should create component from default factory")
	})
}
