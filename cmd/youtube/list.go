package youtube

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/api/youtube/v3"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists playlists and videos in playlist",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		youtubeService, err := youtube.NewService(ctx)
		if err != nil {
			log.Fatalf("Error initializing youtube service: %v", err)
		}
		fmt.Printf("%v", youtubeService)
	},
}
