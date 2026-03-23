package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	allFlag       bool
	dirtyFlag     bool
	cleanFlag     bool
	limitFlag     int
	showFilesFlag bool
)

// DirtyReason is a bitmask describing why a repo is considered dirty
type DirtyReason uint

const (
	DirtyUncommittedChanges DirtyReason = 1 << iota // has staged, unstaged, or untracked changes
	DirtyUnpushedCommits                            // has commits not yet pushed to upstream
	DirtyNoUpstream                                 // current branch has no remote tracking branch
)

func (r DirtyReason) Labels() []string {
	var labels []string
	if r&DirtyUncommittedChanges != 0 {
		labels = append(labels, "uncommitted changes")
	}
	if r&DirtyUnpushedCommits != 0 {
		labels = append(labels, "unpushed commits")
	}
	if r&DirtyNoUpstream != 0 {
		labels = append(labels, "no upstream")
	}
	return labels
}

func (r DirtyReason) String() string {
	labels := r.Labels()
	if len(labels) == 0 {
		return ""
	}
	return "[" + strings.Join(labels, ", ") + "]"
}

// ProjectsCmd represents the projects command
var ProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Display git repositories in ~/src",
	Long: `Display a summary of git repositories found in ~/src directory.

By default, shows the same summary as the 'Projects:' section in 'allbctl status'.
Dirty repos are marked with an asterisk (*).

Examples:
  allbctl status projects                        # Show summary (default, same as status)
  allbctl status projects --all                  # Show all repos
  allbctl status projects --dirty                # Show only dirty repos
  allbctl status projects --clean                # Show only clean repos
  allbctl status projects --dirty --show-files   # Show dirty repos with their changed files`,
	Run: func(cmd *cobra.Command, args []string) {
		if allFlag || dirtyFlag || cleanFlag || showFilesFlag {
			printProjectsSummary()
		} else {
			// Default: show all projects (no limit), unless --limit is specified
			printProjectsInline(limitFlag)
		}
	},
}

func init() {
	ProjectsCmd.Flags().BoolVar(&allFlag, "all", false, "Explicitly show all detected git repos (same as default, useful for clarity in scripts)")
	ProjectsCmd.Flags().BoolVar(&dirtyFlag, "dirty", false, "Show only dirty repos")
	ProjectsCmd.Flags().BoolVar(&cleanFlag, "clean", false, "Show only clean repos")
	ProjectsCmd.Flags().IntVar(&limitFlag, "limit", 0, "Limit the number of projects shown (0 = no limit, show all)")
	ProjectsCmd.Flags().BoolVar(&showFilesFlag, "show-files", false, "Show changed files (tracked and untracked) under each dirty repo")
}

// RepoInfo contains information about a git repository
type RepoInfo struct {
	Path             string
	ModTime          time.Time
	Dirty            bool
	DirtyReasons     DirtyReason
	RemoteRepo       string // e.g., "aallbrig/allbctl" or "godotengine/godot"
	StatusOutput     string // populated when --show-files is set; full `git status --untracked-files=all` output
	UncommittedFiles int    // staged + unstaged file count (excludes untracked)
	UntrackedFiles   int    // untracked file count
	UnpushedCommits  int    // number of commits ahead of upstream
}

// printProjectsSummary prints a summary of git repositories
func printProjectsSummary() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	srcDir := filepath.Join(home, "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		fmt.Printf("~/src directory does not exist\n")
		return
	}

	repos := findGitRepos(srcDir)
	if len(repos) == 0 {
		fmt.Printf("No git repositories found in ~/src\n")
		return
	}

	// Get repos with their info
	repoInfos := getReposByModTime(repos)

	// Filter based on flags
	var displayMode string
	if dirtyFlag {
		displayMode = "dirty"
	} else if cleanFlag {
		displayMode = "clean"
	} else if allFlag {
		displayMode = "all"
	} else {
		displayMode = "summary"
	}

	filtered := filterRepos(repoInfos, displayMode)

	// Display based on mode
	if displayMode == "summary" {
		// Count dirty repos
		dirtyCount := 0
		for _, repo := range repoInfos {
			if repo.Dirty {
				dirtyCount++
			}
		}

		// Format: "Total repos: 4 (2 dirty)" or "Total repos: 4 (2 dirty, 7 files)"
		if showFilesFlag && dirtyCount > 0 {
			totalFiles := countTotalStatusFiles(repoInfos)
			fmt.Printf("Total repos: %d (%d dirty, %d files)\n", len(repos), dirtyCount, totalFiles)
		} else if dirtyCount > 0 {
			fmt.Printf("Total repos: %d (%d dirty)\n", len(repos), dirtyCount)
		} else {
			fmt.Printf("Total repos: %d\n", len(repos))
		}

		count := 5
		if len(filtered) < count {
			count = len(filtered)
		}
		fmt.Printf("\nLast %d recently touched:\n", count)
		printRepoTable(filtered[:count], "  ", showFilesFlag, true)
	} else {
		fmt.Println(buildSummaryLine(filtered))
		fmt.Println()
		printRepoTable(filtered, "  ", showFilesFlag, dirtyFlag || allFlag)
	}
}

func buildSummaryLine(repos []RepoInfo) string {
	unpushedCommits := 0
	uncommittedFiles := 0
	untrackedFiles := 0
	for _, repo := range repos {
		unpushedCommits += repo.UnpushedCommits
		uncommittedFiles += repo.UncommittedFiles
		untrackedFiles += repo.UntrackedFiles
	}
	parts := []string{fmt.Sprintf("Total projects: %d", len(repos))}
	if unpushedCommits > 0 {
		parts = append(parts, fmt.Sprintf("Total unpushed commits: %d", unpushedCommits))
	}
	if uncommittedFiles > 0 {
		parts = append(parts, fmt.Sprintf("Total modified files: %d", uncommittedFiles))
	}
	if untrackedFiles > 0 {
		parts = append(parts, fmt.Sprintf("Total untracked files: %d", untrackedFiles))
	}
	return strings.Join(parts, "  ")
}

// printProjectsInline prints a summary for the status command.
// limit controls how many recently-touched projects to show; 0 means no limit (show all).
func printProjectsInline(limit int) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	srcDir := filepath.Join(home, "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return
	}

	repos := findGitRepos(srcDir)
	if len(repos) == 0 {
		return
	}

	repoInfos := getReposByModTime(repos)
	dirtyCount := 0
	for _, repo := range repoInfos {
		if repo.Dirty {
			dirtyCount++
		}
	}

	// Format: "Projects: 4 total (2 dirty)"
	if dirtyCount > 0 {
		fmt.Printf("Projects: %d total (%d dirty)\n", len(repos), dirtyCount)
	} else {
		fmt.Printf("Projects: %d total\n", len(repos))
	}

	// Show recently touched projects; limit=0 means show all
	count := len(repoInfos)
	if limit > 0 && limit < count {
		count = limit
	}
	if limit > 0 {
		fmt.Printf("  Last %d recently touched:\n", count)
	} else {
		fmt.Printf("  Recently touched (%d):\n", count)
	}
	printRepoTable(repoInfos[:count], "    ", false, true)
}

// findGitRepos recursively finds all git repositories in the given directory
func findGitRepos(rootDir string) []string {
	var repos []string

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Check if this directory is a git repo
		if d.IsDir() && d.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return []string{}
	}

	return repos
}

// getDirtyReasons returns a bitmask describing why a repo is dirty
func getDirtyReasons(repoPath string) DirtyReason {
	var reasons DirtyReason

	if output, err := exec.Command("git", "-C", repoPath, "status", "--porcelain").Output(); err == nil {
		if len(strings.TrimSpace(string(output))) > 0 {
			reasons |= DirtyUncommittedChanges
		}
	}

	// If repo has no commits yet, upstream checks don't apply
	if exec.Command("git", "-C", repoPath, "rev-parse", "HEAD").Run() != nil {
		return reasons
	}

	// Check whether the current branch has an upstream tracking branch
	if exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "@{u}").Run() != nil {
		reasons |= DirtyNoUpstream
		return reasons
	}

	// Has upstream — check for unpushed commits
	if output, err := exec.Command("git", "-C", repoPath, "log", "@{u}..HEAD", "--oneline").Output(); err == nil {
		if len(strings.TrimSpace(string(output))) > 0 {
			reasons |= DirtyUnpushedCommits
		}
	}

	return reasons
}

// countUnpushedCommits returns the number of commits ahead of upstream.
func countUnpushedCommits(repoPath string) int {
	output, err := exec.Command("git", "-C", repoPath, "log", "@{u}..HEAD", "--oneline").Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			count++
		}
	}
	return count
}

// isGitRepoDirty returns true if the repo has any dirty reasons
func isGitRepoDirty(repoPath string) bool {
	return getDirtyReasons(repoPath) != 0
}

// filterStatusLines removes noise from git status output for display:
// strips hint/noise lines and any trailing blank lines.
func filterStatusLines(lines []string) []string {
	noisePatterns := []string{
		"nothing to commit",
		`(use "git push"`,
		`(use "git add`,
		`(use "git restore`,
		"no changes added to commit",
	}
	var filtered []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		noise := false
		for _, pattern := range noisePatterns {
			if strings.Contains(line, pattern) {
				noise = true
				break
			}
		}
		if !noise {
			filtered = append(filtered, line)
		}
	}
	// Strip trailing blank lines
	for len(filtered) > 0 && strings.TrimSpace(filtered[len(filtered)-1]) == "" {
		filtered = filtered[:len(filtered)-1]
	}
	return filtered
}

func getGitStatusOutput(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "status", "--untracked-files=all")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(output), "\n")
}

// countStatusFiles counts file entries in human-readable git status output.
// File entries are tab-indented lines (e.g., "\tmodified:   foo.go" or "\tfoo.go").
func countStatusFiles(statusOutput string) int {
	count := 0
	for _, line := range strings.Split(statusOutput, "\n") {
		if strings.HasPrefix(line, "\t") {
			count++
		}
	}
	return count
}

// countTotalStatusFiles returns the total file count across all repos
func countTotalStatusFiles(repos []RepoInfo) int {
	total := 0
	for _, repo := range repos {
		total += countStatusFiles(repo.StatusOutput)
	}
	return total
}

// countPorcelainFiles parses `git status --porcelain` output and returns
// the count of staged/unstaged files and the count of untracked files.
func countPorcelainFiles(repoPath string) (uncommitted, untracked int) {
	output, err := exec.Command("git", "-C", repoPath, "status", "--porcelain", "--untracked-files=all").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
		if len(line) < 2 {
			continue
		}
		if strings.HasPrefix(line, "??") {
			untracked++
		} else {
			uncommitted++
		}
	}
	return
}

// getReposByModTime gets repository info sorted by modification time (most recent first)
func getReposByModTime(repos []string) []RepoInfo {
	var repoInfos []RepoInfo

	for _, repo := range repos {
		info, err := os.Stat(repo)
		if err != nil {
			continue
		}

		reasons := getDirtyReasons(repo)
		repoInfo := RepoInfo{
			Path:         repo,
			ModTime:      info.ModTime(),
			Dirty:        reasons != 0,
			DirtyReasons: reasons,
			RemoteRepo:   getRemoteRepo(repo),
		}
		if reasons != 0 {
			repoInfo.UncommittedFiles, repoInfo.UntrackedFiles = countPorcelainFiles(repo)
			if reasons&DirtyUnpushedCommits != 0 {
				repoInfo.UnpushedCommits = countUnpushedCommits(repo)
			}
			if showFilesFlag {
				repoInfo.StatusOutput = getGitStatusOutput(repo)
			}
		}
		repoInfos = append(repoInfos, repoInfo)
	}

	// Sort by modification time (most recent first)
	sort.Slice(repoInfos, func(i, j int) bool {
		return repoInfos[i].ModTime.After(repoInfos[j].ModTime)
	})

	return repoInfos
}

// formatRepoPath formats a repository path for display
func formatRepoPath(path string, dirty bool) string {
	home, err := os.UserHomeDir()
	if err == nil {
		path = strings.Replace(path, home, "~", 1)
	}

	// Normalize to forward slashes for consistent cross-platform display
	path = filepath.ToSlash(path)

	if dirty {
		return path + "*"
	}
	return path
}

// filterRepos filters repositories based on the display mode
func filterRepos(repos []RepoInfo, mode string) []RepoInfo {
	switch mode {
	case "dirty":
		var filtered []RepoInfo
		for _, repo := range repos {
			if repo.Dirty {
				filtered = append(filtered, repo)
			}
		}
		return filtered
	case "clean":
		var filtered []RepoInfo
		for _, repo := range repos {
			if !repo.Dirty {
				filtered = append(filtered, repo)
			}
		}
		return filtered
	default:
		return repos
	}
}

// getRemoteRepo gets the remote repository (user/repo) from git remote origin
func getRemoteRepo(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(output))
	return parseRemoteRepo(url)
}

// parseRemoteRepo parses a git remote URL to extract user/repo
func parseRemoteRepo(url string) string {
	if url == "" {
		return ""
	}

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Handle different URL formats:
	// - https://github.com/user/repo
	// - git@github.com:user/repo
	// - https://gitlab.com/user/repo
	// - git@gitlab.com:user/repo

	// For SSH format (git@...)
	if strings.HasPrefix(url, "git@") {
		// Split by : to get the path part
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			path := parts[1]
			// Extract user/repo from path
			pathParts := strings.Split(path, "/")
			if len(pathParts) >= 2 {
				return strings.Join(pathParts[len(pathParts)-2:], "/")
			}
		}
	}

	// For HTTPS format
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// Remove protocol
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimPrefix(url, "http://")
		// Split by / and get last two parts
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			return strings.Join(parts[len(parts)-2:], "/")
		}
	}

	return ""
}

// formatRepoLine formats a repository line with path, dirty marker, remote, and date
func formatRepoLine(repo RepoInfo) string {
	// Format the path with dirty marker
	path := formatRepoPath(repo.Path, repo.Dirty)

	// Format the date/time (e.g., "2024-12-23 15:30 MST -0700")
	dateTime := repo.ModTime.Format("2006-01-02 15:04 MST -0700")

	// Build the line with proper spacing
	// Format: "  ~/src/project*  user/repo  2024-12-23 15:30"
	var parts []string
	parts = append(parts, path)

	if repo.RemoteRepo != "" {
		parts = append(parts, repo.RemoteRepo)
	}

	parts = append(parts, dateTime)

	// Join with multiple spaces for visual separation
	return strings.Join(parts, "  ")
}

// printRepoTable prints repositories in a table format with aligned columns
func printRepoTable(repos []RepoInfo, indent string, showFiles bool, showReasons bool) {
	if len(repos) == 0 {
		return
	}

	// Calculate column widths
	maxPathLen := 0
	maxRemoteLen := 0

	for _, repo := range repos {
		path := formatRepoPath(repo.Path, repo.Dirty)
		if len(path) > maxPathLen {
			maxPathLen = len(path)
		}
		if len(repo.RemoteRepo) > maxRemoteLen {
			maxRemoteLen = len(repo.RemoteRepo)
		}
	}

	// Print each row with aligned columns
	for _, repo := range repos {
		path := formatRepoPath(repo.Path, repo.Dirty)
		remote := repo.RemoteRepo
		if remote == "" {
			remote = "-"
		}
		dateTime := repo.ModTime.Format("2006-01-02 15:04 MST -0700")

		line := fmt.Sprintf("%s%-*s  %-*s  %s", indent, maxPathLen, path, maxRemoteLen, remote, dateTime)
		if showReasons && repo.Dirty {
			line += "  " + repo.DirtyReasons.String()
		}
		fmt.Println(line)

		if showFiles && repo.StatusOutput != "" {
			lines := filterStatusLines(strings.Split(repo.StatusOutput, "\n"))
			for _, l := range lines {
				fmt.Printf("%s  %s\n", indent, l)
			}
			fmt.Println()
		}
	}
}
