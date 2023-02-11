package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
)

func TestNewCassava(t *testing.T) {
	rootModel := testModel{}
	cassava := New(rootModel)
	require.NotNil(t, cassava.tree)
	require.Nil(t, cassava.Parent())
	require.NotNil(t, cassava.Children())
	require.Len(t, cassava.Children(), 0)
}

func TestAddChild(t *testing.T) {
	rootModel := testModel{}
	cassava := New(rootModel)

	// Add a child
	cassava.AddChild(testModel{})
	require.Len(t, cassava.Children(), 1)
	require.NotNil(t, cassava.Child(0))
	// FIXME: flakey test/not a good one because the content can be the same, and we're not using pointers
	require.Equal(t, rootModel, cassava.Child(0).Value())
	// Add sub children
	cassava.Child(0).AddChild(testModel{})
	cassava.Child(0).AddChild(testModel{})
	cassava.Child(0).AddChild(testModel{})
	require.Len(t, cassava.Child(0).Children(), 3)
	require.NotSame(t, cassava.Child(0).Child(0), cassava.Child(0).Child(1))
	require.Same(t, cassava.Child(0), cassava.Child(0).Child(2).Parent())
	// I want to check if the parent of a parent is the node itself
	require.Same(t, cassava.tree, cassava.Child(0).Child(2).Parent().Parent())
}

type testModel struct{}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (t testModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (t testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (t testModel) View() string {
	return ""
}
