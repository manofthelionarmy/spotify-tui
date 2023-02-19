package tui

import (
	"fmt"
	"os"
	"spotify-tui/internal/auth/pcke"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var listStyle = lipgloss.NewStyle().Margin(0, 2).
	Width(25).
	BorderStyle(lipgloss.NormalBorder())
var searchStyle = lipgloss.NewStyle().Margin(0, 2).
	Width(25).
	BorderStyle(lipgloss.NormalBorder())

type searchPrompt struct {
	searching bool
	textInput textinput.Model
}

type displayArtists struct {
	selectedArtist bool
	artists        []*models.Artist
	list           list.Model
}

type displayArtistMenu struct {
	artistMenu     tea.Model
	selectedChoice bool
}

type displayMainMenu struct {
	mainMenu       tea.Model
	selectedChoice bool
}

type displayAlbums struct {
	selectedAlbum bool
	albums        []*models.Album
	list          list.Model
}

type displaySongs struct {
	renderSongs  bool
	selectedSong bool
	songs        []*models.Song
	list         list.Model
}

// searching for artist, browsing artist, browsing song
type appState int

const (
	searchingArtists appState = iota
	selectingFromArtistMenu
	browsingArtists
	browsingSongs
	browsingAlbums
	selectFromMainMenu
)

type composite struct {
	keyMap KeyMap
	displayMainMenu
	searchPrompt
	displayArtists
	displayArtistMenu
	displayAlbums
	displaySongs
	appState
	artistID      spotify.ID
	albumID       spotify.ID
	albumURI      spotify.URI
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
	textInput.PromptStyle.Border(lipgloss.NormalBorder())

	artistList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	artistList.Title = "Artists..."
	artistList.KeyMap.Quit.Unbind()

	songsList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	songsList.Title = "Songs..."
	songsList.KeyMap.Quit.Unbind()

	albumList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	albumList.Title = "Albums..."
	albumList.KeyMap.Quit.Unbind()

	composite := &composite{
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
		displayArtistMenu: displayArtistMenu{
			artistMenu:     NewArtistMenu(),
			selectedChoice: false,
		},
		displayAlbums: displayAlbums{
			selectedAlbum: false,
			albums:        nil,
			list:          albumList,
		},
		displayMainMenu: displayMainMenu{
			mainMenu:       NewMainMenu(),
			selectedChoice: false,
		},
		keyMap:   AppKeyMap(),
		appState: selectFromMainMenu,
	}
	composite.updateKeyBindings()
	return composite
}

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
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case SpotifyAlbumsResponse:
		m.handleAlbumsReponse(SpotifyAlbumsResponse(msg))
	case SpotifySearchSongsRespMsg:
		m.handleSearchSongsResponse(SongsResponse(msg))
	case SpotifySearchArtistsMsg:
		m.handleSearchArtistResponse(ArtistsResponse(msg))
	case tea.WindowSizeMsg:
		h, v := searchStyle.GetFrameSize()
		// m.list.SetSize(msg.Width-h, msg.Height-v)
		// for some reason, this was pushing this up
		m.displayArtists.list.SetSize(msg.Width-h, msg.Height-5-v)
		m.displaySongs.list.SetSize(msg.Width-h, msg.Height-v)
		m.displayAlbums.list.SetSize(msg.Width-h, msg.Height-v) // there is rendering issues when we do max height
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ForceQuit) {
			return m, tea.Quit
		}
	}

	if m.appState == searchingArtists {
		cmds = append(cmds, m.handleSearchingArtists(msg))
	} else if m.appState == browsingArtists {
		cmds = append(cmds, m.handleBrowsingArtists(msg))
	} else if m.appState == browsingSongs {
		cmds = append(cmds, m.handleBrowsingSongs(msg))
	} else if m.appState == selectingFromArtistMenu {
		cmds = append(cmds, m.handleSelectingAlbumsOrTopTracks(msg))
	} else if m.appState == browsingAlbums {
		cmds = append(cmds, m.handleSelectAlbum(msg))
	} else if m.appState == selectFromMainMenu {
		cmds = append(cmds, m.handleSelectFromMainMenu(msg))
	}

	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *composite) View() string {
	var view string
	switch m.appState {
	case browsingAlbums:
		view = m.displayAlbums.list.View()
	case selectingFromArtistMenu:
		view = m.artistMenu.View()
	case searchingArtists, browsingArtists:
		view = lipgloss.JoinVertical(lipgloss.Top, searchStyle.Render(m.searchPrompt.View()),
			listStyle.Render(m.displayArtists.list.View()))
	case browsingSongs:
		view = m.displaySongs.list.View()
	case selectFromMainMenu:
		view = m.displayMainMenu.mainMenu.View()
	}
	return view
}

// It looks like we are currently searching for artists
func (m *composite) handleSearchingForArtists() {
	// should we send this current value as a msg?
	m.searchPrompt.searching = true
	m.searchPrompt.textInput.Blur()
}

func (m *composite) resetDisplayArtist() {
	m.selectedArtist = false
	// remove all of the items
	for len(m.displayArtists.list.Items()) > 0 {
		m.displayArtists.list.RemoveItem(0)
	}
	// remove all of the artists
	for len(m.displayArtists.artists) > 0 {
		m.displayArtists.artists = m.displayArtists.artists[:len(m.displayArtists.artists)-1]
	}
}

func (m *composite) setAppState(state appState) {
	m.appState = state
}

func (m *composite) resetPickFromChoices() {
	m.displayArtistMenu.selectedChoice = false
}

func (m *composite) resetSearchPrompt() {
	m.searchPrompt.searching = false
	m.searchPrompt.textInput.Reset()
	m.searchPrompt.textInput.Focus()
}

// I should be making tests for these kind of things
func (m *composite) resetSongsList() {
	m.displaySongs.list.SetItems([]list.Item{})
}

func (m *composite) resetAlbumsList() {
	// I'm not sure what the consequences of this are, will go's garbage collector take care of it?
	m.displayAlbums.list.SetItems([]list.Item{})
}

func (m *composite) updateKeyBindings() {
	switch m.appState {
	case searchingArtists:
		m.keyMap.SubmitSearch.SetEnabled(true)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(false)
		m.keyMap.SelectedFromMainMenu.SetEnabled(false)
		m.keyMap.SelectAblum.SetEnabled(false)
	case browsingArtists:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(true)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(false)
		m.keyMap.SelectedFromMainMenu.SetEnabled(false)
		m.keyMap.SelectAblum.SetEnabled(false)
	case browsingSongs:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(false)
		m.keyMap.SelectedFromMainMenu.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(true)
		m.keyMap.SelectAblum.SetEnabled(false)
	case selectingFromArtistMenu:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(true)
		m.keyMap.SelectedFromMainMenu.SetEnabled(false)
		m.keyMap.SelectAblum.SetEnabled(false)
	case browsingAlbums:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(false)
		m.keyMap.SelectedFromMainMenu.SetEnabled(false)
		m.keyMap.SelectAblum.SetEnabled(true)
	case selectFromMainMenu:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedFromArtistMenu.SetEnabled(false)
		m.keyMap.SelectedFromMainMenu.SetEnabled(true)
		m.keyMap.SelectAblum.SetEnabled(false)
	}
}

func (m *composite) returnNewAppStateFromMainMenuChoice() appState {
	menu, _ := m.displayMainMenu.mainMenu.(*menu)
	switch menu.choices[menu.cursor] {
	case search:
		return searchingArtists
	case artists:
		// TODO: add the other states we'll handle based on the flows
	case albums:
	case songs:
	case playlists:
	}
	// TODO update
	return searchingArtists
}
