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
	SpanFromContext(ctx).End()
}

func AddAttribute(ctx context.Context, kv ...attribute.KeyValue) {
	SpanFromContext(ctx).SetAttributes(kv...)
}

func AddEvent(ctx context.Context, name string, options ...trace.EventOption) {
	SpanFromContext(ctx).AddEvent(name, options...)
}

func AddException(ctx context.Context, err error, options ...trace.EventOption) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, options...)
	span.SetStatus(codes.Error, err.Error())
}
