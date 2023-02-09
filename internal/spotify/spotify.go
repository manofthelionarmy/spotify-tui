package spotify

import (
	"context"
	"fmt"
	"net/url"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	spotify "github.com/zmb3/spotify/v2"
)

type model struct {
	choices       []string
	songs         []string
	artists       []string
	spotifyClient *spotify.Client
	cursor        int
	selected      int
}

// New returns a bubbletea Model
func New() tea.Model {
	// use this to login to spotify
	client, err := Auth()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &model{
		choices:       []string{"Artists", "Playlists"},
		songs:         []string{},
		spotifyClient: client,
		artists:       make([]string, 0),
		cursor:        0,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *model) Init() tea.Cmd {
	return func() tea.Msg {
		query := url.QueryEscape("artist=playboi carti")
		res, err := m.spotifyClient.Search(context.Background(), query, spotify.SearchTypeArtist)
		if err != nil {
			return err.Error()
		}
		for _, artist := range res.Artists.Artists {
			m.artists = append(m.artists, artist.Name)
		}
		return 0
	}
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// if it is a key press
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.artists)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *model) View() string {

	// query := url.QueryEscape("artist=playboi carti")
	// res, err := m.spotifyClient.Search(context.Background(), query, spotify.SearchTypeArtist)
	// if err != nil {
	// 	return fmt.Sprintf("%+v", err)
	// }

	s := "Artists\n\n"
	for i, artist := range m.artists {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		s += fmt.Sprintf("%s %s\n", cursor, artist)
	}

	// The fotter
	s += "\nPress q to quit.\n"

	return s
}
