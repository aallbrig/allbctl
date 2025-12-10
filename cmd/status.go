package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
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

	// Print Computer Setup Status
	fmt.Println("Computer Setup:")
	printComputerSetupStatus()
	fmt.Println()

	// Print Package Managers Available
	fmt.Println("Package Managers:")
	printPackageManagers()
	fmt.Println()

	// Print Package Manager Summary
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

// printPackageManagers displays available package managers
func printPackageManagers() {
	osType := runtime.GOOS
	systemAvailable := []string{}
	systemUnavailable := []string{}
	runtimeAvailable := []string{}

	// System package managers by OS
	systemManagers := map[string]string{}
	
	switch osType {
	case "linux":
		systemManagers["apt"] = "apt-get"
		systemManagers["dnf"] = "dnf"
		systemManagers["yum"] = "yum"
		systemManagers["pacman"] = "pacman"
		systemManagers["snap"] = "snap"
		systemManagers["flatpak"] = "flatpak"
		systemManagers["zypper"] = "zypper"
		systemManagers["apk"] = "apk"
		systemManagers["nix"] = "nix-env"
		// Check for Homebrew on Linux
		if exists("brew") {
			systemAvailable = append(systemAvailable, "homebrew")
		}
	case "darwin":
		systemManagers["homebrew"] = "brew"
		systemManagers["macports"] = "port"
		systemManagers["nix"] = "nix-env"
	case "windows":
		systemManagers["chocolatey"] = "choco"
		systemManagers["winget"] = "winget"
		systemManagers["scoop"] = "scoop"
	}

	// Check system package managers
	for name, cmd := range systemManagers {
		if exists(cmd) {
			systemAvailable = append(systemAvailable, name)
		} else {
			systemUnavailable = append(systemUnavailable, name)
		}
	}

	// Programming runtime package managers (cross-platform)
	runtimeManagers := map[string]string{
		"npm":      "npm",
		"pip":      "pip",
		"gem":      "gem",
		"cargo":    "cargo",
		"composer": "composer",
		"maven":    "mvn",
		"gradle":   "gradle",
	}

	// Check runtime package managers
	for name, cmd := range runtimeManagers {
		if exists(cmd) || (name == "pip" && exists("pip3")) {
			runtimeAvailable = append(runtimeAvailable, name)
		}
	}

	// Check for WSL on Windows
	var wslAvailable []string
	if osType == "windows" {
		wslAvailable = checkWSLPackageManagers()
	}

	// Print system package managers
	if len(systemAvailable) > 0 {
		fmt.Printf("  System:    %s\n", strings.Join(systemAvailable, ", "))
	} else {
		fmt.Printf("  System:    none\n")
	}

	// Print runtime package managers
	if len(runtimeAvailable) > 0 {
		fmt.Printf("  Runtime:   %s\n", strings.Join(runtimeAvailable, ", "))
	}

	// Print WSL package managers on Windows
	if len(wslAvailable) > 0 {
		fmt.Printf("  WSL:       %s\n", strings.Join(wslAvailable, ", "))
	} else if osType == "windows" {
		fmt.Printf("  WSL:       not available\n")
	}
}

// checkWSLPackageManagers checks for WSL availability and its package managers on Windows
func checkWSLPackageManagers() []string {
	if runtime.GOOS != "windows" {
		return nil
	}

	var wslAvailable []string

	// Check if WSL is available by trying to run a simple command
	wslTestCmd := exec.Command("wsl", "--", "echo", "test")
	if err := wslTestCmd.Run(); err != nil {
		return nil
	}

	// WSL is available, check for package managers inside WSL
	wslManagers := map[string]string{
		"apt":    "apt-get",
		"dnf":    "dnf",
		"yum":    "yum",
		"pacman": "pacman",
		"zypper": "zypper",
		"apk":    "apk",
	}

	for name, cmd := range wslManagers {
		// Check if the command exists in WSL
		checkCmd := exec.Command("wsl", "which", cmd)
		if err := checkCmd.Run(); err == nil {
			wslAvailable = append(wslAvailable, name)
		}
	}

	return wslAvailable
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
		if exists("dnf") {
			managers = append(managers, "dnf")
		}
		if exists("yum") {
			managers = append(managers, "yum")
		}
		if exists("pacman") {
			managers = append(managers, "pacman")
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
		if exists("scoop") {
			managers = append(managers, "scoop")
		}
	}

	// Programming runtime package managers (cross-platform)
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
