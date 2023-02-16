package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type albumTracks struct {
	choices  []string
	cursor   int
	selected bool
}

// NewAlbumTracks returns a new menu with choices to select albums or top tracks
func NewAlbumTracks() tea.Model {
	return &albumTracks{
		choices: []string{"Albums", "Top Tracks"},
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *albumTracks) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *albumTracks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *albumTracks) View() string {
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

func (m *albumTracks) SelectedItem() string {
	return m.choices[m.cursor]
}
