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

	// Check for expected sections (neofetch-style output)
	if !strings.Contains(output, "@") {
		t.Error("Output missing user@hostname header")
	}
	if !strings.Contains(output, "OS:") {
		t.Error("Output missing OS field")
	}
	if !strings.Contains(output, "Network:") {
		t.Error("Output missing Network section")
	}
	if !strings.Contains(output, "Workstation Bootstrap Status:") {
		t.Error("Output missing Workstation Bootstrap Status section")
	}
	if !strings.Contains(output, "Package Managers:") {
		t.Error("Output missing Package Managers section")
	}
	if !strings.Contains(output, "Packages:") {
		t.Error("Output missing Packages section")
	}
}
