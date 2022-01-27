package macbook

import (
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

func (s ScreenshotDirectory) Validate() error {
	return screenshotDirectoryExpectedState.Validate()
}

func (s ScreenshotDirectory) Install() error {
	cmd := exec.Command("mkdir", "-p", "~/Desktop/ScreenShots")
	err := cmd.Run()
	if err != nil {
		return err
	}
	err = screenshotDirectoryExpectedState.Install()
	return err
}

func (s ScreenshotDirectory) Uninstall() error {
	return screenshotDirectoryExpectedState.Uninstall()
}
