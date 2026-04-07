package dotfiles

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// DefaultModularDotfiles lists the dotfiles that, if present in the dotfiles
// repo, are expected to be symlinked into $HOME by the install script.
// Entries not present in the repo are silently skipped — this list is
// forward-looking and will grow as the dotfiles repo adds modules.
var DefaultModularDotfiles = []string{
	".zshrc",
	".profile",
	".aliases",
	".gitconfig",
	".tmux.conf",
	".vimrc",
	".ideavimrc",
	".bashrc",
	".git-shellrc",
	".node-shellrc",
	".go-shellrc",
}

// fetchTimeout is the hard cap for fetching origin during a status check.
// Kept short so `bootstrap status` stays responsive on flaky networks.
const fetchTimeout = 5 * time.Second

// appendGitWarnings appends warning lines to out for git-state problems:
// uncommitted changes, unpushed local commits, unpulled upstream commits.
// It never returns an error and never modifies the caller's err semantics —
// it only writes advisory lines when something is wrong.
func appendGitWarnings(d DotfilesSetup, out *bytes.Buffer) {
	repo, err := git.PlainOpen(d.LocalPath)
	if err != nil {
		writeWarning(out, fmt.Sprintf("git open failed: %v", err), "")
		return
	}

	// Uncommitted changes (dirty worktree).
	if wt, wtErr := repo.Worktree(); wtErr == nil {
		if st, stErr := wt.Status(); stErr == nil && !st.IsClean() {
			writeWarning(out,
				fmt.Sprintf("%d uncommitted change(s) in dotfiles", len(st)),
				fmt.Sprintf("cd %s && git status", d.LocalPath))
		}
	}

	// Fetch origin with a hard timeout so a stuck transport can't freeze us.
	if fetchErr := fetchWithHardTimeout(repo, fetchTimeout); fetchErr != nil {
		writeWarning(out,
			fmt.Sprintf("could not reach origin: %v", fetchErr),
			"")
		// Fall through to check cached refs anyway.
	}

	head, err := repo.Head()
	if err != nil {
		writeWarning(out, fmt.Sprintf("could not read local HEAD: %v", err), "")
		return
	}

	branch := head.Name().Short()
	upstreamRef, err := repo.Reference(
		plumbing.NewRemoteReferenceName("origin", branch), true)
	if err != nil {
		writeWarning(out,
			fmt.Sprintf("no cached upstream ref for origin/%s", branch),
			fmt.Sprintf("cd %s && git fetch", d.LocalPath))
		return
	}

	ahead, behind, err := commitDistance(repo, head.Hash(), upstreamRef.Hash())
	if err != nil {
		writeWarning(out, fmt.Sprintf("could not compute ahead/behind: %v", err), "")
		return
	}
	if ahead > 0 {
		writeWarning(out,
			fmt.Sprintf("%d local commit(s) not pushed to origin", ahead),
			fmt.Sprintf("cd %s && git push", d.LocalPath))
	}
	if behind > 0 {
		writeWarning(out,
			fmt.Sprintf("%d upstream commit(s) not pulled into local", behind),
			fmt.Sprintf("cd %s && git pull", d.LocalPath))
	}
}

// appendModularDotfileWarnings checks every expected modular file and
// appends a warning when a file exists in the dotfiles repo but is not
// correctly symlinked into $HOME.
func appendModularDotfileWarnings(d DotfilesSetup, out *bytes.Buffer) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}
	repoAbs, err := filepath.Abs(d.LocalPath)
	if err != nil {
		repoAbs = d.LocalPath
	}

	for _, name := range d.ModularFiles {
		repoFile := filepath.Join(repoAbs, name)
		if _, statErr := os.Stat(repoFile); statErr != nil {
			// Silent: the user hasn't added this modular file yet.
			continue
		}

		homeFile := filepath.Join(home, name)
		info, lstatErr := os.Lstat(homeFile)
		if lstatErr != nil {
			writeWarning(out,
				fmt.Sprintf("%s not symlinked into $HOME", name),
				fmt.Sprintf("cd %s && ./fresh.sh", d.LocalPath))
			continue
		}
		if info.Mode()&os.ModeSymlink == 0 {
			writeWarning(out,
				fmt.Sprintf("%s in $HOME is a regular file (shadows dotfiles)", name),
				fmt.Sprintf("rm %s && cd %s && ./fresh.sh", homeFile, d.LocalPath))
			continue
		}
		target, readErr := os.Readlink(homeFile)
		if readErr != nil {
			writeWarning(out,
				fmt.Sprintf("%s symlink unreadable: %v", name, readErr),
				"")
			continue
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(homeFile), target)
		}
		targetAbs, absErr := filepath.Abs(target)
		if absErr != nil {
			targetAbs = target
		}
		if targetAbs != repoFile {
			writeWarning(out,
				fmt.Sprintf("%s symlink does not point into dotfiles repo", name),
				fmt.Sprintf("cd %s && ./fresh.sh", d.LocalPath))
		}
	}
}

// writeWarning appends a single indented warning line to out, with an
// optional remediation hint shown after the message.
func writeWarning(out *bytes.Buffer, message, hint string) {
	out.WriteString("\n    ")
	_, _ = color.New(color.FgYellow).Fprintf(out, "⚠ %s", message)
	if hint != "" {
		_, _ = color.New(color.Faint).Fprintf(out, "  (run: %s)", hint)
	}
}

// commitDistance returns (ahead, behind) between local and upstream by
// walking both commit graphs and computing the exclusive counts.
func commitDistance(repo *git.Repository, local, upstream plumbing.Hash) (ahead, behind int, err error) {
	if local == upstream {
		return 0, 0, nil
	}
	localSet, err := ancestorSet(repo, local)
	if err != nil {
		return 0, 0, err
	}
	upstreamSet, err := ancestorSet(repo, upstream)
	if err != nil {
		return 0, 0, err
	}
	for h := range localSet {
		if _, ok := upstreamSet[h]; !ok {
			ahead++
		}
	}
	for h := range upstreamSet {
		if _, ok := localSet[h]; !ok {
			behind++
		}
	}
	return ahead, behind, nil
}

func ancestorSet(repo *git.Repository, from plumbing.Hash) (map[plumbing.Hash]struct{}, error) {
	set := make(map[plumbing.Hash]struct{})
	iter, err := repo.Log(&git.LogOptions{From: from})
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	err = iter.ForEach(func(c *object.Commit) error {
		set[c.Hash] = struct{}{}
		return nil
	})
	return set, err
}

// fetchWithHardTimeout wraps FetchContext in a goroutine with a belt-and-braces
// time.After fallback — go-git's SSH transports don't always honor context
// cancellation cleanly. Returns nil on success or NoErrAlreadyUpToDate.
func fetchWithHardTimeout(repo *git.Repository, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- repo.FetchContext(ctx, &git.FetchOptions{RemoteName: "origin"})
	}()

	select {
	case err := <-done:
		if err == nil || err == git.NoErrAlreadyUpToDate {
			return nil
		}
		return err
	case <-time.After(timeout + time.Second):
		return fmt.Errorf("fetch timed out after %s", timeout)
	}
}
