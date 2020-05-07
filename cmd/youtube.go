package cmd

import (
	"github.com/aallbrig/allbctl/cmd/youtube"
	"github.com/aallbrig/allbctl/pkg"

	"github.com/spf13/cobra"
)

var YoutubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "root command for code youtube commands",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.HelpTextIfEmpty(cmd, args)
	},
}

func init() {
	YoutubeCmd.AddCommand(youtube.ListCmd)

	rootCmd.AddCommand(YoutubeCmd)
}
