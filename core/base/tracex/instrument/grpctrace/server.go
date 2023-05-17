package grpctrace

import (
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/opt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...opt.Option[options]) grpc.UnaryServerInterceptor {
	o := new(options)
	opt.Each(o, opts...)
	return otelgrpc.UnaryServerInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
}

func StreamServerInterceptor(opts ...opt.Option[options]) grpc.StreamServerInterceptor {
	o := new(options)
	opt.Each(o, opts...)
	return otelgrpc.StreamServerInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
}
