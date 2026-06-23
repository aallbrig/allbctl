// Package telemetry wires up OpenTelemetry SDK providers and a structured
// slog logger for allbctl. When debug is false every provider is a no-op so
// instrumented code compiles and runs with zero overhead. When debug is true
// traces are written to stderr via the stdouttrace exporter, metrics are
// collected in-memory and flushed as JSON slog records on shutdown, and the
// slog logger writes JSON to stderr.
package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

const instrumentationScope = "github.com/aallbrig/allbctl"

// Logger is the process-wide structured logger. It is set to a discard handler
// until Setup is called.
var Logger = slog.New(slog.NewTextHandler(io.Discard, nil))

// Setup initialises the OpenTelemetry SDK. Call the returned shutdown function
// (exactly once) to flush and stop all providers. If debug is false, no-op
// providers are installed and Logger discards all records.
func Setup(ctx context.Context, debug bool) (shutdown func(context.Context) error, err error) {
	if !debug {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
		return func(context.Context) error { return nil }, nil
	}

	Logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("allbctl"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("build otel resource: %w", err)
	}

	// ── Traces ──────────────────────────────────────────────────────────────
	traceExp, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stderr),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, fmt.Errorf("stdouttrace exporter: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// ── Metrics ─────────────────────────────────────────────────────────────
	// ManualReader lets us collect and log metrics as structured JSON on shutdown
	// instead of needing a separate stdout metric exporter in the vendor tree.
	mReader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(mReader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return func(ctx context.Context) error {
		var errs []error

		// Flush metrics → slog JSON records on stderr
		if mErr := flushMetrics(ctx, mReader); mErr != nil {
			errs = append(errs, fmt.Errorf("flush metrics: %w", mErr))
		}
		if mErr := mp.Shutdown(ctx); mErr != nil {
			errs = append(errs, fmt.Errorf("metric provider shutdown: %w", mErr))
		}

		// Flush traces
		if tErr := tp.Shutdown(ctx); tErr != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", tErr))
		}

		return errors.Join(errs...)
	}, nil
}

// flushMetrics collects all in-memory metrics and writes them as slog records.
func flushMetrics(ctx context.Context, reader *sdkmetric.ManualReader) error {
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		return err
	}

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			logMetric(m)
		}
	}
	return nil
}

// logMetric writes a single OTel metric as a structured slog record.
func logMetric(m metricdata.Metrics) {
	base := []slog.Attr{
		slog.String("name", m.Name),
		slog.String("description", m.Description),
	}

	log := func(extra ...slog.Attr) func(attrs attribute.Set) {
		return func(attrs attribute.Set) {
			all := append(append(base, extra...), attrsToSlogAttr(attrs)...) //nolint:gocritic
			Logger.LogAttrs(context.Background(), slog.LevelInfo, "otel.metric", all...)
		}
	}

	switch data := m.Data.(type) {
	case metricdata.Sum[int64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "sum"), slog.Int64("value", dp.Value), slog.Time("time", dp.Time))(dp.Attributes)
		}
	case metricdata.Sum[float64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "sum"), slog.Float64("value", dp.Value), slog.Time("time", dp.Time))(dp.Attributes)
		}
	case metricdata.Gauge[int64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "gauge"), slog.Int64("value", dp.Value), slog.Time("time", dp.Time))(dp.Attributes)
		}
	case metricdata.Gauge[float64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "gauge"), slog.Float64("value", dp.Value), slog.Time("time", dp.Time))(dp.Attributes)
		}
	case metricdata.Histogram[int64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "histogram"), slog.Int64("sum", dp.Sum), slog.Uint64("count", dp.Count), slog.Time("time", dp.Time))(dp.Attributes)
		}
	case metricdata.Histogram[float64]:
		for _, dp := range data.DataPoints {
			log(slog.String("kind", "histogram"), slog.Float64("sum", dp.Sum), slog.Uint64("count", dp.Count), slog.Time("time", dp.Time))(dp.Attributes)
		}
	default:
		raw, _ := json.Marshal(m) //nolint:errcheck
		Logger.Info("otel.metric", "name", m.Name, "raw", string(raw))
	}
}

// attrsToSlogAttr converts an OTel attribute set to a slice of slog.Attr.
func attrsToSlogAttr(attrs attribute.Set) []slog.Attr {
	var out []slog.Attr
	iter := attrs.Iter()
	for iter.Next() {
		kv := iter.Attribute()
		out = append(out, slog.Any(string(kv.Key), kv.Value.AsInterface()))
	}
	return out
}

// RecordCommandMetrics records invocation count and duration for a command.
// It is called by the root PersistentPostRunE after each command completes.
func RecordCommandMetrics(ctx context.Context, commandPath string, duration time.Duration, success bool) {
	meter := otel.Meter(instrumentationScope)

	invocations, err := meter.Int64Counter(
		"allbctl.command.invocations",
		metric.WithDescription("Total number of allbctl command invocations"),
	)
	if err == nil {
		invocations.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("command", commandPath),
				attribute.Bool("success", success),
			),
		)
	}

	durationHist, err := meter.Int64Histogram(
		"allbctl.command.duration_ms",
		metric.WithDescription("Duration of allbctl command execution in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err == nil {
		durationHist.Record(ctx, duration.Milliseconds(),
			metric.WithAttributes(
				attribute.String("command", commandPath),
				attribute.Bool("success", success),
			),
		)
	}
}
