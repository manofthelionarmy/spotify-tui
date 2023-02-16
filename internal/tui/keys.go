package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap represents my key bindings
type KeyMap struct {
	GoBack                   key.Binding
	ClearSearch              key.Binding
	SelectedArtist           key.Binding
	SelectedSong             key.Binding
	SelectedAlbumOrTopTracks key.Binding
	SubmitSearch             key.Binding
	ForceQuit                key.Binding
	Quit                     key.Binding
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
		SelectedAlbumOrTopTracks: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedArtist: key.NewBinding(
			key.WithKeys("enter"),
		),
		SelectedSong: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}
