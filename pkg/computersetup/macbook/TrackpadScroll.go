package macbook

import "bytes"

type TrackpadScroll struct {
	ConfigName string
}

func (t TrackpadScroll) Name() string {
	return t.ConfigName
}

var scrollingEnabled = DefaultsCommand{
	Domain:        "com.apple.AppleMultitouchTrackpad",
	Key:           "TrackpadScroll",
	ExpectedValue: "1",
	DefaultValue:  "1",
}

var scrollDirection = DefaultsCommand{
	Domain:        "'Apple Global Domain'",
	Key:           "com.apple.swipescrolldirection",
	ExpectedValue: "0",
	DefaultValue:  "1",
}

func NewTrackpadScrolling() *TrackpadScroll {
	return &TrackpadScroll{ConfigName: "Trackpad Scroll"}
}

func (t TrackpadScroll) Validate() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")

	err, validateOut := scrollDirection.Validate()
	out.WriteString(validateOut.String() + "\n")
	if err != nil {
		return err, out
	}

	err, validateOut = scrollingEnabled.Validate()
	out.WriteString(validateOut.String())
	if err != nil {
		return err, out
	}

	return nil, out
}

func (t TrackpadScroll) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, installOut := scrollDirection.Install()
	out.WriteString(installOut.String() + "\n")
	if err != nil {
		return err, out
	}

	err, installOut = scrollingEnabled.Uninstall()
	out.WriteString(installOut.String())
	if err != nil {
		return err, out
	}
	return nil, out
}

func (t TrackpadScroll) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := scrollDirection.Uninstall()
	out.WriteString(uninstallOut.String() + "\n")
	if err != nil {
		return err, out
	}

	err, uninstallOut = scrollingEnabled.Uninstall()
	out.WriteString(uninstallOut.String())
	if err != nil {
		return err, out
	}

	return nil, out
}
