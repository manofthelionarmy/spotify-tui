package models

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/zmb3/spotify/v2"
)

// Song represents the Tracks object from Spotify
type Song struct {
	Name string
}

// GetSongs gets all of the artists songs based on their spotify id
func GetSongs(client *spotify.Client, artistID spotify.ID) ([]*Song, error) {
	// I see what I'm doing wrong, this is based on track id, but I'm using spotify id
	albumPage, err := client.GetArtistAlbums(context.Background(),
		artistID, []spotify.AlbumType{spotify.AlbumTypeAlbum},
	)

	album := albumPage.Albums[0]

	albumTracksPage, err := client.GetAlbumTracks(context.Background(), album.ID)
	if err != nil {
		log.Fatal(err)
	}

	songs := make([]*Song, len(albumTracksPage.Tracks))
	for i := range albumTracksPage.Tracks {
		// BUG: I don't know what's going on but it seems like we get this as null
		s := &Song{
			Name: albumTracksPage.Tracks[i].Name,
		}

		songs[i] = s
	}
	return songs, nil
}

// FilterValue is the value we use when filtering against this item when we're filtering the list.
func (s *Song) FilterValue() string {
	return s.Name
}

// Title returns the song title
func (s *Song) Title() string {
	return s.Name
}

// Description returns the song description
func (s *Song) Description() string {
	return ""
}

var _ list.DefaultItem = (*Song)(nil)
