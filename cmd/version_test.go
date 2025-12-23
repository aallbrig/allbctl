package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Create a fresh command to avoid state issues
	cmd := &cobra.Command{
		Use: "allbctl",
	}
	cmd.AddCommand(versionCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "allbctl") {
		t.Errorf("version output should contain 'allbctl', got: %s", output)
	}
	if !strings.Contains(output, "commit") {
		t.Errorf("version output should contain 'commit', got: %s", output)
	}
}

func TestVersionFlag(t *testing.T) {
	// Create a fresh command to avoid state issues
	cmd := &cobra.Command{
		Use:     "allbctl",
		Version: Version,
	}
	cmd.SetVersionTemplate("allbctl " + Version + " (commit " + Commit + ")\n")

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("--version flag failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "allbctl") {
		t.Errorf("version output should contain 'allbctl', got: %s", output)
	}
	if !strings.Contains(output, "commit") {
		t.Errorf("version output should contain 'commit', got: %s", output)
	}
}
