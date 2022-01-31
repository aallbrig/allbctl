package macbook

import (
	"bytes"
	"os/exec"
)

type TrackpadTwoFingerRightClick struct {
	ExpectedState DefaultsCommand
}

func (t TrackpadTwoFingerRightClick) Name() string {
	return "Trackpad Two Finger Right Click"
}

var rightClickExpectedState = DefaultsCommand{
	Domain:        "com.apple.AppleMultitouchTrackpad",
	Key:           "TrackpadRightClick",
	ExpectedValue: "1",
	ValueType:     DefaultsInt,
}

func NewTrackpadTwoFingerRightClick() *TrackpadTwoFingerRightClick {
	return &TrackpadTwoFingerRightClick{ExpectedState: rightClickExpectedState}
}

func (t TrackpadTwoFingerRightClick) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := t.ExpectedState.Validate()
	out.WriteString(validateOut.String())
	return err, out
}

func (t TrackpadTwoFingerRightClick) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := t.Validate()
	if err == nil {
		out.WriteString(validateOut.String())
		return nil, out
	}
	err, installOut := t.ExpectedState.WriteExpectedValue()
	out.WriteString(installOut.String())

	cmd := exec.Command("killall", "SystemUIServer")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return err, out
}

func (t TrackpadTwoFingerRightClick) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := t.ExpectedState.Delete()
	out.WriteString(uninstallOut.String())

	cmd := exec.Command("killall", "SystemUIServer")
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return err, out
}
