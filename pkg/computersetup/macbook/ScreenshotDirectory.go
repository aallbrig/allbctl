package macbook

import (
	"bytes"
	"os/exec"
	"path/filepath"
)

var screenshotsDirectory = filepath.Join("~", "Desktop", "Screenshots")
var screenshotDirectoryExpectedState = &DefaultsCommand{
	Domain:        "com.apple.screencapture",
	Key:           "location",
	ExpectedValue: screenshotsDirectory,
	ValueType:     DefaultsString,
}

type ScreenshotDirectory struct{}

func (s ScreenshotDirectory) Name() string {
	return "Screenshots Directory"
}

func NewScreenshotDirectory() *ScreenshotDirectory {
	return &ScreenshotDirectory{}
}

func (s ScreenshotDirectory) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := screenshotDirectoryExpectedState.Validate()
	out.WriteString(validateOut.String())
	return err, out
}

func (s ScreenshotDirectory) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	cmd := exec.Command("mkdir", "-p", screenshotsDirectory)
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	err, installOut := screenshotDirectoryExpectedState.WriteExpectedValue()
	out.WriteString(installOut.String())

	cmd = exec.Command("killall", "SystemUIServer")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return err, out
}

func (s ScreenshotDirectory) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := screenshotDirectoryExpectedState.Delete()
	out.WriteString(uninstallOut.String())

	cmd := exec.Command("killall", "SystemUIServer")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return err, out
}
