package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/computersetup/osagnostic"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/model"
	"log"
	"path/filepath"
)

var os = osagnostic.OperatingSystem{}

type MacbookConfigurationProvider struct{}

func (m MacbookConfigurationProvider) GetConfiguration() []model.IMachineConfiguration {
	homeDir, _ := os.HomeDir()
	return []model.IMachineConfiguration{
		model.MachineConfigurationGroup{
			GroupName: "Expected Directories",
			Configs: []model.IMachineConfiguration{
				osagnostic.NewExpectedDirectory(filepath.Join(homeDir, "src")),
				osagnostic.NewExpectedDirectory(filepath.Join(homeDir, "bin")),
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
		dotfiles.NewDotfilesGremlin(),
	}
}
