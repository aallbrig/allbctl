package providers

import (
	"github.com/aallbrig/allbctl/pkg/computersetup/dotfiles"
	"github.com/aallbrig/allbctl/pkg/externalapi"
	"github.com/aallbrig/allbctl/pkg/model"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"log"
	"path/filepath"
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
