package dotfiles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ---------------------------------------------------------------------------
// Fixture helpers
// ---------------------------------------------------------------------------

// initBareRepo creates a bare git repo at a fresh temp dir and returns its path.
// Used as the fake "remote" so tests stay offline.
func initBareRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if _, err := git.PlainInit(dir, true); err != nil {
		t.Fatalf("PlainInit(bare) failed: %v", err)
	}
	return dir
}

// initWorktreeWithRemote creates a non-bare repo with one initial commit and
// "origin" pointing at the given bare remote. Returns the worktree path.
func initWorktreeWithRemote(t *testing.T, remoteURL string) string {
	t.Helper()
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("PlainInit(worktree) failed: %v", err)
	}
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	}); err != nil {
		t.Fatalf("CreateRemote failed: %v", err)
	}
	commitFile(t, repo, dir, "README.md", "initial")
	return dir
}

func commitFile(t *testing.T, repo *git.Repository, worktreePath, name, content string) plumbing.Hash {
	t.Helper()
	if err := os.WriteFile(filepath.Join(worktreePath, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Worktree(): %v", err)
	}
	if _, err := wt.Add(name); err != nil {
		t.Fatalf("Add(%s): %v", name, err)
	}
	hash, err := wt.Commit("commit "+name, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Commit(%s): %v", name, err)
	}
	return hash
}

// pushCurrentBranch pushes refs/heads/<branch> to the configured origin.
// Used by tests that need a populated origin remote.
func pushCurrentBranch(t *testing.T, repoPath, branch string) {
	t.Helper()
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("PlainOpen(%s): %v", repoPath, err)
	}
	refSpec := config.RefSpec("refs/heads/" + branch + ":refs/heads/" + branch)
	if err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{refSpec},
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		t.Fatalf("Push: %v", err)
	}
}

// fetchOrigin updates the local cached refs/remotes/origin/* refs.
func fetchOrigin(t *testing.T, repoPath string) {
	t.Helper()
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("PlainOpen: %v", err)
	}
	if err := repo.Fetch(&git.FetchOptions{RemoteName: "origin"}); err != nil && err != git.NoErrAlreadyUpToDate {
		t.Fatalf("Fetch: %v", err)
	}
}

// currentBranch returns the short name of HEAD's branch in the given repo.
func currentBranch(t *testing.T, repoPath string) string {
	t.Helper()
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("PlainOpen: %v", err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatalf("Head: %v", err)
	}
	return head.Name().Short()
}

// newDotfilesUnderHome wires a fresh fake $HOME and returns a DotfilesSetup
// pointing at a worktree synced with a bare remote. The worktree is empty
// of modular dotfiles by default — tests add the ones they need.
func newDotfilesUnderHome(t *testing.T) *DotfilesSetup {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	remote := initBareRepo(t)
	worktree := initWorktreeWithRemote(t, remote)
	branch := currentBranch(t, worktree)
	pushCurrentBranch(t, worktree, branch)
	fetchOrigin(t, worktree)

	d := NewDotfilesSetup("https://example.invalid/dotfiles.git", worktree, "./fresh.sh")
	// Disable modular checks by default; per-test overrides set them.
	d.ModularFiles = nil
	return d
}

// ---------------------------------------------------------------------------
// Validate() — existing branches still pass
// ---------------------------------------------------------------------------

func TestValidate_MissingDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	d := NewDotfilesSetup("https://example.invalid/dotfiles.git",
		filepath.Join(t.TempDir(), "missing"), "./fresh.sh")
	out, err := d.Validate()
	if err == nil {
		t.Fatal("expected err for missing directory")
	}
	if !strings.Contains(out.String(), "NOT CLONED") {
		t.Errorf("expected NOT CLONED in output, got: %q", out.String())
	}
	if strings.Contains(out.String(), "⚠") {
		t.Errorf("missing dir should not produce health warnings, got: %q", out.String())
	}
}

func TestValidate_NotAGitRepo(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	d := NewDotfilesSetup("https://example.invalid/dotfiles.git", dir, "./fresh.sh")
	out, err := d.Validate()
	if err == nil {
		t.Fatal("expected err for non-git directory")
	}
	if !strings.Contains(out.String(), "EXISTS BUT NOT A GIT REPO") {
		t.Errorf("expected 'EXISTS BUT NOT A GIT REPO', got: %q", out.String())
	}
	if strings.Contains(out.String(), "⚠") {
		t.Errorf("non-git dir should not produce health warnings, got: %q", out.String())
	}
}

// ---------------------------------------------------------------------------
// Validate() — clean state stays quiet
// ---------------------------------------------------------------------------

func TestValidate_CleanInSync_NoWarnings(t *testing.T) {
	d := newDotfilesUnderHome(t)
	out, err := d.Validate()
	if err != nil {
		t.Fatalf("expected nil err, got: %v", err)
	}
	if !strings.Contains(out.String(), "CLONED") {
		t.Errorf("expected CLONED line, got: %q", out.String())
	}
	if strings.Contains(out.String(), "⚠") {
		t.Errorf("clean repo should produce no warnings, got: %q", out.String())
	}
}

// ---------------------------------------------------------------------------
// Validate() — git health checks
// ---------------------------------------------------------------------------

func TestValidate_DirtyWorktree(t *testing.T) {
	d := newDotfilesUnderHome(t)
	if err := os.WriteFile(filepath.Join(d.LocalPath, "scratch"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "uncommitted change") {
		t.Errorf("expected uncommitted change warning, got: %q", out.String())
	}
}

func TestValidate_AheadOfOrigin(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	commitFile(t, repo, d.LocalPath, "extra.txt", "extra")
	// Note: do NOT push, do NOT re-fetch — local is now ahead by 1.

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "1 local commit(s) not pushed") {
		t.Errorf("expected ahead warning, got: %q", out.String())
	}
}

func TestValidate_BehindOrigin(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	// Make a new commit, push it, but then reset the local HEAD so the
	// cached refs/remotes/origin/<branch> is ahead of HEAD.
	branch := currentBranch(t, d.LocalPath)
	originalHead, _ := repo.Head()
	commitFile(t, repo, d.LocalPath, "remote-only.txt", "remote-only")
	pushCurrentBranch(t, d.LocalPath, branch)
	fetchOrigin(t, d.LocalPath)

	// Move local HEAD back one commit so cached origin/<branch> is ahead.
	wt, _ := repo.Worktree()
	if err := wt.Reset(&git.ResetOptions{
		Commit: originalHead.Hash(),
		Mode:   git.HardReset,
	}); err != nil {
		t.Fatalf("reset: %v", err)
	}

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "1 upstream commit(s) not pulled") {
		t.Errorf("expected behind warning, got: %q", out.String())
	}
}

func TestValidate_NoCachedUpstream(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	remote := initBareRepo(t)
	worktree := initWorktreeWithRemote(t, remote)
	// Intentionally do NOT push or fetch — there's no cached origin/<branch>.
	d := NewDotfilesSetup("https://example.invalid/dotfiles.git", worktree, "./fresh.sh")
	d.ModularFiles = nil

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "no cached upstream ref") {
		t.Errorf("expected 'no cached upstream ref' warning, got: %q", out.String())
	}
}

// ---------------------------------------------------------------------------
// Validate() — modular dotfile checks
// ---------------------------------------------------------------------------

func TestValidate_Modular_FileNotInRepo_Silent(t *testing.T) {
	d := newDotfilesUnderHome(t)
	d.ModularFiles = []string{".not-in-repo"}
	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if strings.Contains(out.String(), ".not-in-repo") {
		t.Errorf("file missing from repo should be silent, got: %q", out.String())
	}
}

func TestValidate_Modular_NotInHome(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	commitFile(t, repo, d.LocalPath, ".zshrc", "# zsh")
	pushCurrentBranch(t, d.LocalPath, currentBranch(t, d.LocalPath))
	fetchOrigin(t, d.LocalPath)
	d.ModularFiles = []string{".zshrc"}

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), ".zshrc not symlinked into $HOME") {
		t.Errorf("expected 'not symlinked' warning, got: %q", out.String())
	}
}

func TestValidate_Modular_RegularFileInHome(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	commitFile(t, repo, d.LocalPath, ".zshrc", "# zsh")
	pushCurrentBranch(t, d.LocalPath, currentBranch(t, d.LocalPath))
	fetchOrigin(t, d.LocalPath)
	// Place a regular file (not a symlink) at $HOME/.zshrc.
	homeFile := filepath.Join(os.Getenv("HOME"), ".zshrc")
	if err := os.WriteFile(homeFile, []byte("manual"), 0644); err != nil {
		t.Fatal(err)
	}
	d.ModularFiles = []string{".zshrc"}

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "regular file (shadows dotfiles)") {
		t.Errorf("expected 'shadows dotfiles' warning, got: %q", out.String())
	}
}

func TestValidate_Modular_SymlinkPointsElsewhere(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	commitFile(t, repo, d.LocalPath, ".zshrc", "# zsh")
	pushCurrentBranch(t, d.LocalPath, currentBranch(t, d.LocalPath))
	fetchOrigin(t, d.LocalPath)

	// Symlink $HOME/.zshrc to a different path entirely.
	other := filepath.Join(t.TempDir(), "other-zshrc")
	if err := os.WriteFile(other, []byte("other"), 0644); err != nil {
		t.Fatal(err)
	}
	homeFile := filepath.Join(os.Getenv("HOME"), ".zshrc")
	if err := os.Symlink(other, homeFile); err != nil {
		t.Fatal(err)
	}
	d.ModularFiles = []string{".zshrc"}

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if !strings.Contains(out.String(), "symlink does not point into dotfiles repo") {
		t.Errorf("expected wrong-target warning, got: %q", out.String())
	}
}

func TestValidate_Modular_CorrectSymlink(t *testing.T) {
	d := newDotfilesUnderHome(t)
	repo, _ := git.PlainOpen(d.LocalPath)
	commitFile(t, repo, d.LocalPath, ".zshrc", "# zsh")
	pushCurrentBranch(t, d.LocalPath, currentBranch(t, d.LocalPath))
	fetchOrigin(t, d.LocalPath)

	homeFile := filepath.Join(os.Getenv("HOME"), ".zshrc")
	repoFile := filepath.Join(d.LocalPath, ".zshrc")
	if err := os.Symlink(repoFile, homeFile); err != nil {
		t.Fatal(err)
	}
	d.ModularFiles = []string{".zshrc"}

	out, err := d.Validate()
	if err != nil {
		t.Fatalf("nil err expected, got: %v", err)
	}
	if strings.Contains(out.String(), ".zshrc") {
		t.Errorf("correct symlink should be silent, got: %q", out.String())
	}
}

// ---------------------------------------------------------------------------
// commitDistance unit test
// ---------------------------------------------------------------------------

func TestCommitDistance(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	a := commitFile(t, repo, dir, "a", "a")
	b := commitFile(t, repo, dir, "b", "b")
	c := commitFile(t, repo, dir, "c", "c")

	tests := []struct {
		name              string
		local, upstream   plumbing.Hash
		wantAhead, wantBe int
	}{
		{"same", c, c, 0, 0},
		{"ahead 1", c, b, 1, 0},
		{"ahead 2", c, a, 2, 0},
		{"behind 1", b, c, 0, 1},
		{"behind 2", a, c, 0, 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ahead, behind, err := commitDistance(repo, tc.local, tc.upstream)
			if err != nil {
				t.Fatal(err)
			}
			if ahead != tc.wantAhead || behind != tc.wantBe {
				t.Errorf("got (%d,%d), want (%d,%d)", ahead, behind, tc.wantAhead, tc.wantBe)
			}
		})
	}
}
