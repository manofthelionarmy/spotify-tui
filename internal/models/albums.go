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
	// TODO: fix bug, getting duplicate albums back, they say I need to specify the market code
	// UPDATE: I did the same query via the developer console and got back duplicate results. the problem is that
	// there are duplicate entries, via US Market Code, with different spotify uris, so I don't know if there was some kind of migration issue...
	artistsAlbums, err := client.GetArtistAlbums(context.Background(), artistID,
		[]spotify.AlbumType{spotify.AlbumTypeAlbum}, spotify.Market(spotify.CountryUSA), spotify.Limit(50), spotify.Country(spotify.MarketFromToken))
	if err != nil {
		return nil, err
	}

	albums := make([]*Album, len(artistsAlbums.Albums))
	for i, album := range artistsAlbums.Albums {
		// TODO: filter out by matching name or differing spotify uri
		// cool feature, display image in terminal
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
