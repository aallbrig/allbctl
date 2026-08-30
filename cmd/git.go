package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/aallbrig/allbctl/pkg/telemetry"
)

var GitConfigCmd = &cobra.Command{
	Use:   "git",
	Short: "Display git global configuration",
	Long:  `Display git global configuration including user name, email, and editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, span := otel.Tracer("github.com/aallbrig/allbctl").Start(ctx, "status.git")
		defer span.End()
		PrintGitConfigInfo(ctx)
	},
}

type GitConfigInfo struct {
	UserName   string
	UserEmail  string
	CoreEditor string
}

func PrintGitConfigInfo(ctx context.Context) {
	if !exists("git") {
		fmt.Println("Git is not installed")
		telemetry.Logger.InfoContext(ctx, "status.git", "installed", false)
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

	telemetry.Logger.InfoContext(ctx, "status.git",
		"installed", true,
		"user_name_set", info.UserName != "",
		"user_email_set", info.UserEmail != "",
		"editor_set", info.CoreEditor != "",
	)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.Bool("git.user_name_set", info.UserName != ""),
		attribute.Bool("git.user_email_set", info.UserEmail != ""),
		attribute.Bool("git.editor_set", info.CoreEditor != ""),
	)
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
