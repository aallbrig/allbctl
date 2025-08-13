package cmd

import (
	"runtime"
	"testing"
)

func TestExists_KnownCommands(t *testing.T) {
	// This test checks detection logic for known package managers
	osType := runtime.GOOS
	var cmds []string
	if osType == "linux" {
		cmds = []string{"apt", "snap", "flatpak"}
	} else if osType == "darwin" {
		cmds = []string{"brew"}
	} else if osType == "windows" {
		cmds = []string{"choco"}
	}
	for _, cmd := range cmds {
		// We can't guarantee all are installed, but exists() should not panic
		exists(cmd)
	}
}

func TestGetPackages_UnknownManager(t *testing.T) {
	result := getPackages("unknown")
	if result != "Unknown package manager." {
		t.Errorf("Expected 'Unknown package manager.', got '%s'", result)
	}
}

func TestRunCmd_InvalidCommand(t *testing.T) {
	output := runCmd("nonexistentcommand123")
	if output == "" {
		t.Error("Expected error output for invalid command")
	}
}

func TestExists_AllSupportedCommands(t *testing.T) {
	osType := runtime.GOOS
	var cmds []string
	if osType == "linux" {
		cmds = []string{"apt", "snap", "flatpak", "dnf", "yum", "zypper", "pacman", "rpm", "dpkg", "apk", "emerge"}
	} else if osType == "darwin" {
		cmds = []string{"brew", "port", "pkgin"}
	} else if osType == "windows" {
		cmds = []string{"choco", "winget", "scoop"}
	}
	for _, cmd := range cmds {
		exists(cmd) // Should not panic
	}
}

func TestGetPackages_AllSupportedManagers(t *testing.T) {
	managers := []string{"apt", "snap", "flatpak", "dnf", "yum", "zypper", "pacman", "rpm", "dpkg", "apk", "emerge", "brew", "macports", "pkgsrc", "choco", "winget", "scoop"}
	for _, m := range managers {
		_ = getPackages(m) // Should not panic or error
	}
}
