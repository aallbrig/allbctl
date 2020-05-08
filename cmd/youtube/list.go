package youtube

import (
	yt "github.com/aallbrig/allbctl/pkg/youtube"
	"github.com/spf13/cobra"
	"log"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists playlists and videos in playlist",
	Run: func(cmd *cobra.Command, args []string) {
		err := yt.List()
		if err != nil {
			log.Fatalf("error executing list command: %v", err)
		}
	},
}
