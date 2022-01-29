package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/macbook"
	"github.com/aallbrig/allbctl/pkg/model"
)

type MacbookConfigurationProvider struct{}

func (m MacbookConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
	return []model.IMachineConfiguration{
		macbook.NewScreenshotDirectory(),
		macbook.NewTrackpadScrolling(),
		macbook.NewTrackpadTwoFingerRightClick(),
	}
}
