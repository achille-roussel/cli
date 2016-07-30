package cli

import (
	"io"
	"strings"
)

type TreeView interface {
	Cell() string

	Nodes() []TreeView
}

func RenderTreeView(w io.Writer, t TreeView) (err error) {
	return renderTreeView(w, t, NewTreeIndent())
}

func renderTreeView(w io.Writer, tree TreeView, indent *TreeIndent) (err error) {
	nodes := tree.Nodes()
	lines := strings.Split(tree.Cell(), "\n")

	for index, line := range lines {
		indent.Clear(index, indent.Depth())

		if err = renderTreeLine(w, line, indent.String()); err != nil {
			return
		}
	}

	depth := indent.Push()
	count := len(nodes)

	for index, node := range nodes {
		indent.Next(index, count, depth)

		if err = renderTreeView(w, node, indent); err != nil {
			return
		}
	}

	indent.Pop()
	return
}

func renderTreeLine(w io.Writer, line string, indent string) (err error) {
	if _, err = io.WriteString(w, indent); err == nil {
		if _, err = io.WriteString(w, line); err == nil {
			_, err = io.WriteString(w, "\n")
		}
	}
	return
}
