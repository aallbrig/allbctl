package languages

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// LanguageBreakdown represents one language's share of a repository.
type LanguageBreakdown struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Percent int    `json:"percent"` // floor of (size/total * 100)
}

// vendoredPrefixes lists directory prefixes to exclude, similar to GitHub's linguist.
var vendoredPrefixes = []string{
	"vendor/",
	"node_modules/",
	"third_party/",
	"third-party/",
	"extern/",
	"deps/",
	"Pods/",
	".bundle/",
	"bower_components/",
	"_vendor/",
}

// isVendored returns true if the file path begins with a known vendored prefix.
func isVendored(path string) bool {
	// Normalize path separators (Windows may use backslash in some contexts)
	normalized := filepath.ToSlash(path)
	for _, prefix := range vendoredPrefixes {
		if strings.HasPrefix(normalized, prefix) {
			return true
		}
		// Also check if any parent directory matches
		if strings.Contains(normalized, "/"+prefix) {
			return true
		}
	}
	return false
}

// fileExt returns the file extension including the leading dot.
// Returns empty string for files with no extension.
func fileExt(path string) string {
	return filepath.Ext(filepath.Base(path))
}

// fileBase returns just the filename without any directory path.
func fileBase(path string) string {
	return filepath.Base(path)
}

// DetectLanguages analyzes a git repository at the given path and returns
// a sorted list of language breakdowns (most bytes first).
// It uses `git ls-tree -r -l HEAD` to enumerate tracked files with their sizes.
func DetectLanguages(repoPath string) ([]LanguageBreakdown, error) {
	output, err := exec.Command("git", "-C", repoPath, "ls-tree", "-r", "-l", "HEAD").Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree failed: %w", err)
	}
	return ParseLsTree(string(output))
}

// ParseLsTree parses the output of `git ls-tree -r -l HEAD` and returns
// language breakdowns. Exported for testability.
//
// Each line has the format:
//
//	<mode> <type> <hash> <size>\t<path>
//
// Example:
//
//	100644 blob abc123  1234\tmain.go
func ParseLsTree(output string) ([]LanguageBreakdown, error) {
	langSizes := make(map[string]int64)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split on tab to separate metadata from path
		tabIdx := strings.IndexByte(line, '\t')
		if tabIdx < 0 {
			continue
		}
		meta := line[:tabIdx]
		path := line[tabIdx+1:]

		if isVendored(path) {
			continue
		}

		// Parse metadata: "<mode> <type> <hash> <size>"
		fields := strings.Fields(meta)
		if len(fields) < 4 {
			continue
		}

		// Skip non-blob entries (submodules show as "commit" type with size "-")
		if fields[1] != "blob" {
			continue
		}

		sizeStr := fields[3]
		if sizeStr == "-" {
			continue
		}

		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			continue
		}

		lang := LanguageForFile(path)
		if lang == "" {
			continue
		}

		langSizes[lang] += size
	}

	return buildBreakdown(langSizes), nil
}

// buildBreakdown converts a language→size map into a sorted slice of breakdowns.
func buildBreakdown(langSizes map[string]int64) []LanguageBreakdown {
	if len(langSizes) == 0 {
		return nil
	}

	var total int64
	for _, s := range langSizes {
		total += s
	}

	result := make([]LanguageBreakdown, 0, len(langSizes))
	for name, size := range langSizes {
		pct := 0
		if total > 0 {
			pct = int(float64(size) / float64(total) * 100)
		}
		result = append(result, LanguageBreakdown{
			Name:    name,
			Size:    size,
			Percent: pct,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	return result
}

// FormatBreakdown formats a language breakdown slice into a human-readable string.
// Example output: "Go: 12345 bytes (67%) | Python: 5678 bytes (33%)"
func FormatBreakdown(breakdown []LanguageBreakdown) string {
	if len(breakdown) == 0 {
		return ""
	}

	parts := make([]string, len(breakdown))
	for i, b := range breakdown {
		parts[i] = fmt.Sprintf("%s: %s (%d%%)", b.Name, formatBytes(b.Size), b.Percent)
	}
	return strings.Join(parts, " | ")
}

// formatBytes returns a human-readable byte size string.
func formatBytes(bytes int64) string {
	switch {
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1<<10))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// GetHeadCommit returns the HEAD commit SHA for a git repository.
func GetHeadCommit(repoPath string) (string, error) {
	output, err := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD failed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
