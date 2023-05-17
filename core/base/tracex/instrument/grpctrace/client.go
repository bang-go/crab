package grpctrace

import (
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/opt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(opts ...opt.Option[Options]) grpc.UnaryClientInterceptor {
	o := new(Options)
	opt.Each(o, opts...)
	return otelgrpc.UnaryClientInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
}

func StreamClientInterceptor(opts ...opt.Option[Options]) grpc.StreamClientInterceptor {
	o := new(Options)
	opt.Each(o, opts...)
	return otelgrpc.StreamClientInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
}
