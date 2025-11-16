package registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rumpl/rb/pkg/tui/components/registry"
)

func TestRegistry(t *testing.T) {
	t.Run("Register and Get with string keys", func(t *testing.T) {
		reg := registry.New[string, int]()

		reg.Register("foo", 42)
		reg.Register("bar", 100)

		val, ok := reg.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, 42, val)

		val, ok = reg.Get("bar")
		assert.True(t, ok)
		assert.Equal(t, 100, val)
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		reg := registry.New[string, int]()

		val, ok := reg.Get("missing")
		assert.False(t, ok)
		assert.Equal(t, 0, val) // zero value for int
	})

	t.Run("Register with custom type keys", func(t *testing.T) {
		type MessageType string
		type Builder func(string) string

		reg := registry.New[MessageType, Builder]()

		builder := func(msg string) string {
			return "built: " + msg
		}

		reg.Register(MessageType("user"), builder)

		retrieved, ok := reg.Get(MessageType("user"))
		assert.True(t, ok)
		assert.Equal(t, "built: test", retrieved("test"))
	})

	t.Run("Concurrent access", func(t *testing.T) {
		reg := registry.New[int, string]()

		done := make(chan bool)

		// Writer goroutine
		go func() {
			for i := range 100 {
				reg.Register(i, "value")
			}
			done <- true
		}()

		// Reader goroutine
		go func() {
			for i := range 100 {
				_, _ = reg.Get(i)
			}
			done <- true
		}()

		<-done
		<-done
	})
}
