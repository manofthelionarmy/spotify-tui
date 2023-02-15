package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
)

var (
	// ErrTokenNotFound is returned when we couldn't find the token at the default or set location
	ErrTokenNotFound = errors.New("token not found. Specify the correct location of the token")

	// ErrSpotifyTuiDirNotFound is returned when we couldn't find the directory
	ErrSpotifyTuiDirNotFound = errors.New("~/.spotify-tui not found")

	// ErrSpotifyTokenExpired is returned when the token expires
	ErrSpotifyTokenExpired = errors.New("token expired")
)

// StoreToken stores the token in a file
func StoreToken(token oauth2.Token) error {
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	appFolder := filepath.Join(homeDir, ".spotify-tui")
	_, err = os.Stat(filepath.Join(appFolder))

	// check if the directory exists
	p := filepath.Join(appFolder, "token")
	if os.IsNotExist(err) {
		err := os.MkdirAll(appFolder, 0755)
		if err != nil {
			return err
		}
	}

	// This will allow the file to be overwritten
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveToken gets the token from file
func RetrieveToken() (*oauth2.Token, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	appFolder := filepath.Join(homeDir, ".spotify-tui")

	// check if the directory exists
	_, err = os.Stat(appFolder)
	if os.IsNotExist(err) {
		return nil, ErrSpotifyTuiDirNotFound
	}

	p := filepath.Join(appFolder, "token")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		return nil, ErrTokenNotFound
	}

	// read the file
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	// unmarshal the json encoded contents into the oauth2 token
	var tok oauth2.Token
	err = json.Unmarshal(b, &tok)
	if err != nil {
		return nil, err
	}

	// this is confusing because if the token expired, they can still use their refresh token
	// but doing this treats it like a static token. I'm not sure if this is safer...
	if !valid(tok) {
		return nil, ErrSpotifyTokenExpired
	}
	return &tok, nil
}

// If it's invalid, we can still use it, I think we can pass the token and it will use the refresh token
func valid(tok oauth2.Token) bool {
	if time.Now().Sub(tok.Expiry) > 0 {
		return false
	}
	return true
}
