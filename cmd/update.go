package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateDryRun   bool
	updateManagers []string
)

// packageManagerUpdate defines how to update a specific package manager
type packageManagerUpdate struct {
	Name        string     // e.g., "apt", "brew"
	NeedsSudo   bool       // whether commands need sudo prefix
	Commands    [][]string // each inner slice is one command to run
	Description string     // human-readable description
}

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update and upgrade all detected package managers",
	Long: `Update and upgrade packages from all detected package managers on the system.

Runs update/upgrade commands for each detected package manager sequentially.
Some managers require sudo and will prompt for your password.

Supported managers: apt, flatpak, snap, dnf, yum, pacman, brew, choco, winget, npm, pipx, gem

Intentionally skipped:
  pip   - Risky to auto-upgrade all pip packages (can break system Python)
  cargo - No global upgrade mechanism
  go    - No global upgrade mechanism

Use --dry-run to preview what commands would be executed.
Use --managers to limit which package managers are updated.

Examples:
  allbctl update                    # Update everything detected
  allbctl update --dry-run          # Preview what would happen
  allbctl update --managers apt,npm # Only update apt and npm`,
	Aliases: []string{"up", "upgrade"},
	Run: func(cmd *cobra.Command, args []string) {
		runUpdate()
	},
}

func init() {
	UpdateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Preview update commands without executing them")
	UpdateCmd.Flags().StringSliceVar(&updateManagers, "managers", nil, "Comma-separated list of package managers to update (default: all detected)")
}

// getUpdatableManagers returns the registry of all supported package manager update definitions
func getUpdatableManagers() []packageManagerUpdate {
	return []packageManagerUpdate{
		{Name: "apt", NeedsSudo: true, Commands: [][]string{
			{"apt-get", "update"},
			{"apt-get", "upgrade", "-y"},
		}, Description: "Update apt package lists and upgrade all packages"},

		{Name: "flatpak", NeedsSudo: false, Commands: [][]string{
			{"flatpak", "update", "-y"},
		}, Description: "Update all Flatpak applications"},

		{Name: "snap", NeedsSudo: true, Commands: [][]string{
			{"snap", "refresh"},
		}, Description: "Refresh all snap packages"},

		{Name: "dnf", NeedsSudo: true, Commands: [][]string{
			{"dnf", "upgrade", "-y"},
		}, Description: "Upgrade all dnf packages"},

		{Name: "yum", NeedsSudo: true, Commands: [][]string{
			{"yum", "update", "-y"},
		}, Description: "Update all yum packages"},

		{Name: "pacman", NeedsSudo: true, Commands: [][]string{
			{"pacman", "-Syu", "--noconfirm"},
		}, Description: "Synchronize and upgrade all pacman packages"},

		{Name: "brew", NeedsSudo: false, Commands: [][]string{
			{"brew", "update"},
			{"brew", "upgrade"},
		}, Description: "Update Homebrew and upgrade all formulae and casks"},

		{Name: "choco", NeedsSudo: false, Commands: [][]string{
			{"choco", "upgrade", "all", "-y"},
		}, Description: "Upgrade all Chocolatey packages"},

		{Name: "winget", NeedsSudo: false, Commands: [][]string{
			{"winget", "upgrade", "--all", "--accept-source-agreements", "--accept-package-agreements"},
		}, Description: "Upgrade all winget packages"},

		{Name: "npm", NeedsSudo: false, Commands: [][]string{
			{"npm", "update", "-g"},
		}, Description: "Update all globally installed npm packages"},

		{Name: "pipx", NeedsSudo: false, Commands: [][]string{
			{"pipx", "upgrade-all"},
		}, Description: "Upgrade all pipx-installed applications"},

		{Name: "gem", NeedsSudo: false, Commands: [][]string{
			{"gem", "update"},
		}, Description: "Update all installed Ruby gems"},
	}
}

// filterUpdatableManagers returns only the managers that are both detected on the system
// and present in the updatable registry, optionally filtered by the --managers flag
func filterUpdatableManagers() []packageManagerUpdate {
	detected := getDetectedPackageManagers()
	detectedSet := make(map[string]bool, len(detected))
	for _, m := range detected {
		detectedSet[m] = true
	}

	// Build set of requested managers if --managers flag is set
	var requestedSet map[string]bool
	if len(updateManagers) > 0 {
		requestedSet = make(map[string]bool, len(updateManagers))
		for _, m := range updateManagers {
			requestedSet[strings.TrimSpace(strings.ToLower(m))] = true
		}
	}

	var result []packageManagerUpdate
	for _, mgr := range getUpdatableManagers() {
		if !detectedSet[mgr.Name] {
			continue
		}
		if requestedSet != nil && !requestedSet[mgr.Name] {
			continue
		}
		result = append(result, mgr)
	}

	// Warn about requested managers that weren't found
	if requestedSet != nil {
		matched := make(map[string]bool)
		for _, mgr := range result {
			matched[mgr.Name] = true
		}
		for m := range requestedSet {
			if !matched[m] {
				if !detectedSet[m] {
					fmt.Printf("Skipping %s: not detected on this system\n", m)
				} else {
					fmt.Printf("Skipping %s: not a supported updatable manager\n", m)
				}
			}
		}
	}

	return result
}

// runUpdateCommand executes a single update command, optionally with sudo
func runUpdateCommand(args []string, needsSudo bool) error {
	if needsSudo {
		if !exists("sudo") {
			return fmt.Errorf("sudo is required but not found on PATH")
		}
		args = append([]string{"sudo"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // args are from a hardcoded registry, not user input
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runUpdate() {
	managers := filterUpdatableManagers()

	if len(managers) == 0 {
		fmt.Println("No updatable package managers detected on this system.")
		return
	}

	// Print summary
	fmt.Println("Package managers to update:")
	for _, mgr := range managers {
		updateCount, _ := checkPackageUpdates(mgr.Name) //nolint:errcheck
		if updateCount > 0 {
			fmt.Printf("  %-12s %s (%d updates available)\n", mgr.Name+":", mgr.Description, updateCount)
		} else {
			fmt.Printf("  %-12s %s\n", mgr.Name+":", mgr.Description)
		}
	}
	fmt.Println()

	// Dry run: show commands and stop
	if updateDryRun {
		fmt.Println("Dry run — commands that would be executed:")
		fmt.Println()
		for _, mgr := range managers {
			fmt.Printf("  # %s\n", mgr.Description)
			for _, cmdArgs := range mgr.Commands {
				if mgr.NeedsSudo {
					fmt.Printf("  sudo %s\n", strings.Join(cmdArgs, " "))
				} else {
					fmt.Printf("  %s\n", strings.Join(cmdArgs, " "))
				}
			}
			fmt.Println()
		}
		return
	}

	// Execute updates
	var succeeded, failed []string

	for _, mgr := range managers {
		fmt.Printf("==> Updating %s...\n", mgr.Name)
		mgrFailed := false

		for _, cmdArgs := range mgr.Commands {
			if mgr.NeedsSudo {
				fmt.Printf("  Running: sudo %s\n", strings.Join(cmdArgs, " "))
			} else {
				fmt.Printf("  Running: %s\n", strings.Join(cmdArgs, " "))
			}

			if err := runUpdateCommand(cmdArgs, mgr.NeedsSudo); err != nil {
				fmt.Printf("  Error: %v\n", err)
				mgrFailed = true
				break
			}
		}

		if mgrFailed {
			failed = append(failed, mgr.Name)
		} else {
			succeeded = append(succeeded, mgr.Name)
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println("---")
	if len(succeeded) > 0 {
		fmt.Printf("Updated successfully: %s\n", strings.Join(succeeded, ", "))
	}
	if len(failed) > 0 {
		fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
	}
}
