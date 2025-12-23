package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of allbctl",
	Long:  `Print the version number and commit hash of allbctl`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("allbctl %s (commit %s)\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
