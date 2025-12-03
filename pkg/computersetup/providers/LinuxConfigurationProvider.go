package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/model"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"path/filepath"
)

type LinuxConfigurationProvider struct{}

func (l LinuxConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
	os := osagnostic.NewOperatingSystem()
	
	return []model.IMachineConfiguration{
		model.MachineConfigurationGroup{
			GroupName: "Expected Directories",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewExpectedDirectory(filepath.Join(os.HomeDirectoryPath, "src")),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Required Tools",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewExpectedCommand("git"),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Dotfiles",
			Configs: []model.IMachineConfiguration{
				dotfiles.NewDotfilesSetup(
					"https://github.com/aallbrig/dotfiles",
					filepath.Join(os.HomeDirectoryPath, "src", "dotfiles"),
					"./fresh.sh",
				),
			},
		},
	}
}
