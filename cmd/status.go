package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
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
	gpuList := getGPUInfoList()

	// Git Info
	gitInstalled := isGitInstalled()

	// Internet Type
	internetType := getInternetType()

	// Get current user for neofetch-style header
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	// Compose header
	header := fmt.Sprintf("%s@%s", user, hostname)

	// Memory Info
	memInfo, memErr := mem.VirtualMemory()
	memStr := "Unknown"
	if memErr == nil {
		memStr = fmt.Sprintf("%.1f GiB / %.1f GiB", float64(memInfo.Used)/1e9, float64(memInfo.Total)/1e9)
	}

	// Hardware Info (from host.Info)
	hwStr := "Unknown"
	if hostInfo != nil {
		if hostInfo.Platform != "" {
			hwStr = hostInfo.Platform
		}
		if hostInfo.Hostname != "" && !strings.Contains(hwStr, hostInfo.Hostname) {
			hwStr = hostInfo.Hostname + " " + hwStr
		}
	}

	// Network Interfaces
	netIfaces, _ := net.Interfaces()
	var netSection []string
	for _, iface := range netIfaces {
		if len(iface.Addrs) == 0 {
			continue
		}
		for _, addr := range iface.Addrs {
			if strings.Contains(addr.Addr, ":") { // skip IPv6
				continue
			}
			netSection = append(netSection, fmt.Sprintf("%s: %s", iface.Name, addr.Addr))
		}
	}
	// Try to get router IP (Linux only, using 'ip route')
	routerIP := "Unknown"
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "ip route | grep default | awk '{print $3}' | head -n1")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			routerIP = strings.TrimSpace(string(out))
		}
	}

	// Print header and blank line
	fmt.Println(header)
	fmt.Println()
	// Print Host section
	fmt.Println("Host:")
	fmt.Printf("  OS:        %s\n", osStr)
	fmt.Printf("  Hostname:  %s\n", hostname)
	fmt.Printf("  Shell:     %s\n", shell)
	fmt.Printf("  Terminal:  %s\n", terminal)
	fmt.Printf("  CPU:       %s\n", cpuStr)
	fmt.Printf("  GPU(s):\n")
	if len(gpuList) == 0 {
		fmt.Printf("    Unavailable\n")
	} else {
		for _, gpu := range gpuList {
			fmt.Printf("    %s\n", gpu)
		}
	}
	fmt.Printf("  Memory:    %s\n", memStr)
	fmt.Printf("  Hardware:  %s\n", hwStr)
	fmt.Println()
	// Print Network section
	fmt.Println("Network:")
	for _, line := range netSection {
		fmt.Printf("  %s\n", line)
	}
	fmt.Printf("  Router:    %s\n", routerIP)
	fmt.Printf("  Type:      %s\n", internetType)
	fmt.Println()
	// Print Installed Software section
	fmt.Println("Installed Software:")
	fmt.Printf("  Git:       %s\n", gitInstalled)
	fmt.Println()
	fmt.Println("Tip: Run 'list-packages' to see all detected package managers and a list of installed packages.")
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

// isGitInstalled checks if git is available in PATH
func isGitInstalled() string {
	if _, err := exec.LookPath("git"); err == nil {
		return "Installed"
	}
	return "Not installed"
}

// getGPUInfoList returns a slice of GPU model names using platform-specific commands
func getGPUInfoList() []string {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		// Include both VGA and 3D controller for hybrid/NVIDIA GPUs
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
		// Try lshw as fallback (may list multiple products)
		cmd = exec.Command("sh", "-c", "lshw -C display | grep 'product:' | cut -d ':' -f2")
		out, err = cmd.Output()
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

// getInternetType tries to determine the type of internet connection (WiFi standard or Ethernet)
func getInternetType() string {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		// Try nmcli first
		cmd := exec.Command("sh", "-c", "nmcli -t -f TYPE,DEVICE,STATE,CONNECTION dev | grep ':connected' | grep '^wifi:'")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			// Get WiFi device name
			fields := strings.SplitN(strings.TrimSpace(string(out)), ":", 4)
			if len(fields) >= 2 {
				iface := fields[1]
				// Try iw to get protocol
				cmd2 := exec.Command("iw", "dev", iface, "link")
				out2, err2 := cmd2.Output()
				if err2 == nil {
					for _, line := range strings.Split(string(out2), "\n") {
						if strings.Contains(line, "tx bitrate") {
							// Not always present, but can parse
						}
						if strings.Contains(line, "connected to") {
							// Not always present
						}
						if strings.Contains(line, "freq:") {
							// Not always present
						}
						if strings.Contains(line, "802.11") {
							// e.g. 802.11ax (WiFi 6/6E), 802.11ac (WiFi 5)
							proto := strings.TrimSpace(line)
							if strings.Contains(proto, "802.11ax") {
								return "WiFi 6/6E (802.11ax)"
							} else if strings.Contains(proto, "802.11ac") {
								return "WiFi 5 (802.11ac)"
							} else if strings.Contains(proto, "802.11n") {
								return "WiFi 4 (802.11n)"
							}
						}
					}
				}
			}
			// Fallback: just say WiFi
			return "WiFi (unknown standard)"
		}
		// Check for Ethernet
		cmd = exec.Command("sh", "-c", "nmcli -t -f TYPE,STATE dev | grep '^ethernet:connected'")
		out, err = cmd.Output()
		if err == nil && len(out) > 0 {
			return "Ethernet"
		}
		return "Unknown"
	case "darwin":
		// Not implemented: would require parsing airport output
		return "Unknown"
	case "windows":
		// Not implemented: would require parsing netsh output
		return "Unknown"
	}
	return "Unknown"
}
