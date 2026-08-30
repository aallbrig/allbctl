package telemetry_test

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/aallbrig/allbctl/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup_NoopWhenDebugFalse(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), false)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Providers should be set and usable without panicking
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	span.End()
	_ = ctx

	meter := otel.Meter("test")
	counter, err := meter.Int64Counter("test.counter")
	require.NoError(t, err)
	counter.Add(context.Background(), 1)

	require.NoError(t, shutdown(context.Background()))
}

func TestSetup_DebugTrue(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), true)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Logger should be non-nil and usable — it routes to the stderr JSON handler
	assert.NotNil(t, telemetry.Logger)
	telemetry.Logger.Info("test structured log", "key", "value")
	telemetry.Logger.InfoContext(context.Background(), "test.with_context",
		"command", "status",
		"count", 3,
	)

	// Create a span and some metrics
	tracer := otel.Tracer("allbctl/test")
	ctx, span := tracer.Start(context.Background(), "test-command")
	span.SetAttributes(attribute.String("command", "test"))
	span.End()
	_ = ctx

	meter := otel.Meter("allbctl/test")
	counter, err := meter.Int64Counter("test.invocations",
		metric.WithDescription("test counter"),
	)
	require.NoError(t, err)
	counter.Add(context.Background(), 3, metric.WithAttributes(attribute.String("cmd", "status")))

	hist, err := meter.Int64Histogram("test.duration_ms",
		metric.WithDescription("test histogram"),
		metric.WithUnit("ms"),
	)
	require.NoError(t, err)
	hist.Record(context.Background(), 42, metric.WithAttributes(attribute.String("cmd", "status")))

	require.NoError(t, shutdown(context.Background()))
}

// TestSetup_LoggerActiveWhenOTLPSet verifies Logger is non-discard when
// OTEL_EXPORTER_OTLP_ENDPOINT is set (even without --debug). The Logger should
// be the otelslog bridge handler; we just verify it doesn't panic.
func TestSetup_LoggerActiveWhenOTLPSet(t *testing.T) {
	// Point at a dummy endpoint — the exporter will fail to connect but Setup
	// should still succeed (non-fatal) and return a working Logger.
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:19999")

	shutdown, err := telemetry.Setup(context.Background(), false)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	assert.NotNil(t, telemetry.Logger)
	// Should not panic even if the OTLP endpoint is unreachable at record time.
	telemetry.Logger.InfoContext(context.Background(), "test.otlp_path",
		"command", "status",
		"runtime_count", 5,
	)

	// Shutdown may return an error due to failed OTLP flush — that's acceptable.
	_ = shutdown(context.Background())
}

func TestRecordCommandMetrics_NoopProviders(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), false)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown(context.Background())) })

	// Should not panic with noop providers
	telemetry.RecordCommandMetrics(context.Background(), "allbctl status", 123*time.Millisecond, true)
	telemetry.RecordCommandMetrics(context.Background(), "allbctl update", 456*time.Millisecond, false)
}

func TestRecordCommandMetrics_RealProviders(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), true)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown(context.Background())) })

	// Should not panic and metrics should be recorded
	telemetry.RecordCommandMetrics(context.Background(), "allbctl status", 200*time.Millisecond, true)
}
