package cmd

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

// captureOutput redirects os.Stdout to a pipe, runs fn, and returns what was printed.
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic("captureOutput: failed to create pipe: " + err.Error())
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var sb strings.Builder
	if _, err := io.Copy(&sb, r); err != nil {
		panic(err)
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// Command registration smoke tests
// ---------------------------------------------------------------------------

// TestAllStatusSubcommandsRegistered verifies every status subcommand is wired up.
func TestAllStatusSubcommandsRegistered(t *testing.T) {
	expected := []string{
		"runtimes",
		"list-packages",
		"projects",
		"db",
		"network",
		"containers",
		"security",
		"systemctl",
		"git",
		"ports",
		"cloud-native",
	}

	registered := map[string]bool{}
	for _, c := range StatusCmd.Commands() {
		registered[c.Name()] = true
	}

	for _, name := range expected {
		if !registered[name] {
			t.Errorf("status subcommand %q is not registered", name)
		}
	}
}

// ---------------------------------------------------------------------------
// containers
// ---------------------------------------------------------------------------

// TestPrintContainersInfo verifies the function produces output without panicking.
func TestPrintContainersInfo(t *testing.T) {
	output := captureOutput(PrintContainersInfo)
	if !strings.Contains(output, "Containers") {
		t.Errorf("PrintContainersInfo() output missing 'Containers'\noutput:\n%s", output)
	}
}

// TestContainersCmdRegistered verifies ContainersCmd is a valid cobra command.
func TestContainersCmdRegistered(t *testing.T) {
	if ContainersCmd.Use == "" {
		t.Error("ContainersCmd.Use is empty")
	}
	if ContainersCmd.Short == "" {
		t.Error("ContainersCmd.Short is empty")
	}
}

// ---------------------------------------------------------------------------
// db
// ---------------------------------------------------------------------------

// TestPrintDatabaseSummaryForStatus verifies the function doesn't panic.
func TestPrintDatabaseSummaryForStatus(t *testing.T) {
	output := captureOutput(PrintDatabaseSummaryForStatus)
	// May print nothing if no databases found — just verify no panic.
	t.Logf("PrintDatabaseSummaryForStatus() output length: %d", len(output))
}

// TestDbCmdFlags verifies --detail / -d flags exist.
func TestDbCmdFlags(t *testing.T) {
	flag := DbCmd.Flags().Lookup("detail")
	if flag == nil {
		t.Error("--detail flag not found on DbCmd")
	}
	flagD := DbCmd.Flags().ShorthandLookup("d")
	if flagD == nil {
		t.Error("-d shorthand not found on DbCmd")
	}
}

// TestDbCmdRegistered verifies DbCmd is a valid cobra command.
func TestDbCmdRegistered(t *testing.T) {
	if DbCmd.Use == "" {
		t.Error("DbCmd.Use is empty")
	}
}

// ---------------------------------------------------------------------------
// ports
// ---------------------------------------------------------------------------

// TestPrintPortsInfo verifies the function produces output without panicking.
func TestPrintPortsInfo(t *testing.T) {
	output := captureOutput(PrintPortsInfo)
	if !strings.Contains(output, "Ports") {
		t.Errorf("PrintPortsInfo() output missing 'Ports'\noutput:\n%s", output)
	}
}

// TestGatherPortsInfo verifies gatherPortsInfo returns a non-nil struct.
func TestGatherPortsInfo(t *testing.T) {
	info := gatherPortsInfo()
	if info == nil {
		t.Fatal("gatherPortsInfo() returned nil")
	}
	if info.TCPPorts < 0 {
		t.Errorf("gatherPortsInfo().TCPPorts = %d, want >= 0", info.TCPPorts)
	}
	if info.UDPPorts < 0 {
		t.Errorf("gatherPortsInfo().UDPPorts = %d, want >= 0", info.UDPPorts)
	}
	t.Logf("gatherPortsInfo(): TCP=%d UDP=%d total=%d", info.TCPPorts, info.UDPPorts, len(info.Ports))
}

// ---------------------------------------------------------------------------
// systemctl
// ---------------------------------------------------------------------------

// TestPrintSystemctlInfo verifies the function handles non-Linux gracefully.
func TestPrintSystemctlInfo(t *testing.T) {
	output := captureOutput(PrintSystemctlInfo)
	if runtime.GOOS != "linux" {
		if !strings.Contains(output, "only available on Linux") {
			t.Errorf("PrintSystemctlInfo() on non-Linux should say 'only available on Linux'\noutput:\n%s", output)
		}
	} else {
		// On Linux: either shows services or "not found" — just must not be empty.
		if output == "" {
			t.Error("PrintSystemctlInfo() produced no output on Linux")
		}
	}
	t.Logf("PrintSystemctlInfo() output length: %d", len(output))
}

// TestGatherSystemctlInfo verifies gatherSystemctlInfo never returns nil.
func TestGatherSystemctlInfo(t *testing.T) {
	info := gatherSystemctlInfo()
	if info == nil {
		t.Fatal("gatherSystemctlInfo() returned nil")
	}
	if info.SystemRunning < 0 {
		t.Errorf("gatherSystemctlInfo().SystemRunning = %d, want >= 0", info.SystemRunning)
	}
	t.Logf("gatherSystemctlInfo(): sysRunning=%d sysFailed=%d userRunning=%d userFailed=%d",
		info.SystemRunning, info.SystemFailed, info.UserRunning, info.UserFailed)
}

// ---------------------------------------------------------------------------
// git config
// ---------------------------------------------------------------------------

// TestPrintGitConfigInfo verifies the function produces output without panicking.
func TestPrintGitConfigInfo(t *testing.T) {
	output := captureOutput(PrintGitConfigInfo)
	// Either shows git config or "Git is not installed" — must produce something.
	if output == "" {
		t.Error("PrintGitConfigInfo() produced no output")
	}
	t.Logf("PrintGitConfigInfo() output length: %d", len(output))
}

// TestGatherGitConfigInfo verifies gatherGitConfigInfo never returns nil.
func TestGatherGitConfigInfo(t *testing.T) {
	info := gatherGitConfigInfo()
	if info == nil {
		t.Fatal("gatherGitConfigInfo() returned nil")
	}
	// Fields may be empty strings if not configured — that's fine.
	t.Logf("gatherGitConfigInfo(): name=%q email=%q editor=%q",
		info.UserName, info.UserEmail, info.CoreEditor)
}

// ---------------------------------------------------------------------------
// security
// ---------------------------------------------------------------------------

// TestPrintSecurityInfo verifies the function produces output without panicking.
func TestPrintSecurityInfo(t *testing.T) {
	output := captureOutput(PrintSecurityInfo)
	if !strings.Contains(output, "SSH Keys") {
		t.Errorf("PrintSecurityInfo() output missing 'SSH Keys'\noutput:\n%s", output)
	}
}

// TestGatherSecurityInfo verifies gatherSecurityInfo returns an initialized struct.
func TestGatherSecurityInfo(t *testing.T) {
	info := gatherSecurityInfo()
	if info == nil {
		t.Fatal("gatherSecurityInfo() returned nil")
	}
	// Slices may be empty if no keys are loaded — that's valid.
	if info.SSHKeys == nil {
		t.Error("gatherSecurityInfo().SSHKeys is nil, expected initialized slice")
	}
	if info.GPGKeys == nil {
		t.Error("gatherSecurityInfo().GPGKeys is nil, expected initialized slice")
	}
	t.Logf("gatherSecurityInfo(): sshKeys=%d gpgKeys=%d keyring=%q",
		len(info.SSHKeys), len(info.GPGKeys), info.KeyringInfo)
}
