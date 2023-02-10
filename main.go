package main

import (
	"fmt"
	"os"
	"spotify-tui/model/spotify"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(spotify.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
