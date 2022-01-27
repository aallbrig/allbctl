package macbook

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

func (t TrackpadScroll) Validate() error {
	err := scrollDirection.Validate()
	if err != nil {
		return err
	}

	err = scrollingEnabled.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (t TrackpadScroll) Install() error {
	err := scrollDirection.Install()
	if err != nil {
		return err
	}

	err = scrollingEnabled.Uninstall()
	if err != nil {
		return err
	}
	return nil
}

func (t TrackpadScroll) Uninstall() error {
	err := scrollDirection.Uninstall()
	if err != nil {
		return err
	}

	err = scrollingEnabled.Uninstall()
	if err != nil {
		return err
	}
	return nil
}
