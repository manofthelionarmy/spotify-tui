package tui

import (
	"context"
	"log"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// this is a list model that displays songs belonging to an artist
type artistsSongsModel struct {
	songs         list.Model
	selectedSong  bool
	spotifyClient *spotify.Client
	prevModel     tea.Model
}

func newArtistsSongsModel(client *spotify.Client, prevModel tea.Model,
	artistID spotify.ID, width int, height int) artistsSongsModel {

	songs, err := models.GetSongs(client, artistID)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, len(songs))
	for i := range items {
		items[i] = songs[i]
	}

	trackList := list.New(items, list.DefaultDelegate{}, width, height)
	trackList.Title = "Songs..."

	return artistsSongsModel{
		songs:         trackList,
		selectedSong:  false,
		spotifyClient: client,
		prevModel:     prevModel,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m artistsSongsModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m artistsSongsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.selectedSong = true
		case tea.KeyEsc:
			if m.songs.FilterState() == list.FilterApplied ||
				m.songs.FilterState() == list.Filtering {
				m.songs.FilterInput.Reset()
			} else {
				// is there a more efficient way to do this? will it take up memory?
				// I can store a prev model as part of the state
				// I think I can do a tree kind of thing
				composite := NewSearchArtistModel(m.spotifyClient, m.songs.Width(), m.songs.Height())
				return composite, tea.EnterAltScreen
			}
		}
	}
	if m.selectedSong {
		selectedItem := m.songs.SelectedItem()
		song, _ := selectedItem.(*models.Song)
		devices, _ := m.spotifyClient.PlayerDevices(context.Background())
		// devices is nil, because it was a permissions issue
		playOpt := &spotify.PlayOptions{
			URIs:     []spotify.URI{song.SongURI},
			DeviceID: &devices[0].ID,
		}
		m.spotifyClient.PlayOpt(context.Background(), playOpt)
		m.selectedSong = false
	}
	m.songs, cmd = m.songs.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m artistsSongsModel) View() string {
	return m.songs.View()
}
