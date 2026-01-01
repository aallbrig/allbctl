package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "allbctl",
	Short: "allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks",
	Long: `allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks.

Example commands for allbctl:

$ allbctl bootstrap status
$ allbctl bootstrap install
$ allbctl status
$ allbctl status runtimes              # Show detected programming runtimes
$ allbctl status projects              # Show git repositories in ~/src
$ allbctl status list-packages         # Show package counts from all package managers
`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() //nolint:errcheck // Help errors are not critical
	},
}

// Execute comment for execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(BootstrapCmd)
	rootCmd.AddCommand(StatusCmd)

	// Add subcommands to status
	StatusCmd.AddCommand(RuntimesCmd)
	StatusCmd.AddCommand(ListPackagesCmd)
	StatusCmd.AddCommand(ProjectsCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.allbctl.yaml)")

	rootCmd.SetVersionTemplate(fmt.Sprintf("allbctl %s (commit %s)\n", Version, Commit))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".allbctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
