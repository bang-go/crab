package tracex

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type SpanCreator interface {
	NewSpan(string) (context.Context, trace.Span)
	NewSpanWithContext(context.Context, string) (context.Context, trace.Span)
	ChildSpan(context.Context, string) (context.Context, trace.Span)
}

type spanCreatorEntity struct {
	tracer trace.Tracer
}

func Tracer(name string, opts ...trace.TracerOption) SpanCreator {
	return &spanCreatorEntity{
		tracer: otel.Tracer(name, opts...),
	}
}

func (t *spanCreatorEntity) NewSpan(spanName string) (context.Context, trace.Span) {
	return t.tracer.Start(context.Background(), spanName)
}

func (t *spanCreatorEntity) NewSpanWithContext(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName)
}

func (t *spanCreatorEntity) ChildSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName)
}
