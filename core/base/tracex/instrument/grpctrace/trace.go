package grpctrace

import (
	"github.com/bang-go/opt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type Filter = otelgrpc.Filter

type options struct {
	filter Filter
}

func WithFilter(f Filter) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.filter = f
	})
}
