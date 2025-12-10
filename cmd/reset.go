package cmd

import (
	"bytes"
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/aallbrig/allbctl/pkg/status"
	"log"

	"github.com/spf13/cobra"
)

// ResetCmd represents status command
var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets machine",
	Run: func(cmd *cobra.Command, args []string) {
		out := bytes.NewBufferString("")
		out.WriteString("System Info\n")
		out.WriteString("-----\n")
		err := status.SystemInfo(out)
		if err != nil {
			log.Fatalf("Issues getting operating system identifier")
		}
		out.WriteString("\n")

		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			log.Fatalf("No configuration provider found for operationg system %s", os.Name)
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, statusOut := tweaker.ResetConfiguration()
		out.WriteString(statusOut.String())

		log.Print(out)
	},
}
