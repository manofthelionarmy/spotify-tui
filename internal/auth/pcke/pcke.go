package pcke

import (
	"context"
	"net/http"
	"os/exec"
	tk "spotify-tui/internal/auth/token"

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
	ch            chan msg
	spotifyAuth   *spotifyauth.Authenticator
}

type msg struct {
	tok    *oauth2.Token
	client *spotify.Client
}

// Auth will create the server to host our redirectURI and get our access token
func Auth() (*spotify.Client, error) {
	// Create the spotify Authenticator and pass in values to set up our redirectURI,
	// scope permissions, and client id
	spotifyAuth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopeStreaming,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
		spotifyauth.WithClientID(clientID))

	// Retrieve the token if it exists
	tok, err := tk.RetrieveToken()
	if err == nil {
		// No error and we have our token
		client := spotify.New(spotifyAuth.Client(context.Background(), tok))
		return client, nil
	} else if err != tk.ErrTokenNotFound &&
		err != tk.ErrSpotifyTuiDirNotFound &&
		err != tk.ErrSpotifyTokenExpired {
		// there was an error and it doens't match our expected errors
		return nil, err
	}

	// The token wasn't found, so we will do our oauth2 flow

	// create the code_verifier and code_challenge
	codeVerifier, _ := generateRandomString(128)
	codeChallenge := sha256UrlEncode(codeVerifier)

	// check the api documenation for why we need a state
	// we need the state to prevent CSRF
	state := generateRandomState()

	// TODO: figure out other permissions
	p := pcke{
		state:         state,
		codeVerifier:  string(codeVerifier),
		codeChallenge: codeChallenge,
		spotifyAuth:   spotifyAuth,
		ch:            make(chan msg),
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
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// a tok will be sent to this channel via the http handler handling our redirectURI as a route
	msg := <-p.ch
	err = tk.StoreToken(*msg.tok)
	if err != nil {
		return nil, err
	}
	return msg.client, nil
}

// newServer creates a new server with our mux that handles our redirectURI as a route
func (p pcke) newServer() http.Server {
	authServer := http.Server{
		Addr:    ":8080",
		Handler: p.routes(),
	}
	return authServer
}
