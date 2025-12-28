package cmd

import (
	"strings"
	"testing"
)

func TestDetectGamingPlatforms(t *testing.T) {
	checks := detectGamingPlatforms()

	// Test should always return a map (empty if no gaming platforms found)
	if checks == nil {
		t.Error("detectGamingPlatforms returned nil, expected a map")
	}
}

func TestDetectSteamCommand(t *testing.T) {
	cmd := detectSteamCommand()

	// This test is environment-dependent, so we just check the return type
	// If Steam is installed, cmd should be non-nil and have elements
	// If Steam is not installed, cmd should be nil
	if cmd != nil {
		if len(cmd) == 0 {
			t.Error("detectSteamCommand returned empty slice, expected nil or non-empty slice")
		}
	}
}

func TestFormatRuntimesOutputWithGaming(t *testing.T) {
	runtimes := []RuntimeInfo{
		{Name: "Python", Version: "3.12.3", Category: "language"},
		{Name: "Steam", Version: "Steam (installed)", Category: "gaming"},
		{Name: "nvm", Version: "0.40.3", Category: "version-manager"},
	}

	output := formatRuntimesOutput(runtimes)

	if !strings.Contains(output, "Languages:") {
		t.Error("Output should contain 'Languages:' section")
	}
	if !strings.Contains(output, "Python") {
		t.Error("Output should contain Python")
	}
	if !strings.Contains(output, "Gaming Platforms:") {
		t.Error("Output should contain 'Gaming Platforms:' section")
	}
	if !strings.Contains(output, "Steam") {
		t.Error("Output should contain Steam")
	}
	if !strings.Contains(output, "Version Managers:") {
		t.Error("Output should contain 'Version Managers:' section")
	}
	if !strings.Contains(output, "nvm") {
		t.Error("Output should contain nvm")
	}
}

func TestFormatRuntimesOutputNoGaming(t *testing.T) {
	runtimes := []RuntimeInfo{
		{Name: "Python", Version: "3.12.3", Category: "language"},
	}

	output := formatRuntimesOutput(runtimes)

	if strings.Contains(output, "Gaming Platforms:") {
		t.Error("Output should not contain 'Gaming Platforms:' section when no gaming platforms detected")
	}
}
