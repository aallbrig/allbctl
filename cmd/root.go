package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	computerSetup "github.com/aallbrig/allbctl/cmd/computer-setup"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "allbctl",
	Short: "allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks",
	Long: `allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks.

Example commands for allbctl:

$ allbctl computersetup
`,
	Run: func(cmd *cobra.Command, args []string) {
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
	rootCmd.AddCommand(computerSetup.RootCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.allbctl.yaml)")
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
