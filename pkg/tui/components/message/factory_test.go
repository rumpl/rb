package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rumpl/rb/pkg/tui/components/message"
	"github.com/rumpl/rb/pkg/tui/core/layout"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

func TestMessageFactory(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)

	t.Run("Create with registered builder", func(t *testing.T) {
		registry := message.NewRegistry()
		factory := message.NewFactory(registry)

		called := false
		customBuilder := func(msg *types.Message, _ *styles.Manager) layout.Model {
			called = true
			return nil
		}

		registry.Register(types.MessageTypeUser, customBuilder)

		msg := &types.Message{Type: types.MessageTypeUser}
		factory.Create(msg, themeManager)

		assert.True(t, called, "Custom builder should be called")
	})

	t.Run("Create with unregistered type falls back to default", func(t *testing.T) {
		registry := message.NewRegistry()
		factory := message.NewFactory(registry)

		msg := &types.Message{Type: types.MessageTypeUser}
		component := factory.Create(msg, themeManager)

		assert.NotNil(t, component, "Should create default component")
	})

	t.Run("Default factory creates components", func(t *testing.T) {
		msg := &types.Message{Type: types.MessageTypeUser}
		component := message.New(msg, themeManager)

		assert.NotNil(t, component, "Should create component from default factory")
	})
}
