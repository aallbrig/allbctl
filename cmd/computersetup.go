package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			log.Fatal(fmt.Sprintf("No configuration provider found for operationg system %s", os.Name))
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ApplyConfiguration()
		log.Print(out)
	},
}
