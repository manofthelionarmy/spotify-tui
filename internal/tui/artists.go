package tui

import (
	"fmt"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// this is a composite model for a search input and a list display
type searchArtistModel struct {
	textInput      textinput.Model
	artistList     list.Model
	artists        []*models.Artist
	err            error
	spotifyClient  *spotify.Client
	windowWidth    int
	windowHeight   int
	selectedArtist bool
}

// NewSearchArtistModel returns a new artists search layout
func NewSearchArtistModel(cli *spotify.Client, windowWidth, windowHeight int) tea.Model {
	textInput := textinput.New()
	textInput.Placeholder = "Enter Artists Name"
	textInput.Focus()
	textInput.CharLimit = 156
	textInput.TextStyle.Height(10)
	textInput.Width = 20

	return &searchArtistModel{
		textInput:      textInput,
		spotifyClient:  cli,
		windowWidth:    windowWidth,
		windowHeight:   windowHeight,
		selectedArtist: false,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *searchArtistModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *searchArtistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// I think it's better to update state in this switch block
	// then do stuff after we exit it
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if !m.textInput.Focused() {
				m.selectedArtist = false
				m.textInput.Focus()
				m.textInput.Reset()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			// if it isn't focused that means we've selected an Artist
			if !m.textInput.Focused() {
				m.selectedArtist = true
			} else {
				m.textInput.Blur()
			}
		}

	// We handle errors just like any other message
	// See the docs in how to properly handle errors
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	if !m.textInput.Focused() && !m.selectedArtist {
		query := m.textInput.Value()
		// how do we know if it's empty?
		// I need to set the focus to my results view
		// oh dang, how am I gonna pass the spotify client among these?
		resultsModel, _ := newResults(m.spotifyClient, query,
			m.windowWidth, m.windowHeight-10)

		m.artistList = resultsModel.list
		// but each update returns a new one, so would it hurt to create a new one?
		m.artistList, cmd = m.artistList.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.selectedArtist {
		selectedItem := m.artistList.SelectedItem()
		artist, _ := selectedItem.(*models.Artist)
		songsList := newArtistsSongsModel(m.spotifyClient, m,
			artist.ID, m.windowWidth, m.windowHeight)

		return songsList, nil
	}
	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *searchArtistModel) View() string {
	return fmt.Sprintf(
		"Which artist have you been vibin' to?\n\n%s\n\n",
		m.textInput.View(),
	) + "\n" +
		m.artistList.View() +
		"\n(esc to quit)"
}

type artistSearchResultsModel struct {
	list list.Model
}

func newResults(cli *spotify.Client, query string, width, height int) (*artistSearchResultsModel, error) {
	artists, err := models.GetArtists(cli, query)
	if err != nil {
		return nil, err
	}
	items := make([]list.Item, len(artists))
	for i := range items {
		items[i] = artists[i]
	}
	list := list.New(items, list.DefaultDelegate{}, width, height)
	list.Title = "Artist Results:"
	return &artistSearchResultsModel{
		list: list,
	}, nil
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *artistSearchResultsModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *artistSearchResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// I was forgetting to call update
	// I can then do additional logic in here while wrapping the original
	// The only downside is overlap of keys
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *artistSearchResultsModel) View() string {
	return m.list.View()
}
