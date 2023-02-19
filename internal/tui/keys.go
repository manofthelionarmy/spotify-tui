package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap represents my key bindings
type KeyMap struct {
	GoBack                 key.Binding
	ClearSearch            key.Binding
	SelectedArtist         key.Binding
	SelectAblum            key.Binding
	SelectedSong           key.Binding // TODO: maintain only one selecting key
	SelectedFromArtistMenu key.Binding
	SelectedFromMainMenu   key.Binding
	SubmitSearch           key.Binding
	ForceQuit              key.Binding
	Quit                   key.Binding
	// TODO: add a keybinding to easily switch between search, artists, albums, playlists search (all have their own flows)
}

// AppKeyMap returns the apps key mapping
func AppKeyMap() KeyMap {
	return KeyMap{
		ClearSearch: key.NewBinding(
			key.WithKeys("esc"),
		),
		GoBack: key.NewBinding(
			key.WithKeys("esc"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
		),
		SubmitSearch: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedFromArtistMenu: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedFromMainMenu: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedArtist: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedSong: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectAblum: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}
