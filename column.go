package cli

import "io"

type column string

func (c column) alignment() (a CellAlign) {
	if n := len(c); n != 0 && c[0] == ':' {
		a += RightAlign
	}

	if n := len(c); n != 0 && c[n-1] == ':' {
		a += LeftAlign
	}

	return
}

func (c column) string() string {
	s := string(c)

	if n := len(s); n != 0 && s[0] == ':' {
		s = s[1:]
	}

	if n := len(s); n != 0 && s[n-1] == ':' {
		s = s[:n-1]
	}

	return s
}

func RenderColumn(w io.Writer, col string, width int) (err error) {
	return RenderCell(w, column(col).string(), width, column(col).alignment())
}
