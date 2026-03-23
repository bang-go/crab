package crab

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunStartupTimeoutRollsBackStartedHooksOnly(t *testing.T) {
	app := New(WithStartupTimeout(50*time.Millisecond), WithShutdownTimeout(time.Second))
	defer func() { _ = app.Unregister() }()

	var mu sync.Mutex
	var events []string
	record := func(event string) {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, event)
	}

	app.Add(Hook{
		Name: "first",
		OnStart: func(ctx context.Context) error {
			record("start:first")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			record("stop:first")
			return nil
		},
	})

	app.Add(Hook{
		Name: "second",
		OnStart: func(ctx context.Context) error {
			record("start:second")
			<-ctx.Done()
			return ctx.Err()
		},
		OnStop: func(ctx context.Context) error {
			record("stop:second")
			return nil
		},
	})

	app.Add(Hook{
		Name: "third",
		OnStart: func(ctx context.Context) error {
			record("start:third")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			record("stop:third")
			return nil
		},
	})

	err := app.Run()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}

	mu.Lock()
	got := append([]string(nil), events...)
	mu.Unlock()

	want := []string{"start:first", "start:second", "stop:first"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected events: got %v want %v", got, want)
	}
}

func TestStopRunsAllHooksInLIFOEvenAfterDeadline(t *testing.T) {
	app := New(WithShutdownTimeout(10 * time.Millisecond))
	defer func() { _ = app.Unregister() }()

	var mu sync.Mutex
	var stopped []string
	record := func(name string) {
		mu.Lock()
		defer mu.Unlock()
		stopped = append(stopped, name)
	}

	app.Add(Hook{Name: "first", OnStop: func(ctx context.Context) error {
		record("first")
		return nil
	}})
	app.Add(Hook{Name: "second", OnStop: func(ctx context.Context) error {
		record("second")
		<-ctx.Done()
		return ctx.Err()
	}})
	app.Add(Hook{Name: "third", OnStop: func(ctx context.Context) error {
		record("third")
		return nil
	}})

	app.mu.Lock()
	app.state = stateRunning
	app.startedHooks = append([]Hook(nil), app.hooks...)
	app.mu.Unlock()

	err := app.Stop(context.Background())
	if err == nil {
		t.Fatal("expected shutdown error")
	}

	mu.Lock()
	got := append([]string(nil), stopped...)
	mu.Unlock()

	want := []string{"third", "second", "first"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected stop order: got %v want %v", got, want)
	}
}

func TestConcurrentStopWaitsForSharedShutdown(t *testing.T) {
	app := New(WithShutdownTimeout(time.Second))
	defer func() { _ = app.Unregister() }()

	release := make(chan struct{})
	var stopCalls atomic.Int32

	app.Add(Hook{Name: "service", OnStop: func(ctx context.Context) error {
		stopCalls.Add(1)
		<-release
		return nil
	}})

	app.mu.Lock()
	app.state = stateRunning
	app.startedHooks = append([]Hook(nil), app.hooks...)
	app.mu.Unlock()

	errCh := make(chan error, 2)
	go func() { errCh <- app.Stop(context.Background()) }()
	go func() { errCh <- app.Stop(context.Background()) }()

	time.Sleep(30 * time.Millisecond)
	close(release)

	for range 2 {
		if err := <-errCh; err != nil {
			t.Fatalf("unexpected stop error: %v", err)
		}
	}

	if got := stopCalls.Load(); got != 1 {
		t.Fatalf("expected OnStop to run once, got %d", got)
	}
}

func TestOnStopCanCallIsRunningWithoutDeadlock(t *testing.T) {
	app := New(WithShutdownTimeout(time.Second))
	defer func() { _ = app.Unregister() }()

	called := make(chan bool, 1)
	app.Add(Hook{Name: "service", OnStop: func(ctx context.Context) error {
		called <- app.IsRunning()
		return nil
	}})

	app.mu.Lock()
	app.state = stateRunning
	app.startedHooks = append([]Hook(nil), app.hooks...)
	app.mu.Unlock()

	done := make(chan error, 1)
	go func() { done <- app.Stop(context.Background()) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected stop error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("stop deadlocked")
	}

	select {
	case running := <-called:
		if running {
			t.Fatal("expected app to report not running during stop")
		}
	default:
		t.Fatal("expected OnStop to run")
	}
}

func TestStopDuringStartupOnlyStopsStartedHooks(t *testing.T) {
	app := New(WithShutdownTimeout(time.Second))
	defer func() { _ = app.Unregister() }()

	var mu sync.Mutex
	var events []string
	record := func(event string) {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, event)
	}

	startedSecond := make(chan struct{})

	app.Add(Hook{
		Name: "first",
		OnStart: func(ctx context.Context) error {
			record("start:first")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			record("stop:first")
			return nil
		},
	})

	app.Add(Hook{
		Name: "second",
		OnStart: func(ctx context.Context) error {
			record("start:second")
			close(startedSecond)
			<-ctx.Done()
			return ctx.Err()
		},
		OnStop: func(ctx context.Context) error {
			record("stop:second")
			return nil
		},
	})

	app.Add(Hook{
		Name: "third",
		OnStart: func(ctx context.Context) error {
			record("start:third")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			record("stop:third")
			return nil
		},
	})

	runDone := make(chan error, 1)
	go func() { runDone <- app.Run() }()

	<-startedSecond
	stopDone := make(chan error, 1)
	go func() { stopDone <- app.Stop(context.Background()) }()

	if err := <-stopDone; err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}
	if err := <-runDone; err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	mu.Lock()
	got := append([]string(nil), events...)
	mu.Unlock()

	want := []string{"start:first", "start:second", "stop:first"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected events: got %v want %v", got, want)
	}
}

func TestShutdownContinuesWithBoundedBestEffortContexts(t *testing.T) {
	app := New(WithShutdownTimeout(200*time.Millisecond), WithBestEffortShutdownTimeout(80*time.Millisecond))
	defer func() { _ = app.Unregister() }()

	var mu sync.Mutex
	var events []string
	record := func(event string) {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, event)
	}

	app.Add(Hook{Name: "first", OnStop: func(ctx context.Context) error {
		record("stop:first")
		return nil
	}})
	app.Add(Hook{Name: "second", OnStop: func(ctx context.Context) error {
		record("stop:second")
		<-ctx.Done()
		return ctx.Err()
	}})
	app.Add(Hook{Name: "third", OnStop: func(ctx context.Context) error {
		record("stop:third")
		return nil
	}})

	app.mu.Lock()
	app.state = stateRunning
	app.startedHooks = append([]Hook(nil), app.hooks...)
	app.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := app.Stop(ctx)
	if err == nil {
		t.Fatal("expected shutdown error")
	}

	mu.Lock()
	got := append([]string(nil), events...)
	mu.Unlock()

	want := []string{"stop:third", "stop:second", "stop:first"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected events: got %v want %v", got, want)
	}
}

func TestIsStoppingReportsDuringShutdown(t *testing.T) {
	app := New(WithShutdownTimeout(time.Second))
	defer func() { _ = app.Unregister() }()

	release := make(chan struct{})
	seenStopping := make(chan bool, 1)
	app.Add(Hook{Name: "service", OnStop: func(ctx context.Context) error {
		seenStopping <- app.IsStopping()
		<-release
		return nil
	}})

	app.mu.Lock()
	app.state = stateRunning
	app.startedHooks = append([]Hook(nil), app.hooks...)
	app.mu.Unlock()

	done := make(chan error, 1)
	go func() { done <- app.Stop(context.Background()) }()

	select {
	case stopping := <-seenStopping:
		if !stopping {
			t.Fatal("expected app to report stopping during OnStop")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for OnStop")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}
}

func TestGlobalShutdownAggregatesErrors(t *testing.T) {
	app1 := New(WithShutdownTimeout(time.Second))
	app2 := New(WithShutdownTimeout(time.Second))
	defer func() { _ = app1.Unregister() }()
	defer func() { _ = app2.Unregister() }()

	app1.Add(Hook{Name: "ok", OnStop: func(ctx context.Context) error { return nil }})
	app2.Add(Hook{Name: "bad", OnStop: func(ctx context.Context) error { return errors.New("boom") }})

	app1.mu.Lock()
	app1.state = stateRunning
	app1.startedHooks = append([]Hook(nil), app1.hooks...)
	app1.mu.Unlock()

	app2.mu.Lock()
	app2.state = stateRunning
	app2.startedHooks = append([]Hook(nil), app2.hooks...)
	app2.mu.Unlock()

	err := Shutdown(context.Background())
	if err == nil {
		t.Fatal("expected shutdown error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected aggregated shutdown error to include boom, got %v", err)
	}
}
