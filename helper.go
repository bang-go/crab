package crab

import (
	"context"
	"fmt"
	"time"

	"github.com/bang-go/crab/pkg/types"
)

// OnStart creates a hook with only a start action.
func OnStart(fn types.Runner) Hook {
	return Hook{OnStart: fn}
}

// OnStop creates a hook with only a stop action.
func OnStop(fn types.Stopper) Hook {
	return Hook{OnStop: fn}
}

// Close creates a hook from a simple close function (no context).
// This is useful for wrapping standard library closers (like sql.DB, os.File).
//
// Example:
//
//	lc.Append(crab.Close(db.Close))
func Close(fn func() error) Hook {
	return Hook{
		OnStop: func(ctx context.Context) error {
			return fn()
		},
	}
}

// CloseWithContext creates a hook from a close function that accepts context.
func CloseWithContext(fn func(context.Context) error) Hook {
	return Hook{
		OnStop: fn,
	}
}

func formatCost(d time.Duration) string {
	if d >= time.Second {
		secs := d.Seconds()
		if secs < 10 {
			return fmt.Sprintf("%.3fs", secs)
		}
		if secs < 100 {
			return fmt.Sprintf("%.2fs", secs)
		}
		return fmt.Sprintf("%.0fs", secs)
	}

	ms := d.Milliseconds()
	if ms > 0 {
		return fmt.Sprintf("%dms", ms)
	}

	us := d.Microseconds()
	if us > 0 {
		return fmt.Sprintf("%dµs", us)
	}

	return fmt.Sprintf("%dns", d.Nanoseconds())
}
