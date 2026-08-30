package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/aallbrig/allbctl/pkg/telemetry"
)

var SystemctlCmd = &cobra.Command{
	Use:   "systemctl",
	Short: "Display systemd service status",
	Long:  `Display count of running system and user services, and any failed services.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, span := otel.Tracer("github.com/aallbrig/allbctl").Start(ctx, "status.systemctl")
		defer span.End()
		PrintSystemctlInfo(ctx)
	},
}

type SystemctlInfo struct {
	SystemRunning int
	SystemFailed  int
	UserRunning   int
	UserFailed    int
}

func PrintSystemctlInfo(ctx context.Context) {
	if runtime.GOOS != "linux" {
		fmt.Println("Systemctl is only available on Linux systems")
		telemetry.Logger.InfoContext(ctx, "status.systemctl", "available", false, "reason", "non-linux")
		return
	}

	if !exists("systemctl") {
		fmt.Println("Systemctl not found on this system")
		telemetry.Logger.InfoContext(ctx, "status.systemctl", "available", false, "reason", "not-installed")
		return
	}

	fmt.Println("Systemd Services:")
	fmt.Println()

	info := gatherSystemctlInfo()

	// System services
	fmt.Printf("  System Services:\n")
	if info.SystemFailed > 0 {
		fmt.Printf("    Running: %d (%d failed)\n", info.SystemRunning, info.SystemFailed)
	} else {
		fmt.Printf("    Running: %d\n", info.SystemRunning)
	}
	fmt.Println()

	// User services
	fmt.Printf("  User Services:\n")
	if info.UserFailed > 0 {
		fmt.Printf("    Running: %d (%d failed)\n", info.UserRunning, info.UserFailed)
	} else {
		fmt.Printf("    Running: %d\n", info.UserRunning)
	}

	telemetry.Logger.InfoContext(ctx, "status.systemctl",
		"available", true,
		"system_running", info.SystemRunning,
		"system_failed", info.SystemFailed,
		"user_running", info.UserRunning,
		"user_failed", info.UserFailed,
	)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.Int("systemctl.system_running", info.SystemRunning),
		attribute.Int("systemctl.system_failed", info.SystemFailed),
		attribute.Int("systemctl.user_running", info.UserRunning),
		attribute.Int("systemctl.user_failed", info.UserFailed),
	)
}

func gatherSystemctlInfo() *SystemctlInfo {
	info := &SystemctlInfo{}

	// Count system running services
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--state=running", "--no-pager", "--no-legend")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.SystemRunning = 0
		} else {
			info.SystemRunning = len(lines)
		}
	}

	// Count system failed services
	cmd = exec.Command("systemctl", "list-units", "--type=service", "--state=failed", "--no-pager", "--no-legend")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.SystemFailed = 0
		} else {
			info.SystemFailed = len(lines)
		}
	}

	// Count user running services
	cmd = exec.Command("systemctl", "--user", "list-units", "--type=service", "--state=running", "--no-pager", "--no-legend")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.UserRunning = 0
		} else {
			info.UserRunning = len(lines)
		}
	}

	// Count user failed services
	cmd = exec.Command("systemctl", "--user", "list-units", "--type=service", "--state=failed", "--no-pager", "--no-legend")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.UserFailed = 0
		} else {
			info.UserFailed = len(lines)
		}
	}

	return info
}
