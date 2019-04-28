package playlist

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

type playlist struct {
	ID   spotify.ID
	Name string
}

func ChoosePlaylistPrompt(client *spotify.Client) (playlist, error) {

	promptItems := getPromptItems(client)
	templates := getPlaylistTemplates()
	prompt := promptui.Select{
		Label:     "Select target playlist",
		Items:     promptItems,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return playlist{}, err
	}

	return promptItems[i], nil
}

func SaveToPlaylist(client *spotify.Client, playlistID spotify.ID, itemID spotify.ID) error {
	var err error
	if playlistID != "" {
		_, err = client.AddTracksToPlaylist(playlistID, itemID)
	} else {
		err = client.AddTracksToLibrary(itemID)
	}

	if err == nil {
		return err
	}

	return nil
}

func getPromptItems(client *spotify.Client) []playlist {

	user, _ := client.CurrentUser()

	playlists, err := client.GetPlaylistsForUser(user.ID)

	if err != nil {
		log.Fatalf("could not get playlists: %v", err)
	}

	promptItems := []playlist{}
	promptItems = append(promptItems, playlist{"", "[Favorites Library]"})

	for _, p := range playlists.Playlists {
		promptItems = append(promptItems, playlist{p.ID, p.Name})
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
