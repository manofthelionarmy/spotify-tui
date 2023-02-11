package models

import tea "github.com/charmbracelet/bubbletea"

// TeaTree is a model tree interface
type TeaTree interface {
	Parent() TeaTree
	Children() []TeaTree
	Child(int) TeaTree
	AddChild(tea.Model)
	Value() tea.Model
}
