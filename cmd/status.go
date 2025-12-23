package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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

	// AI Agents section
	fmt.Println("AI Agents:")
	printAIAgents()
	fmt.Println()

	// Package Managers section
	fmt.Println("Package Managers:")
	printPackageManagers()
	fmt.Println()

	// Packages section
	fmt.Println("Packages:")
	printPackageSummary()
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

// AIAgent represents an AI coding assistant
type AIAgent struct {
	Name    string
	Version string
}

// detectAIAgents detects available AI coding assistants
func detectAIAgents() []AIAgent {
	var agents []AIAgent

	// GitHub Copilot CLI
	if exists("copilot") {
		version := getAIAgentVersion("copilot")
		agents = append(agents, AIAgent{Name: "copilot", Version: version})
	}

	// Claude Code (if it exists as a CLI)
	if exists("claude") {
		version := getAIAgentVersion("claude")
		agents = append(agents, AIAgent{Name: "claude", Version: version})
	}

	// Cursor AI
	if exists("cursor") {
		version := getAIAgentVersion("cursor")
		agents = append(agents, AIAgent{Name: "cursor", Version: version})
	}

	// Aider
	if exists("aider") {
		version := getAIAgentVersion("aider")
		agents = append(agents, AIAgent{Name: "aider", Version: version})
	}

	// Continue.dev (if it has a CLI)
	if exists("continue") {
		version := getAIAgentVersion("continue")
		agents = append(agents, AIAgent{Name: "continue", Version: version})
	}

	// Cody (Sourcegraph)
	if exists("cody") {
		version := getAIAgentVersion("cody")
		agents = append(agents, AIAgent{Name: "cody", Version: version})
	}

	// Tabby (local AI)
	if exists("tabby") {
		version := getAIAgentVersion("tabby")
		agents = append(agents, AIAgent{Name: "tabby", Version: version})
	}

	// Amazon CodeWhisperer
	if exists("codewhisperer") {
		version := getAIAgentVersion("codewhisperer")
		agents = append(agents, AIAgent{Name: "codewhisperer", Version: version})
	}

	// Ollama (local LLM runner)
	if exists("ollama") {
		version := getAIAgentVersion("ollama")
		agents = append(agents, AIAgent{Name: "ollama", Version: version})
	}

	return agents
}

// getAIAgentVersion returns the version of an AI agent
func getAIAgentVersion(agent string) string {
	var cmd *exec.Cmd

	switch agent {
	case "copilot":
		cmd = exec.Command("copilot", "--version")
	case "claude":
		cmd = exec.Command("claude", "--version")
	case "cursor":
		cmd = exec.Command("cursor", "--version")
	case "aider":
		cmd = exec.Command("aider", "--version")
	case "continue":
		cmd = exec.Command("continue", "--version")
	case "cody":
		cmd = exec.Command("cody", "--version")
	case "tabby":
		cmd = exec.Command("tabby", "--version")
	case "codewhisperer":
		cmd = exec.Command("codewhisperer", "--version")
	case "ollama":
		cmd = exec.Command("ollama", "--version")
	default:
		return ""
	}

	if cmd == nil {
		return ""
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return extractAIAgentVersion(agent, version)
}

// extractAIAgentVersion extracts clean version from AI agent output
func extractAIAgentVersion(agent, output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	// Take first line only
	lines := strings.Split(output, "\n")
	firstLine := strings.TrimSpace(lines[0])

	switch agent {
	case "copilot":
		// "0.0.365" - just the version number
		return firstLine
	case "aider":
		// "aider 0.50.0" -> "0.50.0"
		if strings.HasPrefix(firstLine, "aider ") {
			return strings.TrimPrefix(firstLine, "aider ")
		}
		return firstLine
	case "ollama":
		// "ollama version is 0.13.5" -> "0.13.5"
		if strings.Contains(firstLine, "version is ") {
			parts := strings.Split(firstLine, "version is ")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		// Fallback to generic extraction
		fields := strings.Fields(firstLine)
		for _, field := range fields {
			if strings.Contains(field, ".") {
				field = strings.Trim(field, "()[]{}\"',vV")
				if len(field) > 0 && (field[0] >= '0' && field[0] <= '9') {
					return field
				}
			}
		}
		return firstLine
	case "cursor", "claude", "continue", "cody", "tabby", "codewhisperer":
		// Try to extract version number
		fields := strings.Fields(firstLine)
		for _, field := range fields {
			if strings.Contains(field, ".") {
				field = strings.Trim(field, "()[]{}\"',vV")
				if len(field) > 0 && (field[0] >= '0' && field[0] <= '9') {
					return field
				}
			}
		}
		return firstLine
	}

	return firstLine
}

// printAIAgents displays available AI coding assistants
func printAIAgents() {
	agents := detectAIAgents()

	if len(agents) == 0 {
		fmt.Printf("  No AI agents detected\n")
		return
	}

	var agentStrings []string
	for _, agent := range agents {
		if agent.Version != "" {
			agentStrings = append(agentStrings, fmt.Sprintf("%s (%s)", agent.Name, agent.Version))
		} else {
			agentStrings = append(agentStrings, agent.Name)
		}
	}

	fmt.Printf("  %s\n", strings.Join(agentStrings, ", "))
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
			version := getPackageManagerVersion("apt")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("apt (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "apt")
			}
		}
		if exists("flatpak") {
			version := getPackageManagerVersion("flatpak")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("flatpak (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "flatpak")
			}
		}
		if exists("snap") {
			version := getPackageManagerVersion("snap")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("snap (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "snap")
			}
		}
		if exists("dnf") {
			version := getPackageManagerVersion("dnf")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("dnf (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "dnf")
			}
		}
		if exists("yum") {
			version := getPackageManagerVersion("yum")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("yum (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "yum")
			}
		}
		if exists("pacman") {
			version := getPackageManagerVersion("pacman")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("pacman (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "pacman")
			}
		}
	case "darwin":
		if exists("brew") {
			version := getPackageManagerVersion("brew")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("homebrew (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "homebrew")
			}
		}
	case "windows":
		if exists("choco") {
			version := getPackageManagerVersion("choco")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("chocolatey (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "chocolatey")
			}
		}
		if exists("winget") {
			version := getPackageManagerVersion("winget")
			if version != "" {
				systemAvailable = append(systemAvailable, fmt.Sprintf("winget (%s)", version))
			} else {
				systemAvailable = append(systemAvailable, "winget")
			}
		}
	}

	// Programming runtime package managers
	if exists("npm") {
		version := getPackageManagerVersion("npm")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("npm (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "npm")
		}
	}
	if exists("pip") || exists("pip3") {
		version := getPackageManagerVersion("pip")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("pip (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "pip")
		}
	}
	if exists("pipx") {
		version := getPackageManagerVersion("pipx")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("pipx (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "pipx")
		}
	}
	if exists("gem") {
		version := getPackageManagerVersion("gem")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("gem (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "gem")
		}
	}
	if exists("cargo") {
		version := getPackageManagerVersion("cargo")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("cargo (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "cargo")
		}
	}
	if exists("go") {
		version := getPackageManagerVersion("go")
		if version != "" {
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("go (%s)", version))
		} else {
			runtimeAvailable = append(runtimeAvailable, "go")
		}
	}

	// Print
	if len(systemAvailable) > 0 {
		fmt.Printf("  System:    %s\n", strings.Join(systemAvailable, ", "))
	}
	if len(runtimeAvailable) > 0 {
		fmt.Printf("  Runtime:   %s\n", strings.Join(runtimeAvailable, ", "))
	}
}

// getPackageManagerVersion returns the version of a package manager
func getPackageManagerVersion(manager string) string {
	var cmd *exec.Cmd

	switch manager {
	case "apt":
		cmd = exec.Command("apt-get", "--version")
	case "flatpak":
		cmd = exec.Command("flatpak", "--version")
	case "snap":
		cmd = exec.Command("snap", "--version")
	case "dnf":
		cmd = exec.Command("dnf", "--version")
	case "yum":
		cmd = exec.Command("yum", "--version")
	case "pacman":
		cmd = exec.Command("pacman", "--version")
	case "brew":
		cmd = exec.Command("brew", "--version")
	case "choco":
		cmd = exec.Command("choco", "--version")
	case "winget":
		cmd = exec.Command("winget", "--version")
	case "npm":
		cmd = exec.Command("npm", "--version")
	case "pip":
		if exists("pip3") {
			cmd = exec.Command("pip3", "--version")
		} else {
			cmd = exec.Command("pip", "--version")
		}
	case "pipx":
		cmd = exec.Command("pipx", "--version")
	case "gem":
		cmd = exec.Command("gem", "--version")
	case "cargo":
		cmd = exec.Command("cargo", "--version")
	case "go":
		cmd = exec.Command("go", "version")
	default:
		return ""
	}

	if cmd == nil {
		return ""
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return extractPackageManagerVersion(manager, version)
}

// extractPackageManagerVersion extracts clean version from package manager output
func extractPackageManagerVersion(manager, output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	// Take first line only
	if idx := strings.Index(output, "\n"); idx >= 0 {
		output = output[:idx]
	}

	switch manager {
	case "apt":
		// "apt 2.8.3 (amd64)" -> "2.8.3"
		if strings.HasPrefix(output, "apt ") {
			parts := strings.Fields(output)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	case "flatpak":
		// "Flatpak 1.14.6" -> "1.14.6"
		if strings.HasPrefix(output, "Flatpak ") {
			return strings.TrimPrefix(output, "Flatpak ")
		}
	case "snap":
		// "snap    2.63" -> "2.63"
		if strings.HasPrefix(output, "snap") {
			parts := strings.Fields(output)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	case "pip":
		// "pip 24.0 from /usr/lib..." -> "24.0"
		if strings.HasPrefix(output, "pip ") {
			parts := strings.Fields(output)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	case "npm", "pipx", "gem":
		// These usually return just version number
		return output
	case "cargo":
		// "cargo 1.70.0 (7c2f85da6 2023-05-31)" -> "1.70.0"
		if strings.HasPrefix(output, "cargo ") {
			parts := strings.Fields(output)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	case "go":
		// "go version go1.25.5 linux/amd64" -> "1.25.5"
		if strings.HasPrefix(output, "go version go") {
			parts := strings.Fields(output)
			if len(parts) >= 3 {
				return strings.TrimPrefix(parts[2], "go")
			}
		}
	case "brew":
		// "Homebrew 4.0.0" or just "4.0.0"
		if strings.HasPrefix(output, "Homebrew ") {
			return strings.TrimPrefix(output, "Homebrew ")
		}
		return output
	case "dnf", "yum":
		// Usually just version number
		return output
	case "pacman":
		// "Pacman v6.0.1 - libalpm v13.0.1" -> "6.0.1"
		if strings.Contains(output, "Pacman v") {
			parts := strings.Fields(output)
			for _, part := range parts {
				if strings.HasPrefix(part, "v") && len(part) > 1 {
					return strings.TrimPrefix(part, "v")
				}
			}
		}
	case "choco", "winget":
		// Usually just version number
		return output
	}

	// Generic: try to extract version-like pattern
	fields := strings.Fields(output)
	for _, field := range fields {
		if strings.Contains(field, ".") {
			field = strings.Trim(field, "()[]{}\"',")
			if len(field) > 0 && (field[0] >= '0' && field[0] <= '9') {
				return field
			}
		}
	}

	return output
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
	if exists("pipx") {
		managers = append(managers, "pipx")
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
