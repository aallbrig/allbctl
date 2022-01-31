package cmd

import (
	"bytes"
	"fmt"
	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/computersetup/os_agnostic"
	"github.com/aallbrig/allbctl/pkg/status"
	"log"

	"github.com/spf13/cobra"
)

// StatusCmd represents status command
var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets machine",
	Run: func(cmd *cobra.Command, args []string) {
		out := bytes.NewBufferString("")
		out.WriteString("System Info\n")
		out.WriteString("-----\n")
		err := status.SystemInfo(out)
		out.WriteString("\n")

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
		_, statusOut := tweaker.ResetConfiguration()
		out.WriteString(statusOut.String())

		log.Print(out)
	},
}