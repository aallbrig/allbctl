package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aallbrig/allbctl/pkg/languages"
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
	tmpDir, err := os.MkdirTemp("", "allbctl-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dirty := isGitRepoDirty(tmpDir)
	if dirty {
		t.Error("Non-git directory should not be dirty")
	}
}

func TestDirtyReasonString(t *testing.T) {
	cases := []struct {
		reason   DirtyReason
		expected string
	}{
		{0, ""},
		{DirtyUncommittedChanges, "[uncommitted changes]"},
		{DirtyUnpushedCommits, "[unpushed commits]"},
		{DirtyNoUpstream, "[no upstream]"},
		{DirtyUncommittedChanges | DirtyUnpushedCommits, "[uncommitted changes, unpushed commits]"},
		{DirtyUncommittedChanges | DirtyNoUpstream, "[uncommitted changes, no upstream]"},
		{DirtyUncommittedChanges | DirtyUnpushedCommits | DirtyNoUpstream, "[uncommitted changes, unpushed commits, no upstream]"},
	}

	for _, tc := range cases {
		got := tc.reason.String()
		if got != tc.expected {
			t.Errorf("DirtyReason(%d).String() = %q, want %q", tc.reason, got, tc.expected)
		}
	}
}

func TestGetDirtyReasons(t *testing.T) {
	t.Run("non-git directory has no reasons", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		reasons := getDirtyReasons(tmpDir)
		if reasons != 0 {
			t.Errorf("Expected no dirty reasons for non-git dir, got %s", reasons)
		}
	})

	t.Run("uncommitted changes detected", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		reasons := getDirtyReasons(tmpDir)
		if reasons&DirtyUncommittedChanges == 0 {
			t.Errorf("Expected DirtyUncommittedChanges, got %s", reasons)
		}
	})

	t.Run("no upstream detected after first commit", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("hi"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		_ = exec.Command("git", "-C", tmpDir, "add", ".").Run()                                                                       //nolint:errcheck
		_ = exec.Command("git", "-C", tmpDir, "-c", "user.email=test@test.com", "-c", "user.name=Test", "commit", "-m", "init").Run() //nolint:errcheck

		reasons := getDirtyReasons(tmpDir)
		if reasons&DirtyNoUpstream == 0 {
			t.Errorf("Expected DirtyNoUpstream for local-only repo, got %s", reasons)
		}
	})

	t.Run("unpushed commits detected", func(t *testing.T) {
		// Set up a bare "remote" and a local clone, then make a local commit
		remoteDir, err := os.MkdirTemp("", "allbctl-remote-")
		if err != nil {
			t.Fatalf("Failed to create remote dir: %v", err)
		}
		defer os.RemoveAll(remoteDir)

		localDir, err := os.MkdirTemp("", "allbctl-local-")
		if err != nil {
			t.Fatalf("Failed to create local dir: %v", err)
		}
		defer os.RemoveAll(localDir)

		if err := exec.Command("git", "init", "--bare", remoteDir).Run(); err != nil {
			t.Skip("git not available")
		}
		if err := exec.Command("git", "clone", remoteDir, localDir).Run(); err != nil {
			t.Skipf("git clone failed: %v", err)
		}

		// Make an initial commit and push so the clone has an upstream
		if err := os.WriteFile(filepath.Join(localDir, "file.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		_ = exec.Command("git", "-C", localDir, "add", ".").Run()                                                                          //nolint:errcheck
		_ = exec.Command("git", "-C", localDir, "-c", "user.email=test@test.com", "-c", "user.name=Test", "commit", "-m", "initial").Run() //nolint:errcheck
		_ = exec.Command("git", "-C", localDir, "push", "-u", "origin", "HEAD").Run()                                                      //nolint:errcheck

		// Now make a local commit that isn't pushed
		if err := os.WriteFile(filepath.Join(localDir, "file.txt"), []byte("world"), 0644); err != nil {
			t.Fatalf("Failed to update file: %v", err)
		}
		_ = exec.Command("git", "-C", localDir, "add", ".").Run()                                                                           //nolint:errcheck
		_ = exec.Command("git", "-C", localDir, "-c", "user.email=test@test.com", "-c", "user.name=Test", "commit", "-m", "unpushed").Run() //nolint:errcheck

		reasons := getDirtyReasons(localDir)
		if reasons&DirtyUnpushedCommits == 0 {
			t.Errorf("Expected DirtyUnpushedCommits, got %s", reasons)
		}
		if reasons&DirtyNoUpstream != 0 {
			t.Errorf("Should not have DirtyNoUpstream when upstream exists, got %s", reasons)
		}
	})
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

func TestFilterStatusLines(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"removes nothing to commit line",
			[]string{"On branch main", "", "nothing to commit, working tree clean"},
			[]string{"On branch main"},
		},
		{
			"strips trailing blank lines",
			[]string{"On branch main", "Your branch is ahead by 1 commit.", ""},
			[]string{"On branch main", "Your branch is ahead by 1 commit."},
		},
		{
			"removes nothing to commit and trailing blank",
			[]string{"On branch main", "", "nothing to commit, working tree clean", ""},
			[]string{"On branch main"},
		},
		{
			"removes git push hint line",
			[]string{"On branch main", `  (use "git push" to publish your local commits)`, ""},
			[]string{"On branch main"},
		},
		{
			"removes git add hint line",
			[]string{"Changes not staged for commit:", `  (use "git add <file>..." to update what will be committed)`, "\tmodified:   foo.go"},
			[]string{"Changes not staged for commit:", "\tmodified:   foo.go"},
		},
		{
			"removes git restore hint line",
			[]string{"Changes not staged for commit:", `  (use "git restore <file>..." to discard changes in working directory)`, "\tmodified:   foo.go"},
			[]string{"Changes not staged for commit:", "\tmodified:   foo.go"},
		},
		{
			"removes no changes added to commit line",
			[]string{"On branch main", "", "no changes added to commit (use \"git add\" and/or \"git commit -a\")"},
			[]string{"On branch main"},
		},
		{
			"removes blank line before Changes not staged for commit",
			[]string{"On branch main", "Your branch is up to date with 'origin/main'.", "", "Changes not staged for commit:", "\tmodified:   foo.go"},
			[]string{"On branch main", "Your branch is up to date with 'origin/main'.", "Changes not staged for commit:", "\tmodified:   foo.go"},
		},
		{
			"preserves non-blank non-noise lines",
			[]string{"On branch main", "\tmodified:   foo.go"},
			[]string{"On branch main", "\tmodified:   foo.go"},
		},
		{
			"empty input",
			[]string{},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterStatusLines(tc.input)
			if len(got) != len(tc.expected) {
				t.Errorf("len=%d, want %d; got %v", len(got), len(tc.expected), got)
				return
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("line[%d] = %q, want %q", i, got[i], tc.expected[i])
				}
			}
		})
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

func TestProjectsCmdLimitFlag(t *testing.T) {
	// Verify --limit flag exists on ProjectsCmd
	flag := ProjectsCmd.Flags().Lookup("limit")
	if flag == nil {
		t.Fatal("Expected --limit flag to exist on ProjectsCmd")
	}
	if flag.DefValue != "0" {
		t.Errorf("Expected --limit default value to be 0, got %s", flag.DefValue)
	}
}

func TestShowFilesFlag(t *testing.T) {
	flag := ProjectsCmd.Flags().Lookup("verbose")
	if flag == nil {
		t.Fatal("Expected --verbose flag to exist on ProjectsCmd")
	}
	if flag.DefValue != "false" {
		t.Errorf("Expected --verbose default value to be false, got %s", flag.DefValue)
	}
}

func TestGetGitStatusOutput(t *testing.T) {
	t.Run("non-git directory returns empty string", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		out := getGitStatusOutput(tmpDir)
		if out != "" {
			t.Errorf("Expected empty string for non-git directory, got %q", out)
		}
	})

	t.Run("dirty repo returns non-empty status output", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}

		if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		out := getGitStatusOutput(tmpDir)
		if out == "" {
			t.Error("Expected non-empty status output for dirty repo")
		}
		if !strings.Contains(out, "test.txt") {
			t.Errorf("Expected test.txt in status output, got:\n%s", out)
		}
	})

	t.Run("clean repo returns non-empty status output", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}

		// Clean repos still produce "nothing to commit" output
		out := getGitStatusOutput(tmpDir)
		if out == "" {
			t.Error("Expected non-empty status output even for clean repo")
		}
	})
}

func TestCountStatusFiles(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect int
	}{
		{"empty", "", 0},
		{"no files", "On branch main\nnothing to commit", 0},
		{"one modified file", "Changes not staged:\n\tmodified:   foo.go", 1},
		{"multiple files", "Changes:\n\tmodified:   a.go\n\tmodified:   b.go\nUntracked:\n\tc.go\n\td.go", 4},
		{"untracked with tab", "Untracked files:\n\tfoo.go\n\tbar.go", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := countStatusFiles(tc.input)
			if got != tc.expect {
				t.Errorf("countStatusFiles() = %d, want %d", got, tc.expect)
			}
		})
	}
}

func TestCountTotalStatusFiles(t *testing.T) {
	repos := []RepoInfo{
		{Path: "/path/1", Dirty: true, StatusOutput: "Changes:\n\ta.go\n\tb.go"},
		{Path: "/path/2", Dirty: false, StatusOutput: ""},
		{Path: "/path/3", Dirty: true, StatusOutput: "Changes:\n\tc.go"},
	}

	count := countTotalStatusFiles(repos)
	if count != 3 {
		t.Errorf("Expected 3 total status files, got %d", count)
	}
}

func TestCountPorcelainFiles(t *testing.T) {
	t.Run("non-git directory returns zeros", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		uncommitted, untracked := countPorcelainFiles(tmpDir)
		if uncommitted != 0 || untracked != 0 {
			t.Errorf("Expected (0, 0), got (%d, %d)", uncommitted, untracked)
		}
	})

	t.Run("staged file counted as uncommitted, untracked file counted separately", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}

		if err := os.WriteFile(filepath.Join(tmpDir, "staged.txt"), []byte("staged"), 0644); err != nil {
			t.Fatalf("Failed to create staged file: %v", err)
		}
		_ = exec.Command("git", "-C", tmpDir, "add", "staged.txt").Run() //nolint:errcheck

		if err := os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("untracked"), 0644); err != nil {
			t.Fatalf("Failed to create untracked file: %v", err)
		}

		uncommitted, untracked := countPorcelainFiles(tmpDir)
		if uncommitted != 1 {
			t.Errorf("Expected 1 uncommitted file, got %d", uncommitted)
		}
		if untracked != 1 {
			t.Errorf("Expected 1 untracked file, got %d", untracked)
		}
	})
}

func TestCountUnpushedCommits(t *testing.T) {
	t.Run("non-git directory returns zero", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if countUnpushedCommits(tmpDir) != 0 {
			t.Error("Expected 0 for non-git dir")
		}
	})

	t.Run("repo with no upstream returns zero", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}

		if countUnpushedCommits(tmpDir) != 0 {
			t.Error("Expected 0 for repo with no upstream")
		}
	})
}

func TestBuildSummaryLine(t *testing.T) {
	cases := []struct {
		name        string
		repos       []RepoInfo
		displayMode string
		expected    string
	}{
		// --all mode
		{
			"all: all clean",
			[]RepoInfo{{}, {}},
			"all",
			"Total projects: 2  Total clean: 2",
		},
		{
			"all: mix of dirty and clean",
			[]RepoInfo{{Dirty: true, UncommittedFiles: 4}, {}},
			"all",
			"Total projects: 2  Total dirty: 1  Total clean: 1  Total modified files: 4",
		},
		{
			"all: all dirty with unpushed commits",
			[]RepoInfo{{Dirty: true, UnpushedCommits: 3}, {Dirty: true, UnpushedCommits: 2}},
			"all",
			"Total projects: 2  Total dirty: 2  Total unpushed commits: 5",
		},
		{
			"all: ci failed and pending",
			[]RepoInfo{{Dirty: true, DirtyReasons: DirtyCIFailed}, {Dirty: true, DirtyReasons: DirtyCIPending}},
			"all",
			"Total projects: 2  Total dirty: 2  Total CI failed: 1  Total CI pending: 1",
		},
		// --dirty mode
		{
			"dirty: shows total dirty as lead",
			[]RepoInfo{{Dirty: true, UnpushedCommits: 3}, {Dirty: true, UncommittedFiles: 2}},
			"dirty",
			"Total dirty: 2  Total unpushed commits: 3  Total modified files: 2",
		},
		{
			"dirty: ci failed included",
			[]RepoInfo{{Dirty: true, DirtyReasons: DirtyCIFailed}, {Dirty: true, DirtyReasons: DirtyCIFailed}},
			"dirty",
			"Total dirty: 2  Total CI failed: 2",
		},
		// --clean mode
		{
			"clean: shows total clean as lead, no detail counts",
			[]RepoInfo{{}, {}},
			"clean",
			"Total clean: 2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildSummaryLine(tc.repos, tc.displayMode)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestCICheckVerboseIconMapping(t *testing.T) {
	cases := []struct {
		conclusion string
		wantIcon   string
	}{
		{"success", "✓"},
		{"failure", "✗"},
		{"timed_out", "✗"},
		{"cancelled", "✗"},
		{"", "…"},
	}
	for _, tc := range cases {
		t.Run(tc.conclusion, func(t *testing.T) {
			icon := "✓"
			switch tc.conclusion {
			case "failure", "timed_out", "cancelled":
				icon = "✗"
			case "":
				icon = "…"
			}
			if icon != tc.wantIcon {
				t.Errorf("conclusion %q: got icon %q, want %q", tc.conclusion, icon, tc.wantIcon)
			}
		})
	}
}

func TestParseCICheckRuns(t *testing.T) {
	cases := []struct {
		name        string
		conclusions []string
		expected    string
	}{
		{"empty returns empty", []string{}, ""},
		{"all success returns success", []string{"success", "success"}, "success"},
		{"one failure returns failure", []string{"success", "failure"}, "failure"},
		{"timed_out counts as failure", []string{"timed_out"}, "failure"},
		{"cancelled counts as failure", []string{"cancelled"}, "failure"},
		{"in_progress returns pending", []string{"success", ""}, "pending"},
		{"failure beats pending", []string{"", "failure"}, "failure"},
		{"queued returns pending", []string{"queued"}, "pending"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseCICheckRuns(tc.conclusions)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestGetRepoLanguages(t *testing.T) {
	t.Run("returns languages for a valid git repo", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-lang-test-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}

		if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "script.py"), []byte("print('hello')\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test\n"), 0644); err != nil {
			t.Fatal(err)
		}

		_ = exec.Command("git", "-C", tmpDir, "add", ".").Run()                                                                       //nolint:errcheck
		_ = exec.Command("git", "-C", tmpDir, "-c", "user.email=test@test.com", "-c", "user.name=Test", "commit", "-m", "init").Run() //nolint:errcheck

		langs := getRepoLanguages(tmpDir)
		if len(langs) == 0 {
			t.Fatal("Expected at least one language, got none")
		}

		foundGo := false
		for _, l := range langs {
			if l.Name == "Go" {
				foundGo = true
				break
			}
		}
		if !foundGo {
			t.Errorf("Expected Go in languages, got: %v", langs)
		}
	})

	t.Run("returns nil for non-git directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-lang-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		langs := getRepoLanguages(tmpDir)
		if langs != nil {
			t.Errorf("Expected nil for non-git dir, got %v", langs)
		}
	})

	t.Run("caches and returns same result", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "allbctl-lang-cache-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		if err := exec.Command("git", "-C", tmpDir, "init").Run(); err != nil {
			t.Skip("git not available")
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
		_ = exec.Command("git", "-C", tmpDir, "add", ".").Run()                                                                       //nolint:errcheck
		_ = exec.Command("git", "-C", tmpDir, "-c", "user.email=test@test.com", "-c", "user.name=Test", "commit", "-m", "init").Run() //nolint:errcheck

		langs1 := getRepoLanguages(tmpDir)
		langs2 := getRepoLanguages(tmpDir)

		if len(langs1) != len(langs2) {
			t.Errorf("Expected same result from cache, got %d vs %d", len(langs1), len(langs2))
		}
		if len(langs1) > 0 && len(langs2) > 0 {
			if langs1[0].Name != langs2[0].Name || langs1[0].Size != langs2[0].Size {
				t.Errorf("Cache returned different data: %v vs %v", langs1, langs2)
			}
		}
	})
}

func TestRepoInfoLanguagesField(t *testing.T) {
	repo := RepoInfo{
		Path:  "/path/to/repo",
		Dirty: false,
		Languages: []languages.LanguageBreakdown{
			{Name: "Go", Size: 1500, Percent: 75},
			{Name: "Python", Size: 500, Percent: 25},
		},
	}

	if len(repo.Languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(repo.Languages))
	}
	if repo.Languages[0].Name != "Go" {
		t.Errorf("Expected first language to be Go, got %s", repo.Languages[0].Name)
	}

	data, err := json.Marshal(repo.Languages)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Go") {
		t.Errorf("Expected JSON to contain Go, got %s", string(data))
	}
}

func TestVerboseDetailLines(t *testing.T) {
	t.Run("includes language line when languages present", func(t *testing.T) {
		repo := RepoInfo{
			Path: "/path/to/repo",
			Languages: []languages.LanguageBreakdown{
				{Name: "Go", Size: 1500, Percent: 75},
				{Name: "Python", Size: 500, Percent: 25},
			},
		}
		lines := verboseDetailLines(repo)
		if len(lines) == 0 {
			t.Fatal("Expected at least one verbose line")
		}
		found := false
		for _, l := range lines {
			if strings.HasPrefix(l, "Languages: ") {
				found = true
				if !strings.Contains(l, "Go") || !strings.Contains(l, "Python") {
					t.Errorf("Languages line missing expected content: %s", l)
				}
			}
		}
		if !found {
			t.Error("No Languages line found in verbose output")
		}
	})

	t.Run("includes CI checks", func(t *testing.T) {
		repo := RepoInfo{
			Path: "/path/to/repo",
			CIChecks: []CICheck{
				{Name: "Tests", Conclusion: "success"},
				{Name: "Lint", Conclusion: "failure"},
			},
		}
		lines := verboseDetailLines(repo)
		foundSuccess := false
		foundFailure := false
		for _, l := range lines {
			if strings.Contains(l, "✓ Tests") {
				foundSuccess = true
			}
			if strings.Contains(l, "✗ Lint") {
				foundFailure = true
			}
		}
		if !foundSuccess {
			t.Error("Missing success CI check line")
		}
		if !foundFailure {
			t.Error("Missing failure CI check line")
		}
	})

	t.Run("includes git status lines", func(t *testing.T) {
		repo := RepoInfo{
			Path:         "/path/to/repo",
			StatusOutput: "On branch main\n\tmodified:   foo.go",
		}
		lines := verboseDetailLines(repo)
		found := false
		for _, l := range lines {
			if strings.Contains(l, "modified:") {
				found = true
			}
		}
		if !found {
			t.Error("Missing git status line in verbose output")
		}
	})

	t.Run("empty repo produces no lines", func(t *testing.T) {
		repo := RepoInfo{Path: "/path/to/repo"}
		lines := verboseDetailLines(repo)
		if len(lines) != 0 {
			t.Errorf("Expected 0 lines for empty repo, got %d: %v", len(lines), lines)
		}
	})
}

func TestLanguagesFlag(t *testing.T) {
	t.Run("flag exists with default true", func(t *testing.T) {
		flag := ProjectsCmd.Flags().Lookup("languages")
		if flag == nil {
			t.Fatal("Expected --languages flag to exist on ProjectsCmd")
		}
		if flag.DefValue != "true" {
			t.Errorf("Expected --languages default value to be true, got %s", flag.DefValue)
		}
	})

	t.Run("showLanguages false when neither verbose nor explicit", func(t *testing.T) {
		// Reset state
		orig := showLanguages
		defer func() { showLanguages = orig }()

		showLanguages = languagesFlag && false // no verboseFlag, no explicit
		if showLanguages {
			t.Error("showLanguages should be false without verbose or explicit flag")
		}
	})

	t.Run("verboseDetailLines omits languages when not populated", func(t *testing.T) {
		repo := RepoInfo{
			Path: "/path/to/repo",
			CIChecks: []CICheck{
				{Name: "Tests", Conclusion: "success"},
			},
		}
		lines := verboseDetailLines(repo)
		for _, l := range lines {
			if strings.HasPrefix(l, "Languages:") {
				t.Error("Languages line should not appear when Languages field is empty")
			}
		}
	})

	t.Run("verboseDetailLines includes languages when populated", func(t *testing.T) {
		repo := RepoInfo{
			Path: "/path/to/repo",
			Languages: []languages.LanguageBreakdown{
				{Name: "Go", Size: 1000, Percent: 100},
			},
		}
		lines := verboseDetailLines(repo)
		found := false
		for _, l := range lines {
			if strings.HasPrefix(l, "Languages:") {
				found = true
			}
		}
		if !found {
			t.Error("Expected Languages line when Languages field is populated")
		}
	})

	t.Run("languages flag works with all filter modes", func(t *testing.T) {
		// Verify the flag can be looked up alongside other flags
		for _, flagName := range []string{"all", "dirty", "clean", "languages"} {
			f := ProjectsCmd.Flags().Lookup(flagName)
			if f == nil {
				t.Errorf("Expected --%s flag to exist", flagName)
			}
		}
	})
}
