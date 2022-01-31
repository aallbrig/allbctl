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
	ValueType:     DefaultsInt,
}

var scrollDirection = DefaultsCommand{
	Domain:        "'Apple Global Domain'",
	Key:           "com.apple.swipescrolldirection",
	ExpectedValue: "0",
	ValueType:     DefaultsInt,
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
	out.WriteString(validateOut.String() + "\n")
	if err != nil {
		return err, out
	}

	return nil, out
}

func (t TrackpadScroll) Install() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, validateOut := scrollDirection.Validate()
	if err == nil {
		out.WriteString(validateOut.String() + "\n")
	} else {
		err, installOut := scrollDirection.WriteExpectedValue()
		out.WriteString(installOut.String() + "\n")
		if err != nil {
			return err, out
		}
	}

	err, validateOut = scrollingEnabled.Validate()
	if err == nil {
		out.WriteString(validateOut.String() + "\n")
	} else {
		err, installOut := scrollingEnabled.WriteExpectedValue()
		out.WriteString(installOut.String() + "\n")
		if err != nil {
			return err, out
		}
	}

	return nil, out
}

func (t TrackpadScroll) Uninstall() (error, *bytes.Buffer) {
	out := bytes.NewBufferString("")
	err, uninstallOut := scrollDirection.Delete()
	out.WriteString(uninstallOut.String() + "\n")
	if err != nil {
		return err, out
	}

	err, uninstallOut = scrollingEnabled.Delete()
	out.WriteString(uninstallOut.String() + "\n")
	if err != nil {
		return err, out
	}

	return nil, out
}
