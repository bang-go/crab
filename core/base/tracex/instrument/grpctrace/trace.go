package grpctrace

import (
	"github.com/bang-go/opt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type Filter = otelgrpc.Filter

type Options struct {
	filter Filter
}

func WithFilter(f Filter) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.filter = f
	})
}
