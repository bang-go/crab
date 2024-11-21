package grpctrace

//func UnaryServerInterceptor(opts ...opt.Option[Options]) grpc.UnaryServerInterceptor {
//	o := new(Options)
//	opt.Each(o, opts...)
//	return otelgrpc.UnaryServerInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
//}
//
//func StreamServerInterceptor(opts ...opt.Option[Options]) grpc.StreamServerInterceptor {
//	o := new(Options)
//	opt.Each(o, opts...)
//	return otelgrpc.StreamServerInterceptor(otelgrpc.WithPropagators(tracex.Propagator()), otelgrpc.WithTracerProvider(tracex.Provider()), otelgrpc.WithInterceptorFilter(o.filter))
//}
