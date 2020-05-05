package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var YoutubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "root command for code youtube commands",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(YoutubeCmd)
}
