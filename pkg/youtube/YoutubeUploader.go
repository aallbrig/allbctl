package youtube

import (
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"os"
)

var channelId = os.Getenv("YOUTUBE_CHANNEL_ID")

func getYoutubeService() (*youtube.Service, error) {
	ctx := context.Background()
	yt, err := youtube.NewService(ctx)
	if err != nil {
		fmt.Printf("Error initializing youtube service: %v\n", err)
		return nil, err
	}
	return yt, nil
}

func ListPlaylists() error {
	yt, err := getYoutubeService()
	if err != nil {
		return err
	}

	call := yt.Playlists.List("snippet,contentDetails")
	call.ChannelId(channelId)
	call.MaxResults(25)
	resp, err := call.Do()
	if err != nil {
		fmt.Printf("Error executing playlist listing: %v\n", err)
		return err
	}
	for _, playlist := range resp.Items {
		fmt.Printf("%s\n%s\n", playlist.Snippet.Title, playlist.Snippet.Description)
	}
	return nil
}

func ListVideos() error {
	yt, err := getYoutubeService()
	if err != nil {
		fmt.Printf("Error initializing youtube service: %v\n", err)
		return err
	}
	call := yt.Channels.List("contentDetails")
	call = call.Mine(true)
	response, err := call.Do()
	if err != nil {
		fmt.Printf("error getting client")
	}
	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads
		fmt.Printf("Videos in playlist %s\n", playlistId)

		nextPageToken := ""
		for {
			call := yt.PlaylistItems.List("snippet")
			call = call.PlaylistId(playlistId)
			if nextPageToken != "" {
				call = call.PageToken(nextPageToken)
			}
			playlistResp, err := call.Do()
			if err != nil {
				fmt.Printf("Error getting playlist videos")
				return err
			}
			for _, playlistItem := range playlistResp.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId
				fmt.Printf("%v (%v)\n", title, videoId)
			}

			nextPageToken = playlistResp.NextPageToken
			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}
	}

	return nil
}
