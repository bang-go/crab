package tracex

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

const (
	ExporterKindStdout = iota
	ExporterKindOltpGrpc
	ExporterKindOltpHttp
)

func NewExporterByOltpGrpc(ctx context.Context, endpoint string, additionalOpts ...otlptracegrpc.Option) (exp sdktrace.SpanExporter, err error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}
	opts = append(opts, additionalOpts...)
	client := otlptracegrpc.NewClient(opts...)
	exp, err = otlptrace.New(ctx, client)
	if err != nil {
		return
	}
	return
}

func NewExporterByOltpHttp(ctx context.Context, endpoint string, additionalOpts ...otlptracehttp.Option) (exp sdktrace.SpanExporter, err error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	}
	opts = append(opts, additionalOpts...)
	client := otlptracehttp.NewClient(opts...)
	exp, err = otlptrace.New(ctx, client)
	if err != nil {
		return
	}
	return
}

func NewExporterByStdout(additionalOpts ...stdouttrace.Option) (exp sdktrace.SpanExporter, err error) {
	return stdouttrace.New(additionalOpts...)
}
