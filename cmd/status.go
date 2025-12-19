package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

// StatusCmd represents status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display system information (like neofetch)",
	Run: func(cmd *cobra.Command, args []string) {
		printSystemInfo()
	},
}

// printSystemInfo collects and prints system information in a neofetch-inspired format
func printSystemInfo() {
	// Get current user for neofetch-style header
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// Print header with separator
	header := fmt.Sprintf("%s@%s", user, hostname)
	separator := strings.Repeat("-", len(header))
	fmt.Println(header)
	fmt.Println(separator)

	// Host Info using gopsutil
	hostInfo, err := host.Info()
	if err == nil {
		// OS
		osStr := fmt.Sprintf("%s %s %s", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch)
		fmt.Printf("OS: %s\n", osStr)

		// Host/Hardware
		if hostInfo.VirtualizationSystem != "" {
			fmt.Printf("Host: %s (%s)\n", hostInfo.VirtualizationSystem, hostInfo.VirtualizationRole)
		}

		// Kernel
		fmt.Printf("Kernel: %s\n", hostInfo.KernelVersion)

		// Uptime
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		fmt.Printf("Uptime: %s\n", formatUptime(uptime))
	}

	// Package count
	printPackageCountInline()

	// Shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("COMSPEC")
	}
	if shell != "" {
		fmt.Printf("Shell: %s\n", shell)
	}

	// Terminal
	terminal := detectTerminal()
	if terminal != "Unknown" {
		fmt.Printf("Terminal: %s\n", terminal)
	}

	// CPU Info using gopsutil
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		cpuStr := fmt.Sprintf("%s (%d) @ %.3fGHz", cpuInfo[0].ModelName, runtime.NumCPU(), cpuInfo[0].Mhz/1000)
		fmt.Printf("CPU: %s\n", cpuStr)
	}

	// GPU Info (platform-specific fallback)
	gpuList := getGPUInfoList()
	for _, gpu := range gpuList {
		fmt.Printf("GPU: %s\n", gpu)
	}

	// Memory using gopsutil
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		memStr := fmt.Sprintf("%dMiB / %dMiB", memInfo.Used/(1024*1024), memInfo.Total/(1024*1024))
		fmt.Printf("Memory: %s\n", memStr)
	}

	fmt.Println()

	// Computer Setup Status
	fmt.Println("Computer Setup Status:")
	printComputerSetupStatus()
}

// printComputerSetupStatus runs the computer-setup status logic
func printComputerSetupStatus() {
	os := osagnostic.NewOperatingSystem()
	identifier := computerSetup.MachineIdentifier{}
	configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
	if configProvider == nil {
		fmt.Printf("  No configuration provider for %s\n", os.Name)
		return
	}

	tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
	_, out := tweaker.ConfigurationStatus()

	// Indent the output
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("  %s\n", line)
		}
	}
}

// formatUptime formats a duration into a human-readable uptime string
func formatUptime(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		if hours > 0 {
			return fmt.Sprintf("%d days, %d hours, %d mins", days, hours, mins)
		}
		return fmt.Sprintf("%d days, %d mins", days, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%d hours, %d mins", hours, mins)
	}
	return fmt.Sprintf("%d mins", mins)
}

// printPackageCountInline prints package counts in neofetch style (e.g., "Packages: 2035 (dpkg), 9 (flatpak)")
func printPackageCountInline() {
	osType := runtime.GOOS
	var counts []string

	// System package managers
	switch osType {
	case "linux":
		if exists("dpkg") {
			if pkgs := getPackages("dpkg"); pkgs != "" {
				count := countPackages("dpkg", pkgs)
				counts = append(counts, fmt.Sprintf("%d (dpkg)", count))
			}
		}
		if exists("rpm") {
			if pkgs := getPackages("rpm"); pkgs != "" {
				count := countPackages("rpm", pkgs)
				counts = append(counts, fmt.Sprintf("%d (rpm)", count))
			}
		}
		if exists("pacman") {
			if pkgs := getPackages("pacman"); pkgs != "" {
				count := countPackages("pacman", pkgs)
				counts = append(counts, fmt.Sprintf("%d (pacman)", count))
			}
		}
		if exists("snap") {
			if pkgs := getPackages("snap"); pkgs != "" {
				count := countPackages("snap", pkgs)
				counts = append(counts, fmt.Sprintf("%d (snap)", count))
			}
		}
		if exists("flatpak") {
			if pkgs := getPackages("flatpak"); pkgs != "" {
				count := countPackages("flatpak", pkgs)
				counts = append(counts, fmt.Sprintf("%d (flatpak)", count))
			}
		}
	case "darwin":
		if exists("brew") {
			if pkgs := getPackages("brew"); pkgs != "" {
				count := countPackages("brew", pkgs)
				counts = append(counts, fmt.Sprintf("%d (brew)", count))
			}
		}
	case "windows":
		if exists("choco") {
			if pkgs := getPackages("choco"); pkgs != "" {
				count := countPackages("choco", pkgs)
				counts = append(counts, fmt.Sprintf("%d (choco)", count))
			}
		}
		if exists("winget") {
			if pkgs := getPackages("winget"); pkgs != "" {
				count := countPackages("winget", pkgs)
				counts = append(counts, fmt.Sprintf("%d (winget)", count))
			}
		}
	}

	if len(counts) > 0 {
		fmt.Printf("Packages: %s\n", strings.Join(counts, ", "))
	}
}

// getGPUInfoList returns a slice of GPU model names using platform-specific commands
func getGPUInfoList() []string {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		cmd := exec.Command("sh", "-c", "lspci | grep -Ei 'vga|3d controller' | cut -d ':' -f3")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			var gpus []string
			for _, line := range lines {
				gpu := strings.TrimSpace(line)
				if gpu != "" {
					gpus = append(gpus, gpu)
				}
			}
			if len(gpus) > 0 {
				return gpus
			}
		}
	case "darwin":
		cmd := exec.Command("system_profiler", "SPDisplaysDataType")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			var gpus []string
			for _, line := range lines {
				if strings.Contains(line, "Chipset Model:") {
					gpus = append(gpus, strings.TrimSpace(strings.SplitN(line, ":", 2)[1]))
				}
			}
			if len(gpus) > 0 {
				return gpus
			}
		}
	case "windows":
		cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			var gpus []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(strings.ToLower(line), "name") {
					gpus = append(gpus, line)
				}
			}
			if len(gpus) > 0 {
				return gpus
			}
		}
	}
	return nil
}

// detectTerminal tries to determine the terminal emulator in use
func detectTerminal() string {
	term := os.Getenv("TERM_PROGRAM")
	if term != "" {
		return term
	}
	term = os.Getenv("TERM")
	if term != "" {
		return term
	}
	return "Unknown"
}
