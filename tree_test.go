package cli

import (
	"bytes"
	"testing"
)

func TestTree(t *testing.T) {
	b := &bytes.Buffer{}
	b.Grow(1024)

	s := Style(Bold, Blue)

	RenderTreeView(b, NewTree(s.S("."),
		NewTree(s.S("A"),
			NewTree("Hello\nWorld")),
		NewTree(s.S("B"),
			NewTree("1"),
			NewTree("2"),
			NewTree("3")),
		NewTree(s.S("C"))),
	)

	t.Logf("\n\n%s\n", b.String())
}
