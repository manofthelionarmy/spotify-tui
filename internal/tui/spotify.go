package tui

import (
	"fmt"
	"os"
	"spotify-tui/internal/auth/pcke"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type spotifyModel struct {
	choice        selection
	choices       []string
	cursor        int
	artistsLayout tea.Model
	width         int
	height        int
	spotifyClient *spotify.Client
}

type selection int

var (
	// MAINMENU represents our main menu
	MAINMENU = selection(0)
	// ARTISTS represents artists
	ARTISTS = selection(1)
	// PLAYLISTS represents playlists
	PLAYLISTS = selection(2)
	// SONGS represents songs
	SONGS = selection(3)
)

// New returns a bubbletea Model
func New() tea.Model {
	// use this to login to spotify
	client, err := pcke.Auth()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &spotifyModel{
		choices:       []string{"Artists", "Playlists", "Tracks"},
		cursor:        0,
		choice:        0,
		spotifyClient: client,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *spotifyModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *spotifyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var selectedModel tea.Model = m
	switch msg := msg.(type) {
	// if it is a key press
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			// I need to load a sub-model, I think this is what I need to do
		case "enter":
			// because the mainmenu is ouor default menu, we default to this
			m.choice = selection(m.cursor) + 1

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
		switch m.choice {
		case MAINMENU:
			selectedModel = m
		case ARTISTS:
			artistLayout := NewSearchArtistModel(m.spotifyClient, m.width, m.height)
			selectedModel = artistLayout
		case PLAYLISTS:
			selectedModel = m
		case SONGS:
			selectedModel = m
		}
	}
	return selectedModel, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *spotifyModel) View() string {
	s := "Choices\n\n"
	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	// The footer
	s += "\nPress q to quit.\n"
	return s
}
