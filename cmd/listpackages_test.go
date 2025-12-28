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
	if result != "" {
		t.Errorf("Expected empty string for unknown manager, got '%s'", result)
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
		cmds = []string{"apt-mark", "snap", "flatpak", "dnf", "yum", "pacman"}
	} else if osType == "darwin" {
		cmds = []string{"brew"}
	} else if osType == "windows" {
		cmds = []string{"choco", "winget", "scoop"}
	}
	// Add cross-platform runtime package managers
	cmds = append(cmds, "npm", "pip", "pip3", "gem", "cargo", "go", "pipx")

	for _, cmd := range cmds {
		exists(cmd) // Should not panic
	}
}

func TestGetPackages_AllSupportedManagers(t *testing.T) {
	managers := []string{"apt", "snap", "flatpak", "dnf", "yum", "pacman", "brew", "choco", "winget", "scoop", "npm", "pip", "gem", "cargo", "go", "pipx", "ollama", "vagrant"}
	for _, m := range managers {
		_ = getPackages(m) // Should not panic or error
	}
}

func TestCountPackages_Apt(t *testing.T) {
	output := "package1\npackage2\npackage3\n"
	count := countPackages("apt", output)
	if count != 3 {
		t.Errorf("Expected 3 packages for apt, got %d", count)
	}
}

func TestCountPackages_Npm(t *testing.T) {
	output := "/home/user/.nvm/versions/node/v20.0.0/lib\n├── package1@1.0.0\n├── package2@2.0.0\n└── package3@3.0.0\n"
	count := countPackages("npm", output)
	if count != 3 {
		t.Errorf("Expected 3 packages for npm, got %d", count)
	}
}

func TestCountPackages_Pip(t *testing.T) {
	output := "Package    Version\n---------- -------\npkg1       1.0.0\npkg2       2.0.0\n"
	count := countPackages("pip", output)
	if count != 2 {
		t.Errorf("Expected 2 packages for pip, got %d", count)
	}
}

func TestCountPackages_Error(t *testing.T) {
	output := "Error running command: exit status 1"
	count := countPackages("apt", output)
	if count != 0 {
		t.Errorf("Expected 0 packages for error output, got %d", count)
	}
}

func TestCountPackages_EmptyOutput(t *testing.T) {
	output := ""
	count := countPackages("apt", output)
	if count != 0 {
		t.Errorf("Expected 0 packages for empty output, got %d", count)
	}
}

func TestCountPackages_Pipx(t *testing.T) {
	output := "   package    pkg1\n   package    pkg2\n"
	count := countPackages("pipx", output)
	if count != 2 {
		t.Errorf("Expected 2 packages for pipx, got %d", count)
	}
}

func TestCountPackages_Ollama(t *testing.T) {
	output := "NAME                      ID              SIZE      MODIFIED\nllama3.2:latest          a80c4f17acd5    2.0 GB    3 days ago\ncodellama:latest         8fdf8f752f6e    3.8 GB    2 weeks ago\nqwen2.5-coder:latest     6c701bcd39d9    4.7 GB    1 week ago\n"
	count := countPackages("ollama", output)
	if count != 3 {
		t.Errorf("Expected 3 models for ollama, got %d", count)
	}
}

func TestCountPackages_Ollama_EmptyOutput(t *testing.T) {
	output := "NAME    ID    SIZE    MODIFIED\n"
	count := countPackages("ollama", output)
	if count != 0 {
		t.Errorf("Expected 0 models for empty ollama output, got %d", count)
	}
}

func TestGetPackages_Ollama(t *testing.T) {
	_ = getPackages("ollama") // Should not panic
}

func TestCountPackages_Vagrant(t *testing.T) {
	output := "gusztavvargadr/windows-10 (virtualbox, 2511.0.0, (amd64))\nubuntu/focal64            (virtualbox, 20240821.0.0)\n"
	count := countPackages("vagrant", output)
	if count != 2 {
		t.Errorf("Expected 2 VMs for vagrant, got %d", count)
	}
}

func TestCountPackages_Vagrant_EmptyOutput(t *testing.T) {
	output := ""
	count := countPackages("vagrant", output)
	if count != 0 {
		t.Errorf("Expected 0 VMs for empty vagrant output, got %d", count)
	}
}

func TestGetPackages_Vagrant(t *testing.T) {
	_ = getPackages("vagrant") // Should not panic
}

func TestGetQueryCommand_Vagrant(t *testing.T) {
	cmd := getQueryCommand("vagrant")
	expected := "vagrant box list"
	if cmd != expected {
		t.Errorf("Expected query command '%s', got '%s'", expected, cmd)
	}
}
