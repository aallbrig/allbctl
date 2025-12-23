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
	allFlag   bool
	dirtyFlag bool
	cleanFlag bool
)

// ProjectsCmd represents the projects command
var ProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Display git repositories in ~/src",
	Long: `Display a summary of git repositories found in ~/src directory.

By default, shows a count and the last 5 recently touched repos.
Dirty repos are marked with an asterisk (*).

Examples:
  allbctl projects           # Show summary (default)
  allbctl projects --all     # Show all repos
  allbctl projects --dirty   # Show only dirty repos
  allbctl projects --clean   # Show only clean repos`,
	Run: func(cmd *cobra.Command, args []string) {
		printProjectsSummary()
	},
}

func init() {
	ProjectsCmd.Flags().BoolVar(&allFlag, "all", false, "Show all detected git repos")
	ProjectsCmd.Flags().BoolVar(&dirtyFlag, "dirty", false, "Show only dirty repos")
	ProjectsCmd.Flags().BoolVar(&cleanFlag, "clean", false, "Show only clean repos")
}

// RepoInfo contains information about a git repository
type RepoInfo struct {
	Path       string
	ModTime    time.Time
	Dirty      bool
	RemoteRepo string // e.g., "aallbrig/allbctl" or "godotengine/godot"
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

		// Format: "Total repos: 4 (2 dirty)"
		if dirtyCount > 0 {
			fmt.Printf("Total repos: %d (%d dirty)\n", len(repos), dirtyCount)
		} else {
			fmt.Printf("Total repos: %d\n", len(repos))
		}

		fmt.Printf("\nLast 5 recently touched:\n")
		count := 5
		if len(filtered) < count {
			count = len(filtered)
		}
		printRepoTable(filtered[:count], "  ")
	} else {
		fmt.Printf("Total repos: %d\n\n", len(filtered))
		printRepoTable(filtered, "  ")
	}
}

// printProjectsInline prints a one-line summary for status command
func printProjectsInline() {
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

	// Show last 5 recently touched
	count := 5
	if len(repoInfos) < count {
		count = len(repoInfos)
	}
	fmt.Printf("  Last 5 recently touched:\n")
	printRepoTable(repoInfos[:count], "    ")
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

// isGitRepoDirty checks if a git repository has uncommitted changes
func isGitRepoDirty(repoPath string) bool {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(output))) > 0
}

// getReposByModTime gets repository info sorted by modification time (most recent first)
func getReposByModTime(repos []string) []RepoInfo {
	var repoInfos []RepoInfo

	for _, repo := range repos {
		info, err := os.Stat(repo)
		if err != nil {
			continue
		}

		repoInfo := RepoInfo{
			Path:       repo,
			ModTime:    info.ModTime(),
			Dirty:      isGitRepoDirty(repo),
			RemoteRepo: getRemoteRepo(repo),
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

	// Format the date/time (e.g., "2024-12-23 15:30")
	dateTime := repo.ModTime.Format("2006-01-02 15:04")

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
func printRepoTable(repos []RepoInfo, indent string) {
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
		dateTime := repo.ModTime.Format("2006-01-02 15:04")

		// Format with right-aligned columns
		fmt.Printf("%s%-*s  %-*s  %s\n", indent, maxPathLen, path, maxRemoteLen, remote, dateTime)
	}
}
