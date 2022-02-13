package cmd

import (
	"fmt"
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/spf13/cobra"
	"log"
)

// ComputerSetupCmd defines the root of computer setup
var ComputerSetupCmd = &cobra.Command{
	Use: "computer-setup",
	Aliases: []string{
		"computersetup",
		"cs",
		"setup",
	},
	Short: "Configure host to developer preferences (cross platform)",
	Run: func(cmd *cobra.Command, args []string) {
		os := os_agnostic.OperatingSystem{}
		identifier := computerSetup.MachineIdentifier{}
		err, name := os.GetName()
		if err != nil {
			log.Fatalf("Issues getting operating system identifier")
		}

		configProvider := identifier.ConfigurationProviderForOperatingSystem(name)
		if configProvider == nil {
			log.Fatal(fmt.Sprintf("No configuration provider found for operationg system %s", os))
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ApplyConfiguration()
		log.Print(out)
	},
}
