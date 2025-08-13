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
	r, w, _ := os.Pipe()
	os.Stdout = w

	printSystemInfo()

	w.Close()
	os.Stdout = oldStdout

	var sb strings.Builder
	_, err := io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := sb.String()

	// Accept either single or double quotes for the tip
	if !strings.Contains(output, "Host:") {
		t.Error("Output missing Host section")
	}
	if !strings.Contains(output, "Installed Software:") {
		t.Error("Output missing Installed Software section")
	}
	if !strings.Contains(output, "Tip: Run 'allbctl list-packages'") && !strings.Contains(output, "Tip: Run \"allbctl list-packages\"") {
		t.Error("Output missing tip for list-packages")
	}
}
