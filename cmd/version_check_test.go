package cmd

import (
	"strings"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"v1 greater major", "2.0.0", "1.9.9", 1},
		{"v1 less major", "1.0.0", "2.0.0", -1},
		{"v1 greater minor", "1.5.0", "1.4.9", 1},
		{"v1 less minor", "1.2.0", "1.3.0", -1},
		{"v1 greater patch", "1.2.5", "1.2.3", 1},
		{"v1 less patch", "1.2.1", "1.2.3", -1},
		{"different lengths v1 longer", "1.2.3.4", "1.2.3", 1},
		{"different lengths v2 longer", "1.2.3", "1.2.3.1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestCheckNodeJSUpdate(t *testing.T) {
	tests := []struct {
		name    string
		current string
	}{
		{"current LTS", "24.11.1"},
		{"old version", "18.0.0"},
		{"very old version", "14.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := checkNodeJSUpdate(tt.current)
			if update != nil {
				t.Logf("Node.js update available: %s → %s (LTS: %v)", update.Current, update.Available, update.IsLTS)
			} else {
				t.Logf("Node.js %s is up to date", tt.current)
			}
		})
	}
}

func TestCheckPythonUpdate(t *testing.T) {
	tests := []struct {
		name    string
		current string
	}{
		{"recent version", "3.12.3"},
		{"old version", "3.9.0"},
		{"very old version", "2.7.18"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := checkPythonUpdate(tt.current)
			if update != nil {
				t.Logf("Python update available: %s → %s (LTS: %v)", update.Current, update.Available, update.IsLTS)
			} else {
				t.Logf("Python %s is up to date", tt.current)
			}
		})
	}
}

func TestCheckGoUpdate(t *testing.T) {
	tests := []struct {
		name    string
		current string
	}{
		{"recent version", "1.25.5"},
		{"old version", "1.20.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := checkGoUpdate(tt.current)
			if update != nil {
				t.Logf("Go update available: %s → %s", update.Current, update.Available)
			} else {
				t.Logf("Go %s is up to date", tt.current)
			}
		})
	}
}

func TestCheckJavaUpdate(t *testing.T) {
	tests := []struct {
		name    string
		current string
	}{
		{"LTS version", "21.0.9"},
		{"old LTS", "17.0.1"},
		{"very old", "11.0.1"},
		{"ancient", "8.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := checkJavaUpdate(tt.current)
			if update != nil {
				t.Logf("Java update available: %s → %s (LTS: %v)", update.Current, update.Available, update.IsLTS)
			} else {
				t.Logf("Java %s is up to date", tt.current)
			}
		})
	}
}

func TestCheckVersionUpdate(t *testing.T) {
	tests := []struct {
		name    string
		tool    string
		current string
	}{
		{"Node.js", "Node.js", "24.11.1"},
		{"Python", "Python", "3.12.3"},
		{"Go", "Go", "1.25.5"},
		{"npm", "npm", "11.6.2"},
		{"ollama", "ollama", "0.13.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := checkVersionUpdate(tt.tool, tt.current)
			if update != nil {
				t.Logf("%s update available: %s → %s (LTS: %v)", tt.tool, update.Current, update.Available, update.IsLTS)
			} else {
				t.Logf("%s %s: no update info or up to date", tt.tool, tt.current)
			}
		})
	}
}

func TestFormatVersionWithUpdate_WithUpdate(t *testing.T) {
	// Test with a version that definitely has an update
	result := formatVersionWithUpdate("Node.js", "18.0.0")
	t.Logf("formatVersionWithUpdate(Node.js, 18.0.0) = %s", result)

	// Should either show update arrow or original version
	if !strings.Contains(result, "18.0.0") {
		t.Errorf("Result should contain original version: %s", result)
	}
}

func TestExtractDatabaseVersion(t *testing.T) {
	tests := []struct {
		name     string
		dbName   string
		output   string
		expected string
	}{
		{
			name:     "sqlite3 version",
			dbName:   "sqlite3",
			output:   "3.45.1 2024-01-30 16:01:20 e876e51a0ed5c5b3126f52e532044363a014bc594cfefa87ffb5b82257cc467a",
			expected: "3.45.1",
		},
		{
			name:     "mysql version",
			dbName:   "mysql",
			output:   "mysql  Ver 8.0.39 for Linux on x86_64 (MySQL Community Server - GPL)",
			expected: "8.0.39",
		},
		{
			name:     "postgres version",
			dbName:   "postgres",
			output:   "psql (PostgreSQL) 16.6",
			expected: "16.6",
		},
		{
			name:     "redis version",
			dbName:   "redis",
			output:   "redis-cli 7.0.15",
			expected: "7.0.15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDatabaseVersion(tt.dbName, tt.output)
			if result != tt.expected {
				t.Errorf("extractDatabaseVersion(%s, %q) = %q, want %q", tt.dbName, tt.output, result, tt.expected)
			}
		})
	}
}
