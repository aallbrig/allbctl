package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	computerSetup "github.com/aallbrig/allbctl/pkg/computersetup"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/aallbrig/allbctl/pkg/status"
	"github.com/aallbrig/allbctl/pkg/telemetry"
	"github.com/spf13/cobra"
)

var (
	registerSSHKeys bool
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
		_ = cmd.Help() //nolint:errcheck // Help errors are not critical
	},
}

var bootstrapStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check workstation bootstrap status",
	Long:  `Check the status of workstation bootstrap configuration including directories, tools, SSH keys, and dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		printBootstrapStatus(ctx)
	},
}

var bootstrapInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure workstation bootstrap",
	Long: `Install and configure workstation bootstrap including directories, tools, SSH keys, and dotfiles.

By default, SSH key generation and GitHub registration are SKIPPED.
Use --register-ssh-keys flag to enable SSH key generation and GitHub registration.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		os := osagnostic.NewOperatingSystem()
		identifier := computerSetup.MachineIdentifier{}
		configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
		if configProvider == nil {
			fmt.Printf("No configuration provider for %s\n", os.Name)
			return
		}

		// Get configuration and filter out SSH key registration if flag not set
		configs := configProvider.GetConfiguration()
		if !registerSSHKeys {
			configs = computerSetup.FilterOutSSHKeyRegistration(configs)
		}

		telemetry.Logger.InfoContext(ctx, "bootstrap.install.start",
			"os", os.Name,
			"register_ssh_keys", registerSSHKeys,
			"config_count", len(configs),
		)

		tweaker := computerSetup.NewMachineTweaker(configs)
		_, out := tweaker.ApplyConfiguration()
		fmt.Print(out.String())

		telemetry.Logger.InfoContext(ctx, "bootstrap.install.finish",
			"os", os.Name,
			"config_count", len(configs),
		)
	},
}

var bootstrapResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset workstation bootstrap configuration",
	Long:  `Reset workstation bootstrap configuration, removing directories, tools, SSH keys, and dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

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

		telemetry.Logger.InfoContext(ctx, "bootstrap.reset.start", "os", os.Name)

		_, statusOut := tweaker.ResetConfiguration()
		out.WriteString(statusOut.String())

		telemetry.Logger.InfoContext(ctx, "bootstrap.reset.finish", "os", os.Name)

		log.Print(out)
	},
}

func printBootstrapStatus(ctx context.Context) {
	os := osagnostic.NewOperatingSystem()
	identifier := computerSetup.MachineIdentifier{}
	configProvider := identifier.ConfigurationProviderForOperatingSystem(os.Name)
	if configProvider == nil {
		fmt.Printf("No configuration provider for %s\n", os.Name)
		return
	}

	tweaker := computerSetup.NewMachineTweaker(configProvider.GetConfiguration())
	_, out := tweaker.ConfigurationStatus()

	telemetry.Logger.InfoContext(ctx, "bootstrap.status",
		"os", os.Name,
		"config_count", len(configProvider.GetConfiguration()),
	)

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

	// Add flags to install command
	bootstrapInstallCmd.Flags().BoolVar(&registerSSHKeys, "register-ssh-keys", false, "Generate SSH keys and register with GitHub (requires gh CLI)")
}
