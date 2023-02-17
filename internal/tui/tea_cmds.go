package tui

import (
	"context"
	"spotify-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

func (m *composite) getArtists(query string) tea.Cmd {
	return func() tea.Msg {
		artists, _ := models.GetArtists(m.spotifyClient, query)
		return SpotifySearchArtistsMsg(artists)
	}
}

func (m *composite) handleSelectedArtist() tea.Cmd {
	// return songs
	return func() tea.Msg {
		// set this to true
		// m.displayArtists.selectedArtist = false

		// get the selected item
		seletectedItem := m.displayArtists.list.SelectedItem()
		artist := seletectedItem.(*models.Artist)

		// search the songs belonging to this artist
		songs, _ := models.GetSongs(m.spotifyClient, artist.ID)
		return SpotifySearchSongsRespMsg(songs)
	}
}

// TODO: handle error and send it as a message, or send the selected song as a message and handle it, return an error as its message
// I think that's better
func (m *composite) handleSelectedSong() {
	// m.displaySongs.selectedSong = true
	selectedItem := m.displaySongs.list.SelectedItem()
	song := selectedItem.(*models.Song)

	// bug, if it's paused and try to get the devices, it will return none
	// I need to check the state if it's paused or played
	devices, _ := m.spotifyClient.PlayerDevices(context.Background())
	m.spotifyClient.PlayOpt(context.Background(),
		&spotify.PlayOptions{
			URIs:     []spotify.URI{song.SongURI},
			DeviceID: &devices[0].ID,
		},
	)
}

func (m *composite) handleSearchAlbums() tea.Cmd {
	return func() tea.Msg {
		albums, _ := models.GetAlbums(m.spotifyClient, m.artistID)
		return SpotifyAlbumsResponse(albums)
	}
}
