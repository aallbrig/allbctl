package cmd

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLICommandsExist verifies that all commands documented in README exist
func TestCLICommandsExist(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "help command",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "status command",
			args:        []string{"status", "--help"},
			expectError: false,
		},
		{
			name:        "list-packages command",
			args:        []string{"list-packages", "--help"},
			expectError: false,
		},
		{
			name:        "computer-setup status",
			args:        []string{"computer-setup", "status", "--help"},
			expectError: false,
		},
		{
			name:        "computer-setup install",
			args:        []string{"computer-setup", "install", "--help"},
			expectError: false,
		},
		{
			name:        "cs alias",
			args:        []string{"cs", "--help"},
			expectError: false,
		},
		{
			name:        "runtimes command",
			args:        []string{"runtimes", "--help"},
			expectError: false,
		},
		{
			name:        "reset command",
			args:        []string{"reset", "--help"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the compiled binary if available, otherwise use go run from project root
			binary := "../bin/allbctl"
			var cmd *exec.Cmd
			if _, err := os.Stat(binary); err == nil {
				cmd = exec.Command(binary, tt.args...)
			} else {
				// Try from current directory (when running tests from project root)
				binary = "./bin/allbctl"
				if _, err := os.Stat(binary); err == nil {
					cmd = exec.Command(binary, tt.args...)
				} else {
					// Fallback to using go run
					args := append([]string{"run", "../main.go"}, tt.args...)
					cmd = exec.Command("go", args...)
					cmd.Dir = ".."
				}
			}

			output, err := cmd.CombinedOutput()
			
			if tt.expectError && err == nil {
				t.Errorf("%s: expected error but got none", tt.name)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: unexpected error: %v\nOutput: %s", tt.name, err, string(output))
			}

			// All help commands should produce some output
			if len(output) == 0 {
				t.Errorf("%s: expected output but got none", tt.name)
			}
		})
	}
}

// TestRootCommandOutput verifies the root command shows all available commands
func TestRootCommandOutput(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Root command failed: %v", err)
	}
	
	// Verify all commands are registered
	commands := []string{"computer-setup", "list-packages", "status", "runtimes", "reset"}
	for _, cmd := range commands {
		found := false
		for _, c := range rootCmd.Commands() {
			if c.Name() == cmd || contains(c.Aliases, cmd) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command %s not found in root command", cmd)
		}
	}
}

// TestComputerSetupAliases verifies all documented aliases work
func TestComputerSetupAliases(t *testing.T) {
	aliases := []string{"computersetup", "cs", "setup"}
	for _, alias := range aliases {
		found := false
		if ComputerSetupCmd.Name() == alias {
			found = true
		}
		for _, a := range ComputerSetupCmd.Aliases {
			if a == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Alias %s not found for computer-setup command", alias)
		}
	}
}

// TestListPackagesFlagExists verifies --detail flag exists
func TestListPackagesFlagExists(t *testing.T) {
	flag := ListPackagesCmd.Flags().Lookup("detail")
	if flag == nil {
		t.Error("--detail flag not found for list-packages command")
	}
	
	// Check short flag
	flagD := ListPackagesCmd.Flags().ShorthandLookup("d")
	if flagD == nil {
		t.Error("-d flag not found for list-packages command")
	}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// TestStatusCommandSections verifies status command output has documented sections
func TestStatusCommandSections(t *testing.T) {
	// Capture output
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

	// Verify documented sections exist
	expectedSections := []string{
		"Host:",
		"Network:",
		"Computer Setup:",
		"Package Managers:",
		"Packages:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Status output missing documented section: %s", section)
		}
	}

	// Verify some of the documented fields
	expectedFields := []string{
		"OS:",
		"Hostname:",
		"Shell:",
		"Terminal:",
		"CPU:",
		"Memory:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Status output missing documented field: %s", field)
		}
	}
}
