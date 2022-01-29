package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/macbook"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/aallbrig/allbctl/pkg/model"
	"path/filepath"
)

var os = os_agnostic.OperatingSystem{}

type MacbookConfigurationProvider struct{}

func (m MacbookConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
	_, homeDir := os.HomeDir()
	return []model.IMachineConfiguration{
		os_agnostic.NewExpectedDirectory(filepath.Join(homeDir, "src")),
		os_agnostic.NewExpectedDirectory(filepath.Join(homeDir, "bin")),
		macbook.NewScreenshotDirectory(),
		macbook.NewTrackpadScrolling(),
		macbook.NewTrackpadTwoFingerRightClick(),
	}
}
