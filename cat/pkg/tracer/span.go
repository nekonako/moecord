package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, event, spanName string, opts ...trace.TracerOption) (context.Context, trace.Span) {
	return otel.Tracer(event).Start(ctx, spanName)
}

func SpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func Finish(span trace.Span) {
	span.End()
}
