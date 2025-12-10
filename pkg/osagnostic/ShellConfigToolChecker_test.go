package osagnostic

import (
	"testing"
)

func TestShellConfigToolChecker_ExtractTools(t *testing.T) {
	// This test verifies the tool extraction patterns work correctly
	checker := &ShellConfigToolChecker{}
	
	testCases := []struct {
		name     string
		line     string
		expected []string
	}{
		{
			name:     "Command substitution with which",
			line:     `complete -C "$(which aws_completer)" aws`,
			expected: []string{"which", "aws_completer"}, // Both are extracted, filtering happens later
		},
		{
			name:     "Command -v check",
			line:     `if command -v tmux &> /dev/null; then`,
			expected: []string{"tmux"},
		},
		{
			name:     "Eval with command substitution",
			line:     `eval "$(gh copilot alias -- zsh)"`,
			expected: []string{"gh"},
		},
		{
			name:     "PATH export with command substitution",
			line:     `export PATH=$PATH:$(go env GOPATH)/bin`,
			expected: []string{"go"},
		},
		{
			name:     "Comment line should be ignored",
			line:     `# $(kubectl completion zsh)`,
			expected: []string{},
		},
		{
			name:     "Empty line",
			line:     ``,
			expected: []string{},
		},
		{
			name:     "Multiple commands",
			line:     `export GOPATH=$(go env GOPATH) && which kubectl`,
			expected: []string{"go", "kubectl"},
		},
		{
			name:     "Source with process substitution",
			line:     `source <(kubectl completion zsh)`,
			expected: []string{"kubectl"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checker.extractToolsFromLine(tc.line)
			
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d tools, got %d: %v", len(tc.expected), len(result), result)
				return
			}
			
			for i, tool := range tc.expected {
				found := false
				for _, r := range result {
					if r == tool {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find tool '%s' but it was not found. Got: %v", tc.expected[i], result)
				}
			}
		})
	}
}

func TestShellConfigToolChecker_IsToolAvailable(t *testing.T) {
	checker := &ShellConfigToolChecker{}
	
	// Test with a command that should always exist
	if !checker.isToolAvailable("ls") {
		t.Error("Expected 'ls' to be available")
	}
	
	// Test with a command that should never exist
	if checker.isToolAvailable("this-command-does-not-exist-xyz123") {
		t.Error("Expected non-existent command to return false")
	}
}
