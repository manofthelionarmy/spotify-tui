package spotify

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type authServer struct {
	close  chan bool
	server http.ServeMux
}

type config struct {
	redirectURL string
	clientID    string
	secretKey   string
	tokenDir    string
}

var cfg config

// Auth creates a new auth server and does the authentication
func Auth() (*spotify.Client, error) {
	ctx := context.Background()

	// I won't be able to access the users profile data if I use the client credentials
	// which makes me think they won't be able to access their playlists
	// let's have a agile mindset. Get a working mvp, abstract auth details away, and then implement the other kinds of auth (PCKE or authorization code flow)
	config := &clientcredentials.Config{
		ClientID:     cfg.clientID,
		ClientSecret: cfg.secretKey,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}
	// the doc says that this will renew the token
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client, nil
}

func init() {
	flag.StringVar(&cfg.redirectURL, "redirectURL", "localhost:8080", "Overrides the default redirectURL")
	flag.StringVar(&cfg.clientID, "clientID", "", "Set the client id for your application")
	flag.StringVar(&cfg.secretKey, "secretKey", "", "Set the client id for your application")
	flag.StringVar(&cfg.tokenDir, "tokenDir", "~/.spotify-tui/", "Override the token dir for spotify-tui. Default is ~/.spotify-tui/")

	flag.Parse()
	if cfg.clientID == "" {
		clientID, ok := os.LookupEnv("SPOTIFY_CLIENT_ID")
		if !ok {
			fmt.Println("Either set SPOTIFY_CLIENT_ID or use --clientID")
			os.Exit(1)
		}
		cfg.clientID = clientID
	}

	if cfg.secretKey == "" {
		secretKey, ok := os.LookupEnv("SPOTIFY_SECRET_KEY")
		if !ok {
			fmt.Println("Either set SPOTIFY_SECRET_KEY or use --secretKey")
			os.Exit(1)
		}
		cfg.secretKey = secretKey
	}
}
