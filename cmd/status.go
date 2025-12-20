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
	"github.com/shirou/gopsutil/v4/net"
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

// printSystemInfo collects and prints system information in a structured format
func printSystemInfo() {
	// Get current user for header
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// Print header
	header := fmt.Sprintf("%s@%s", user, hostname)
	fmt.Println(header)
	fmt.Println()

	// Host Info using gopsutil
	hostInfo, err := host.Info()
	osStr := "Unknown"
	if err == nil {
		osStr = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	}

	// Shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("COMSPEC")
	}
	if shell == "" {
		shell = "Unknown"
	}

	// Terminal
	terminal := detectTerminal()

	// CPU Info using gopsutil
	cpuInfo, err := cpu.Info()
	cpuStr := "Unknown"
	if err == nil && len(cpuInfo) > 0 {
		cpuStr = fmt.Sprintf("%s (%d cores)", cpuInfo[0].ModelName, runtime.NumCPU())
	}

	// GPU Info
	gpuList := getGPUInfoList()

	// Memory using gopsutil
	memInfo, err := mem.VirtualMemory()
	memStr := "Unknown"
	if err == nil {
		memStr = fmt.Sprintf("%.1f GiB / %.1f GiB", float64(memInfo.Used)/1e9, float64(memInfo.Total)/1e9)
	}

	// Hardware Info
	hwStr := "Unknown"
	if hostInfo != nil {
		if hostInfo.Platform != "" {
			hwStr = hostInfo.Platform
		}
		if hostInfo.Hostname != "" && !strings.Contains(hwStr, hostInfo.Hostname) {
			hwStr = hostInfo.Hostname + " " + hwStr
		}
	}

	// Runtimes (inline)
	runtimesInline := detectRuntimesInline()

	// Print system information (no "Host:" header)
	fmt.Printf("OS:        %s\n", osStr)
	fmt.Printf("Hostname:  %s\n", hostname)
	fmt.Printf("Shell:     %s\n", shell)
	fmt.Printf("Terminal:  %s\n", terminal)
	fmt.Printf("CPU:       %s\n", cpuStr)
	fmt.Printf("GPU(s):\n")
	if len(gpuList) == 0 {
		fmt.Printf("  Unavailable\n")
	} else {
		for _, gpu := range gpuList {
			fmt.Printf("  %s\n", gpu)
		}
	}
	fmt.Printf("Memory:    %s\n", memStr)
	fmt.Printf("Hardware:  %s\n", hwStr)
	if runtimesInline != "" {
		fmt.Printf("Runtimes:  %s\n", runtimesInline)
	}
	fmt.Println()

	// Network section
	fmt.Println("Network:")
	printNetworkInfo()
	fmt.Println()

	// Workstation Bootstrap Status
	fmt.Println("Workstation Bootstrap Status:")
	printComputerSetupStatus()
	fmt.Println()

	// Package Managers section
	fmt.Println("Package Managers:")
	printPackageManagers()
	fmt.Println()

	// Packages section
	fmt.Println("Packages:")
	printPackageSummary()
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

// printNetworkInfo prints network interface information
func printNetworkInfo() {
	netIfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("  Unable to detect network interfaces\n")
		return
	}

	for _, iface := range netIfaces {
		if len(iface.Addrs) == 0 {
			continue
		}
		for _, addr := range iface.Addrs {
			if strings.Contains(addr.Addr, ":") {
				continue // skip IPv6
			}
			fmt.Printf("  %s: %s\n", iface.Name, addr.Addr)
		}
	}

	// Try to get router IP (Linux only)
	routerIP := "Unknown"
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "ip route | grep default | awk '{print $3}' | head -n1")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			routerIP = strings.TrimSpace(string(out))
		}
	}
	fmt.Printf("  Router:    %s\n", routerIP)

	// Internet type
	internetType := getInternetType()
	fmt.Printf("  Type:      %s\n", internetType)
}

// getInternetType tries to determine the type of internet connection
func getInternetType() string {
	if runtime.GOOS == "linux" {
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
						if strings.Contains(line, "802.11") {
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
			return "WiFi (unknown standard)"
		}
		// Check for Ethernet
		cmd = exec.Command("sh", "-c", "nmcli -t -f TYPE,STATE dev | grep '^ethernet:connected'")
		out, err = cmd.Output()
		if err == nil && len(out) > 0 {
			return "Ethernet"
		}
	}
	return "Unknown"
}

// printPackageManagers displays available package managers
func printPackageManagers() {
	osType := runtime.GOOS
	systemAvailable := []string{}
	runtimeAvailable := []string{}

	// System package managers
	switch osType {
	case "linux":
		if exists("apt-get") {
			systemAvailable = append(systemAvailable, "apt")
		}
		if exists("flatpak") {
			systemAvailable = append(systemAvailable, "flatpak")
		}
		if exists("snap") {
			systemAvailable = append(systemAvailable, "snap")
		}
		if exists("dnf") {
			systemAvailable = append(systemAvailable, "dnf")
		}
		if exists("yum") {
			systemAvailable = append(systemAvailable, "yum")
		}
		if exists("pacman") {
			systemAvailable = append(systemAvailable, "pacman")
		}
	case "darwin":
		if exists("brew") {
			systemAvailable = append(systemAvailable, "homebrew")
		}
	case "windows":
		if exists("choco") {
			systemAvailable = append(systemAvailable, "chocolatey")
		}
		if exists("winget") {
			systemAvailable = append(systemAvailable, "winget")
		}
	}

	// Programming runtime package managers
	if exists("npm") {
		runtimeAvailable = append(runtimeAvailable, "npm")
	}
	if exists("pip") || exists("pip3") {
		runtimeAvailable = append(runtimeAvailable, "pip")
	}
	if exists("gem") {
		runtimeAvailable = append(runtimeAvailable, "gem")
	}
	if exists("cargo") {
		runtimeAvailable = append(runtimeAvailable, "cargo")
	}
	if exists("go") {
		runtimeAvailable = append(runtimeAvailable, "go")
	}

	// Print
	if len(systemAvailable) > 0 {
		fmt.Printf("  System:    %s\n", strings.Join(systemAvailable, ", "))
	}
	if len(runtimeAvailable) > 0 {
		fmt.Printf("  Runtime:   %s\n", strings.Join(runtimeAvailable, ", "))
	}
}

// printPackageSummary runs the list-packages summary logic
func printPackageSummary() {
	osType := runtime.GOOS
	var managers []string

	// System package managers
	switch osType {
	case "linux":
		if exists("apt-mark") {
			managers = append(managers, "apt")
		}
		if exists("snap") {
			managers = append(managers, "snap")
		}
		if exists("flatpak") {
			managers = append(managers, "flatpak")
		}
	case "darwin":
		if exists("brew") {
			managers = append(managers, "brew")
		}
	case "windows":
		if exists("choco") {
			managers = append(managers, "choco")
		}
		if exists("winget") {
			managers = append(managers, "winget")
		}
	}

	// Programming runtime package managers
	if exists("npm") {
		managers = append(managers, "npm")
	}
	if exists("pip") || exists("pip3") {
		managers = append(managers, "pip")
	}
	if exists("gem") {
		managers = append(managers, "gem")
	}
	if exists("cargo") {
		managers = append(managers, "cargo")
	}
	if exists("go") {
		managers = append(managers, "go")
	}

	if len(managers) == 0 {
		fmt.Println("  No package managers detected")
		return
	}

	// Summary mode: just count packages
	for _, m := range managers {
		pkgs := getPackages(m)
		if pkgs != "" {
			count := countPackages(m, pkgs)
			fmt.Printf("  %-15s %d packages\n", m+":", count)
		}
	}
}
