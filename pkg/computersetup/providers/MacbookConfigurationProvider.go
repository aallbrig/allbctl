package providers

import (
	"log"
	"path/filepath"

	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/model"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
)

var os = osagnostic.NewOperatingSystem()

type MacbookConfigurationProvider struct{}

func (m MacbookConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
	return []model.IMachineConfiguration{
		model.MachineConfigurationGroup{
			GroupName: "Expected Directories",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewExpectedDirectory(filepath.Join(os.HomeDirectoryPath, "src")),
				osagnostic.NewExpectedDirectory(filepath.Join(os.HomeDirectoryPath, "bin")),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Required Tools",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewInstallableCommand("git").
					SetMacOSPackage("git"),
				osagnostic.NewInstallableCommand("gh").
					SetMacOSPackage("gh"),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "SSH Configuration",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewSSHKeyGitHubRegistration(),
			},
		},
		model.MachineConfigurationGroup{
			GroupName: "Expected Environment Variables",
			Configs: []model.IMachineConfiguration{
				osagnostic.ExpectedEnvVar{
					Key: externalapi.GithubAuthTokenEnvVar,
					OnInstall: func() error {
						log.Println("Read documentation: https://cli.github.com/manual/gh_help_environment")
						return nil
					},
					OnUninstall: func() error {
						log.Println("❌ It is up to the user to uninstall this environment variable")
						return nil
					},
				},
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
