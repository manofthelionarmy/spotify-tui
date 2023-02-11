package tui

import (
	"log"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// this is a list model that displays songs belonging to an artist
type artistsSongsModel struct {
	songs list.Model
}

func newArtistsSongsModel(client *spotify.Client,
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

	return artistsSongsModel{trackList}
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
	m.songs, cmd = m.songs.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m artistsSongsModel) View() string {
	return m.songs.View()
}
