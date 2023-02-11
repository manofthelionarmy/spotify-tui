package models

import (
	"context"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	"github.com/zmb3/spotify/v2"
)

// Artist represents our artist returned from spotify
type Artist struct {
	Name string
	ID   spotify.ID
}

// GetArtists returns artists by a search query
func GetArtists(cli *spotify.Client, searchQuery string) ([]*Artist, error) {
	urlQuery := url.QueryEscape(searchQuery)
	res, err := cli.Search(context.Background(), urlQuery, spotify.SearchTypeArtist)
	if err != nil {
		return nil, err
	}

	artists := make([]*Artist, len(res.Artists.Artists))
	for i, artist := range res.Artists.Artists {
		artists[i] = &Artist{
			Name: artist.Name,
			ID:   artist.ID,
		}
	}
	return artists, nil
}

// Title belongs to bubbles list.Item
func (a *Artist) Title() string {
	return a.Name
}

// Description belongs to bubbles  list.Item
func (a *Artist) Description() string {
	return ""
}

// FilterValue returns the value we filter by
func (a *Artist) FilterValue() string {
	return a.Name
}

var _ list.DefaultItem = (*Artist)(nil)
