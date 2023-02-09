package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os/exec"
	"strings"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

// Can this be a random clientID?
const clientID = "66934524ee284599bf9862b60b7fac53"

var (
	auth          = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate), spotifyauth.WithClientID(clientID))
	ch            = make(chan *spotify.Client)
	codeVerifier  = ""
	codeChallenge = ""
)

type app struct {
	state         string
	codeVerifier  string
	codeChallenge string
}

// PCKE returns an authenticator and that does the PCKE flow
func PCKE() Authenticator {
	return authFunc(func() (*spotify.Client, error) {
		codeVerifier, _ := generateRandomString(128)
		codeChallenge := sha256UrlEncode(codeVerifier)
		// check the api documenation
		state := generateRandomState()

		// set up server
		server := &app{
			state:         state,
			codeVerifier:  string(codeVerifier),
			codeChallenge: codeChallenge,
		}

		// set up a serve mux
		mux := http.NewServeMux()
		mux.HandleFunc("/callback", server.completeAuth)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got request for:", r.URL.String())
		})

		authServer := http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		go authServer.ListenAndServe()
		defer authServer.Shutdown(context.Background())

		url := auth.AuthURL(state,
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		)

		// open url
		cmd := exec.Command("xdg-open", url)
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		client := <-ch
		return client, nil
	})
}

// They said it's best to use crypto/rand
func generateRandomString(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_.-~"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		ret[i] = letters[num.Int64()]
	}

	return ret, nil
}

func sha256UrlEncode(b []byte) string {
	res := sha256.Sum256(b)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(res[:]), "=")
}

func generateRandomState() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *app) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), s.state, r,
		oauth2.SetAuthURLParam("code_verifier", s.codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	// ah cool, here we're validating the state
	if st := r.FormValue("state"); st != s.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, s.state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
