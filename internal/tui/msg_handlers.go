package tui

import (
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
)

func (m *composite) handleSearchArtistResponse(artists []*models.Artist) {
	// TODO: do we really need the artists here? it will be maintained in the list model.
	m.displayArtists.artists = artists
	items := make([]list.Item, len(m.displayArtists.artists))
	for i := range items {
		items[i] = m.displayArtists.artists[i]
	}
	m.displayArtists.list.SetItems(items)
}

func (m *composite) handleSearchSongsResponse(songs []*models.Song) {
	items := make([]list.Item, len(songs))
	for i := range items {
		items[i] = songs[i]
	}
	m.displaySongs.list.SetItems(items)
}

func (m *composite) handleAlbumsReponse(albums []*models.Album) {
	items := make([]list.Item, len(albums))
	for i := range items {
		items[i] = albums[i]
	}
	m.displayAlbums.list.SetItems(items)
}
