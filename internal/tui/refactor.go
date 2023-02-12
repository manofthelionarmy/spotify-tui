package tui

import (
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

type artistList struct {
	artists    []*models.Artist
	artistList list.Model
}

type composite struct {
	searchPrompt
	artistList
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

	listModel := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 0)

	return &composite{
		spotifyClient: client,
		searchPrompt: searchPrompt{
			textInput: textInput,
			searching: false,
		},
		artistList: artistList{
			artists:    nil,
			artistList: listModel,
		},
	}
}

// SpotifySearchArtistsMsg is a message sent
type SpotifySearchArtistsMsg []*models.Artist

// ArtistsResponse is this
type ArtistsResponse []*models.Artist

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
	case SpotifySearchArtistsMsg:
		// handle the message being sent with the retrieved response sent from the command
		m.artistList.artists = ArtistsResponse(msg)
		items := make([]list.Item, len(m.artistList.artists))
		for i := range items {
			items[i] = m.artistList.artists[i]
		}
		m.artistList.artistList.SetItems(items)
		m.searching = false
	case tea.WindowSizeMsg:
		// m.width, m.height = msg.Width, msg.Height
		m.artistList.artistList.SetSize(msg.Width, msg.Height-10)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			// the esc key is resolved for list.Model... ugggghh
			m.searchPrompt.searching = false
			m.searchPrompt.textInput.Reset()
			m.searchPrompt.textInput.Focus()
			for len(m.artistList.artistList.Items()) > 0 {
				m.artistList.artistList.RemoveItem(0)
			}
			for len(m.artistList.artists) > 0 {
				m.artistList.artists = m.artistList.artists[:len(m.artistList.artists)-1]
			}
		case tea.KeyEnter:
			m.searchPrompt.searching = true
			m.searchPrompt.textInput.Blur()
		}
	}

	if m.searchPrompt.searching == true {
		cmd := m.getArtists(m.searchPrompt.textInput.Value())
		cmds = append(cmds, cmd)
	}

	if m.searchPrompt.textInput.Focused() {
		m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// the esc key will be picked up
		// I only want this to update when the searchPrompt isn't in focus
		m.artistList.artistList, cmd = m.artistList.artistList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *composite) View() string {
	return m.searchPrompt.View() + m.artistList.artistList.View()
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
func (m *artistList) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *artistList) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	panic("not implemented") // TODO: Implement
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *artistList) View() string {
	panic("not implemented") // TODO: Implement
}
