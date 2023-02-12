package tui

import (
	"context"
	"fmt"
	"os"
	"spotify-tui/internal/auth/pcke"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type searchPrompt struct {
	searching bool
	textInput textinput.Model
}

type displayArtists struct {
	selectedArtist bool
	artists        []*models.Artist
	list           list.Model
}

type displaySongs struct {
	renderSongs  bool
	selectedSong bool
	songs        []*models.Song
	list         list.Model
}

type composite struct {
	searchPrompt
	displayArtists
	displaySongs
	spotifyClient *spotify.Client
	width         int
	height        int
}

// NewComposite returns a new composite
func NewComposite() tea.Model {
	client, err := pcke.Auth()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	textInput := textinput.New()
	textInput.Placeholder = "Enter Artists Name"
	textInput.Focus()
	textInput.CharLimit = 156
	textInput.TextStyle.Height(10)
	textInput.Width = 20

	artistList := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 0)
	artistList.Title = "Artists..."

	songsList := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 0)
	songsList.Title = "Songs..."

	return &composite{
		spotifyClient: client,
		searchPrompt: searchPrompt{
			textInput: textInput,
			searching: false,
		},
		displayArtists: displayArtists{
			artists: nil,
			list:    artistList,
		},
		displaySongs: displaySongs{
			songs:        nil,
			list:         songsList,
			renderSongs:  false,
			selectedSong: false,
		},
	}
}

// SpotifySearchArtistsMsg is a message sent
type SpotifySearchArtistsMsg []*models.Artist

// SelectedArtistMsg is message sent when we've selected an artist
type SelectedArtistMsg *models.Artist

// ArtistsResponse is the artists response we sent as a message
type ArtistsResponse []*models.Artist

// SongsResponse is the songs response we sent as a message
type SongsResponse []*models.Song

// SpotifySearchSongsRespMsg is a message signaling we got back an artists songs from spotify api
type SpotifySearchSongsRespMsg []*models.Song

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m searchPrompt) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m searchPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m searchPrompt) View() string {
	return m.textInput.View()
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *composite) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *composite) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// TODO: study bubble tea lists to see how they handle this stuff more cleanly
	case SpotifySearchSongsRespMsg:
		m.displaySongs.songs = SongsResponse(msg)
		items := make([]list.Item, len(m.displaySongs.songs))
		for i := range items {
			items[i] = m.displaySongs.songs[i]
		}
		m.displaySongs.list.SetItems(items)
		m.displaySongs.renderSongs = true
	case SpotifySearchArtistsMsg:
		// handle the message being sent with the retrieved response sent from the command
		m.displayArtists.artists = ArtistsResponse(msg)
		items := make([]list.Item, len(m.displayArtists.artists))
		for i := range items {
			items[i] = m.displayArtists.artists[i]
		}
		m.displayArtists.list.SetItems(items)
		m.searching = false
	case tea.WindowSizeMsg:
		// m.width, m.height = msg.Width, msg.Height
		m.displayArtists.list.SetSize(msg.Width, msg.Height-10)
		m.displaySongs.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			// the esc key is resolved for list.Model... ugggghh
			m.searchPrompt.searching = false
			m.searchPrompt.textInput.Reset()
			m.searchPrompt.textInput.Focus()
			for len(m.displayArtists.list.Items()) > 0 {
				m.displayArtists.list.RemoveItem(0)
			}
			for len(m.displayArtists.artists) > 0 {
				m.displayArtists.artists = m.displayArtists.artists[:len(m.displayArtists.artists)-1]
			}
		case tea.KeyEnter:
			// TODO: handle when we are rendering the displaySongs model
			if m.displaySongs.renderSongs {
				m.handleSelectedSong()
			} else if m.searchPrompt.textInput.Focused() {
				m.handleSearching()
			} else {
				cmd := m.handleSelectedArtist()
				cmds = append(cmds, cmd)
			}
		}
	}

	if m.searchPrompt.searching == true {
		cmd := m.getArtists(m.searchPrompt.textInput.Value())
		cmds = append(cmds, cmd)
	}

	if m.searchPrompt.textInput.Focused() {
		m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if !m.searchPrompt.textInput.Focused() && !m.displaySongs.renderSongs {
		// the esc key will be picked up
		// I only want this to update when the searchPrompt isn't in focus
		m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.displaySongs.renderSongs {
		m.displaySongs.list, cmd = m.displaySongs.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *composite) View() string {
	if !m.displaySongs.renderSongs {
		return m.searchPrompt.View() + "\n" +
			m.displayArtists.list.View()
	}
	return m.displaySongs.list.View()
}

func (m *composite) getArtists(query string) tea.Cmd {

	// return artists
	return func() tea.Msg {
		artists, _ := models.GetArtists(m.spotifyClient, query)
		return SpotifySearchArtistsMsg(artists)
	}
}

// func (m *composite) handleArtistsResponse()

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *displayArtists) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *displayArtists) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	panic("not implemented") // TODO: Implement
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *displayArtists) View() string {
	panic("not implemented") // TODO: Implement
}

func (m *composite) handleSelectedArtist() tea.Cmd {
	// return songs
	return func() tea.Msg {
		// set this to true
		m.displayArtists.selectedArtist = true

		// get the selected item
		seletectedItem := m.displayArtists.list.SelectedItem()
		artist := seletectedItem.(*models.Artist)

		// search the songs belonging to this artist
		songs, _ := models.GetSongs(m.spotifyClient, artist.ID)
		return SpotifySearchSongsRespMsg(songs)
	}
}

func (m *composite) handleSearching() {
	m.searchPrompt.searching = true
	m.searchPrompt.textInput.Blur()
}

// TODO: handle error and send it as a message, or send the selected song as a message and handle it
// I think that's better
func (m *composite) handleSelectedSong() {

	m.displaySongs.selectedSong = true
	selectedItem := m.displaySongs.list.SelectedItem()
	song := selectedItem.(*models.Song)

	devices, _ := m.spotifyClient.PlayerDevices(context.Background())
	m.spotifyClient.PlayOpt(context.Background(),
		&spotify.PlayOptions{
			URIs:     []spotify.URI{song.SongURI},
			DeviceID: &devices[0].ID,
		},
	)
}
