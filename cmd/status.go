package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
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

// printSystemInfo collects and prints system information in a formatted way
func printSystemInfo() {
	// OS Info
	hostInfo, err := host.Info()
	osStr := "Unknown"
	if err == nil {
		osStr = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	}

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// Shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("COMSPEC") // Windows
	}
	if shell == "" {
		shell = "Unknown"
	}

	// Terminal
	terminal := detectTerminal()

	// CPU Info
	cpuInfo, err := cpu.Info()
	cpuStr := "Unknown"
	if err == nil && len(cpuInfo) > 0 {
		cpuStr = fmt.Sprintf("%s (%d cores)", cpuInfo[0].ModelName, runtime.NumCPU())
	}

	// GPU Info
	gpuStr := getGPUInfo()

	fmt.Println("System Information:")
	fmt.Println("-------------------")
	fmt.Printf("OS:      %s\n", osStr)
	fmt.Printf("Host:    %s\n", hostname)
	fmt.Printf("Shell:   %s\n", shell)
	fmt.Printf("Terminal:%s\n", terminal)
	fmt.Printf("CPU:     %s\n", cpuStr)
	fmt.Printf("GPU:     %s\n", gpuStr)
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
	// Try to get parent process name (works on Unix)
	if runtime.GOOS != "windows" {
		ppid := os.Getppid()
		cmd := exec.Command("ps", "-p", fmt.Sprint(ppid), "-o", "comm=")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return "Unknown"
}

// getGPUInfo tries to get GPU model using platform-specific commands
func getGPUInfo() string {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		// Try lspci first
		cmd := exec.Command("sh", "-c", "lspci | grep -i vga | cut -d ':' -f3")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
		// Try lshw as fallback
		cmd = exec.Command("sh", "-c", "lshw -C display | grep 'product:' | head -n1 | cut -d ':' -f2")
		out, err = cmd.Output()
		if err == nil && len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
	case "darwin":
		cmd := exec.Command("system_profiler", "SPDisplaysDataType")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Chipset Model:") {
					return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				}
			}
		}
	case "windows":
		cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(strings.ToLower(line), "name") {
					return line
				}
			}
		}
	}
	return "Unavailable"
}
