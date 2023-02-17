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
			m.resetSearchPrompt()
			m.setAppState(searchingArtists)
			m.updateKeyBindings()
			return nil
		}
		if key.Matches(msg, m.keyMap.SelectedArtist) {
			m.selectedArtist = true
			m.setAppState(selectingAlbumsOrTopTracks)
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
		if key.Matches(msg, m.keyMap.GoBack) {
			m.setAppState(browsingAlbums)
			m.updateKeyBindings()
			m.resetSongsList()
			// reset this because we want to select a new one
			m.albumID = ""
			// reset this because we want to select a new album
			m.selectedAlbum = false
			return nil
		} else if key.Matches(msg, m.keyMap.SelectedSong) {
			m.selectedSong = true
		}
	}

	if m.selectedSong {
		m.handleSelectedSong()
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

			m.searchPrompt.textInput, cmd = m.searchPrompt.textInput.Update(msg)
			cmds = append(cmds, cmd)
			m.displayArtists.list, cmd = m.displayArtists.list.Update(msg)
			cmds = append(cmds, cmd)
		} else if key.Matches(msg, m.keyMap.SelectedAlbumOrTopTracks) {
			m.selectedChoice = true
			m.appState = browsingAlbums
			m.updateKeyBindings()
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

func (m *composite) handleSelectAlbum(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.GoBack) {
			m.resetPickFromChoices()
			m.setAppState(selectingAlbumsOrTopTracks)
			m.updateKeyBindings()

			return nil
		} else if key.Matches(msg, m.keyMap.SelectAblum) {
			m.displayAlbums.selectedAlbum = true
			m.appState = browsingSongs

			album, _ := m.displayAlbums.list.SelectedItem().(*models.Album)
			m.albumID = album.ID
			m.updateKeyBindings()
		}
	}

	if m.displayAlbums.selectedAlbum {
		cmds = append(cmds, m.handleSearchSongsInAlbum())
	}

	m.displayAlbums.list, cmd = m.displayAlbums.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
