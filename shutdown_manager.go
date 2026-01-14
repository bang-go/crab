package crab

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// ShutdownManager 管理所有App实例的全局shutdown
type ShutdownManager struct {
	apps   map[string]*App // key为app ID
	mu     sync.RWMutex
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

var globalShutdown = NewShutdownManager()

func NewShutdownManager() *ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShutdownManager{
		apps:   make(map[string]*App),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (sm *ShutdownManager) Register(app *App) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if app == nil {
		return fmt.Errorf("cannot register nil app")
	}

	if app.id == "" {
		return fmt.Errorf("app ID cannot be empty")
	}

	if _, exists := sm.apps[app.id]; exists {
		return fmt.Errorf("app with ID %s already registered", app.id)
	}

	sm.apps[app.id] = app
	return nil
}

func (sm *ShutdownManager) Unregister(appID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.apps[appID]; !exists {
		return fmt.Errorf("app with ID %s not found", appID)
	}

	delete(sm.apps, appID)
	return nil
}

func (sm *ShutdownManager) Shutdown(ctx context.Context) error {
	sm.mu.RLock()
	apps := make([]*App, 0, len(sm.apps))
	for _, app := range sm.apps {
		apps = append(apps, app)
	}
	sm.mu.RUnlock()

	if len(apps) == 0 {
		return nil
	}

	type result struct {
		appID string
		err   error
	}

	results := make(chan result, len(apps))
	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(a *App) {
			defer wg.Done()
			if err := a.Stop(ctx); err != nil {
				results <- result{appID: a.id, err: err}
			} else {
				results <- result{appID: a.id, err: nil}
			}
		}(app)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var errors []error
	for res := range results {
		if res.err != nil {
			errors = append(errors, fmt.Errorf("app %s shutdown failed: %w", res.appID, res.err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}

	return nil
}

func (sm *ShutdownManager) GetAppIDs() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ids := make([]string, 0, len(sm.apps))
	for id := range sm.apps {
		ids = append(ids, id)
	}
	return ids
}

func (sm *ShutdownManager) GetAppCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.apps)
}

func Shutdown(ctx context.Context) error {
	return globalShutdown.Shutdown(ctx)
}

func ShutdownWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return globalShutdown.Shutdown(ctx)
}

func GetApps() []string {
	return globalShutdown.GetAppIDs()
}

func generateAppID() string {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return fmt.Sprintf("app-%d", time.Now().UnixNano())
	}
	return "app-" + hex.EncodeToString(b[:])
}
