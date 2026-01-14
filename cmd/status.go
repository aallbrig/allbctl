package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

// browserVersionRegex is used to extract version numbers from browser output
var browserVersionRegex = regexp.MustCompile(`\d+\.\d+[\d.]*`)

// StatusCmd represents status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display system information (like neofetch)",
	Run: func(cmd *cobra.Command, args []string) {
		printSystemInfo()
	},
}

// BrowserInfo holds browser information
type BrowserInfo struct {
	Name    string
	Version string
}

// detectBrowsers detects installed web browsers and their versions
func detectBrowsers() []BrowserInfo {
	var browsers []BrowserInfo
	osType := runtime.GOOS

	switch osType {
	case "linux":
		browsers = detectLinuxBrowsers()
	case "darwin":
		browsers = detectMacBrowsers()
	case "windows":
		browsers = detectWindowsBrowsers()
	}

	return browsers
}

// detectLinuxBrowsers detects browsers on Linux
func detectLinuxBrowsers() []BrowserInfo {
	var browsers []BrowserInfo

	// Chrome/Chromium
	if version := getBrowserVersion("google-chrome", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Chrome", Version: version})
	} else if version := getBrowserVersion("chromium", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Chromium", Version: version})
	} else if version := getBrowserVersion("chromium-browser", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Chromium", Version: version})
	} else if version := getFlatpakBrowserVersion("org.chromium.Chromium"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Chromium", Version: version})
	}

	// Firefox
	if version := getBrowserVersion("firefox", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Firefox", Version: version})
	} else if version := getFlatpakBrowserVersion("org.mozilla.firefox"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Firefox", Version: version})
	}

	// Brave
	if version := getBrowserVersion("brave-browser", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Brave", Version: version})
	} else if version := getBrowserVersion("brave", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Brave", Version: version})
	} else if version := getFlatpakBrowserVersion("com.brave.Browser"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Brave", Version: version})
	}

	// Edge
	if version := getBrowserVersion("microsoft-edge", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Edge", Version: version})
	} else if version := getBrowserVersion("microsoft-edge-stable", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Edge", Version: version})
	} else if version := getFlatpakBrowserVersion("com.microsoft.Edge"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Edge", Version: version})
	}

	// Opera
	if version := getBrowserVersion("opera", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Opera", Version: version})
	} else if version := getFlatpakBrowserVersion("com.opera.Opera"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Opera", Version: version})
	}

	// Vivaldi
	if version := getBrowserVersion("vivaldi", "--version"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Vivaldi", Version: version})
	} else if version := getFlatpakBrowserVersion("com.vivaldi.Vivaldi"); version != "" {
		browsers = append(browsers, BrowserInfo{Name: "Vivaldi", Version: version})
	}

	return browsers
}

// detectMacBrowsers detects browsers on macOS
func detectMacBrowsers() []BrowserInfo {
	var browsers []BrowserInfo

	// Check for browsers in /Applications
	appPaths := map[string]string{
		"Chrome":  "/Applications/Google Chrome.app",
		"Firefox": "/Applications/Firefox.app",
		"Safari":  "/Applications/Safari.app",
		"Brave":   "/Applications/Brave Browser.app",
		"Edge":    "/Applications/Microsoft Edge.app",
		"Opera":   "/Applications/Opera.app",
		"Vivaldi": "/Applications/Vivaldi.app",
	}

	for name, appPath := range appPaths {
		if _, err := os.Stat(appPath); err == nil {
			version := getMacAppVersion(appPath)
			if version != "" {
				browsers = append(browsers, BrowserInfo{Name: name, Version: version})
			} else {
				browsers = append(browsers, BrowserInfo{Name: name, Version: "installed"})
			}
		}
	}

	return browsers
}

// detectWindowsBrowsers detects browsers on Windows
func detectWindowsBrowsers() []BrowserInfo {
	var browsers []BrowserInfo

	// Check common browser paths using filepath.Join for cross-platform compatibility
	programFiles := `C:\Program Files`
	programFilesX86 := `C:\Program Files (x86)`

	browserPaths := map[string][]string{
		"Chrome": {
			filepath.Join(programFiles, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(programFilesX86, "Google", "Chrome", "Application", "chrome.exe"),
		},
		"Firefox": {
			filepath.Join(programFiles, "Mozilla Firefox", "firefox.exe"),
			filepath.Join(programFilesX86, "Mozilla Firefox", "firefox.exe"),
		},
		"Edge": {
			filepath.Join(programFilesX86, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(programFiles, "Microsoft", "Edge", "Application", "msedge.exe"),
		},
		"Brave": {
			filepath.Join(programFiles, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
			filepath.Join(programFilesX86, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
		},
		"Opera": {
			filepath.Join(programFiles, "Opera", "launcher.exe"),
			filepath.Join(programFilesX86, "Opera", "launcher.exe"),
		},
		"Vivaldi": {
			filepath.Join(programFiles, "Vivaldi", "Application", "vivaldi.exe"),
			filepath.Join(programFilesX86, "Vivaldi", "Application", "vivaldi.exe"),
		},
	}

	for name, paths := range browserPaths {
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				version := getWindowsBrowserVersion(path)
				if version != "" {
					browsers = append(browsers, BrowserInfo{Name: name, Version: version})
				} else {
					browsers = append(browsers, BrowserInfo{Name: name, Version: "installed"})
				}
				break
			}
		}
	}

	return browsers
}

// getBrowserVersion gets browser version using command line
func getBrowserVersion(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return parseBrowserVersion(version)
}

// getFlatpakBrowserVersion gets browser version from Flatpak
func getFlatpakBrowserVersion(appID string) string {
	cmd := exec.Command("flatpak", "run", appID, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return parseBrowserVersion(version)
}

// parseBrowserVersion extracts version number from browser output
func parseBrowserVersion(output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	// Common patterns:
	// "Google Chrome 120.0.6099.109"
	// "Chromium 120.0.6099.109"
	// "Mozilla Firefox 121.0"
	// "Brave 1.61.109 Chromium: 120.0.6099.109"
	// Multi-line with warnings:
	// "Gtk-Message: ...\nChromium 143.0.7499.169"

	// Process each line to find the browser version line
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip known warning/message lines
		lineLower := strings.ToLower(line)
		if strings.HasPrefix(lineLower, "gtk-message:") ||
			strings.HasPrefix(lineLower, "warning:") ||
			strings.HasPrefix(lineLower, "error:") {
			continue
		}

		// Try to find version pattern in each field
		fields := strings.Fields(line)
		for i, field := range fields {
			// Skip known browser name prefixes
			fieldLower := strings.ToLower(field)
			if fieldLower == "google" || fieldLower == "chrome" || fieldLower == "chromium" ||
				fieldLower == "mozilla" || fieldLower == "firefox" || fieldLower == "brave" ||
				fieldLower == "microsoft" || fieldLower == "edge" || fieldLower == "opera" ||
				fieldLower == "vivaldi" {
				continue
			}

			// Check if this is a version label followed by version
			if fieldLower == "version" && i+1 < len(fields) {
				if version := browserVersionRegex.FindString(fields[i+1]); version != "" {
					return version
				}
			}

			// Look for version pattern directly
			if version := browserVersionRegex.FindString(field); version != "" {
				return version
			}
		}
	}

	return ""
}

// getMacAppVersion gets version from macOS app bundle
func getMacAppVersion(appPath string) string {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	cmd := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString")
	output, err := cmd.Output()
	if err != nil {
		// Try alternative version key
		cmd = exec.Command("defaults", "read", plistPath, "CFBundleVersion")
		output, err = cmd.Output()
		if err != nil {
			return ""
		}
	}

	version := strings.TrimSpace(string(output))
	return version
}

// isVersionString checks if a string looks like a version number
func isVersionString(s string) bool {
	if s == "" {
		return false
	}
	// Check if it starts with a digit and contains a dot
	return s[0] >= '0' && s[0] <= '9' && strings.Contains(s, ".")
}

// getWindowsBrowserVersion gets browser version on Windows
func getWindowsBrowserVersion(browserPath string) string {
	// Try to get version from file properties using wmic
	dir := filepath.Dir(browserPath)

	// Look for version info in the directory (version folders for Chrome/Edge)
	// These browsers store their version as a folder name in the Application directory
	if strings.Contains(browserPath, filepath.Join("Google", "Chrome")) ||
		strings.Contains(browserPath, filepath.Join("Microsoft", "Edge")) {
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					// Check if directory name looks like a version
					name := entry.Name()
					if isVersionString(name) {
						return name
					}
				}
			}
		}
	}

	return ""
}

// printBrowsers displays detected browsers
func printBrowsers(browsers []BrowserInfo) {
	if len(browsers) == 0 {
		return
	}

	var browserStrings []string
	for _, browser := range browsers {
		if browser.Version != "" && browser.Version != "installed" {
			browserStrings = append(browserStrings, fmt.Sprintf("%s (%s)", browser.Name, browser.Version))
		} else {
			browserStrings = append(browserStrings, browser.Name)
		}
	}

	if len(browserStrings) > 0 {
		fmt.Printf("  %s\n", strings.Join(browserStrings, ", "))
	}
}

// printSystemInfo collects and prints system information in a structured format
func printSystemInfo() {
	// Start package detection early (runs in background)
	packagesFuture := StartPackageSummary()

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

	// CPU Info - get detailed information
	cpuDetails := getDetailedCPUInfo()

	// GPU Info - get detailed information
	gpuDetails := getDetailedGPUInfo()

	// Memory using gopsutil - show only total installed
	memInfo, err := mem.VirtualMemory()
	memStr := "Unknown"
	if err == nil {
		memStr = fmt.Sprintf("%.1f GiB", float64(memInfo.Total)/1e9)
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

	// Print system information (no "Host:" header)
	fmt.Printf("OS:        %s\n", osStr)
	fmt.Printf("Hostname:  %s\n", hostname)
	fmt.Printf("Shell:     %s\n", shell)
	fmt.Printf("Terminal:  %s\n", terminal)
	fmt.Printf("CPU:\n")
	printCPUInfo(cpuDetails)
	fmt.Printf("GPU(s):\n")
	printGPUInfo(gpuDetails)
	fmt.Printf("Memory:    %s\n", memStr)

	// Disks
	diskInfo := getDiskInfo()
	if diskInfo != "" {
		fmt.Printf("Disks:     %s\n", diskInfo)
	}

	fmt.Printf("Hardware:  %s\n", hwStr)

	// Runtimes (using shared function)
	runtimesInline := detectRuntimesInline()
	if runtimesInline != "" {
		fmt.Printf("Runtimes:  %s\n", runtimesInline)
	}

	// Databases summary
	PrintDatabaseSummaryForStatus()
	fmt.Println()

	// Network section
	printNetworkInfo()
	fmt.Println()

	// Ports section
	printPortsSummary()
	fmt.Println()

	// Browsers section
	browsers := detectBrowsers()
	if len(browsers) > 0 {
		fmt.Println("Browsers:")
		printBrowsers(browsers)
		fmt.Println()
	}

	// AI Agents section
	fmt.Println("AI Agents:")
	printAIAgents()
	fmt.Println()

	// Package Managers section
	fmt.Println("Package Managers:")
	printPackageManagers()
	fmt.Println()

	// Packages section - wait for background detection to complete
	fmt.Println("Packages:")
	if packagesFuture != nil {
		packagesFuture.PrintResults()
	} else {
		fmt.Println("  No package managers detected")
	}
	fmt.Println()

	// Projects section
	printProjectsInline()
}

// GPUInfo holds detailed GPU information
type GPUInfo struct {
	Name          string
	Vendor        string
	Memory        string
	Driver        string
	ComputeCap    string
	ClockGraphics string
	ClockMemory   string
}

// getDetailedGPUInfo gathers detailed GPU information from multiple sources
func getDetailedGPUInfo() []GPUInfo {
	var gpus []GPUInfo

	osType := runtime.GOOS

	// Try nvidia-smi first for NVIDIA GPUs
	if exists("nvidia-smi") {
		nvidiaGPUs := getNvidiaGPUInfo()
		gpus = append(gpus, nvidiaGPUs...)
	}

	// Always run platform-specific detection to find all GPUs (e.g., integrated Intel GPUs)
	var platformGPUs []GPUInfo
	switch osType {
	case "linux":
		platformGPUs = getLinuxGPUInfo()
	case "darwin":
		platformGPUs = getMacGPUInfo()
	case "windows":
		platformGPUs = getWindowsGPUInfo()
	}

	// Merge platform-specific GPUs with NVIDIA GPUs, avoiding duplicates
	// Skip NVIDIA GPUs from platform detection if we already have them from nvidia-smi
	hasNvidiaFromSmi := false
	for _, gpu := range gpus {
		if gpu.Vendor == "NVIDIA" {
			hasNvidiaFromSmi = true
			break
		}
	}

	for _, platformGPU := range platformGPUs {
		// Skip NVIDIA GPUs from platform detection if nvidia-smi already detected them
		if hasNvidiaFromSmi && platformGPU.Vendor == "NVIDIA" {
			continue
		}

		isDuplicate := false
		for _, gpu := range gpus {
			// Check if this GPU is already in the list (basic duplicate detection by name)
			if gpu.Name == platformGPU.Name {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			gpus = append(gpus, platformGPU)
		}
	}

	return gpus
}

// getNvidiaGPUInfo gets GPU information from nvidia-smi
func getNvidiaGPUInfo() []GPUInfo {
	var gpus []GPUInfo

	cmd := exec.Command("nvidia-smi", "--query-gpu=name,memory.total,driver_version,compute_cap,clocks.current.graphics,clocks.current.memory", "--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) >= 6 {
			gpu := GPUInfo{
				Name:          strings.TrimSpace(fields[0]),
				Vendor:        "NVIDIA",
				Memory:        strings.TrimSpace(fields[1]) + " MiB",
				Driver:        strings.TrimSpace(fields[2]),
				ComputeCap:    strings.TrimSpace(fields[3]),
				ClockGraphics: strings.TrimSpace(fields[4]) + " MHz",
				ClockMemory:   strings.TrimSpace(fields[5]) + " MHz",
			}
			gpus = append(gpus, gpu)
		}
	}

	return gpus
}

// getLinuxGPUInfo gets GPU information on Linux using lspci
func getLinuxGPUInfo() []GPUInfo {
	var gpus []GPUInfo

	cmd := exec.Command("sh", "-c", "lspci | grep -Ei 'vga|3d controller'")
	out, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse lspci output: "00:08.0 VGA compatible controller: Vendor Name Device Name"
		parts := strings.SplitN(line, ":", 3)
		if len(parts) >= 3 {
			name := strings.TrimSpace(parts[2])
			vendor := detectVendor(name)

			gpu := GPUInfo{
				Name:   name,
				Vendor: vendor,
			}

			// Note: Could add AMD-specific detection here with rocm-smi if needed

			gpus = append(gpus, gpu)
		}
	}

	return gpus
}

// getMacGPUInfo gets GPU information on macOS
func getMacGPUInfo() []GPUInfo {
	var gpus []GPUInfo

	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	out, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(string(out), "\n")
	var currentGPU *GPUInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Chipset Model:") {
			if currentGPU != nil {
				gpus = append(gpus, *currentGPU)
			}
			name := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			currentGPU = &GPUInfo{
				Name:   name,
				Vendor: detectVendor(name),
			}
		} else if currentGPU != nil {
			if strings.Contains(line, "VRAM (Total):") || strings.Contains(line, "VRAM (Dynamic, Max):") {
				currentGPU.Memory = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			}
		}
	}

	if currentGPU != nil {
		gpus = append(gpus, *currentGPU)
	}

	return gpus
}

// getWindowsGPUInfo gets GPU information on Windows
func getWindowsGPUInfo() []GPUInfo {
	var gpus []GPUInfo

	cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "Name,AdapterRAM,DriverVersion", "/format:csv")
	out, err := cmd.Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		// Skip header and empty lines
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 3 {
			name := strings.TrimSpace(fields[2])
			if name == "" || strings.HasPrefix(strings.ToLower(name), "name") {
				continue
			}

			gpu := GPUInfo{
				Name:   name,
				Vendor: detectVendor(name),
			}

			// Parse adapter RAM
			if fields[1] != "" {
				var ramBytes int64
				//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
				_, _ = fmt.Sscanf(strings.TrimSpace(fields[1]), "%d", &ramBytes)
				if ramBytes > 0 {
					gpu.Memory = fmt.Sprintf("%.0f MB", float64(ramBytes)/1024/1024)
				}
			}

			// Driver version
			if len(fields) >= 4 && fields[3] != "" {
				gpu.Driver = strings.TrimSpace(fields[3])
			}

			gpus = append(gpus, gpu)
		}
	}

	return gpus
}

// detectVendor detects GPU vendor from the name
func detectVendor(name string) string {
	nameLower := strings.ToLower(name)

	// Check for specific vendor patterns (order matters - check more specific patterns first)
	if strings.Contains(nameLower, "nvidia") || strings.Contains(nameLower, "geforce") || strings.Contains(nameLower, "quadro") || strings.Contains(nameLower, "tesla") {
		return "NVIDIA"
	} else if strings.Contains(nameLower, "radeon") || (strings.Contains(nameLower, "amd") && !strings.Contains(nameLower, "amdahl")) {
		return "AMD"
	} else if strings.Contains(nameLower, "ati technologies") {
		return "AMD" // ATI Technologies is now part of AMD
	} else if matched, err := regexp.MatchString(`\bati\b`, nameLower); err == nil && matched {
		// Match ATI as a whole word to avoid false matches in words like "Corporation"
		return "AMD"
	} else if strings.Contains(nameLower, "intel") {
		return "Intel"
	} else if strings.Contains(nameLower, "apple") {
		return "Apple"
	} else if strings.Contains(nameLower, "microsoft") || strings.Contains(nameLower, "hyper-v") {
		return "Microsoft"
	}
	return "Unknown"
}

// printGPUInfo prints detailed GPU information
func printGPUInfo(gpus []GPUInfo) {
	if len(gpus) == 0 {
		fmt.Printf("  Unavailable\n")
		return
	}

	for i, gpu := range gpus {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("  Name:      %s\n", gpu.Name)
		if gpu.Vendor != "" && gpu.Vendor != "Unknown" {
			fmt.Printf("  Vendor:    %s\n", gpu.Vendor)
		}
		if gpu.Memory != "" {
			fmt.Printf("  Memory:    %s\n", gpu.Memory)
		}
		if gpu.Driver != "" {
			fmt.Printf("  Driver:    %s\n", gpu.Driver)
		}
		if gpu.ComputeCap != "" {
			fmt.Printf("  Compute:   %s\n", gpu.ComputeCap)
		}
		if gpu.ClockGraphics != "" {
			fmt.Printf("  Clock:     %s (graphics)\n", gpu.ClockGraphics)
		}
		if gpu.ClockMemory != "" {
			fmt.Printf("  Clock Mem: %s (memory)\n", gpu.ClockMemory)
		}
	}
}

// CPUDetails holds detailed CPU information
type CPUDetails struct {
	ModelName      string
	Architecture   string
	Cores          int
	ThreadsPerCore int
	CoresPerSocket int
	Sockets        int
	PhysicalCores  int
	LogicalCores   int
	BaseClock      string
	PCores         int
	ECores         int
	HasPECores     bool
}

// getDetailedCPUInfo gathers detailed CPU information from multiple sources
func getDetailedCPUInfo() CPUDetails {
	details := CPUDetails{
		ModelName:      "Unknown",
		Architecture:   runtime.GOARCH,
		LogicalCores:   runtime.NumCPU(),
		ThreadsPerCore: 1,
		CoresPerSocket: runtime.NumCPU(),
		Sockets:        1,
	}

	// Try to get CPU info from gopsutil
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		details.ModelName = cpuInfo[0].ModelName
		if cpuInfo[0].Mhz > 0 {
			details.BaseClock = fmt.Sprintf("%.2f GHz", cpuInfo[0].Mhz/1000.0)
		}
	}

	// On Linux, use lscpu for more detailed information
	if runtime.GOOS == "linux" {
		cmd := exec.Command("lscpu")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Architecture:") {
					details.Architecture = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				} else if strings.HasPrefix(line, "Thread(s) per core:") {
					//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
					_, _ = fmt.Sscanf(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "%d", &details.ThreadsPerCore)
				} else if strings.HasPrefix(line, "Core(s) per socket:") {
					//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
					_, _ = fmt.Sscanf(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "%d", &details.CoresPerSocket)
				} else if strings.HasPrefix(line, "Socket(s):") {
					//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
					_, _ = fmt.Sscanf(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "%d", &details.Sockets)
				} else if strings.HasPrefix(line, "CPU(s):") {
					//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
					_, _ = fmt.Sscanf(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "%d", &details.LogicalCores)
				} else if strings.HasPrefix(line, "CPU max MHz:") {
					var mhz float64
					//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
					_, _ = fmt.Sscanf(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "%f", &mhz)
					if mhz > 0 {
						details.BaseClock = fmt.Sprintf("%.2f GHz", mhz/1000.0)
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		// On macOS, use sysctl for detailed information
		cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
		if out, err := cmd.Output(); err == nil {
			details.ModelName = strings.TrimSpace(string(out))
		}

		// Get core counts
		cmd = exec.Command("sysctl", "-n", "hw.physicalcpu")
		if out, err := cmd.Output(); err == nil {
			//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
			_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &details.PhysicalCores)
		}

		cmd = exec.Command("sysctl", "-n", "hw.logicalcpu")
		if out, err := cmd.Output(); err == nil {
			//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
			_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &details.LogicalCores)
		}

		// Try to get P and E core counts (Apple Silicon)
		cmd = exec.Command("sysctl", "-n", "hw.perflevel0.physicalcpu")
		if out, err := cmd.Output(); err == nil {
			//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
			_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &details.PCores)
			details.HasPECores = true
		}

		cmd = exec.Command("sysctl", "-n", "hw.perflevel1.physicalcpu")
		if out, err := cmd.Output(); err == nil {
			//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
			_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &details.ECores)
			details.HasPECores = true
		}

		// Get base clock
		cmd = exec.Command("sysctl", "-n", "hw.cpufrequency")
		if out, err := cmd.Output(); err == nil {
			var hz int64
			//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
			_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &hz)
			if hz > 0 {
				details.BaseClock = fmt.Sprintf("%.2f GHz", float64(hz)/1e9)
			}
		}
	} else if runtime.GOOS == "windows" {
		// On Windows, use wmic
		cmd := exec.Command("wmic", "cpu", "get", "Name")
		if out, err := cmd.Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				details.ModelName = strings.TrimSpace(lines[1])
			}
		}

		cmd = exec.Command("wmic", "cpu", "get", "NumberOfCores")
		if out, err := cmd.Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
				_, _ = fmt.Sscanf(strings.TrimSpace(lines[1]), "%d", &details.PhysicalCores)
			}
		}

		cmd = exec.Command("wmic", "cpu", "get", "NumberOfLogicalProcessors")
		if out, err := cmd.Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				//nolint:errcheck // Sscanf errors are non-critical for best-effort parsing
				_, _ = fmt.Sscanf(strings.TrimSpace(lines[1]), "%d", &details.LogicalCores)
			}
		}
	}

	// Calculate physical cores if not set
	if details.PhysicalCores == 0 {
		details.PhysicalCores = details.CoresPerSocket * details.Sockets
	}

	// Ensure logical cores is at least physical cores
	if details.LogicalCores < details.PhysicalCores {
		details.LogicalCores = details.PhysicalCores
	}

	return details
}

// printCPUInfo prints detailed CPU information
func printCPUInfo(details CPUDetails) {
	fmt.Printf("  Model:     %s\n", details.ModelName)
	fmt.Printf("  Arch:      %s\n", details.Architecture)

	if details.BaseClock != "" {
		fmt.Printf("  Clock:     %s\n", details.BaseClock)
	}

	// Show physical vs logical cores
	if details.PhysicalCores > 0 && details.LogicalCores > 0 {
		fmt.Printf("  Cores:     %d physical, %d logical", details.PhysicalCores, details.LogicalCores)
		if details.ThreadsPerCore > 1 {
			fmt.Printf(" (%d threads/core)", details.ThreadsPerCore)
		}
		fmt.Println()
	} else {
		fmt.Printf("  Cores:     %d\n", details.LogicalCores)
	}

	// Show P/E core breakdown if available (Apple Silicon)
	if details.HasPECores && (details.PCores > 0 || details.ECores > 0) {
		fmt.Printf("  P-cores:   %d (performance)\n", details.PCores)
		fmt.Printf("  E-cores:   %d (efficiency)\n", details.ECores)
	}

	// Show socket/core organization if multiple sockets or meaningful
	if details.Sockets > 1 || (details.CoresPerSocket > 0 && details.CoresPerSocket != details.PhysicalCores) {
		fmt.Printf("  Layout:    %d socket(s), %d core(s) per socket\n", details.Sockets, details.CoresPerSocket)
	}
}

// detectTerminal tries to determine the terminal emulator in use
func detectTerminal() string {
	osType := runtime.GOOS

	// Windows-specific terminal detection
	if osType == "windows" {
		// Check for Windows Terminal
		if os.Getenv("WT_SESSION") != "" {
			return "Windows Terminal"
		}

		// Check for PowerShell (look for PSModulePath which is PowerShell-specific)
		if os.Getenv("PSModulePath") != "" {
			// Distinguish between PowerShell Core and Windows PowerShell
			if os.Getenv("POWERSHELL_DISTRIBUTION_CHANNEL") != "" {
				return "PowerShell Core"
			}
			return "PowerShell"
		}

		// Check for Git Bash (MSYSTEM is set by Git Bash)
		if msys := os.Getenv("MSYSTEM"); msys != "" {
			return "Git Bash"
		}

		// Check ConEmu
		if os.Getenv("ConEmuPID") != "" {
			return "ConEmu"
		}

		// Check for cmd.exe by examining COMSPEC
		if comspec := os.Getenv("COMSPEC"); comspec != "" && strings.Contains(strings.ToLower(comspec), "cmd.exe") {
			// If none of the above matched, we're likely in cmd.exe
			return "Command Prompt (cmd.exe)"
		}

		// Fallback for Windows
		return "Unknown Windows Terminal"
	}

	// Unix-like systems (Linux, macOS, etc.)
	// Try TERM_PROGRAM first (set by many modern terminals)
	if term := os.Getenv("TERM_PROGRAM"); term != "" {
		return term
	}

	// Check for specific terminal indicators
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return "kitty"
	}
	if os.Getenv("ALACRITTY_SOCKET") != "" || os.Getenv("ALACRITTY_LOG") != "" {
		return "alacritty"
	}
	if os.Getenv("KONSOLE_VERSION") != "" {
		return "konsole"
	}
	if os.Getenv("GNOME_TERMINAL_SERVICE") != "" {
		return "gnome-terminal"
	}

	// Check TERM variable as fallback
	if term := os.Getenv("TERM"); term != "" {
		return term
	}

	return "Unknown"
}

// printNetworkInfo prints network information for status command
func printNetworkInfo() {
	// Use the enhanced network details function from network.go
	PrintNetworkInfo()
}

// getRouterIP gets the default gateway IP (exported for network subcommand)
func getRouterIP() string {
	osType := runtime.GOOS

	switch osType {
	case "linux":
		cmd := exec.Command("sh", "-c", "ip route | grep default | awk '{print $3}' | head -n1")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
	case "darwin":
		cmd := exec.Command("route", "-n", "get", "default")
		out, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "gateway:") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						return parts[1]
					}
				}
			}
		}
	case "windows":
		cmd := exec.Command("ipconfig")
		out, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "Default Gateway") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						gateway := strings.TrimSpace(parts[1])
						if gateway != "" {
							return gateway
						}
					}
				}
			}
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
			versionStr := formatVersionWithUpdate(agent.Name, agent.Version)
			agentStrings = append(agentStrings, fmt.Sprintf("%s (%s)", agent.Name, versionStr))
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
	languageAvailable := []string{}
	runtimeAvailable := []string{}
	infrastructureAvailable := []string{}

	// System package managers
	switch osType {
	case "linux":
		if exists("apt-get") {
			version := getPackageManagerVersion("apt")
			if version != "" {
				versionStr := formatVersionWithUpdate("apt", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("apt (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "apt")
			}
		}
		if exists("flatpak") {
			version := getPackageManagerVersion("flatpak")
			if version != "" {
				versionStr := formatVersionWithUpdate("flatpak", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("flatpak (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "flatpak")
			}
		}
		if exists("snap") {
			version := getPackageManagerVersion("snap")
			if version != "" {
				versionStr := formatVersionWithUpdate("snap", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("snap (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "snap")
			}
		}
		if exists("dnf") {
			version := getPackageManagerVersion("dnf")
			if version != "" {
				versionStr := formatVersionWithUpdate("dnf", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("dnf (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "dnf")
			}
		}
		if exists("yum") {
			version := getPackageManagerVersion("yum")
			if version != "" {
				versionStr := formatVersionWithUpdate("yum", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("yum (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "yum")
			}
		}
		if exists("pacman") {
			version := getPackageManagerVersion("pacman")
			if version != "" {
				versionStr := formatVersionWithUpdate("pacman", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("pacman (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "pacman")
			}
		}
	case "darwin":
		if exists("brew") {
			version := getPackageManagerVersion("brew")
			if version != "" {
				versionStr := formatVersionWithUpdate("brew", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("homebrew (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "homebrew")
			}
		}
	case "windows":
		if exists("choco") {
			version := getPackageManagerVersion("choco")
			if version != "" {
				versionStr := formatVersionWithUpdate("choco", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("chocolatey (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "chocolatey")
			}
		}
		if exists("winget") {
			version := getPackageManagerVersion("winget")
			if version != "" {
				versionStr := formatVersionWithUpdate("winget", version)
				systemAvailable = append(systemAvailable, fmt.Sprintf("winget (%s)", versionStr))
			} else {
				systemAvailable = append(systemAvailable, "winget")
			}
		}
	}

	// Language version managers
	if checkNvmInstalled() {
		version := getVersionManagerVersion("nvm")
		if version != "" {
			versionStr := formatVersionWithUpdate("nvm", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("nvm (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "nvm")
		}
	}
	if exists("pyenv") {
		version := getVersionManagerVersion("pyenv")
		if version != "" {
			versionStr := formatVersionWithUpdate("pyenv", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("pyenv (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "pyenv")
		}
	}
	if exists("rbenv") {
		version := getVersionManagerVersion("rbenv")
		if version != "" {
			versionStr := formatVersionWithUpdate("rbenv", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("rbenv (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "rbenv")
		}
	}
	if exists("jenv") {
		version := getVersionManagerVersion("jenv")
		if version != "" {
			versionStr := formatVersionWithUpdate("jenv", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("jenv (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "jenv")
		}
	}
	if exists("rustup") {
		version := getVersionManagerVersion("rustup")
		if version != "" {
			versionStr := formatVersionWithUpdate("rustup", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("rustup (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "rustup")
		}
	}
	if exists("asdf") {
		version := getVersionManagerVersion("asdf")
		if version != "" {
			versionStr := formatVersionWithUpdate("asdf", version)
			languageAvailable = append(languageAvailable, fmt.Sprintf("asdf (%s)", versionStr))
		} else {
			languageAvailable = append(languageAvailable, "asdf")
		}
	}
	// Check for sdkman
	home, err := os.UserHomeDir()
	if err == nil {
		sdkmanInit := filepath.Join(home, ".sdkman", "bin", "sdkman-init.sh")
		if _, err := os.Stat(sdkmanInit); err == nil {
			version := getVersionManagerVersion("sdkman")
			if version != "" {
				versionStr := formatVersionWithUpdate("sdkman", version)
				languageAvailable = append(languageAvailable, fmt.Sprintf("sdkman (%s)", versionStr))
			} else {
				languageAvailable = append(languageAvailable, "sdkman")
			}
		}
	}

	// Programming runtime package managers
	if exists("npm") {
		version := getPackageManagerVersion("npm")
		if version != "" {
			versionStr := formatVersionWithUpdate("npm", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("npm (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "npm")
		}
	}
	if exists("pip") || exists("pip3") {
		version := getPackageManagerVersion("pip")
		if version != "" {
			versionStr := formatVersionWithUpdate("pip", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("pip (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "pip")
		}
	}
	if exists("pipx") {
		version := getPackageManagerVersion("pipx")
		if version != "" {
			versionStr := formatVersionWithUpdate("pipx", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("pipx (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "pipx")
		}
	}
	if exists("gem") {
		version := getPackageManagerVersion("gem")
		if version != "" {
			versionStr := formatVersionWithUpdate("gem", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("gem (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "gem")
		}
	}
	if exists("cargo") {
		version := getPackageManagerVersion("cargo")
		if version != "" {
			versionStr := formatVersionWithUpdate("cargo", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("cargo (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "cargo")
		}
	}
	if exists("go") {
		version := getPackageManagerVersion("go")
		if version != "" {
			versionStr := formatVersionWithUpdate("go", version)
			runtimeAvailable = append(runtimeAvailable, fmt.Sprintf("go (%s)", versionStr))
		} else {
			runtimeAvailable = append(runtimeAvailable, "go")
		}
	}

	// Infrastructure package managers
	if exists("VBoxManage") {
		version := getPackageManagerVersion("vboxmanage")
		if version != "" {
			versionStr := formatVersionWithUpdate("vboxmanage", version)
			infrastructureAvailable = append(infrastructureAvailable, fmt.Sprintf("VBoxManage (%s)", versionStr))
		} else {
			infrastructureAvailable = append(infrastructureAvailable, "VBoxManage")
		}
	}

	// Print
	if len(systemAvailable) > 0 {
		fmt.Printf("  System:         %s\n", strings.Join(systemAvailable, ", "))
	}
	if len(languageAvailable) > 0 {
		fmt.Printf("  Language:       %s\n", strings.Join(languageAvailable, ", "))
	}
	if len(runtimeAvailable) > 0 {
		fmt.Printf("  Runtime:        %s\n", strings.Join(runtimeAvailable, ", "))
	}
	if len(infrastructureAvailable) > 0 {
		fmt.Printf("  Infrastructure: %s\n", strings.Join(infrastructureAvailable, ", "))
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
	case "vboxmanage":
		cmd = exec.Command("VBoxManage", "--version")
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
	case "vboxmanage":
		// "7.0.14r161095" - just the version number
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

// getVersionManagerVersion returns the version of a language version manager
func getVersionManagerVersion(manager string) string {
	var cmd *exec.Cmd

	switch manager {
	case "nvm":
		cmd = exec.Command("bash", "-c", ". ~/.nvm/nvm.sh 2>/dev/null && nvm --version || echo ''")
	case "pyenv":
		cmd = exec.Command("pyenv", "--version")
	case "rbenv":
		cmd = exec.Command("rbenv", "--version")
	case "jenv":
		cmd = exec.Command("jenv", "--version")
	case "rustup":
		cmd = exec.Command("rustup", "--version")
	case "asdf":
		cmd = exec.Command("asdf", "--version")
	case "sdkman":
		cmd = exec.Command("bash", "-c", "source ~/.sdkman/bin/sdkman-init.sh 2>/dev/null && sdk version || echo ''")
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
	return extractVersionManagerVersion(manager, version)
}

// extractVersionManagerVersion extracts clean version from version manager output
func extractVersionManagerVersion(manager, output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	// Take first line only
	lines := strings.Split(output, "\n")
	firstLine := strings.TrimSpace(lines[0])

	switch manager {
	case "nvm":
		// "0.40.3" - just the version number
		return firstLine
	case "pyenv":
		// "pyenv 2.3.0" -> "2.3.0"
		if strings.HasPrefix(firstLine, "pyenv ") {
			return strings.TrimPrefix(firstLine, "pyenv ")
		}
		return firstLine
	case "rbenv":
		// "rbenv 1.2.0" -> "1.2.0"
		if strings.HasPrefix(firstLine, "rbenv ") {
			return strings.TrimPrefix(firstLine, "rbenv ")
		}
		return firstLine
	case "jenv":
		// "jenv 0.5.6" -> "0.5.6"
		if strings.HasPrefix(firstLine, "jenv ") {
			return strings.TrimPrefix(firstLine, "jenv ")
		}
		return firstLine
	case "rustup":
		// "rustup 1.26.0 (5af9b9484 2023-04-05)" -> "1.26.0"
		if strings.HasPrefix(firstLine, "rustup ") {
			parts := strings.Fields(firstLine)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
		return firstLine
	case "asdf":
		// "v0.14.0" -> "0.14.0"
		if strings.HasPrefix(firstLine, "v") {
			return strings.TrimPrefix(firstLine, "v")
		}
		return firstLine
	case "sdkman":
		// "SDKMAN 5.18.2" or "script: 5.18.2" -> "5.18.2"
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

	// Generic: try to extract version-like pattern
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

// getDiskInfo gets disk information
func getDiskInfo() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return ""
	}

	diskCount := 0
	totalSpace := uint64(0)
	seen := make(map[string]bool)

	for _, partition := range partitions {
		// Skip special filesystems
		if strings.HasPrefix(partition.Mountpoint, "/dev") ||
			strings.HasPrefix(partition.Mountpoint, "/sys") ||
			strings.HasPrefix(partition.Mountpoint, "/proc") ||
			strings.HasPrefix(partition.Mountpoint, "/run") {
			continue
		}

		// Get usage stats
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		// Track unique devices
		if !seen[partition.Device] {
			seen[partition.Device] = true
			diskCount++
		}

		totalSpace += usage.Total
	}

	if diskCount == 0 {
		return ""
	}

	totalGB := float64(totalSpace) / 1e9
	if diskCount == 1 {
		return fmt.Sprintf("%d disk (%.1f GB)", diskCount, totalGB)
	}
	return fmt.Sprintf("%d disks (%.1f GB total)", diskCount, totalGB)
}

// printPortsSummary prints a summary of listening ports
func printPortsSummary() {
	info := gatherPortsInfo()
	total := info.TCPPorts + info.UDPPorts

	if total > 0 {
		fmt.Printf("Ports:     %d listening (TCP: %d, UDP: %d)\n", total, info.TCPPorts, info.UDPPorts)
	}
}
