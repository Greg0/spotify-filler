package playlist

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

type Playlist struct {
	ID   spotify.ID
	Name string
}

var currentPlaylistID spotify.ID
var currentPlaylistTrackIDList []spotify.ID

func ChoosePlaylistPrompt(client *spotify.Client) (Playlist, error) {

	promptItems := getPromptItems(client)
	templates := getPlaylistTemplates()
	prompt := promptui.Select{
		Label:     "Select target playlist",
		Items:     promptItems,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return Playlist{}, err
	}

	loadPlaylistItems(client, promptItems[i])
	currentPlaylistID = promptItems[i].ID

	return promptItems[i], nil
}

func HasItem(itemID spotify.ID) bool {
	for _, trackID := range currentPlaylistTrackIDList {
		if trackID == itemID {
			return true
		}
	}

	return false
}

func SaveToPlaylist(client *spotify.Client, itemID spotify.ID) error {
	var err error
	if currentPlaylistID != "" {
		_, err = client.AddTracksToPlaylist(currentPlaylistID, itemID)
	} else {
		err = client.AddTracksToLibrary(itemID)
	}

	if err == nil {
		return err
	}

	currentPlaylistTrackIDList = append(currentPlaylistTrackIDList, itemID)

	return nil
}

func getPromptItems(client *spotify.Client) []Playlist {

	user, _ := client.CurrentUser()

	playlists, err := client.GetPlaylistsForUser(user.ID)

	if err != nil {
		log.Fatalf("could not get playlists: %v", err)
	}

	promptItems := []Playlist{}
	promptItems = append(promptItems, Playlist{"", "[Favorites Library]"})

	for _, p := range playlists.Playlists {
		promptItems = append(promptItems, Playlist{p.ID, p.Name})
	}

	return promptItems
}

func getPlaylistTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    fmt.Sprintf("%s {{.Name}}: ", promptui.IconInitial),
		Active:   fmt.Sprintf("%s {{ .Name | underline }}", promptui.IconSelect),
		Inactive: "  {{.Name}}",
		Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}`, promptui.IconGood),
	}
}

func loadPlaylistItems(client *spotify.Client, playlist Playlist) {
	fmt.Println("Loading playlist items to memory...")
	currentPlaylistTrackIDList = []spotify.ID{}
	var limit int
	offset := 0
	opts := spotify.Options{Limit: &limit, Offset: &offset}
	if playlist.ID != "" {
		limit = 100
		for {
			response, err := client.GetPlaylistTracksOpt(playlist.ID, &opts, "next,total,items(track(id))")
			if err != nil {
				log.Fatal(err)
			}
			for _, track := range response.Tracks {
				currentPlaylistTrackIDList = append(currentPlaylistTrackIDList, track.Track.ID)
			}
			if response.Next == "" {
				break
			} else {
				offset = offset + limit
			}
		}
	} else {
		limit = 50
		for {
			response, err := client.CurrentUsersTracksOpt(&opts)
			if err != nil {
				log.Fatal(err)
			}
			for _, track := range response.Tracks {
				currentPlaylistTrackIDList = append(currentPlaylistTrackIDList, track.ID)
			}
			if response.Next == "" {
				break
			} else {
				offset = offset + limit
			}
		}
	}

	fmt.Println("Loading complete!")
}
