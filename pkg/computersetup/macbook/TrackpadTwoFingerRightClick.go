package macbook

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
	rightClickExpectedState.SyncCurrentValue()
	return &TrackpadTwoFingerRightClick{ConfigName: "Trackpad Two Finger Right Click", ExpectedState: rightClickExpectedState}
}

func (t TrackpadTwoFingerRightClick) Validate() error {
	return t.ExpectedState.Validate()
}

func (t TrackpadTwoFingerRightClick) Install() error {
	return t.ExpectedState.Install()
}

func (t TrackpadTwoFingerRightClick) Uninstall() error {
	return t.ExpectedState.Uninstall()
}
