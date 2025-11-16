package registry

// Factory creates components using a registry to look up builders.
// K is the key type used for lookups
// B is the builder function type
type Factory[K comparable, B any] struct {
	registry *Registry[K, B]
}

// NewFactory creates a new Factory instance.
func NewFactory[
	K comparable,
	B any,
](
	registry *Registry[K, B],
) *Factory[K, B] {
	return &Factory[K, B]{
		registry: registry,
	}
}

// GetBuilder retrieves a builder by key.
// Returns the builder and a boolean indicating if it was found.
func (f *Factory[K, B]) GetBuilder(key K) (B, bool) {
	return f.registry.Get(key)
}
