package main

import (
	"fmt"
	"log"
	"strings"

	// "github.com/rapito/go-spotify/spotify"
	"github.com/zmb3/spotify"
	"spotify-filler/auth"
	"spotify-filler/playlist"
)

func main() {

	client := auth.GetClient()
	playlistItem, err := playlist.ChoosePlaylistPrompt(client)

	if err != nil {
		log.Fatal(err)
	}

	title := "Jack White - Blunderbuss"

	// search for playlists and albums containing "holiday"
	results, err := client.Search(title, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Tracks != nil {
		fmt.Println("Tracks:")
		for _, item := range results.Tracks.Tracks {
			if strings.Contains(title, item.Name) {
				err := playlist.SaveToPlaylist(client, playlistItem.ID, item.ID)
				if err != nil {
					log.Fatalf("Failed save '%s' to '%s': %s", item.Name, playlistItem.Name, err)
				}
				fmt.Printf("Saved '%s' to %s", item.Name, playlistItem.Name)

				break
			}
		}
	}
}
