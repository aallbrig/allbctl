package status

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

// CheckForDirectory adds to output buffer based on if directory exists or not
func CheckForDirectory(w *bytes.Buffer, dir string, dirToFind string) (err error) {
	targetDirectory := filepath.Join(dir, dirToFind)
	w.WriteString(fmt.Sprintf("%-30s", targetDirectory))
	w.WriteString(" ")

	dirStat, err := os.Stat(targetDirectory)
	if os.IsNotExist(err) {
		color.New(color.FgRed).Fprint(w, "Missing")
	} else if dirStat.IsDir() {
		color.New(color.FgGreen).Fprint(w, "Present")
	} else {
		color.New(color.FgRed).Fprint(w, "Not Directory")
	}

	w.WriteString("\n")
	return
}

var goos = runtime.GOOS

// SystemInfo outputs basic system information
func SystemInfo(buf *bytes.Buffer) (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	buf.WriteString(fmt.Sprintf("Hostname:         %s\n", hostname))
	buf.WriteString("Operating System: ")

	switch goos {
	case "windows":
		buf.WriteString("Windows")
	case "darwin":
		buf.WriteString("MAC OS")
	case "linux":
		buf.WriteString("Linux")
	default:
		buf.WriteString(goos)
	}
	buf.WriteString("\n")

	return
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// PackageManager adds package manager information to output
func PackageManager(buf *bytes.Buffer) (err error) {
	var packageManagerCommand string
	buf.WriteString("Package Manager ")

	switch goos {
	case "windows":
		buf.WriteString("Chocolatey: ")
		packageManagerCommand = "choco"
	case "darwin":
		buf.WriteString("Homebrew: ")
		packageManagerCommand = "brew"
	case "linux":
		// Assume Ubuntu
		buf.WriteString("Apt-Get: ")
		packageManagerCommand = "apt-get"
	default:
		buf.WriteString(fmt.Sprintf("%s: ", goos))
	}

	if commandExists(packageManagerCommand) {
		_, err = color.New(color.FgGreen).Fprint(buf, "Present")
	} else {
		_, err = color.New(color.FgRed).Fprint(buf, "Missing")
	}

	buf.WriteString("\n")
	return
}
