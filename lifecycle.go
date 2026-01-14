package crab

import "sync"

// Lifecycle defines the interface for managing application lifecycle hooks.
// Providers should depend on this interface to register their startup/shutdown logic
// instead of returning cleanup closures or Hooks directly.
//
// Example Provider:
//
//	func NewRedis(conf *Config, lc crab.Lifecycle) (*redis.Client, error) {
//	    rdb := redis.NewClient(...)
//	    lc.Append(crab.Hook{
//	        OnStop: func(ctx context.Context) error {
//	            return rdb.Close()
//	        },
//	    })
//	    return rdb, nil
//	}
type Lifecycle interface {
	Append(Hook)
}

// Registry is a thread-safe implementation of Lifecycle.
// It is designed to be used as a singleton in the dependency injection container.
type Registry struct {
	hooks []Hook
	mu    sync.Mutex
}

// NewRegistry creates a new hook registry.
func NewRegistry() *Registry {
	return &Registry{
		hooks: make([]Hook, 0),
	}
}

// Append registers a new hook safely.
func (r *Registry) Append(h Hook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks = append(r.hooks, h)
}

// Hooks returns a copy of all registered hooks.
// This is typically called in main.go to pass hooks to app.Add(...).
func (r *Registry) Hooks() []Hook {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Return a copy to ensure thread safety during read
	result := make([]Hook, len(r.hooks))
	copy(result, r.hooks)
	return result
}
