package cmd

import (
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/spf13/cobra"
	"log"
)

var ComputerSetupCmd = &cobra.Command{
	Use: "computer-setup",
	Aliases: []string{
		"computersetup",
		"cs",
		"setup",
	},
	Short: "Configure host to developer preferences (cross platform)",
	Long: `Configure host to developer preferences.

Available subcommands:
  status - Check current setup status
  install - Apply configuration to setup this machine`,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check current setup status",
	Run: func(cmd *cobra.Command, args []string) {
		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			log.Fatalf("No configuration provider found for operating system %s", os.Name)
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ConfigurationStatus()
		log.Print(out)
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Apply configuration to setup this machine",
	Run: func(cmd *cobra.Command, args []string) {
		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			log.Fatalf("No configuration provider found for operating system %s", os.Name)
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ApplyConfiguration()
		log.Print(out)
	},
}

func init() {
	ComputerSetupCmd.AddCommand(statusCmd)
	ComputerSetupCmd.AddCommand(installCmd)
}
