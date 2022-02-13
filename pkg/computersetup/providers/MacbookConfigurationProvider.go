package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/model"
	"log"
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
		dotfiles.NewDotfilesGremlin(),
	}
}
