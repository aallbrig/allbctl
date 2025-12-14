package providers

import (
	"path/filepath"

	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/model"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
)

type WindowsConfigurationProvider struct{}

func (w WindowsConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
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
				osagnostic.NewInstallableCommand("git").
					SetWindowsPackage("winget", "Git.Git").
					SetWindowsPackage("choco", "git").
					SetWindowsPackage("scoop", "git"),
				osagnostic.NewInstallableCommand("gh").
					SetWindowsPackage("winget", "GitHub.cli").
					SetWindowsPackage("choco", "gh").
					SetWindowsPackage("scoop", "gh"),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "SSH Configuration",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewSSHKeyGitHubRegistration(),
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
		osagnostic.NewShellConfigTools(),
	}
}
