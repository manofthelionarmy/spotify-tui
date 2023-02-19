package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// these are the choices available in these menus
const (
	search    = "Search"
	artists   = "Artists"
	albums    = "Albums"
	songs     = "Songs"
	playlists = "Playlists"
	topTracks = "Top Tracks"
)

type menu struct {
	choices  []string
	cursor   int
	selected bool
}

// NewArtistMenu returns a new menu with choices to select albums or top tracks
func NewArtistMenu() tea.Model {
	return &menu{
		choices: []string{albums, topTracks},
	}
}

// NewMainMenu returns the main menu to select search, artists, albums, etc
func NewMainMenu() tea.Model {
	return &menu{
		choices: []string{search, artists, albums, songs, playlists},
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *menu) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}

	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *menu) View() string {
	s := ""
	for i := range m.choices {
		var selected string
		if i == m.cursor {
			selected = ">"
		} else {
			selected = " "
		}
		s += fmt.Sprintf("%s %s\n", selected, m.choices[i])
	}
	return s
}

func (m *menu) SelectedItem() string {
	return m.choices[m.cursor]
}
