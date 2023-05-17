package tracex

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func SpanEnd(ctx context.Context) {
	span := SpanFromContext(ctx)
	span.End()
}

func AddAttribute(ctx context.Context, kv ...attribute.KeyValue) {
	span := SpanFromContext(ctx)
	span.SetAttributes(kv...)
}

func AddEvent(ctx context.Context, name string, options ...trace.EventOption) {
	span := SpanFromContext(ctx)
	span.AddEvent(name, options...)
}

func AddException(ctx context.Context, err error, options ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, options...)
	span.SetStatus(codes.Error, err.Error())
}
