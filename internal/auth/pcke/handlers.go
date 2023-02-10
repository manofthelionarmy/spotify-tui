package pcke

import (
	"fmt"
	"log"
	"net/http"

	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (p pcke) completeAuth(w http.ResponseWriter, r *http.Request) {
	// here we pass the codeVerifier in order to get the token
	tok, err := p.spotifyAuth.Token(r.Context(), p.state, r,
		oauth2.SetAuthURLParam("code_verifier", p.codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	// ah cool, here we're validating the state, I read we need to do this
	if st := r.FormValue("state"); st != p.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, p.state)
	}
	// use the token to get an authenticated client
	client := spotify.New(p.spotifyAuth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	// signal we have a client
	p.ch <- client
}
