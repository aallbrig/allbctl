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
	err, createOut := CreateScreenshotsDirectory()
	out.WriteString(createOut.String())
	if err != nil {
		return err, out
	}

	err, installOut := screenshotDirectoryExpectedState.WriteExpectedValue()
	out.WriteString(installOut.String())
	if err != nil {
		return err, out
	}

	err, _ = RestartSystemUIServer()

	err, validateOut := s.Validate()
	out.WriteString(validateOut.String())
	if err != nil {
		return err, out
	}
	return err, out
}

func CreateScreenshotsDirectory() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("Creating screenshots directory")
	cmd := exec.Command("mkdir", "-p", screenshotsDirectory)
	err := cmd.Run()
	if err != nil {
		out.WriteString(" create fail ❌\n")
		return err, out
	}
	out.WriteString(" success ✅\n")
	return nil, out
}

func (s ScreenshotDirectory) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, _ := screenshotDirectoryExpectedState.Delete()

	err, _ = RestartSystemUIServer()
	return err, out
}
