package tui

import (
	"context"
	"fmt"
	"os"
	"spotify-tui/internal/auth/pcke"
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/key"
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

type pickFromChoices struct {
	albumTracks    tea.Model
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
	selectingAlbumsOrTopTracks
	browsingArtists
	browsingSongs
	browsingAlbums
)

type composite struct {
	keyMap KeyMap
	searchPrompt
	displayArtists
	pickFromChoices
	displayAlbums
	displaySongs
	appState
	artistID      spotify.ID
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
	artistList.KeyMap.Quit.Unbind()

	songsList := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 0)
	songsList.Title = "Songs..."
	songsList.KeyMap.Quit.Unbind()

	albumList := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 0)
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
		pickFromChoices: pickFromChoices{
			albumTracks:    NewAlbumTracks(),
			selectedChoice: false,
		},
		displayAlbums: displayAlbums{
			selectedAlbum: false,
			albums:        nil,
			list:          albumList,
		},
		keyMap:   AppKeyMap(),
		appState: searchingArtists,
	}
	composite.updateKeyBindings()
	return composite
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

// SpotifyAlbumsResponse is the alumb response we sent as a message
type SpotifyAlbumsResponse []*models.Album

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
	// TODO: study bubble tea lists to see how they handle this stuff more cleanly
	case SpotifyAlbumsResponse:
		m.handleAlbumsReponse(SpotifyAlbumsResponse(msg))
	case SpotifySearchSongsRespMsg:
		m.handleSearchSongsResponse(SongsResponse(msg))
	case SpotifySearchArtistsMsg:
		// handle the message being sent with the retrieved response sent from the tea.Cmd
		m.handleSearchArtistResponse(ArtistsResponse(msg))
	case tea.WindowSizeMsg:
		m.displayArtists.list.SetSize(msg.Width, msg.Height-10)
		m.displaySongs.list.SetSize(msg.Width, msg.Height)
		m.displayAlbums.list.SetSize(msg.Width, msg.Height)
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
	} else if m.appState == selectingAlbumsOrTopTracks {
		cmds = append(cmds, m.handleSelectingAlbumsOrTopTracks(msg))
	} else if m.appState == browsingAlbums {
		var cmd tea.Cmd
		m.displayAlbums.list, cmd = m.displayAlbums.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *composite) View() string {
	if m.appState == browsingAlbums {
		// why isn't this displaying?
		return m.displayAlbums.list.View()
	} else if m.appState == selectingAlbumsOrTopTracks {
		return m.albumTracks.View()
	} else if m.appState != browsingSongs {
		return m.searchPrompt.View() + "\n" +
			m.displayArtists.list.View()
	}
	return m.displaySongs.list.View()
}

func (m *composite) getArtists(query string) tea.Cmd {
	return func() tea.Msg {
		artists, _ := models.GetArtists(m.spotifyClient, query)
		return SpotifySearchArtistsMsg(artists)
	}
}

func (m *composite) handleSelectedArtist() tea.Cmd {
	// return songs
	return func() tea.Msg {
		// set this to true
		// m.displayArtists.selectedArtist = false

		// get the selected item
		seletectedItem := m.displayArtists.list.SelectedItem()
		artist := seletectedItem.(*models.Artist)

		// search the songs belonging to this artist
		songs, _ := models.GetSongs(m.spotifyClient, artist.ID)
		return SpotifySearchSongsRespMsg(songs)
	}
}

// It looks like we are currently searching for artists
func (m *composite) handleSearchingForArtists() {
	// should we send this current value as a msg?
	m.searchPrompt.searching = true
	m.searchPrompt.textInput.Blur()
}

// TODO: handle error and send it as a message, or send the selected song as a message and handle it
// I think that's better
func (m *composite) handleSelectedSong() {
	// m.displaySongs.selectedSong = true
	selectedItem := m.displaySongs.list.SelectedItem()
	song := selectedItem.(*models.Song)

	// bug, if it's paused and try to get the devices, it will return none
	// I need to check the state if it's paused or played
	devices, _ := m.spotifyClient.PlayerDevices(context.Background())
	m.spotifyClient.PlayOpt(context.Background(),
		&spotify.PlayOptions{
			URIs:     []spotify.URI{song.SongURI},
			DeviceID: &devices[0].ID,
		},
	)
}

func (m *composite) resetDisplayArtistList() {
	// remove all of the items
	for len(m.displayArtists.list.Items()) > 0 {
		m.displayArtists.list.RemoveItem(0)
	}
	// remove all of the artists
	for len(m.displayArtists.artists) > 0 {
		m.displayArtists.artists = m.displayArtists.artists[:len(m.displayArtists.artists)-1]
	}
}

func (m *composite) resetSearchPrompt() {
	// the esc key is resolved for list.Model... ugggghh
	m.searchPrompt.searching = false
	m.searchPrompt.textInput.Reset()
	m.searchPrompt.textInput.Focus()
}

func (m *composite) resetSongsList() {
	// remove all of the items
	for len(m.displayArtists.list.Items()) > 0 {
		m.displaySongs.list.RemoveItem(0)
	}

	// remove all of the artists
	for len(m.displaySongs.songs) > 0 {
		m.displaySongs.songs = m.displaySongs.songs[:len(m.displaySongs.songs)-1]
	}
}

func (m *composite) populateSongsList(songs []*models.Song) {
	m.displaySongs.songs = songs
	items := make([]list.Item, len(m.displaySongs.songs))
	for i := range items {
		items[i] = m.displaySongs.songs[i]
	}
	m.displaySongs.list.SetItems(items)
}

func (m *composite) populateArtists(artists []*models.Artist) {
	m.displayArtists.artists = artists
	items := make([]list.Item, len(m.displayArtists.artists))
	for i := range items {
		items[i] = m.displayArtists.artists[i]
	}
	m.displayArtists.list.SetItems(items)
}

func (m *composite) handleSearchArtistResponse(artists []*models.Artist) {
	m.populateArtists(artists)
}

func (m *composite) handleSearchSongsResponse(songs []*models.Song) {
	m.populateSongsList(songs)
}

func (m *composite) handleAlbumsReponse(albums []*models.Album) {
	m.displayAlbums.albums = albums
	items := make([]list.Item, len(m.displayAlbums.albums))
	for i := range items {
		items[i] = m.displayAlbums.albums[i]
	}
	m.displayAlbums.list.SetItems(items)
}

func (m *composite) handleSearchingArtists(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keyMap.SubmitSearch):
			m.searchPrompt.textInput.Blur()
			m.appState = browsingArtists
			m.searching = true
			m.updateKeyBindings()
		}
	}
	newTextInputInput, cmd := m.searchPrompt.textInput.Update(msg)
	m.searchPrompt.textInput = newTextInputInput
	cmds = append(cmds, cmd)

	if m.searching == true {
		cmds = append(cmds, m.getArtists(m.searchPrompt.textInput.Value()))
	}
	return tea.Batch(cmds...)
}

func (m *composite) handleBrowsingArtists(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ClearSearch) {
			// a bit weird but we're saying we want to clear the search
			m.searching = false
			m.appState = searchingArtists
			m.updateKeyBindings()
			m.searchPrompt.textInput.Focus()
			m.searchPrompt.textInput.Reset()
			// I need to reset the list
			return nil
		}
		if key.Matches(msg, m.keyMap.SelectedArtist) {
			m.selectedArtist = true
			// m.appState = browsingSongs
			m.appState = selectingAlbumsOrTopTracks
			artist, _ := m.displayArtists.list.SelectedItem().(*models.Artist)
			m.artistID = artist.ID
			m.updateKeyBindings()
		}
	}
	m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
	cmds = append(cmds, cmd)
	if m.selectedArtist {
		cmds = append(cmds, m.handleSelectedArtist())
	}
	return tea.Batch(cmds...)
}

func (m *composite) handleBrowsingSongs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// I can figure out how to do go back, this is feature creep, but overall I have a cleaner way of
		// coding this project
		if key.Matches(msg, m.keyMap.GoBack) {
			// a bit weird but we're saying we want to clear the search
			m.searching = false
			m.searchPrompt.textInput.Focus()
			m.searchPrompt.textInput.Reset()
			m.selectedArtist = false
			m.appState = searchingArtists
			m.updateKeyBindings()
			m.resetDisplayArtistList()
			m.resetSongsList()

			m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
			cmds = append(cmds, cmd)
			m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
			cmds = append(cmds, cmd)
		} else if key.Matches(msg, m.keyMap.SelectedSong) {
			m.handleSelectedSong()
		}
	}
	// something weird is happening, this doesn't go away
	// this is taking a while to update
	m.displaySongs.list, cmd = m.displaySongs.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *composite) updateKeyBindings() {
	switch m.appState {
	case searchingArtists:
		m.keyMap.SubmitSearch.SetEnabled(true)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedAlbumOrTopTracks.SetEnabled(false)
	case browsingArtists:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(true)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedAlbumOrTopTracks.SetEnabled(false)
	case browsingSongs:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedAlbumOrTopTracks.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(true)
	case selectingAlbumsOrTopTracks:
		m.keyMap.SubmitSearch.SetEnabled(false)
		m.keyMap.SelectedArtist.SetEnabled(false)
		m.keyMap.SelectedSong.SetEnabled(false)
		m.keyMap.SelectedAlbumOrTopTracks.SetEnabled(true)
	}
}

func (m *composite) handleSelectingAlbumsOrTopTracks(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// I can figure out how to do go back, this is feature creep, but overall I have a cleaner way of
		// coding this project
		if key.Matches(msg, m.keyMap.GoBack) {
			// a bit weird but we're saying we want to clear the search
			m.searching = false
			m.searchPrompt.textInput.Focus()
			m.searchPrompt.textInput.Reset()
			m.selectedArtist = false
			m.selectedChoice = false
			m.appState = searchingArtists
			m.artistID = ""
			m.updateKeyBindings()
			m.resetDisplayArtistList()
			m.resetSongsList()

			m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
			cmds = append(cmds, cmd)
			m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
			cmds = append(cmds, cmd)
		} else if key.Matches(msg, m.keyMap.SelectedAlbumOrTopTracks) {
			m.selectedChoice = true
			m.appState = browsingAlbums
		}
	}

	m.pickFromChoices.albumTracks, cmd = m.albumTracks.Update(msg)
	cmds = append(cmds, cmd)
	if m.selectedChoice {
		albumTrcks, _ := m.albumTracks.(*albumTracks)
		switch albumTrcks.SelectedItem() {
		case "Albums":
			cmds = append(cmds, m.handleSearchAlbums())
		case "Top Tracks":
		}
	}
	return tea.Batch(cmds...)
}

func (m *composite) handleSearchAlbums() tea.Cmd {
	return func() tea.Msg {
		albums, _ := models.GetAlbums(m.spotifyClient, m.artistID)
		return SpotifyAlbumsResponse(albums)
	}
}
