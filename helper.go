package crab

import (
	"context"

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
