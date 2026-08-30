// Package telemetry wires up OpenTelemetry SDK providers and a structured
// slog logger for allbctl.
//
// Three independent signal paths exist and can be active simultaneously:
//
//   - Console path (--debug flag): traces → stdouttrace pretty-print on stderr;
//     metrics → JSON slog records on stderr at shutdown; slog logger → JSON on stderr.
//
//   - OTLP path (OTEL_EXPORTER_OTLP_ENDPOINT env var): traces, metrics, and logs are
//     shipped via OTLP HTTP to the configured endpoint (e.g. a local Grafana
//     LGTM stack). This is always active when the env var is set, regardless of
//     the --debug flag.
//
// When neither path is active, no-op providers are installed so instrumented
// code compiles and runs with zero overhead.
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

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlploghttp "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otlpmetrichttp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

// multiHandler fans slog records out to multiple slog.Handler implementations.
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

// ensure multiHandler satisfies the interface at compile time
var _ slog.Handler = (*multiHandler)(nil)

const instrumentationScope = "github.com/aallbrig/allbctl"

// otlpEndpointEnv is the standard OTel env var for the OTLP HTTP endpoint.
const otlpEndpointEnv = "OTEL_EXPORTER_OTLP_ENDPOINT"

// Logger is the process-wide structured logger. It discards all records until
// Setup is called with debug=true or OTEL_EXPORTER_OTLP_ENDPOINT is set.
var Logger = slog.New(slog.NewTextHandler(io.Discard, nil))

// Setup initialises the OpenTelemetry SDK. Call the returned shutdown function
// (exactly once) to flush and stop all providers.
//
//   - debug=true enables the console (stderr) signal path.
//   - OTEL_EXPORTER_OTLP_ENDPOINT being set enables the OTLP signal path.
//
// If neither is active, no-op providers are installed.
func Setup(ctx context.Context, debug bool) (shutdown func(context.Context) error, err error) {
	otlpEndpoint := os.Getenv(otlpEndpointEnv)
	otlpEnabled := otlpEndpoint != ""

	if !debug && !otlpEnabled {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
		return func(context.Context) error { return nil }, nil
	}

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

	// ── Trace provider ───────────────────────────────────────────────────────
	var traceOpts []sdktrace.TracerProviderOption
	traceOpts = append(traceOpts, sdktrace.WithResource(res))

	if debug {
		consoleExp, cErr := stdouttrace.New(
			stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithPrettyPrint(),
		)
		if cErr != nil {
			return nil, fmt.Errorf("stdouttrace exporter: %w", cErr)
		}
		traceOpts = append(traceOpts, sdktrace.WithBatcher(consoleExp))
	}

	if otlpEnabled {
		otlpExp, oErr := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpointURL(otlpEndpoint+"/v1/traces"),
			otlptracehttp.WithInsecure(),
		)
		if oErr != nil {
			// Non-fatal: log and continue without OTLP traces
			fmt.Fprintf(os.Stderr, "OTLP trace exporter warning: %v\n", oErr)
		} else {
			traceOpts = append(traceOpts, sdktrace.WithBatcher(otlpExp))
		}
	}

	tp := sdktrace.NewTracerProvider(traceOpts...)
	otel.SetTracerProvider(tp)

	// ── Metric provider ──────────────────────────────────────────────────────
	var metricReaders []sdkmetric.Option
	metricReaders = append(metricReaders, sdkmetric.WithResource(res))

	// ManualReader lets us flush metrics as slog JSON records at shutdown.
	mReader := sdkmetric.NewManualReader()
	metricReaders = append(metricReaders, sdkmetric.WithReader(mReader))

	if otlpEnabled {
		otlpMetricExp, oErr := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpointURL(otlpEndpoint+"/v1/metrics"),
			otlpmetrichttp.WithInsecure(),
		)
		if oErr != nil {
			fmt.Fprintf(os.Stderr, "OTLP metric exporter warning: %v\n", oErr)
		} else {
			metricReaders = append(metricReaders,
				sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpMetricExp)),
			)
		}
	}

	mp := sdkmetric.NewMeterProvider(metricReaders...)
	otel.SetMeterProvider(mp)

	// ── Log provider ─────────────────────────────────────────────────────────
	// Build a fan-out slog.Logger that routes to every active log sink:
	//   - stderr JSON handler (debug mode)
	//   - OTel log provider → OTLP HTTP exporter (when endpoint is set)
	var logHandlers []slog.Handler

	if debug {
		logHandlers = append(logHandlers,
			slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	var lp *sdklog.LoggerProvider
	if otlpEnabled {
		otlpLogExp, oErr := otlploghttp.New(ctx,
			otlploghttp.WithEndpointURL(otlpEndpoint+"/v1/logs"),
			otlploghttp.WithInsecure(),
		)
		if oErr != nil {
			fmt.Fprintf(os.Stderr, "OTLP log exporter warning: %v\n", oErr)
		} else {
			lp = sdklog.NewLoggerProvider(
				sdklog.WithResource(res),
				sdklog.WithProcessor(sdklog.NewBatchProcessor(otlpLogExp)),
			)
			global.SetLoggerProvider(lp)
			logHandlers = append(logHandlers,
				otelslog.NewHandler(instrumentationScope, otelslog.WithLoggerProvider(lp)),
			)
		}
	}

	switch len(logHandlers) {
	case 0:
		// both debug and otlp paths inactive — keep the discard logger
	case 1:
		Logger = slog.New(logHandlers[0])
	default:
		Logger = slog.New(&multiHandler{handlers: logHandlers})
	}

	return func(ctx context.Context) error {
		var errs []error

		if debug {
			if mErr := flushMetrics(ctx, mReader); mErr != nil {
				errs = append(errs, fmt.Errorf("flush metrics: %w", mErr))
			}
		}
		if mErr := mp.Shutdown(ctx); mErr != nil {
			errs = append(errs, fmt.Errorf("metric provider shutdown: %w", mErr))
		}
		if tErr := tp.Shutdown(ctx); tErr != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", tErr))
		}
		if lp != nil {
			if lErr := lp.Shutdown(ctx); lErr != nil {
				errs = append(errs, fmt.Errorf("log provider shutdown: %w", lErr))
			}
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
