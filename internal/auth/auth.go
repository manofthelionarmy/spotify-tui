package auth

import (
	spotify "github.com/zmb3/spotify/v2"
)

// Authenticator is one that implements authentication logic for spotify
type Authenticator interface {
	Auth() (*spotify.Client, error)
}

type authFunc func() (*spotify.Client, error)

func (f authFunc) Auth() (*spotify.Client, error) {
	return f()
}
