package pubsub

import (
	"context"
	"sync"
)

// Subscriber is a function that receives events of type T.
type Subscriber[T any] func(ctx context.Context, event T)

// Hub is a generic pub-sub hub that allows publishing events of type T and subscribing to them.
type Hub[T any] struct {
	mu          sync.RWMutex
	subscribers []Subscriber[T]
	closed      bool
}

// New creates a new pub-sub hub for events of type T.
func New[T any]() *Hub[T] {
	return &Hub[T]{
		subscribers: make([]Subscriber[T], 0),
	}
}

// Subscribe registers a new subscriber that will receive all published events.
// Returns an unsubscribe function that can be called to remove the subscriber.
func (h *Hub[T]) Subscribe(subscriber Subscriber[T]) func() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.closed {
		// Return a no-op unsubscribe function if hub is closed
		return func() {}
	}

	h.subscribers = append(h.subscribers, subscriber)
	index := len(h.subscribers) - 1

	// Return unsubscribe function
	return func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		if index < len(h.subscribers) {
			// Remove subscriber by setting it to nil and compacting the slice
			h.subscribers[index] = nil
			h.compactSubscribers()
		}
	}
}

// Publish sends an event to all registered subscribers.
// This is a non-blocking operation - subscribers are notified concurrently.
func (h *Hub[T]) Publish(ctx context.Context, event T) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.closed {
		return
	}

	// Create a copy of subscribers to avoid holding the lock during notification
	subscribersCopy := make([]Subscriber[T], 0, len(h.subscribers))
	for _, sub := range h.subscribers {
		if sub != nil {
			subscribersCopy = append(subscribersCopy, sub)
		}
	}

	// Notify all subscribers concurrently
	var wg sync.WaitGroup
	for _, sub := range subscribersCopy {
		wg.Add(1)
		go func(s Subscriber[T]) {
			defer wg.Done()
			defer func() {
				// Recover from panics in subscribers to prevent one bad subscriber
				// from affecting others
				if r := recover(); r != nil {
					// Log the panic if needed, but don't propagate it
				}
			}()
			s(ctx, event)
		}(sub)
	}

	wg.Wait()
}

// PublishAsync sends an event to all registered subscribers asynchronously.
// This returns immediately without waiting for subscribers to process the event.
func (h *Hub[T]) PublishAsync(ctx context.Context, event T) {
	go h.Publish(ctx, event)
}

// Close closes the hub and prevents new subscriptions or publications.
func (h *Hub[T]) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.closed = true
	h.subscribers = nil
}

// SubscriberCount returns the number of active subscribers.
func (h *Hub[T]) SubscriberCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, sub := range h.subscribers {
		if sub != nil {
			count++
		}
	}
	return count
}

// compactSubscribers removes nil entries from the subscribers slice.
// Must be called with h.mu locked.
func (h *Hub[T]) compactSubscribers() {
	n := 0
	for _, sub := range h.subscribers {
		if sub != nil {
			h.subscribers[n] = sub
			n++
		}
	}
	h.subscribers = h.subscribers[:n]
}
