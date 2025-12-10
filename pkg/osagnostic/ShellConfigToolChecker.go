package osagnostic

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ShellConfigToolChecker struct {
	configPaths []string
}

type ToolStatus struct {
	Tool      string
	Available bool
	Source    string // Which config file referenced it
}

func NewShellConfigToolChecker() *ShellConfigToolChecker {
	osInfo := NewOperatingSystem()
	homeDir := osInfo.HomeDirectoryPath

	configFiles := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".profile"),
		filepath.Join(homeDir, ".config", "fish", "config.fish"),
	}

	var existingConfigs []string
	for _, path := range configFiles {
		if _, err := os.Stat(path); err == nil {
			existingConfigs = append(existingConfigs, path)
		}
	}

	return &ShellConfigToolChecker{
		configPaths: existingConfigs,
	}
}

func (c *ShellConfigToolChecker) ExtractTools() []ToolStatus {
	toolMap := make(map[string]string) // tool -> source file (full path)
	
	// Common shell builtins and utilities to exclude
	commonCommands := map[string]bool{
		// Core shell builtins
		"echo": true, "cd": true, "ls": true, "cat": true, "grep": true,
		"sed": true, "awk": true, "find": true, "sort": true, "uniq": true,
		"head": true, "tail": true, "tr": true, "cut": true, "wc": true,
		"chmod": true, "chown": true, "mkdir": true, "rm": true, "cp": true,
		"mv": true, "ln": true, "touch": true, "which": true, "command": true,
		"type": true, "alias": true, "export": true, "source": true, "test": true,
		"true": true, "false": true, "printf": true, "read": true, "eval": true,
		"exec": true, "bash": true, "zsh": true, "sh": true, "env": true,
		"set": true, "unset": true, "shift": true, "return": true, "exit": true,
		"local": true, "declare": true, "typeset": true, "readonly": true,
		// Zsh specific
		"autoload": true, "compinit": true, "bashcompinit": true,
		// Bash specific
		"shopt": true, "complete": true, "history": true, "bind": true,
		// Common utilities
		"tput": true, "dircolors": true, "yes": true, "watch": true,
		// Keywords and variables (common false positives)
		"if": true, "then": true, "else": true, "elif": true, "fi": true,
		"case": true, "esac": true, "for": true, "while": true, "do": true,
		"done": true, "function": true, "select": true, "time": true, "until": true,
		"SHELL": true, "PATH": true, "HOME": true, "USER": true, "TERM": true,
	}

	for _, configPath := range c.configPaths {
		file, err := os.Open(configPath)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			tools := c.extractToolsFromLine(line)
			
			// Deduplicate tools from this line
			seenInLine := make(map[string]bool)
			for _, tool := range tools {
				// Skip common commands, already found tools, and duplicates in the same line
				if !commonCommands[tool] && toolMap[tool] == "" && !seenInLine[tool] {
					toolMap[tool] = configPath // Store full path
					seenInLine[tool] = true
				}
			}
		}
	}

	// Convert map to slice and check availability
	var results []ToolStatus
	for tool, source := range toolMap {
		available := c.isToolAvailable(tool)
		results = append(results, ToolStatus{
			Tool:      tool,
			Available: available,
			Source:    source,
		})
	}

	return results
}

func (c *ShellConfigToolChecker) extractToolsFromLine(line string) []string {
	toolSet := make(map[string]bool) // Use a set to avoid duplicates
	
	// Skip comments
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return nil
	}

	// Pattern 1: $(command ...) - command substitution
	cmdSubstPattern := regexp.MustCompile(`\$\(([a-zA-Z0-9_-]+)`)
	matches := cmdSubstPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 2: `command` - backtick command substitution
	backtickPattern := regexp.MustCompile("`([a-zA-Z0-9_-]+)")
	matches = backtickPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 3: which command, command -v command, type command
	whichPattern := regexp.MustCompile(`(?:which|type)\s+([a-zA-Z0-9_-]+)`)
	matches = whichPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 4: command -v command (more specific)
	commandVPattern := regexp.MustCompile(`command\s+-v\s+([a-zA-Z0-9_-]+)`)
	matches = commandVPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 5: eval "$(command ...)"
	evalPattern := regexp.MustCompile(`eval\s+["']?\$\(([a-zA-Z0-9_-]+)`)
	matches = evalPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 6: export PATH=$PATH:$(command)
	pathExportPattern := regexp.MustCompile(`PATH=.*\$\(([a-zA-Z0-9_-]+)`)
	matches = pathExportPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Pattern 7: source <(command ...)
	sourcePattern := regexp.MustCompile(`source\s+<\(([a-zA-Z0-9_-]+)`)
	matches = sourcePattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			toolSet[match[1]] = true
		}
	}

	// Convert set to slice
	var tools []string
	for tool := range toolSet {
		tools = append(tools, tool)
	}
	return tools
}

func (c *ShellConfigToolChecker) isToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
