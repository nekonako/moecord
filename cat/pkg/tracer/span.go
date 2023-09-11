package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, spanName string, opts ...trace.TracerOption) (context.Context, trace.Span) {
	return otel.Tracer("", opts...).Start(ctx, spanName)
}

func SpanFromContext(ctx context.Context, name string) trace.Span {
	span := trace.SpanFromContext(ctx)
	span.SetName(name)
	return span
}

func SpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func Finish(span trace.Span) {
	span.End()
}
