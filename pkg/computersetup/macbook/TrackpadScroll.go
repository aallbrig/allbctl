package macbook

type TrackpadScroll struct {
	ConfigName    string
	ExpectedState DefaultsCommand
}

func (t TrackpadScroll) Name() string {
	return t.ConfigName
}

var scrollingExpectedState = DefaultsCommand{
	Domain:        "com.apple.AppleMultitouchTrackpad",
	Key:           "TrackpadScroll",
	ExpectedValue: "0",
	DefaultValue:  "1",
}

func NewTrackpadScrolling() *TrackpadScroll {
	return &TrackpadScroll{ConfigName: "Trackpad Scroll", ExpectedState: scrollingExpectedState}
}

func (t TrackpadScroll) Validate() error {
	return t.ExpectedState.Validate()
}

func (t TrackpadScroll) Install() error {
	return t.ExpectedState.Install()
}

func (t TrackpadScroll) Uninstall() error {
	return t.ExpectedState.Uninstall()
}
