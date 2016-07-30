package cli

type Tree struct {
	cell  string
	nodes []TreeView
}

func NewTree(cell string, nodes ...TreeView) *Tree {
	return &Tree{cell: cell, nodes: nodes}
}

func (t *Tree) Cell() string {
	return t.cell
}

func (t *Tree) Nodes() []TreeView {
	nodes := make([]TreeView, len(t.nodes))
	copy(nodes, t.nodes)
	return nodes
}

type TreeIndent struct {
	s []rune
}

func NewTreeIndent() *TreeIndent {
	return &TreeIndent{
		s: make([]rune, 0, 20),
	}
}

func (ti *TreeIndent) Push() int {
	for i, c := range ti.s {
		switch c {
		case '├':
			ti.s[i] = '│'
		case '└', '─':
			ti.s[i] = ' '
		}
	}
	ti.s = append(ti.s, '├', '─', '─', ' ')
	return len(ti.s) / 4
}

func (ti *TreeIndent) Pop() {
	ti.s = ti.s[:len(ti.s)-4]
}

func (ti *TreeIndent) Next(index int, count int, depth int) {
	depth = (depth - 1) * 4

	if (index + 1) != count {
		ti.s[depth] = '├'
	} else {
		ti.s[depth] = '└'
	}

	ti.s[depth+1] = '─'
	ti.s[depth+2] = '─'
}

func (ti *TreeIndent) Clear(index int, depth int) {
	if depth = (depth - 1) * 4; depth >= 0 && index != 0 {
		switch ti.s[depth] {
		case '└':
			ti.s[depth] = ' '
		case '├':
			ti.s[depth] = '│'
		}
		ti.s[depth+1] = ' '
		ti.s[depth+2] = ' '
	}
}

func (ti *TreeIndent) Depth() int {
	return len(ti.s) / 4
}

func (ti *TreeIndent) String() string {
	return string(ti.s)
}
