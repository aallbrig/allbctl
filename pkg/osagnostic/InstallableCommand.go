package osagnostic

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

type InstallableCommand struct {
	CommandName     string
	LinuxPackages   map[string]string // package manager -> package name
	MacOSPackage    string            // homebrew package name
	WindowsPackages map[string]string // package manager -> package name
}

func NewInstallableCommand(commandName string) *InstallableCommand {
	return &InstallableCommand{
		CommandName:     commandName,
		LinuxPackages:   make(map[string]string),
		MacOSPackage:    "",
		WindowsPackages: make(map[string]string),
	}
}

func (i *InstallableCommand) SetLinuxPackage(packageManager, packageName string) *InstallableCommand {
	i.LinuxPackages[packageManager] = packageName
	return i
}

func (i *InstallableCommand) SetMacOSPackage(packageName string) *InstallableCommand {
	i.MacOSPackage = packageName
	return i
}

func (i *InstallableCommand) SetWindowsPackage(packageManager, packageName string) *InstallableCommand {
	i.WindowsPackages[packageManager] = packageName
	return i
}

func (i InstallableCommand) Name() string {
	return fmt.Sprintf("Installable Command: %s", i.CommandName)
}

func (i InstallableCommand) Validate() (out *bytes.Buffer, err error) {
	out = bytes.NewBufferString("")

	_, err = exec.LookPath(i.CommandName)
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(out, "NOT FOUND")
	} else {
		_, _ = color.New(color.FgGreen).Fprint(out, "INSTALLED")
	}
	out.WriteString(fmt.Sprintf(" %s", i.CommandName))

	return
}

func (i InstallableCommand) Install() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")

	// Check if already installed
	_, err := exec.LookPath(i.CommandName)
	if err == nil {
		_, _ = color.New(color.FgGreen).Fprint(out, "✅ Already installed")
		out.WriteString(fmt.Sprintf(" %s\n", i.CommandName))
		return out, nil
	}

	out.WriteString(fmt.Sprintf("Installing %s...\n", i.CommandName))

	switch runtime.GOOS {
	case "linux":
		return i.installOnLinux(out)
	case "darwin":
		return i.installOnMacOS(out)
	case "windows":
		return i.installOnWindows(out)
	default:
		out.WriteString(fmt.Sprintf("❌ Unsupported OS: %s\n", runtime.GOOS))
		return out, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (i InstallableCommand) installOnLinux(out *bytes.Buffer) (*bytes.Buffer, error) {
	// Try package managers in order of preference
	packageManagers := []string{"apt", "dnf", "yum", "pacman", "zypper", "apk"}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm); err == nil {
			if packageName, exists := i.LinuxPackages[pm]; exists {
				return i.installWithPackageManager(out, pm, packageName)
			}
		}
	}

	// Fallback to generic package name if available
	if genericPackage, exists := i.LinuxPackages["generic"]; exists {
		// Try common package managers with generic package name
		for _, pm := range packageManagers {
			if _, err := exec.LookPath(pm); err == nil {
				return i.installWithPackageManager(out, pm, genericPackage)
			}
		}
	}

	out.WriteString("❌ No supported package manager found for Linux\n")
	return out, fmt.Errorf("no supported package manager found")
}

func (i InstallableCommand) installOnMacOS(out *bytes.Buffer) (*bytes.Buffer, error) {
	if i.MacOSPackage == "" {
		out.WriteString("❌ No macOS package configured\n")
		return out, fmt.Errorf("no macOS package configured for %s", i.CommandName)
	}

	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		return i.installWithPackageManager(out, "brew install", i.MacOSPackage)
	}

	out.WriteString("❌ Homebrew not found. Please install Homebrew first: https://brew.sh/\n")
	return out, fmt.Errorf("homebrew not available")
}

func (i InstallableCommand) installOnWindows(out *bytes.Buffer) (*bytes.Buffer, error) {
	// Try package managers in order of preference
	packageManagers := []struct {
		command string
		lookup  string
	}{
		{"winget install", "winget"},
		{"choco install", "choco"},
		{"scoop install", "scoop"},
	}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm.lookup); err == nil {
			if packageName, exists := i.WindowsPackages[pm.lookup]; exists {
				return i.installWithPackageManager(out, pm.command, packageName)
			}
		}
	}

	out.WriteString("❌ No supported package manager found for Windows (winget, choco, scoop)\n")
	return out, fmt.Errorf("no supported package manager found")
}

func (i InstallableCommand) installWithPackageManager(out *bytes.Buffer, pmCommand, packageName string) (*bytes.Buffer, error) {
	// Add non-interactive flags for package managers that require them
	fullCommand := fmt.Sprintf("%s %s", pmCommand, packageName)

	// winget requires --accept-source-agreements to avoid interactive prompts
	if strings.HasPrefix(pmCommand, "winget") {
		fullCommand = fmt.Sprintf("%s --accept-source-agreements %s", pmCommand, packageName)
	}

	out.WriteString(fmt.Sprintf("Using: %s\n", fullCommand))

	// Split the command for exec
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", fullCommand)
	} else {
		cmd = exec.Command("sh", "-c", fullCommand)
	}

	output, err := cmd.CombinedOutput()
	out.WriteString(string(output))

	if err != nil {
		_, _ = color.New(color.FgRed).Fprintf(out, "❌ Failed to install %s\n", i.CommandName)
		return out, err
	}

	_, _ = color.New(color.FgGreen).Fprintf(out, "✅ Successfully installed %s\n", i.CommandName)
	return out, nil
}

func (i InstallableCommand) Uninstall() (*bytes.Buffer, error) {
	out := bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("❌ Cannot auto-uninstall %s - please uninstall manually\n", i.CommandName))
	return out, nil
}
