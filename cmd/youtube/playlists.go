package youtube

import (
	yt "github.com/aallbrig/allbctl/pkg/youtube"
	"github.com/spf13/cobra"
	"log"
)

var PlaylistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "lists playlists",
	Run: func(cmd *cobra.Command, args []string) {
		if err := yt.ListPlaylists(); err != nil {
			log.Fatalf("error executing playlist list command: %v", err)
		}
	},
}

