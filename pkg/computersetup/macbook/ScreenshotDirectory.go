package macbook

import (
	"bytes"
	"os/exec"
)

var screenshotsDirectory = "~/Desktop/Screenshots"
var screenshotDirectoryExpectedState = &DefaultsCommand{
	Domain:        "com.apple.screencapture",
	Key:           "location",
	ExpectedValue: screenshotsDirectory,
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
	err := cmd.Run()
	if err != nil {
		return err, out
	}

	err, installOut := screenshotDirectoryExpectedState.Install()
	out.WriteString(installOut.String())
	return err, out
}

func (s ScreenshotDirectory) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := screenshotDirectoryExpectedState.Uninstall()
	out.WriteString(uninstallOut.String())
	return err, out
}
