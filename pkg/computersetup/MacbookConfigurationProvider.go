package computersetup

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/macbook"
)

type MacbookConfigurationProvider struct{}

func (m MacbookConfigurationProvider) GetConfiguration() []IMachineConfiguration {
	return []IMachineConfiguration{
		macbook.NewScreenshotDirectory(),
		macbook.NewTrackpadScrolling(),
		macbook.NewTrackpadTwoFingerRightClick(),
	}
}
