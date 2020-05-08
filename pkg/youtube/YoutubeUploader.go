package youtube

import (
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
)

func List() error {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		fmt.Printf("Error initializing youtube service: %v\n", err)
		return err
	}
	fmt.Printf("%v\n", youtubeService)
	return nil
}
