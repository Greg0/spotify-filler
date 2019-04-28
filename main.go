package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// "github.com/rapito/go-spotify/spotify"
	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserLibraryModify,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {

	client := getClient()
	playlistID := choosePlaylist(client)

	fmt.Println(playlistID)
	fmt.Println(spotify.ID(playlistID))

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
				var err error
				if playlistID != "" {
					_, err = client.AddTracksToPlaylist(spotify.ID(playlistID), item.ID)
				} else {
					err = client.AddTracksToLibrary(item.ID)
				}

				if err == nil {
					log.Fatal(err)
				}
				fmt.Print(item.ExternalURLs)
				fmt.Println("   ", item.ID, " ", item.Name)
				break
			}
		}
	}
}

func getClient() *spotify.Client {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	return client
}

type playlist struct {
	ID   string
	Name string
}

func choosePlaylist(client *spotify.Client) string {
	user, _ := client.CurrentUser()

	playlists, err := client.GetPlaylistsForUser(user.ID)

	if err != nil {
		log.Fatalf("could not get playlists: %v", err)
	}

	promptItems := []playlist{}
	promptItems = append(promptItems, playlist{"", "[Favorites Library]"})

	for _, p := range playlists.Playlists {
		promptItems = append(promptItems, playlist{string(p.ID), p.Name})
	}

	templates := &promptui.SelectTemplates{
		Label:    fmt.Sprintf("%s {{.Name}}: ", promptui.IconInitial),
		Active:   fmt.Sprintf("%s {{ .Name | underline }}", promptui.IconSelect),
		Inactive: "  {{.Name}}",
		Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}`, promptui.IconGood),
	}

	prompt := promptui.Select{
		Label:     "Select target playlist",
		Items:     promptItems,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
		return ""
	}

	return promptItems[i].ID
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	fmt.Printf("Token: %s\n", tok.AccessToken)

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
