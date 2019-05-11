package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "github.com/rapito/go-spotify/spotify"
	"spotify-filler/auth"
	"spotify-filler/playlist"

	"github.com/zmb3/spotify"
)

var supportedExtensions = []string{
	".3gp",
	".mp3",
	".aa",
	".aac",
	".aiff",
	".au",
	".flac",
	".m4a",
	".ogg",
	".oga",
	".wav",
	".aax",
	".mp4",
}

var currentPlaylist playlist.Playlist

func main() {
	var root string
	if len(os.Args) == 1 {
		log.Fatal("No path given, Please specify path.")
		return
	}
	if root = os.Args[1]; root == "" {
		log.Fatal("No path given, Please specify path.")
		return
	}

	// filepath.Walk
	files, err := IOReadDir(root)
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		log.Fatalf("Directory '%s' does not contain files with extensions %s", root, strings.Join(supportedExtensions, ", "))
	}

	client := auth.GetClient()
	currentPlaylist, err = playlist.ChoosePlaylistPrompt(client)

	if err != nil {
		log.Fatal(err)
	}

	for _, title := range files {
		Search(client, title)
	}
}

func Search(client *spotify.Client, title string) {
	// search for playlists and albums containing "holiday"
	results, err := client.Search(title, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	saved := false
	// handle album results
	if results.Tracks != nil {
		for _, item := range results.Tracks.Tracks {
			if strings.ContainsAny(title, item.Name) {

				if playlist.HasItem(item.ID) {
					// fmt.Printf("[Exists] Track '%s' on playlist '%s'\n", title, currentPlaylist.Name)
					saved = true
					break
				}

				err := playlist.SaveToPlaylist(client, item.ID)
				if err != nil {
					log.Fatalf("[Fail] Track '%s' to playlist '%s': %s\n", title, currentPlaylist.Name, err)
				}
				fmt.Printf("[Saved] Track '%s' to playlist %s\n", item.Name, currentPlaylist.Name)
				saved = true
				break
			}
		}

	}
	if saved == false {
		fmt.Printf("[Not Found] Track '%s'\n", title)
	}
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	var extension string
	for _, file := range fileInfo {
		if file.IsDir() {
			files, _ = IOReadDir(root + "/" + file.Name())
		}

		extension = filepath.Ext(file.Name())
		if Contains(supportedExtensions, extension) {
			files = append(files, strings.TrimSuffix(file.Name(), extension))
		}

	}
	return files, nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
