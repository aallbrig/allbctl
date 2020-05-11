package youtube

import (
	yt "github.com/aallbrig/allbctl/pkg/youtube"
	"github.com/spf13/cobra"
	"log"
)

var VideosCmd = &cobra.Command{
	Use:   "videos",
	Short: "lists videos of playlists",
	Run: func(cmd *cobra.Command, args []string) {
		if err := yt.ListVideos(); err != nil {
			log.Fatalf("error executing video list command:\n%v", err)
		}
	},
}

