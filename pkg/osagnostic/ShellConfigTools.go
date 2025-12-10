package osagnostic

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	
	"github.com/fatih/color"
)

type ShellConfigTools struct {
	checker *ShellConfigToolChecker
}

func NewShellConfigTools() *ShellConfigTools {
	return &ShellConfigTools{
		checker: NewShellConfigToolChecker(),
	}
}

func (s *ShellConfigTools) Name() string {
	return "Shell Config Tools"
}

func (s *ShellConfigTools) Validate() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	
	tools := s.checker.ExtractTools()
	
	if len(tools) == 0 {
		out.WriteString("No tools found in shell config files\n")
		return out, nil
	}

	// Group tools by source file
	toolsByFile := make(map[string][]ToolStatus)
	for _, tool := range tools {
		toolsByFile[tool.Source] = append(toolsByFile[tool.Source], tool)
	}

	// Get sorted list of files
	var files []string
	for file := range toolsByFile {
		files = append(files, file)
	}
	sort.Strings(files)

	homeDir := os.Getenv("HOME")
	hasErrors := false

	// Output grouped by file
	for _, file := range files {
		// Replace home directory with $HOME for display
		displayPath := file
		if homeDir != "" && strings.HasPrefix(file, homeDir) {
			displayPath = "$HOME" + strings.TrimPrefix(file, homeDir)
		}
		
		out.WriteString(fmt.Sprintf("%s:\n", displayPath))
		
		// Sort tools within each file
		fileTools := toolsByFile[file]
		sort.Slice(fileTools, func(i, j int) bool {
			return fileTools[i].Tool < fileTools[j].Tool
		})
		
		for _, tool := range fileTools {
			if tool.Available {
				color.New(color.FgGreen).Fprint(out, "INSTALLED")
				out.WriteString(fmt.Sprintf(" %s\n", tool.Tool))
			} else {
				color.New(color.FgRed).Fprint(out, "MISSING")
				out.WriteString(fmt.Sprintf("   %s\n", tool.Tool))
				hasErrors = true
			}
		}
	}

	if hasErrors {
		// Count missing tools
		missingCount := 0
		for _, tool := range tools {
			if !tool.Available {
				missingCount++
			}
		}
		return out, fmt.Errorf("%d tool(s) referenced in shell config are not available", missingCount)
	}

	return out, nil
}

func (s *ShellConfigTools) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	
	tools := s.checker.ExtractTools()
	var missingTools []string
	
	for _, tool := range tools {
		if !tool.Available {
			missingTools = append(missingTools, tool.Tool)
		}
	}

	if len(missingTools) == 0 {
		out.WriteString("All shell config tools are available\n")
		return out, nil
	}

	out.WriteString(fmt.Sprintf("The following tools are referenced in your shell config but not installed:\n"))
	out.WriteString(fmt.Sprintf("  %s\n", strings.Join(missingTools, ", ")))
	out.WriteString("\nPlease install them manually using your package manager.\n")
	
	return out, fmt.Errorf("cannot automatically install shell config tools")
}

func (s *ShellConfigTools) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString("Shell config tools checker does not support uninstall\n")
	return out, nil
}
