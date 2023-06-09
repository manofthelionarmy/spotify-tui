package tui

import (
	"spotify-tui/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *composite) handleSearchingArtists(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keyMap.SubmitSearch):
			m.searchPrompt.textInput.Blur()
			m.appState = browsingArtists
			m.searching = true
			m.updateKeyBindings()
		case key.Matches(msg, m.keyMap.GoBack):
			m.appState = selectFromMainMenu
			m.searching = false
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
			// TODO: handle logic for filtering
			m.resetSearchPrompt()
			m.setAppState(searchingArtists)
			m.updateKeyBindings()
			return nil
		}
		if key.Matches(msg, m.keyMap.SelectedArtist) {
			m.selectedArtist = true
			m.setAppState(selectingFromArtistMenu)
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
		// FIXME: because I disabled esc, we need a way to undo a filter
		// or I need to handle stuff better, such as checking if we are filtering and then go back
		// when we are not filtering
		if key.Matches(msg, m.keyMap.GoBack) {
			m.setAppState(browsingAlbums)
			m.updateKeyBindings()
			m.resetSongsList()
			// reset this because we want to select a new one
			m.albumID = ""
			m.albumURI = ""
			// reset this because we want to select a new album
			m.selectedAlbum = false
		} else if key.Matches(msg, m.keyMap.SelectedSong) {
			// TODO: how to pass the album uri?
			m.handleSelectedSong(m.albumURI)
		}
	}

	// something weird is happening, this doesn't go away
	// this is taking a while to update
	m.displaySongs.list, cmd = m.displaySongs.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *composite) handleSelectingAlbumsOrTopTracks(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// I can figure out how to do go back, this is feature creep, but overall I have a cleaner way of
		// coding this project
		if key.Matches(msg, m.keyMap.GoBack) {
			// Should we reinitialize to a bunch of new components instead?
			m.resetSearchPrompt()
			m.resetDisplayArtist()
			m.setAppState(searchingArtists)
			m.artistID = ""
			m.updateKeyBindings()

			// update these on a reset...?
			m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
			cmds = append(cmds, cmd)
			m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
			cmds = append(cmds, cmd)
		} else if key.Matches(msg, m.keyMap.SelectedFromArtistMenu) {
			m.displayArtistMenu.selectedChoice = true
			m.appState = browsingAlbums
			m.updateKeyBindings()
		}
	}
	if m.displayArtistMenu.selectedChoice {
		albumTrcks, _ := m.artistMenu.(*menu)
		switch albumTrcks.SelectedItem() {
		case "Albums":
			cmds = append(cmds, m.handleSearchAlbums())
		case "Top Tracks":
		}
	}

	m.displayArtistMenu.artistMenu, cmd = m.artistMenu.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *composite) handleSelectAlbum(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.GoBack) {
			m.resetPickFromChoices()
			m.setAppState(selectingFromArtistMenu)
			m.updateKeyBindings()

			return nil
		} else if key.Matches(msg, m.keyMap.SelectAblum) {
			m.displayAlbums.selectedAlbum = true
			m.appState = browsingSongs

			album, _ := m.displayAlbums.list.SelectedItem().(*models.Album)
			m.albumID = album.ID
			m.albumURI = album.AlbumURI
			m.updateKeyBindings()
		}
	}

	if m.displayAlbums.selectedAlbum {
		cmds = append(cmds, m.handleSearchSongsInAlbum(m.albumID))
	}

	m.displayAlbums.list, cmd = m.displayAlbums.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *composite) handleSelectFromMainMenu(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.SelectedFromMainMenu) {
			m.displayMainMenu.selectedChoice = true
		}
	}
	if m.displayMainMenu.selectedChoice {
		m.appState = m.returnNewAppStateFromMainMenuChoice()
		m.displayMainMenu.selectedChoice = false
		m.updateKeyBindings()
	}
	m.displayMainMenu.mainMenu, cmd = m.displayMainMenu.mainMenu.Update(msg)
	return cmd
}
