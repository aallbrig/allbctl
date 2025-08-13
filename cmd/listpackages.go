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
		if exists("dnf") {
			managers = append(managers, "dnf")
		}
		if exists("yum") {
			managers = append(managers, "yum")
		}
		if exists("zypper") {
			managers = append(managers, "zypper")
		}
		if exists("pacman") {
			managers = append(managers, "pacman")
		}
		if exists("rpm") {
			managers = append(managers, "rpm")
		}
		if exists("dpkg") {
			managers = append(managers, "dpkg")
		}
		if exists("apk") {
			managers = append(managers, "apk")
		}
		if exists("emerge") {
			managers = append(managers, "emerge")
		}
	case "darwin":
		if exists("brew") {
			managers = append(managers, "brew")
		}
		if exists("port") {
			managers = append(managers, "macports")
		}
		if exists("pkgin") {
			managers = append(managers, "pkgsrc")
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
		// Only return the count of installed packages
		return runCmd("bash -c \"dpkg-query -f '${binary:Package}\n' -W | wc -l\"")
	case "snap":
		return runCmd("snap list")
	case "flatpak":
		return runCmd("flatpak list")
	case "brew":
		return runCmd("brew list --formula --cask")
	case "choco":
		return runCmd("choco list --local-only")
	case "dnf":
		return runCmd("dnf list installed")
	case "yum":
		return runCmd("yum list installed")
	case "zypper":
		return runCmd("zypper se --installed-only")
	case "pacman":
		return runCmd("pacman -Q")
	case "rpm":
		return runCmd("rpm -qa")
	case "dpkg":
		return runCmd("dpkg --get-selections")
	case "apk":
		return runCmd("apk info")
	case "emerge":
		return runCmd("emerge --quiet --nocolor --pretend --tree @world")
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
