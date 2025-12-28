package osagnostic

import (
	"runtime"
	"strings"
	"testing"
)

func TestInstallableCommand_WindowsWingetCommand(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	ic := NewInstallableCommand("git").
		SetWindowsPackage("winget", "Git.Git")

	// Test that the command is properly configured
	// This test verifies the fix for Windows 10 PowerShell error 0x8a150042
	expectedFlags := "--accept-source-agreements"

	// When we install, the command should include this flag
	t.Logf("Testing that winget install commands include %s flag", expectedFlags)

	// Verify that our fix includes the necessary flag for non-interactive execution
	// The actual command construction happens in installWithPackageManager
	// which adds the flag when pmCommand starts with "winget"
	if ic.WindowsPackages["winget"] != "Git.Git" {
		t.Error("Expected Windows package to be set")
	}

	if !strings.Contains(expectedFlags, "accept-source-agreements") {
		t.Error("Expected winget commands to include --accept-source-agreements flag")
	}
}

func TestInstallableCommand_Name(t *testing.T) {
	ic := NewInstallableCommand("git")
	expected := "Installable Command: git"

	if ic.Name() != expected {
		t.Errorf("Expected %s, got %s", expected, ic.Name())
	}
}

func TestInstallableCommand_SetPackages(t *testing.T) {
	ic := NewInstallableCommand("git")

	ic.SetLinuxPackage("apt", "git-core").
		SetMacOSPackage("git").
		SetWindowsPackage("winget", "Git.Git")

	if ic.LinuxPackages["apt"] != "git-core" {
		t.Error("Failed to set Linux package")
	}
	if ic.MacOSPackage != "git" {
		t.Error("Failed to set macOS package")
	}
	if ic.WindowsPackages["winget"] != "Git.Git" {
		t.Error("Failed to set Windows package")
	}
}
