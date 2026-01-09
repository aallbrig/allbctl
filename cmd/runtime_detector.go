package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type RuntimeInfo struct {
	Name     string
	Version  string
	Category string
}

type RuntimeCheck struct {
	Command  []string
	Category string
}

func getAllRuntimeChecks() map[string]RuntimeCheck {
	checks := map[string]RuntimeCheck{
		// Programming Languages
		"Python":  {[]string{"python3", "--version"}, "language"},
		"Node.js": {[]string{"node", "--version"}, "language"},
		"Go":      {[]string{"go", "version"}, "language"},
		"Java":    {[]string{"java", "-version"}, "language"},
		"Ruby":    {[]string{"ruby", "--version"}, "language"},
		"Rust":    {[]string{"rustc", "--version"}, "language"},
		"PHP":     {[]string{"php", "--version"}, "language"},
		"Perl":    {[]string{"perl", "--version"}, "language"},
		"R":       {[]string{"R", "--version"}, "language"},
		"Scala":   {[]string{"scala", "-version"}, "language"},
		"Kotlin":  {[]string{"kotlin", "-version"}, "language"},
		"Swift":   {[]string{"swift", "--version"}, "language"},
		"Elixir":  {[]string{"elixir", "--version"}, "language"},
		"Erlang":  {[]string{"erl", "-version"}, "language"},
		"Haskell": {[]string{"ghc", "--version"}, "language"},
		"Lua":     {[]string{"lua", "-v"}, "language"},
		"Dart":    {[]string{"dart", "--version"}, "language"},
		"Zig":     {[]string{"zig", "version"}, "language"},
		"C#":      {[]string{"dotnet", "--version"}, "language"},

		// Version Managers
		"nvm":    {[]string{"bash", "-c", ". ~/.nvm/nvm.sh 2>/dev/null && nvm --version || echo ''"}, "version-manager"},
		"pyenv":  {[]string{"pyenv", "--version"}, "version-manager"},
		"rbenv":  {[]string{"rbenv", "--version"}, "version-manager"},
		"jenv":   {[]string{"jenv", "--version"}, "version-manager"},
		"rustup": {[]string{"rustup", "--version"}, "version-manager"},
		"sdkman": {[]string{"bash", "-c", "source ~/.sdkman/bin/sdkman-init.sh 2>/dev/null && sdk version || echo ''"}, "version-manager"},
		"asdf":   {[]string{"asdf", "--version"}, "version-manager"},
	}

	// Add gaming platforms
	gamingChecks := detectGamingPlatforms()
	for name, check := range gamingChecks {
		checks[name] = check
	}

	return checks
}

func detectRuntimes() []RuntimeInfo {
	var runtimes []RuntimeInfo
	checks := getAllRuntimeChecks()

	for name, check := range checks {
		version := checkRuntime(check.Command)
		if version != "" {
			runtimes = append(runtimes, RuntimeInfo{
				Name:     name,
				Version:  version,
				Category: check.Category,
			})
		}
	}

	return runtimes
}

func checkRuntime(cmdArgs []string) string {
	if len(cmdArgs) == 0 {
		return ""
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return ""
	}

	return parseVersion(version)
}

func parseVersion(output string) string {
	return firstLine(output)
}

func firstLine(s string) string {
	for i, c := range s {
		if c == '\n' || c == '\r' {
			return s[:i]
		}
	}
	return s
}

func formatRuntimesOutput(runtimes []RuntimeInfo) string {
	if len(runtimes) == 0 {
		return "No runtimes detected."
	}

	// Group by category
	languages := []RuntimeInfo{}
	versionManagers := []RuntimeInfo{}
	gaming := []RuntimeInfo{}

	for _, rt := range runtimes {
		if rt.Category == "language" {
			languages = append(languages, rt)
		} else if rt.Category == "version-manager" {
			versionManagers = append(versionManagers, rt)
		} else if rt.Category == "gaming" {
			gaming = append(gaming, rt)
		}
	}

	var output strings.Builder

	if len(languages) > 0 {
		output.WriteString("Languages:\n")
		for _, rt := range languages {
			output.WriteString(fmt.Sprintf("  %-15s %s\n", rt.Name+":", rt.Version))
		}
		output.WriteString("\n")
	}

	if len(versionManagers) > 0 {
		output.WriteString("Version Managers:\n")
		for _, rt := range versionManagers {
			output.WriteString(fmt.Sprintf("  %-15s %s\n", rt.Name+":", rt.Version))
		}
		output.WriteString("\n")
	}

	if len(gaming) > 0 {
		output.WriteString("Gaming Platforms:\n")
		for _, rt := range gaming {
			output.WriteString(fmt.Sprintf("  %-15s %s\n", rt.Name+":", rt.Version))
		}
	}

	return output.String()
}

func detectRuntimesInline() string {
	runtimes := detectRuntimes()
	if len(runtimes) == 0 {
		return ""
	}

	// Return comma-separated list of language names with versions in parentheses
	var parts []string
	for _, rt := range runtimes {
		if rt.Category == "language" {
			// Extract just the version number for cleaner display
			version := extractVersionNumber(rt.Version)
			if version != "" {
				versionStr := formatVersionWithUpdate(rt.Name, version)
				parts = append(parts, fmt.Sprintf("%s (%s)", rt.Name, versionStr))
			} else {
				parts = append(parts, rt.Name)
			}
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}

// extractVersionNumber extracts just the version number from version output
func extractVersionNumber(versionOutput string) string {
	// Handle common version output patterns
	output := strings.TrimSpace(versionOutput)

	// For "Python 3.9.0" or "python 3.9.0"
	if strings.HasPrefix(strings.ToLower(output), "python ") {
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	// For "go version go1.20.0 linux/amd64"
	if strings.HasPrefix(output, "go version go") {
		parts := strings.Fields(output)
		if len(parts) >= 3 {
			return strings.TrimPrefix(parts[2], "go")
		}
	}

	// For "rustc 1.70.0 (90c541806 2023-05-31)"
	if strings.HasPrefix(output, "rustc ") {
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	// For "node v18.0.0" or "Node.js v18.0.0"
	if strings.Contains(strings.ToLower(output), "node") || strings.HasPrefix(output, "v") {
		parts := strings.Fields(output)
		for _, part := range parts {
			if strings.HasPrefix(part, "v") && len(part) > 1 {
				return strings.TrimPrefix(part, "v")
			}
		}
	}

	// For "ruby 3.0.0p0 (2020-12-25 revision 95aff21468)"
	if strings.HasPrefix(output, "ruby ") {
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			// Remove patch level suffix like "p0"
			version := parts[1]
			if idx := strings.Index(version, "p"); idx > 0 {
				return version[:idx]
			}
			return version
		}
	}

	// For "This is perl 5, version 38, subversion 2 (v5.38.2)"
	if strings.Contains(output, "perl") {
		// Look for version in parentheses like (v5.38.2)
		if idx := strings.Index(output, "(v"); idx >= 0 {
			rest := output[idx+2:]
			if endIdx := strings.Index(rest, ")"); endIdx >= 0 {
				return rest[:endIdx]
			}
		}
		// Try "version X, subversion Y" pattern
		if strings.Contains(output, "version") {
			parts := strings.Split(output, ",")
			var versionParts []string
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "version ") {
					versionParts = append(versionParts, strings.TrimPrefix(part, "version "))
				} else if strings.HasPrefix(part, "subversion ") {
					versionParts = append(versionParts, strings.TrimPrefix(part, "subversion "))
				}
			}
			if len(versionParts) >= 2 {
				return versionParts[0] + "." + versionParts[1]
			}
		}
	}

	// For "PHP 8.1.0 (cli) (built: Nov 23 2021)"
	if strings.HasPrefix(strings.ToUpper(output), "PHP ") {
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	// For "javac 17.0.1" or "java version \"17.0.1\""
	if strings.Contains(strings.ToLower(output), "java") {
		parts := strings.Fields(output)
		for i, part := range parts {
			if (part == "version" || strings.HasPrefix(part, "version")) && i+1 < len(parts) {
				version := parts[i+1]
				// Remove quotes
				version = strings.Trim(version, "\"'")
				return version
			}
		}
		// Try second field
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	// For C# / dotnet: "6.0.100"
	if !strings.Contains(output, " ") {
		// Just a version number
		return output
	}

	// Generic: try to find version-like pattern (e.g., "1.2.3")
	fields := strings.Fields(output)
	for _, field := range fields {
		// Check if field looks like a version (contains digits and dots)
		if strings.Contains(field, ".") {
			// Clean up any surrounding characters
			field = strings.Trim(field, "()[]{}\"',")
			if len(field) > 0 && (field[0] >= '0' && field[0] <= '9') {
				return field
			}
		}
	}

	// If no pattern matches, return the first line as-is
	return output
}

func checkNvmInstalled() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	nvmDir := filepath.Join(home, ".nvm")
	if _, err := os.Stat(nvmDir); err == nil {
		return true
	}
	return false
}

// detectGamingPlatforms returns gaming platform checks based on the OS
func detectGamingPlatforms() map[string]RuntimeCheck {
	checks := make(map[string]RuntimeCheck)

	// Steam detection - cross-platform
	steamCmd := detectSteamCommand()
	if len(steamCmd) > 0 {
		checks["Steam"] = RuntimeCheck{steamCmd, "gaming"}
	}

	return checks
}

// detectSteamCommand returns the appropriate command to check Steam installation
func detectSteamCommand() []string {
	osType := runtime.GOOS
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch osType {
	case "linux":
		// Try command-line first
		if exists("steam") {
			return []string{"bash", "-c", "steam --version 2>/dev/null | head -1 || echo 'Steam (installed)'"}
		}

		// Check common Linux installation paths
		linuxPaths := []string{
			filepath.Join(home, ".steam", "steam.sh"),
			filepath.Join(home, ".local", "share", "Steam", "steam.sh"),
			"/usr/bin/steam",
			"/usr/games/steam",
		}
		for _, path := range linuxPaths {
			if _, err := os.Stat(path); err == nil {
				return []string{"bash", "-c", fmt.Sprintf("%s --version 2>/dev/null | head -1 || echo 'Steam (installed)'", path)}
			}
		}

	case "darwin":
		// Check for Steam.app on macOS
		macPaths := []string{
			"/Applications/Steam.app",
			filepath.Join(home, "Applications", "Steam.app"),
		}
		for _, appPath := range macPaths {
			if _, err := os.Stat(appPath); err == nil {
				plistPath := filepath.Join(appPath, "Contents", "Info.plist")
				return []string{"bash", "-c", fmt.Sprintf("defaults read '%s' CFBundleShortVersionString 2>/dev/null || echo 'Steam (installed)'", plistPath)}
			}
		}

	case "windows":
		// Check Windows registry for Steam installation
		if exists("reg") {
			return []string{"cmd", "/c", "reg query \"HKCU\\Software\\Valve\\Steam\" /v SteamPath >nul 2>&1 && echo Steam (installed) || echo"}
		}

		// Fallback: check Program Files
		windowsPaths := []string{
			"C:\\Program Files (x86)\\Steam\\steam.exe",
			"C:\\Program Files\\Steam\\steam.exe",
		}
		for _, path := range windowsPaths {
			if _, err := os.Stat(path); err == nil {
				return []string{"cmd", "/c", "echo Steam (installed)"}
			}
		}
	}

	return nil
}
