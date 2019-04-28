package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

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

func GetClient() *spotify.Client {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:\n", url)

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
