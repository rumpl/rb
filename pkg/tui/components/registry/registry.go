package registry

import "sync"

// Registry manages component builders with type-safe key-value pairs.
// K is the key type (e.g., MessageType, string for tool names)
// V is the value type (e.g., ComponentBuilder function type)
type Registry[K comparable, V any] struct {
	mu      sync.RWMutex
	entries map[K]V
}

// New creates a new Registry instance.
func New[K comparable, V any]() *Registry[K, V] {
	return &Registry[K, V]{
		entries: make(map[K]V),
	}
}

// Register adds a new key-value pair to the registry.
func (r *Registry[K, V]) Register(key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[key] = value
}

// Get retrieves a value by key, returning the value and a boolean indicating existence.
func (r *Registry[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, exists := r.entries[key]
	return value, exists
}
