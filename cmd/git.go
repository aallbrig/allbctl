package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var GitConfigCmd = &cobra.Command{
	Use:   "git",
	Short: "Display git global configuration",
	Long:  `Display git global configuration including user name, email, and editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintGitConfigInfo()
	},
}

type GitConfigInfo struct {
	UserName   string
	UserEmail  string
	CoreEditor string
}

func PrintGitConfigInfo() {
	if !exists("git") {
		fmt.Println("Git is not installed")
		return
	}

	fmt.Println("Git Global Configuration:")
	fmt.Println()

	info := gatherGitConfigInfo()

	if info.UserName != "" {
		fmt.Printf("  User Name:  %s\n", info.UserName)
	} else {
		fmt.Printf("  User Name:  (not set)\n")
	}

	if info.UserEmail != "" {
		fmt.Printf("  User Email: %s\n", info.UserEmail)
	} else {
		fmt.Printf("  User Email: (not set)\n")
	}

	if info.CoreEditor != "" {
		fmt.Printf("  Editor:     %s\n", info.CoreEditor)
	} else {
		fmt.Printf("  Editor:     (not set)\n")
	}
}

func gatherGitConfigInfo() *GitConfigInfo {
	info := &GitConfigInfo{}

	// Get user.name
	cmd := exec.Command("git", "config", "--global", "user.name")
	out, err := cmd.Output()
	if err == nil {
		info.UserName = strings.TrimSpace(string(out))
	}

	// Get user.email
	cmd = exec.Command("git", "config", "--global", "user.email")
	out, err = cmd.Output()
	if err == nil {
		info.UserEmail = strings.TrimSpace(string(out))
	}

	// Get core.editor
	cmd = exec.Command("git", "config", "--global", "core.editor")
	out, err = cmd.Output()
	if err == nil {
		info.CoreEditor = strings.TrimSpace(string(out))
	}

	return info
}
