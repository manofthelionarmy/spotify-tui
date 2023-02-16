package models

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// Album contains data we need from spotify about albums
type Album struct {
	Name     string
	ID       spotify.ID
	AlbumURI spotify.URI
}

// GetAlbums gets the albums belonging to an artist given their artist id
func GetAlbums(client *spotify.Client, artistID spotify.ID) ([]*Album, error) {
	artistsAlbums, err := client.GetArtistAlbums(context.Background(), artistID,
		[]spotify.AlbumType{spotify.AlbumTypeAlbum})
	if err != nil {
		return nil, err
	}

	albums := make([]*Album, len(artistsAlbums.Albums))
	for i, album := range artistsAlbums.Albums {
		albums[i] = &Album{
			Name:     album.Name,
			ID:       album.ID,
			AlbumURI: album.URI,
		}
	}
	return albums, nil
}

// FilterValue is the value we use when filtering against this item when
// we're filtering the list.
func (a *Album) FilterValue() string {
	return a.Name
}

// Title displays the albums name as the title
func (a *Album) Title() string {
	return a.Name
}

// Description is empty for now
func (a *Album) Description() string {
	return ""
}
