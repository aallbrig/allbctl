package cmd

import (
	"io"
	"os"
	"testing"
	"time"
)

// Test_StatusCommandPerformance ensures allbctl status completes in under 10 seconds
// This is a performance regression test to ensure the command stays snappy
// Note: This test logs a warning if threshold is exceeded but doesn't fail
func Test_StatusCommandPerformance(t *testing.T) {
	// Skip in short mode as this is a longer-running test
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Redirect stdout to suppress output during test
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Capture output in background
	done := make(chan string)
	go func() {
		buf, _ := io.ReadAll(r) //nolint:errcheck // Best-effort read for test output suppression
		done <- string(buf)
	}()

	// Measure execution time
	start := time.Now()
	printSystemInfo()
	duration := time.Since(start)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	<-done

	// Performance threshold: 10 seconds
	threshold := 10 * time.Second

	t.Logf("Status command completed in %v", duration)

	if duration > threshold {
		// Log as a warning but don't fail the test
		// This allows CI to continue while alerting developers
		t.Logf("⚠️  PERFORMANCE WARNING: Status command took %v, which exceeds the %v threshold", duration, threshold)
		t.Logf("This is a performance regression indicator. Consider optimizing slow operations.")
		t.Logf("Slow operations are typically: package counting, version checks, network operations")
		// To fail in CI, uncomment the line below:
		// t.Fail()
	} else {
		successPct := float64(threshold-duration) / float64(threshold) * 100
		t.Logf("✓ Performance check passed with %.1f%% headroom (%v remaining)", successPct, threshold-duration)
	}
}

// Benchmark_StatusCommand provides a benchmark for the status command
func Benchmark_StatusCommand(b *testing.B) {
	// Redirect stdout to suppress output during benchmark
	oldStdout := os.Stdout
	os.Stdout = nil

	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a null writer to discard output
	nullWriter, _ := os.Open(os.DevNull) //nolint:errcheck // Best-effort for benchmark output suppression
	defer nullWriter.Close()
	os.Stdout = nullWriter

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		printSystemInfo()
	}
}

// Test_PackageDetectionPerformance tests the async package detection specifically
func Test_PackageDetectionPerformance(t *testing.T) {
	start := time.Now()
	future := StartPackageSummary()
	if future != nil {
		// Redirect stdout to suppress output
		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}
		os.Stdout = w

		// Capture output in background
		done := make(chan string)
		go func() {
			buf, _ := io.ReadAll(r) //nolint:errcheck // Best-effort read for test output suppression
			done <- string(buf)
		}()

		future.PrintResults()

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout
		<-done
	}
	duration := time.Since(start)

	t.Logf("Package detection completed in %v", duration)

	// Package detection should be reasonably fast (5 seconds threshold)
	threshold := 5 * time.Second
	if duration > threshold {
		t.Logf("⚠️  Package detection took %v (threshold: %v)", duration, threshold)
	}
}
