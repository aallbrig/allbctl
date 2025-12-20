package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	return map[string]RuntimeCheck{
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

	for _, rt := range runtimes {
		if rt.Category == "language" {
			languages = append(languages, rt)
		} else if rt.Category == "version-manager" {
			versionManagers = append(versionManagers, rt)
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
	}

	return output.String()
}

func detectRuntimesInline() string {
	runtimes := detectRuntimes()
	if len(runtimes) == 0 {
		return ""
	}

	// Just return comma-separated list of language names
	var names []string
	for _, rt := range runtimes {
		if rt.Category == "language" {
			names = append(names, rt.Name)
		}
	}

	if len(names) == 0 {
		return ""
	}

	return strings.Join(names, ", ")
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
