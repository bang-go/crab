package vars

import "go.uber.org/atomic"

var (
	DefaultAppName = atomic.NewString("crab service")
)
