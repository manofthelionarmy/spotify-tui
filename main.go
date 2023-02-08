package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func main() {

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func initialModel() model {
	return model{
		choices:  []string{"Buy carrots", "Buy Celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m model) Init() tea.Cmd {
	return nil
}

var keyMap map[string]int = make(map[string]int)

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// if it is a key press

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "G":
			m.cursor = len(m.choices) - 1
		case "g":
			if numTimesPressed, ok := keyMap["g"]; !ok {
				keyMap["g"] = 1
			} else {
				numTimesPressed = numTimesPressed + 1
				keyMap["g"] = numTimesPressed
				// why it's one, idk
				if numTimesPressed == 2 {
					keyMap["g"] = 0
					m.cursor = 0
				}
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m model) View() string {
	// the header
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The fotter
	s += "\nPress q to quit.\n"

	// send the ui for rendering
	return s
}
