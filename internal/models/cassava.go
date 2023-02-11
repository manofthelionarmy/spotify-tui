package models

// Though cassava isn't a tree, this is where they get boba from to make bubbletea:
// https://en.wikipedia.org/wiki/Cassava

import tea "github.com/charmbracelet/bubbletea"

var _ TeaTree = (*tree)(nil)

// tree implements our tree
type tree struct {
	parent   *tree
	children []TeaTree
	value    tea.Model
}

func (t *tree) Parent() TeaTree {
	return t.parent
}

func (t *tree) Children() []TeaTree {
	return t.children
}

func (t *tree) Child(idx int) TeaTree {
	return t.children[idx]
}

func (t *tree) AddChild(m tea.Model) {
	child := &tree{
		parent:   t,
		children: make([]TeaTree, 0),
		value:    m,
	}
	t.children = append(t.children, child)
}

func (t *tree) Value() tea.Model {
	return t.value
}

func newTree(rootModel tea.Model) *tree {
	return &tree{
		parent:   nil,
		children: make([]TeaTree, 0),
		value:    rootModel,
	}
}

// Cassava is a TeaTree implementation
type Cassava struct {
	*tree
}

// New returns a new Cassava
func New(rootModel tea.Model) Cassava {
	return Cassava{
		tree: newTree(rootModel),
	}
}

// Parent returns the parent node
func (c *Cassava) Parent() TeaTree {
	return c.tree.Parent()
}

// Children returns the children belonging to the current node
func (c *Cassava) Children() []TeaTree {
	return c.tree.Children()
}

// Child returns the child node at the given idx
func (c *Cassava) Child(idx int) TeaTree {
	return c.tree.Child(idx)
}

// AddChild adds a child node to the current node
func (c *Cassava) AddChild(m tea.Model) {
	c.tree.AddChild(m)
}

// Value returns the tea.Model
func (c *Cassava) Value() tea.Model {
	return c.tree.Value()
}
