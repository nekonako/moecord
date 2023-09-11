package infra

import (
	"os"
	"time"

	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type zerologHook struct{}

func (t zerologHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	sCtx := span.SpanContext()
	if sCtx.HasTraceID() {
		e.Str("trace_id", sCtx.TraceID().String())
	}
	if sCtx.HasSpanID() {
		e.Str("span_id", sCtx.SpanID().String())
	}

	attrs := make([]attribute.KeyValue, 0)
	logSeverityKey := attribute.Key("log.severity")
	logMessageKey := attribute.Key("log.message")
	attrs = append(attrs, logSeverityKey.String(level.String()))
	attrs = append(attrs, logMessageKey.String(message))

	span.AddEvent("log", trace.WithAttributes(attrs...))
	if level <= zerolog.ErrorLevel {
		span.SetStatus(codes.Error, message)
	} else {
		span.SetStatus(codes.Ok, message)
	}

}

func initLogger(c *config.Config) {

	now := time.Now().Format(time.DateOnly)
	file, err := os.OpenFile(
		"logs/"+now+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	level, err := zerolog.ParseLevel(c.Apm.LogLevel)
	if err != nil {
		panic(err)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(os.Stdout, file)).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(level).
		Hook(&zerologHook{})

}
