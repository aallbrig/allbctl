package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Bootstrap command registration tests
// ---------------------------------------------------------------------------

func TestBootstrapCmdRegistered(t *testing.T) {
	if BootstrapCmd == nil {
		t.Fatal("BootstrapCmd is nil")
	}
	if BootstrapCmd.Use != "bootstrap" {
		t.Errorf("BootstrapCmd.Use = %q, want %q", BootstrapCmd.Use, "bootstrap")
	}
}

func TestBootstrapSubcommandsRegistered(t *testing.T) {
	want := []string{"status", "install", "reset"}
	registered := make(map[string]bool)
	for _, sub := range BootstrapCmd.Commands() {
		registered[sub.Use] = true
	}
	for _, name := range want {
		if !registered[name] {
			t.Errorf("bootstrap subcommand %q not registered", name)
		}
	}
}

// ---------------------------------------------------------------------------
// Bootstrap status smoke test
// ---------------------------------------------------------------------------

func TestPrintBootstrapStatus_Runs(t *testing.T) {
	output := captureOutput(func() {
		printBootstrapStatus()
	})
	// Should produce some output — either the status table or "No configuration provider"
	if len(output) == 0 {
		t.Error("printBootstrapStatus() produced no output")
	}
	t.Logf("bootstrap status output:\n%s", output)
}

func TestPrintBootstrapStatus_ContainsHeader(t *testing.T) {
	output := captureOutput(func() {
		printBootstrapStatus()
	})
	// On a supported OS, should print "Workstation Bootstrap Status:"
	// On an unsupported OS, should print "No configuration provider for ..."
	if !strings.Contains(output, "Bootstrap") && !strings.Contains(output, "configuration provider") {
		t.Errorf("printBootstrapStatus() output missing expected header; got:\n%s", output)
	}
}

// ---------------------------------------------------------------------------
// Bootstrap command aliases
// ---------------------------------------------------------------------------

func TestBootstrapCmdAliases(t *testing.T) {
	aliases := BootstrapCmd.Aliases
	found := false
	for _, a := range aliases {
		if a == "bs" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("BootstrapCmd.Aliases = %v, want to contain %q", aliases, "bs")
	}
}

// ---------------------------------------------------------------------------
// Install flag registration
// ---------------------------------------------------------------------------

func TestBootstrapInstallFlagRegistered(t *testing.T) {
	var installCmd *cobra.Command
	for _, sub := range BootstrapCmd.Commands() {
		if sub.Use == "install" {
			installCmd = sub
			break
		}
	}
	if installCmd == nil {
		t.Fatal("bootstrap install subcommand not found")
	}
	flag := installCmd.Flags().Lookup("register-ssh-keys")
	if flag == nil {
		t.Error("bootstrap install --register-ssh-keys flag not registered")
	}
}
