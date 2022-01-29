package macbook

import "bytes"

type TrackpadTwoFingerRightClick struct {
	ConfigName    string
	ExpectedState DefaultsCommand
}

func (t TrackpadTwoFingerRightClick) Name() string {
	return t.ConfigName
}

var rightClickExpectedState = DefaultsCommand{
	Domain:        "com.apple.AppleMultitouchTrackpad",
	Key:           "TrackpadRightClick",
	ExpectedValue: "1",
	DefaultValue:  "0",
}

func NewTrackpadTwoFingerRightClick() *TrackpadTwoFingerRightClick {
	rightClickExpectedState.ReadCurrentValue()
	return &TrackpadTwoFingerRightClick{ConfigName: "Trackpad Two Finger Right Click", ExpectedState: rightClickExpectedState}
}

func (t TrackpadTwoFingerRightClick) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := t.ExpectedState.Validate()
	out.WriteString(validateOut.String())
	return err, out
}

func (t TrackpadTwoFingerRightClick) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, installOut := t.ExpectedState.Install()
	out.WriteString(installOut.String())
	return err, out
}

func (t TrackpadTwoFingerRightClick) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := t.ExpectedState.Uninstall()
	out.WriteString(uninstallOut.String())
	return err, out
}
