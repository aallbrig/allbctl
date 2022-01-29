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
		model.MachineConfigurationGroup{
			GroupName: "Expected Directories",
			Configs: []model.IMachineConfiguration{
				os_agnostic.NewExpectedDirectory(filepath.Join(homeDir, "src")),
				os_agnostic.NewExpectedDirectory(filepath.Join(homeDir, "bin")),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Expected Environment Variables",
			Configs: []model.IMachineConfiguration{
				os_agnostic.ExpectedEnvVar{
					Key: "GH_AUTH_TOKEN",
					OnInstall: func() error {
						// log out instructions/link to documentation for what this envvar is
						return nil
					},
					OnUninstall: func() error {
						return nil
					},
				},
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Macbook Configuration",
			Configs: []model.IMachineConfiguration{
				macbook.NewScreenshotDirectory(),
				macbook.NewTrackpadScrolling(),
				macbook.NewTrackpadTwoFingerRightClick(),
			},
		},
	}
}
