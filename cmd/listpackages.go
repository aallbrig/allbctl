package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var detailFlag bool

// ListPackagesCmd represents the list-packages command
var ListPackagesCmd = &cobra.Command{
	Use:   "list-packages [package-manager]",
	Short: "Show package count summary from all detected package managers",
	Long: `Show package count summary from all detected package managers.

By default, shows the same summary as the 'Packages:' section in 'allbctl status'.
	
For system package managers (apt, yum, brew, chocolatey, etc.), only explicitly 
installed packages are counted, not their dependencies.

For programming runtime package managers (npm, pip, gem, etc.), only globally 
installed packages are counted.

If a package manager is not detected on the system, it will not be displayed.

Use --detail flag to see the full list of all installed packages.

You can also specify a specific package manager to list only its packages:
  allbctl list-packages apt
  allbctl list-packages npm
  allbctl list-packages flatpak`,
	Run: func(cmd *cobra.Command, args []string) {
		listInstalledPackages(args)
	},
}

func init() {
	ListPackagesCmd.Flags().BoolVarP(&detailFlag, "detail", "d", false, "Show detailed list of all packages instead of just counts")
}

// getDetectedPackageManagers returns a list of all detected package managers
func getDetectedPackageManagers() []string {
	osType := runtime.GOOS

	var managers []string

	// System package managers
	switch osType {
	case "linux":
		if exists("dpkg") {
			managers = append(managers, "dpkg")
		}
		if exists("rpm") {
			managers = append(managers, "rpm")
		}
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
	if exists("ollama") {
		managers = append(managers, "ollama")
	}
	if exists("vagrant") {
		managers = append(managers, "vagrant")
	}
	if exists("VBoxManage") {
		managers = append(managers, "vboxmanage")
	}

	return managers
}

// PackageResult holds package count results for a single package manager
type PackageResult struct {
	Manager     string
	Count       int
	UpdateCount int
	Index       int
}

// PackageSummaryFuture represents an ongoing package detection operation
type PackageSummaryFuture struct {
	resultChan chan PackageResult
	managers   []string
}

// StartPackageSummary initiates package detection in the background
// Returns a future that can be used to retrieve results later
func StartPackageSummary() *PackageSummaryFuture {
	managers := getDetectedPackageManagers()

	if len(managers) == 0 {
		return nil
	}

	resultChan := make(chan PackageResult, len(managers))

	// Launch goroutines to count packages in parallel
	for i, m := range managers {
		go func(manager string, idx int) {
			pkgs := getPackages(manager)
			var count, updateCount int
			if pkgs != "" {
				count = countPackages(manager, pkgs)
				updateCount, _ = checkPackageUpdates(manager) //nolint:errcheck
			}
			resultChan <- PackageResult{
				Manager:     manager,
				Count:       count,
				UpdateCount: updateCount,
				Index:       idx,
			}
		}(m, i)
	}

	return &PackageSummaryFuture{
		resultChan: resultChan,
		managers:   managers,
	}
}

// PrintResults waits for all package detection to complete and prints the results
func (f *PackageSummaryFuture) PrintResults() {
	if f == nil {
		fmt.Println("  No package managers detected")
		return
	}

	// Collect all results
	resultMap := make(map[string]PackageResult)
	for i := 0; i < len(f.managers); i++ {
		result := <-f.resultChan
		resultMap[result.Manager] = result
	}
	close(f.resultChan)

	// Print in original order with results
	for _, m := range f.managers {
		result := resultMap[m]
		if result.Count > 0 {
			var output string
			if m == "ollama" {
				if result.UpdateCount > 0 {
					output = fmt.Sprintf("  %-15s %d models (%d want updates)\n", m+":", result.Count, result.UpdateCount)
				} else {
					output = fmt.Sprintf("  %-15s %d models\n", m+":", result.Count)
				}
			} else if m == "vagrant" || m == "vboxmanage" {
				output = fmt.Sprintf("  %-15s %d VMs\n", m+":", result.Count)
			} else {
				if result.UpdateCount > 0 {
					output = fmt.Sprintf("  %-15s %d packages (%d want updates)\n", m+":", result.Count, result.UpdateCount)
				} else {
					output = fmt.Sprintf("  %-15s %d packages\n", m+":", result.Count)
				}
			}
			fmt.Print(output)
		}
	}
}

// PrintPackageSummary prints package counts for all detected package managers (for status command)
// This is the synchronous version for backward compatibility
func PrintPackageSummary() {
	future := StartPackageSummary()
	if future == nil {
		fmt.Println("  No package managers detected")
		return
	}
	future.PrintResults()
}

func listInstalledPackages(args []string) {
	// If a specific package manager is requested
	if len(args) > 0 {
		manager := args[0]
		if !exists(getCommandForManager(manager)) {
			fmt.Printf("Package manager '%s' not found on this system.\n", manager)
			return
		}
		pkgs := getPackages(manager)
		if pkgs != "" {
			fmt.Printf("Packages installed via %s:\n", manager)
			fmt.Println(pkgs)
			fmt.Printf("\nCommand: %s\n", getQueryCommand(manager))

			// Show recent installations if --detail flag is used
			if detailFlag {
				recent := getRecentPackages(manager)
				if recent != "" {
					fmt.Printf("\nLast 5 installed packages:\n")
					fmt.Print(recent)
					recentCmd := getRecentPackagesCommand(manager)
					if recentCmd != "" {
						fmt.Printf("\nCommand: %s\n", recentCmd)
					}
				}
			}
		} else {
			fmt.Printf("No packages found for %s\n", manager)
			fmt.Printf("\nCommand: %s\n", getQueryCommand(manager))
		}
		return
	}

	// Otherwise, list all detected package managers
	managers := getDetectedPackageManagers()

	if len(managers) == 0 {
		fmt.Println("No known package managers detected.")
		return
	}

	if detailFlag {
		// Detail mode: show full listing
		for _, m := range managers {
			pkgs := getPackages(m)
			if pkgs != "" {
				fmt.Printf("Packages installed via %s:\n", m)
				fmt.Println(pkgs)
				fmt.Println()
			}
		}
	} else {
		// Summary mode (default): just count packages (no indentation for direct command)
		for _, m := range managers {
			pkgs := getPackages(m)
			if pkgs != "" {
				count := countPackages(m, pkgs)
				if m == "ollama" {
					fmt.Printf("%-15s %d models\n", m+":", count)
				} else if m == "vagrant" || m == "vboxmanage" {
					fmt.Printf("%-15s %d VMs\n", m+":", count)
				} else {
					fmt.Printf("%-15s %d packages\n", m+":", count)
				}
			}
		}
		fmt.Println("\nUse --detail flag to see the full list of all installed packages.")
		fmt.Println("Or specify a package manager: allbctl status list-packages <manager>")
	}
}

func getCommandForManager(manager string) string {
	switch manager {
	case "apt":
		return "apt-mark"
	case "pip":
		if exists("pip3") {
			return "pip3"
		}
		return "pip"
	case "vboxmanage":
		return "VBoxManage"
	default:
		return manager
	}
}

func getQueryCommand(manager string) string {
	switch manager {
	case "dpkg":
		return "dpkg --get-selections"
	case "rpm":
		return "rpm -qa"
	case "apt":
		return "apt-mark showmanual"
	case "snap":
		return "snap list --color=never"
	case "flatpak":
		return "flatpak list --app --columns=name,application"
	case "brew":
		return "brew leaves && brew list --cask"
	case "choco":
		return "choco list"
	case "dnf":
		return "dnf repoquery --userinstalled --qf '%{name}'"
	case "yum":
		return "yum history userinstalled"
	case "pacman":
		return "pacman -Qe"
	case "winget":
		return "winget list"
	case "scoop":
		return "scoop list"
	case "npm":
		return "npm list -g --depth=0"
	case "pip":
		cmd := "pip3"
		if !exists("pip3") {
			cmd = "pip"
		}
		return cmd + " list --format=columns"
	case "pipx":
		return "pipx list"
	case "gem":
		return "gem list --local"
	case "cargo":
		return "cargo install --list"
	case "go":
		return "ls -1 $(go env GOPATH)/bin"
	case "ollama":
		return "ollama list"
	case "vagrant":
		return "vagrant box list"
	case "vboxmanage":
		return "VBoxManage list vms"
	default:
		return ""
	}
}

func exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getPackages(manager string) string {
	var output string
	switch manager {
	case "dpkg":
		output = runCmd("dpkg --get-selections")
	case "rpm":
		output = runCmd("rpm -qa")
	case "apt":
		// List only manually installed packages (not auto-installed dependencies)
		output = runCmd("apt-mark showmanual")
	case "snap":
		// Snap doesn't track dependencies separately, list all
		output = runCmd("snap list --color=never")
	case "flatpak":
		// List user-installed apps (columns: name, app-id, version, branch, origin)
		output = runCmd("flatpak list --app --columns=name,application")
	case "brew":
		// List only top-level formulae and casks (explicitly installed, not dependencies)
		output = runCmd("brew leaves") + "\n" + runCmd("brew list --cask")
	case "choco":
		// List only explicitly installed packages
		output = runCmd("choco list")
	case "dnf":
		// List user-installed packages (not dependencies)
		output = runCmd("dnf repoquery --userinstalled --qf '%{name}'")
	case "yum":
		// List user-installed packages
		output = runCmd("yum history userinstalled")
	case "pacman":
		// List explicitly installed packages (not dependencies)
		output = runCmd("pacman -Qe")
	case "winget":
		// List installed packages
		output = runCmd("winget list")
	case "scoop":
		// List installed packages
		output = runCmd("scoop list")
	case "npm":
		// List globally installed packages (depth 0 = no dependencies)
		output = runCmd("npm list -g --depth=0")
	case "pip":
		// List globally installed packages
		cmd := "pip3"
		if !exists("pip3") {
			cmd = "pip"
		}
		output = runCmd(cmd + " list --format=columns")
	case "pipx":
		// List packages installed via pipx
		output = runCmd("pipx list")
	case "gem":
		// List globally installed gems (no dependencies shown by default)
		output = runCmd("gem list --local")
	case "cargo":
		// List globally installed cargo binaries
		output = runCmd("cargo install --list")
	case "go":
		// Go doesn't have a traditional global install list
		// List binaries in GOPATH/bin or GOBIN
		output = runCmd("bash -c 'ls -1 $(go env GOPATH)/bin 2>/dev/null || echo \"No Go binaries found\"'")
	case "ollama":
		// List ollama models
		output = runCmd("ollama list")
	case "vagrant":
		// List vagrant boxes
		output = runCmd("vagrant box list")
	case "vboxmanage":
		// List VirtualBox VMs
		output = runCmd("VBoxManage list vms")
	default:
		return ""
	}
	return strings.TrimSpace(output)
}

func runCmd(command string) string {
	// Add non-interactive flags for commands that require them
	if strings.HasPrefix(command, "winget ") {
		// winget requires --accept-source-agreements to avoid interactive prompts
		command = strings.Replace(command, "winget ", "winget --accept-source-agreements ", 1)
	}

	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running %s: %v", command, err)
	}
	return string(output)
}

func getRecentPackages(manager string) string {
	switch manager {
	case "apt", "dpkg":
		// Read dpkg logs to find recent installations
		// Use sh -c to properly handle pipes
		cmd := exec.Command("sh", "-c", "cat /var/log/dpkg.log /var/log/dpkg.log.1 2>/dev/null | grep ' install ' | tail -5")
		output, err := cmd.CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) == "" {
			return ""
		}

		logs := string(output)
		var result strings.Builder
		lines := strings.Split(strings.TrimSpace(logs), "\n")
		for _, line := range lines {
			// Format: "2026-01-06 17:08:49 install sqlite3:amd64 <none> 3.45.1-1ubuntu2.5"
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				date := parts[0]
				time := parts[1]
				pkgName := parts[3]
				// Remove architecture suffix
				if idx := strings.Index(pkgName, ":"); idx > 0 {
					pkgName = pkgName[:idx]
				}
				result.WriteString(fmt.Sprintf("  %s %s - %s\n", date, time, pkgName))
			}
		}
		return result.String()

	case "npm":
		// Use filesystem timestamps for npm packages
		cmd := exec.Command("sh", "-c", "ls -lt $(npm root -g 2>/dev/null) 2>/dev/null | grep '^d' | head -5")
		output, err := cmd.CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) == "" {
			return ""
		}

		var result strings.Builder
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			// Format: "drwxrwxr-x 4 user user 4096 Dec 31 13:54 puppeteer-mcp-server"
			parts := strings.Fields(line)
			if len(parts) >= 9 {
				month := parts[5]
				day := parts[6]
				timeOrYear := parts[7]
				pkgName := parts[8]
				result.WriteString(fmt.Sprintf("  %s %s %s - %s\n", month, day, timeOrYear, pkgName))
			}
		}
		return result.String()

	case "pip":
		// Use Python to get package timestamps
		// Suppress pkg_resources deprecation warning with -W ignore
		cmd := exec.Command("python3", "-W", "ignore::DeprecationWarning", "-c", "import pkg_resources; import os; pkgs = [(p.project_name, p.version, os.stat(p.location).st_mtime) for p in pkg_resources.working_set]; import time; pkgs_sorted = sorted(pkgs, key=lambda x: x[2], reverse=True); [print(f'{time.strftime(\"%Y-%m-%d %H:%M:%S\", time.localtime(p[2]))} - {p[0]} ({p[1]})') for p in pkgs_sorted[:5]]")
		output, err := cmd.CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) == "" {
			return ""
		}

		var result strings.Builder
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			result.WriteString("  " + line + "\n")
		}
		return result.String()

	default:
		return ""
	}
}

func getRecentPackagesCommand(manager string) string {
	switch manager {
	case "apt", "dpkg":
		return "cat /var/log/dpkg.log /var/log/dpkg.log.1 2>/dev/null | grep ' install ' | tail -5"
	case "npm":
		return "ls -lt $(npm root -g) 2>/dev/null | grep '^d' | head -5"
	case "pip":
		return "python3 -c \"import pkg_resources; import os; pkgs = [(p.project_name, p.version, os.stat(p.location).st_mtime) for p in pkg_resources.working_set]; import time; pkgs_sorted = sorted(pkgs, key=lambda x: x[2], reverse=True); [print(f'{time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(p[2]))} - {p[0]} ({p[1]})') for p in pkgs_sorted[:5]]\""
	default:
		return ""
	}
}

func countPackages(manager string, output string) int {
	if strings.HasPrefix(output, "Error running") {
		return 0
	}

	output = strings.TrimSpace(output)
	if output == "" {
		return 0
	}

	lines := strings.Split(output, "\n")

	switch manager {
	case "dpkg":
		// dpkg --get-selections output: "package-name	install"
		count := 0
		for _, line := range lines {
			if strings.Contains(line, "install") {
				count++
			}
		}
		return count
	case "rpm":
		// Simple line count for rpm -qa
		return len(lines)
	case "apt", "cargo", "go":
		// Simple line count (one package per line)
		return len(lines)
	case "snap", "flatpak", "brew", "dnf", "yum", "pacman", "winget", "scoop", "gem":
		// Skip header lines and count data lines
		count := 0
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Skip header lines (varies by manager)
			if manager == "snap" && i == 0 {
				continue // Skip "Name Version Rev Tracking Publisher Notes"
			}
			if manager == "flatpak" && strings.HasPrefix(line, "Name") {
				continue
			}
			if (manager == "winget" || manager == "scoop") && i < 2 {
				continue // Skip headers and separator
			}
			count++
		}
		return count
	case "npm":
		// npm output has a tree structure, count top-level packages
		// Format: /path/to/lib\n├── package@version\n└── package@version
		count := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "├──") || strings.HasPrefix(line, "└──") {
				count++
			}
		}
		return count
	case "pip":
		// pip list has 2 header lines
		if len(lines) <= 2 {
			return 0
		}
		return len(lines) - 2
	case "pipx":
		// pipx list output format: "   package    pkg-name"
		count := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "venvs are in") || strings.HasPrefix(line, "apps are exposed") {
				continue
			}
			// Count lines that start with "package" (each installed app)
			if strings.HasPrefix(line, "package ") {
				count++
			}
		}
		return count
	case "choco":
		// choco list has header and footer, count packages in between
		count := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.Contains(line, "packages installed") {
				continue
			}
			// Package lines typically have a version number
			if strings.Contains(line, " ") && !strings.HasSuffix(line, ":") {
				count++
			}
		}
		return count
	case "ollama":
		// ollama list has 1 header line: "NAME    ID    SIZE    MODIFIED"
		// Count data lines (skip header)
		count := 0
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if i == 0 || line == "" {
				continue // Skip header and empty lines
			}
			count++
		}
		return count
	case "vagrant":
		// vagrant box list output: "box-name (provider, version)"
		// Each line represents one box
		count := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			count++
		}
		return count
	case "vboxmanage":
		// VBoxManage list vms output: "VM-name" {uuid}
		// Each line represents one VM
		count := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			count++
		}
		return count
	default:
		return len(lines)
	}
}
