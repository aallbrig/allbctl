package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFindGitRepos(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "allbctl-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock git repos
	repo1 := filepath.Join(tmpDir, "repo1", ".git")
	repo2 := filepath.Join(tmpDir, "repo2", ".git")
	nestedRepo := filepath.Join(tmpDir, "parent", "nested-repo", ".git")

	//nolint:errcheck // Test setup errors are not critical
	_ = os.MkdirAll(repo1, 0755)
	//nolint:errcheck // Test setup
	_ = os.MkdirAll(repo2, 0755)
	//nolint:errcheck // Test setup
	_ = os.MkdirAll(nestedRepo, 0755)
	//nolint:errcheck // Test setup

	// Create non-repo directory
	_ = os.MkdirAll(filepath.Join(tmpDir, "not-a-repo"), 0755) //nolint:errcheck // Test setup

	repos := findGitRepos(tmpDir)
	if len(repos) != 3 {
		t.Errorf("Expected 3 repos, got %d", len(repos))
	}
}

func TestIsGitRepoDirty(t *testing.T) {
	// This test checks that the function doesn't panic
	// Actual behavior depends on git being available and having a repo
	tmpDir, err := os.MkdirTemp("", "allbctl-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with non-git directory
	dirty := isGitRepoDirty(tmpDir)
	if dirty {
		t.Error("Non-git directory should not be dirty")
	}
}

func TestGetReposByModTime(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "allbctl-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directories with different mod times
	repo1 := filepath.Join(tmpDir, "repo1")
	repo2 := filepath.Join(tmpDir, "repo2")
	repo3 := filepath.Join(tmpDir, "repo3")

	_ = os.MkdirAll(repo1, 0755) //nolint:errcheck // Test setup
	time.Sleep(10 * time.Millisecond)
	_ = os.MkdirAll(repo2, 0755) //nolint:errcheck // Test setup
	time.Sleep(10 * time.Millisecond)
	_ = os.MkdirAll(repo3, 0755) //nolint:errcheck // Test setup

	repos := []string{repo1, repo2, repo3}
	sorted := getReposByModTime(repos)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 repos, got %d", len(sorted))
	}

	// Most recently touched should be first
	if sorted[0].Path != repo3 {
		t.Errorf("Expected repo3 first (most recent), got %s", sorted[0].Path)
	}
}

func TestFormatRepoPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	testPath := filepath.Join(home, "src", "myproject")

	formatted := formatRepoPath(testPath, false)
	expected := "~/src/myproject"
	if formatted != expected {
		t.Errorf("Expected %s, got %s", expected, formatted)
	}

	formattedDirty := formatRepoPath(testPath, true)
	expectedDirty := "~/src/myproject*"
	if formattedDirty != expectedDirty {
		t.Errorf("Expected %s, got %s", expectedDirty, formattedDirty)
	}
}

func TestGetRemoteRepo(t *testing.T) {
	// Test various git remote URL formats
	testCases := []struct {
		url      string
		expected string
	}{
		{"https://github.com/aallbrig/allbctl.git", "aallbrig/allbctl"},
		{"git@github.com:godotengine/godot.git", "godotengine/godot"},
		{"https://github.com/user/repo", "user/repo"},
		{"git@gitlab.com:user/repo.git", "user/repo"},
		{"https://gitlab.com/user/repo.git", "user/repo"},
		{"invalid-url", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		result := parseRemoteRepo(tc.url)
		if result != tc.expected {
			t.Errorf("parseRemoteRepo(%q) = %q, expected %q", tc.url, result, tc.expected)
		}
	}
}

func TestFormatRepoLine(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	testPath := filepath.Join(home, "src", "allbctl")
	testTime := time.Date(2024, 12, 23, 15, 30, 0, 0, time.Local)

	repo := RepoInfo{
		Path:       testPath,
		Dirty:      true,
		RemoteRepo: "aallbrig/allbctl",
		ModTime:    testTime,
	}

	line := formatRepoLine(repo)
	// Should contain path, *, remote, and date
	if !strings.Contains(line, "~/src/allbctl*") {
		t.Errorf("Expected line to contain path with dirty marker, got: %s", line)
	}
	if !strings.Contains(line, "aallbrig/allbctl") {
		t.Errorf("Expected line to contain remote repo, got: %s", line)
	}
	if !strings.Contains(line, "2024-12-23") {
		t.Errorf("Expected line to contain date, got: %s", line)
	}
}

func TestFilterRepos(t *testing.T) {
	repos := []RepoInfo{
		{Path: "/path/1", Dirty: true, RemoteRepo: "user/repo1"},
		{Path: "/path/2", Dirty: false, RemoteRepo: "user/repo2"},
		{Path: "/path/3", Dirty: true, RemoteRepo: "user/repo3"},
	}

	dirty := filterRepos(repos, "dirty")
	if len(dirty) != 2 {
		t.Errorf("Expected 2 dirty repos, got %d", len(dirty))
	}

	clean := filterRepos(repos, "clean")
	if len(clean) != 1 {
		t.Errorf("Expected 1 clean repo, got %d", len(clean))
	}

	all := filterRepos(repos, "all")
	if len(all) != 3 {
		t.Errorf("Expected 3 total repos, got %d", len(all))
	}
}
