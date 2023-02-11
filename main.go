package main

import (
	"fmt"
	"os"
	"spotify-tui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
