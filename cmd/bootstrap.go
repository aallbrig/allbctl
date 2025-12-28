package cmd

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/aallbrig/allbctl/pkg/status"
	"github.com/spf13/cobra"
)

var BootstrapCmd = &cobra.Command{
	Use: "bootstrap",
	Aliases: []string{
		"bs",
	},
	Short: "Manage workstation bootstrap configuration",
	Long: `Manage workstation bootstrap configuration including directories, tools, SSH keys, and dotfiles.

Available subcommands:
  status - Check current bootstrap status
  install - Apply bootstrap configuration to setup this machine
  reset - Reset bootstrap configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() //nolint:errcheck // Help always succeeds
	},
}

var bootstrapStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check workstation bootstrap status",
	Long:  `Check the status of workstation bootstrap configuration including directories, tools, SSH keys, and dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		printBootstrapStatus()
	},
}

var bootstrapInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure workstation bootstrap",
	Long:  `Install and configure workstation bootstrap including directories, tools, SSH keys, and dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			fmt.Printf("No configuration provider for %s\n", os.Name)
			return
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, out := tweaker.ApplyConfiguration()
		fmt.Print(out.String())
	},
}

var bootstrapResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset workstation bootstrap configuration",
	Long:  `Reset workstation bootstrap configuration, removing directories, tools, SSH keys, and dotfiles.`,
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
			log.Fatalf("No configuration provider found for operating system %s", os.Name)
		}

		tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
		_, statusOut := tweaker.ResetConfiguration()
		out.WriteString(statusOut.String())

		log.Print(out)
	},
}

func printBootstrapStatus() {
	os := osagnostic.NewOperatingSystem()
	identifier := computerSetup.MachineIdentifier{}
	configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
	if configProvider == nil {
		fmt.Printf("No configuration provider for %s\n", os.Name)
		return
	}

	tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
	_, out := tweaker.ConfigurationStatus()

	fmt.Println("Workstation Bootstrap Status:")
	fmt.Println()
	// Indent the output
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("  %s\n", line)
		}
	}
}

func init() {
	BootstrapCmd.AddCommand(bootstrapStatusCmd)
	BootstrapCmd.AddCommand(bootstrapInstallCmd)
	BootstrapCmd.AddCommand(bootstrapResetCmd)
}
