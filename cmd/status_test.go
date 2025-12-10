package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintSystemInfo_Output(t *testing.T) {
	// Redirect stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	printSystemInfo()

	w.Close()
	os.Stdout = oldStdout

	var sb strings.Builder
	_, err = io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := sb.String()

	// Check for expected sections
	if !strings.Contains(output, "Host:") {
		t.Error("Output missing Host section")
	}
	if !strings.Contains(output, "Computer Setup:") {
		t.Error("Output missing Computer Setup section")
	}
	if !strings.Contains(output, "Packages:") {
		t.Error("Output missing Packages section")
	}
}
