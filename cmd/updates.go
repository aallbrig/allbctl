package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// UpdateInfo holds information about available updates
type UpdateInfo struct {
	Current   string
	Available string
	IsLTS     bool
}

// Version comparison and update checking cache
var versionCache = make(map[string]*UpdateInfo)

// checkPackageUpdates checks if packages have available updates for a given package manager
func checkPackageUpdates(manager string) (int, error) {
	switch manager {
	case "apt":
		return checkAptUpdates()
	case "flatpak":
		return checkFlatpakUpdates()
	case "snap":
		return checkSnapUpdates()
	case "dnf":
		return checkDnfUpdates()
	case "yum":
		return checkYumUpdates()
	case "pacman":
		return checkPacmanUpdates()
	case "brew":
		return checkBrewUpdates()
	case "choco":
		return checkChocoUpdates()
	case "winget":
		return checkWingetUpdates()
	case "npm":
		return checkNpmUpdates()
	case "pip":
		return checkPipUpdates()
	case "pipx":
		return checkPipxUpdates()
	default:
		return 0, nil
	}
}

func checkAptUpdates() (int, error) {
	// Don't run apt-get update as it requires sudo
	// Instead, check the existing cache
	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return 0, nil // Silently fail if can't check
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// First line is "Listing...", so subtract 1
	count := len(lines) - 1
	if count < 0 {
		count = 0
	}
	return count, nil
}

func checkFlatpakUpdates() (int, error) {
	cmd := exec.Command("flatpak", "remote-ls", "--updates", "--app")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}
	return len(lines), nil
}

func checkSnapUpdates() (int, error) {
	cmd := exec.Command("snap", "refresh", "--list")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Skip header line
	count := len(lines) - 1
	if count < 0 {
		count = 0
	}
	return count, nil
}

func checkDnfUpdates() (int, error) {
	cmd := exec.Command("dnf", "check-update", "-q")
	output, err := cmd.Output()
	// dnf returns exit code 100 if updates are available
	if err != nil && !strings.Contains(err.Error(), "exit status 100") {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return len(lines), nil
}

func checkYumUpdates() (int, error) {
	cmd := exec.Command("yum", "check-update", "-q")
	output, err := cmd.Output()
	// yum returns exit code 100 if updates are available
	if err != nil && !strings.Contains(err.Error(), "exit status 100") {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return len(lines), nil
}

func checkPacmanUpdates() (int, error) {
	cmd := exec.Command("checkupdates")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}
	return len(lines), nil
}

func checkBrewUpdates() (int, error) {
	cmd := exec.Command("brew", "outdated")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}
	return len(lines), nil
}

func checkChocoUpdates() (int, error) {
	cmd := exec.Command("choco", "outdated", "-r")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}
	return len(lines), nil
}

func checkWingetUpdates() (int, error) {
	cmd := exec.Command("winget", "upgrade", "--accept-source-agreements")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "Available") {
			count++
		}
	}
	return count, nil
}

func checkNpmUpdates() (int, error) {
	cmd := exec.Command("npm", "outdated", "-g", "--json")
	output, err := cmd.Output()
	if err != nil && len(output) == 0 {
		return 0, err
	}

	var outdated map[string]interface{}
	if err := json.Unmarshal(output, &outdated); err != nil {
		return 0, nil
	}

	return len(outdated), nil
}

func checkPipUpdates() (int, error) {
	pipCmd := "pip3"
	if !exists("pip3") {
		pipCmd = "pip"
	}

	cmd := exec.Command(pipCmd, "list", "--outdated", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var outdated []interface{}
	if err := json.Unmarshal(output, &outdated); err != nil {
		return 0, nil
	}

	return len(outdated), nil
}

func checkPipxUpdates() (int, error) {
	// pipx doesn't have a built-in outdated check
	// We would need to check each package individually
	return 0, nil
}

// checkVersionUpdate checks if a newer version is available for a tool
func checkVersionUpdate(name, current string) *UpdateInfo {
	// Check cache first
	cacheKey := strings.ToLower(name) + ":" + current
	if cached, ok := versionCache[cacheKey]; ok {
		return cached
	}

	var update *UpdateInfo

	switch strings.ToLower(name) {
	case "node.js", "node":
		update = checkNodeJSUpdate(current)
	case "python":
		update = checkPythonUpdate(current)
	case "go":
		update = checkGoUpdate(current)
	case "java":
		update = checkJavaUpdate(current)
	case "ruby":
		update = checkRubyUpdate(current)
	case "npm":
		update = checkNPMUpdate(current)
	case "pip":
		update = checkPipVersionUpdate(current)
	case "copilot":
		update = checkCopilotUpdate(current)
	case "ollama":
		update = checkOllamaUpdate(current)
	default:
		// Generic GitHub release checker for other tools
		update = checkGitHubRelease(name, current)
	}

	// Cache the result
	if update != nil {
		versionCache[cacheKey] = update
	}

	return update
}

// checkNodeJSUpdate checks for Node.js LTS updates
func checkNodeJSUpdate(current string) *UpdateInfo {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://nodejs.org/dist/index.json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var releases []struct {
		Version string      `json:"version"`
		LTS     interface{} `json:"lts"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		return nil
	}

	// Find latest LTS version
	for _, release := range releases {
		if release.LTS != nil && release.LTS != false {
			latest := strings.TrimPrefix(release.Version, "v")
			if compareVersions(latest, current) > 0 {
				return &UpdateInfo{
					Current:   current,
					Available: latest,
					IsLTS:     true,
				}
			}
			break
		}
	}

	return nil
}

// checkPythonUpdate checks for Python updates
func checkPythonUpdate(current string) *UpdateInfo {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://endoflife.date/api/python.json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var releases []struct {
		Cycle  string `json:"cycle"`
		Latest string `json:"latest"`
		LTS    bool   `json:"lts"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		return nil
	}

	// Find latest stable version
	for _, release := range releases {
		if release.Latest != "" {
			if compareVersions(release.Latest, current) > 0 {
				return &UpdateInfo{
					Current:   current,
					Available: release.Latest,
					IsLTS:     release.LTS,
				}
			}
			break
		}
	}

	return nil
}

// checkGoUpdate checks for Go updates
func checkGoUpdate(current string) *UpdateInfo {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var releases []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		return nil
	}

	// Find latest stable version
	for _, release := range releases {
		if release.Stable {
			latest := strings.TrimPrefix(release.Version, "go")
			if compareVersions(latest, current) > 0 {
				return &UpdateInfo{
					Current:   current,
					Available: latest,
					IsLTS:     false,
				}
			}
			break
		}
	}

	return nil
}

// checkJavaUpdate checks for Java LTS updates
func checkJavaUpdate(current string) *UpdateInfo {
	// Java LTS versions: 8, 11, 17, 21
	ltsVersions := []string{"21", "17", "11", "8"}

	currentMajor := strings.Split(current, ".")[0]

	for _, lts := range ltsVersions {
		if compareVersions(lts, currentMajor) > 0 {
			return &UpdateInfo{
				Current:   current,
				Available: lts,
				IsLTS:     true,
			}
		}
	}

	return nil
}

// checkRubyUpdate checks for Ruby updates
func checkRubyUpdate(current string) *UpdateInfo {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://endoflife.date/api/ruby.json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var releases []struct {
		Cycle  string `json:"cycle"`
		Latest string `json:"latest"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		return nil
	}

	// Find latest version
	for _, release := range releases {
		if release.Latest != "" {
			if compareVersions(release.Latest, current) > 0 {
				return &UpdateInfo{
					Current:   current,
					Available: release.Latest,
					IsLTS:     false,
				}
			}
			break
		}
	}

	return nil
}

// checkNPMUpdate checks for npm updates
func checkNPMUpdate(current string) *UpdateInfo {
	return checkGitHubRelease("npm/cli", current)
}

// checkPipVersionUpdate checks for pip updates
func checkPipVersionUpdate(current string) *UpdateInfo {
	return checkGitHubRelease("pypa/pip", current)
}

// checkCopilotUpdate checks for GitHub Copilot CLI updates
func checkCopilotUpdate(current string) *UpdateInfo {
	// GitHub Copilot CLI uses a different versioning scheme
	// For now, return nil as it auto-updates
	return nil
}

// checkOllamaUpdate checks for Ollama updates
func checkOllamaUpdate(current string) *UpdateInfo {
	return checkGitHubRelease("ollama/ollama", current)
}

// checkGitHubRelease checks GitHub releases for a given repo
func checkGitHubRelease(repo, current string) *UpdateInfo {
	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	// Set User-Agent to avoid rate limiting
	req.Header.Set("User-Agent", "allbctl")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		return nil
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	if compareVersions(latest, current) > 0 {
		return &UpdateInfo{
			Current:   current,
			Available: latest,
			IsLTS:     false,
		}
	}

	return nil
}

// compareVersions compares two semantic version strings
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		var v1Part, v2Part int

		if i < len(v1Parts) {
			v1Part, _ = strconv.Atoi(strings.TrimSpace(v1Parts[i]))
		}

		if i < len(v2Parts) {
			v2Part, _ = strconv.Atoi(strings.TrimSpace(v2Parts[i]))
		}

		if v1Part > v2Part {
			return 1
		} else if v1Part < v2Part {
			return -1
		}
	}

	return 0
}

// formatVersionWithUpdate formats version string with update arrow if available
func formatVersionWithUpdate(name, current string) string {
	update := checkVersionUpdate(name, current)
	if update != nil && update.Available != "" && update.Available != current {
		if update.IsLTS {
			return fmt.Sprintf("%s → %s (LTS)", current, update.Available)
		}
		return fmt.Sprintf("%s → %s", current, update.Available)
	}
	return current
}

// checkOSUpdates checks for operating system updates
func checkOSUpdates() (int, error) {
	osType := runtime.GOOS

	switch osType {
	case "linux":
		return checkLinuxUpdates()
	case "windows":
		return checkWindowsUpdates()
	case "darwin":
		return checkMacOSUpdates()
	default:
		return 0, nil
	}
}

func checkLinuxUpdates() (int, error) {
	// Try to detect the distribution
	cmd := exec.Command("lsb_release", "-is")
	distOutput, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	distro := strings.ToLower(strings.TrimSpace(string(distOutput)))

	switch {
	case strings.Contains(distro, "ubuntu"), strings.Contains(distro, "mint"):
		return checkAptUpdates()
	case strings.Contains(distro, "fedora"):
		return checkDnfUpdates()
	case strings.Contains(distro, "arch"):
		return checkPacmanUpdates()
	default:
		return 0, nil
	}
}

func checkWindowsUpdates() (int, error) {
	// Use PowerShell to check for Windows updates
	cmd := exec.Command("powershell", "-Command",
		"Get-WindowsUpdate -AcceptAll -IgnoreReboot | Measure-Object | Select-Object -ExpandProperty Count")
	output, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func checkMacOSUpdates() (int, error) {
	cmd := exec.Command("softwareupdate", "-l")
	output, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "* Label:") || strings.Contains(line, "* Title:") {
			count++
		}
	}

	return count, nil
}

// checkGPUDriverUpdates checks for GPU driver updates
func checkGPUDriverUpdates() string {
	// Check for NVIDIA GPU driver updates
	if exists("nvidia-smi") {
		return checkNvidiaDriverUpdate()
	}

	// Could add AMD/Intel checks here in the future
	return ""
}

func checkNvidiaDriverUpdate() string {
	// Get current driver version
	cmd := exec.Command("nvidia-smi", "--query-gpu=driver_version", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	currentDriver := strings.TrimSpace(string(output))

	// This is a simplified check - in reality, we'd query NVIDIA's API
	// or check the package manager for nvidia-driver updates
	return fmt.Sprintf("Current: %s", currentDriver)
}
