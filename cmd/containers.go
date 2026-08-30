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

var ContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Display container and virtualization information",
	Long:  `Display information about containers (Docker, Podman) and virtualization status.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, span := otel.Tracer("github.com/aallbrig/allbctl").Start(ctx, "status.containers")
		defer span.End()
		PrintContainersInfo(ctx)
	},
}

type ContainerInfo struct {
	Runtime     string
	Running     int
	Images      int
	ImagesList  []string
	Virtualized bool
	VirtType    string
}

func PrintContainersInfo(ctx context.Context) {
	fmt.Println("Containers/Virtualization:")
	fmt.Println()

	// Check Docker
	dockerInfo := checkDocker()
	if dockerInfo != nil {
		fmt.Printf("  Docker:\n")
		fmt.Printf("    Running Containers: %d\n", dockerInfo.Running)
		fmt.Printf("    Images: %d\n", dockerInfo.Images)
		if len(dockerInfo.ImagesList) > 0 && len(dockerInfo.ImagesList) <= 10 {
			fmt.Printf("    Image List:\n")
			for _, img := range dockerInfo.ImagesList {
				fmt.Printf("      - %s\n", img)
			}
		}
		fmt.Println()
	}

	// Check Podman
	podmanInfo := checkPodman()
	if podmanInfo != nil {
		fmt.Printf("  Podman:\n")
		fmt.Printf("    Running Containers: %d\n", podmanInfo.Running)
		fmt.Printf("    Images: %d\n", podmanInfo.Images)
		if len(podmanInfo.ImagesList) > 0 && len(podmanInfo.ImagesList) <= 10 {
			fmt.Printf("    Image List:\n")
			for _, img := range podmanInfo.ImagesList {
				fmt.Printf("      - %s\n", img)
			}
		}
		fmt.Println()
	}

	// Check virtualization
	virtInfo := checkVirtualization()
	if virtInfo.Virtualized {
		fmt.Printf("  Virtualization: %s\n", virtInfo.VirtType)
	} else {
		fmt.Printf("  Virtualization: None detected (bare metal)\n")
	}

	if dockerInfo == nil && podmanInfo == nil {
		fmt.Println("  No container runtimes detected")
	}

	dockerRunning, dockerImages, podmanRunning, podmanImages := 0, 0, 0, 0
	if dockerInfo != nil {
		dockerRunning = dockerInfo.Running
		dockerImages = dockerInfo.Images
	}
	if podmanInfo != nil {
		podmanRunning = podmanInfo.Running
		podmanImages = podmanInfo.Images
	}
	telemetry.Logger.InfoContext(ctx, "status.containers",
		"docker_available", dockerInfo != nil,
		"docker_running", dockerRunning,
		"docker_images", dockerImages,
		"podman_available", podmanInfo != nil,
		"podman_running", podmanRunning,
		"podman_images", podmanImages,
		"virtualized", virtInfo.Virtualized,
		"virt_type", virtInfo.VirtType,
	)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.Bool("containers.docker_available", dockerInfo != nil),
		attribute.Int("containers.docker_running", dockerRunning),
		attribute.Bool("containers.podman_available", podmanInfo != nil),
		attribute.Int("containers.podman_running", podmanRunning),
		attribute.Bool("containers.virtualized", virtInfo.Virtualized),
	)
}

func checkDocker() *ContainerInfo {
	if !exists("docker") {
		return nil
	}

	info := &ContainerInfo{Runtime: "docker"}

	// Count running containers
	cmd := exec.Command("docker", "ps", "-q")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.Running = 0
		} else {
			info.Running = len(lines)
		}
	}

	// Count images and get list
	cmd = exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.Images = 0
		} else {
			info.Images = len(lines)
			info.ImagesList = lines
		}
	}

	return info
}

func checkPodman() *ContainerInfo {
	if !exists("podman") {
		return nil
	}

	info := &ContainerInfo{Runtime: "podman"}

	// Count running containers
	cmd := exec.Command("podman", "ps", "-q")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.Running = 0
		} else {
			info.Running = len(lines)
		}
	}

	// Count images and get list
	cmd = exec.Command("podman", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			info.Images = 0
		} else {
			info.Images = len(lines)
			info.ImagesList = lines
		}
	}

	return info
}

func checkVirtualization() *ContainerInfo {
	info := &ContainerInfo{Virtualized: false}

	if runtime.GOOS != "linux" {
		return info
	}

	// Use systemd-detect-virt
	if exists("systemd-detect-virt") {
		cmd := exec.Command("systemd-detect-virt")
		out, err := cmd.Output()
		if err == nil {
			virtType := strings.TrimSpace(string(out))
			if virtType != "none" {
				info.Virtualized = true
				info.VirtType = virtType
			}
		}
	}

	return info
}
