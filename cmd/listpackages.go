package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// ListPackagesCmd represents the list-packages command
var ListPackagesCmd = &cobra.Command{
	Use:   "list-packages",
	Short: "List installed packages from all detected package managers",
	Run: func(cmd *cobra.Command, args []string) {
		listInstalledPackages()
	},
}

func listInstalledPackages() {
	osType := runtime.GOOS
	fmt.Printf("Detected OS: %s\n", osType)

	var managers []string

	switch osType {
	case "linux":
		if exists("apt") {
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
	}

	if len(managers) == 0 {
		fmt.Println("No known package managers detected.")
		return
	}

	for _, m := range managers {
		fmt.Printf("\nPackages installed via %s:\n", m)
		pkgs := getPackages(m)
		fmt.Println(pkgs)
	}
}

func exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getPackages(manager string) string {
	switch manager {
	case "apt":
		return runCmd("apt list --installed")
	case "snap":
		return runCmd("snap list")
	case "flatpak":
		return runCmd("flatpak list")
	case "brew":
		return runCmd("brew list --formula --cask")
	case "choco":
		return runCmd("choco list --local-only")
	default:
		return "Unknown package manager."
	}
}

func runCmd(command string) string {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running %s: %v", command, err)
	}
	return string(output)
}
