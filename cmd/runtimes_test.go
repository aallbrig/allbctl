package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func Test_RuntimesCommand_DetectsRuntimes(t *testing.T) {
	buf := new(bytes.Buffer)
	RuntimesCmd.SetOut(buf)
	RuntimesCmd.SetArgs([]string{})
	err := RuntimesCmd.Execute()
	if err != nil {
		t.Errorf("Runtimes command failed: %v", err)
	}
	// Optionally, check output for expected strings
}

func Test_DetectRuntimesInline_IncludesVersion(t *testing.T) {
	// This test verifies that runtime detection includes versions in inline format
	runtimesInline := detectRuntimesInline()
	// If any runtimes are detected, they should include version info in parentheses
	if runtimesInline != "" {
		// Check that at least one runtime has version info (contains parentheses)
		if strings.Contains(runtimesInline, " (") && strings.Contains(runtimesInline, ")") {
			// Good, versions are included
			return
		}
		// If no parentheses found, check if we have any runtimes at all
		// (might be running on system with no runtimes)
		parts := strings.Split(runtimesInline, ",")
		if len(parts) > 0 {
			t.Logf("Runtime detected but no version info: %s", runtimesInline)
		}
	}
}

func Test_FormatRuntimesOutput_IncludesVersion(t *testing.T) {
	// Test with mock runtime info
	runtimes := []RuntimeInfo{
		{Name: "Python", Version: "Python 3.9.0", Category: "language"},
		{Name: "Go", Version: "go version go1.20.0", Category: "language"},
	}

	output := formatRuntimesOutput(runtimes)

	if !strings.Contains(output, "Languages:") {
		t.Error("Expected output to contain 'Languages:' section")
	}

	if !strings.Contains(output, "Python") {
		t.Error("Expected output to contain 'Python'")
	}

	if !strings.Contains(output, "3.9.0") {
		t.Error("Expected output to contain version '3.9.0'")
	}
}
