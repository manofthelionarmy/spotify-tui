package pcke

import (
	"context"
	"net/http"
	"os/exec"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

// Can this be a random clientID?
const clientID = "66934524ee284599bf9862b60b7fac53"

type pcke struct {
	state         string
	codeVerifier  string
	codeChallenge string
	ch            chan *spotify.Client
	spotifyAuth   *spotifyauth.Authenticator
}

// Auth will create the server to host our redirectURI and get our access token
func Auth() (*spotify.Client, error) {
	// create the code_verifier and code_challenge
	codeVerifier, _ := generateRandomString(128)
	codeChallenge := sha256UrlEncode(codeVerifier)

	// check the api documenation for why we need a state
	// we need the state to prevent CSRF
	state := generateRandomState()

	// Create the spotify Authenticator and pass in these optional values
	spotifyAuth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(clientID))

	p := pcke{
		state:         state,
		codeVerifier:  string(codeVerifier),
		codeChallenge: codeChallenge,
		spotifyAuth:   spotifyAuth,
		ch:            make(chan *spotify.Client),
	}

	// create a new server that handles our redirectURI as a route. there will be hander that gets a token
	server := p.newServer()
	go server.ListenAndServe()
	// Shutdown the server after we get the token and create a client
	defer server.Shutdown(context.Background())

	// build the url by passing in the codeChallenge and state
	url := p.spotifyAuth.AuthURL(p.state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", p.codeChallenge),
	)

	// open the spotify auth url in the browser
	cmd := exec.Command("xdg-open", url)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// a client will be sent to this channel via the http handler handling our redirectURI as a route
	client := <-p.ch
	return client, nil
}

// newServer creates a new server with our mux that handles our redirectURI as a route
func (p pcke) newServer() http.Server {
	authServer := http.Server{
		Addr:    ":8080",
		Handler: p.routes(),
	}
	return authServer
}
