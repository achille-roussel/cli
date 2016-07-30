package cli

import "io"

type CellAlign int

const (
	LeftAlign  CellAlign = -1
	RightAlign CellAlign = 1
	Centered   CellAlign = 0
)

func RenderCellLeftAlign(w io.Writer, cell string, width int) (err error) {
	if _, err = io.WriteString(w, cell); err != nil {
		return
	}

	if n := RuneCountInString(cell); n < width {
		err = writeSpaces(w, width-n)
	}

	return
}

func RenderCellRightAlign(w io.Writer, cell string, width int) (err error) {
	if n := RuneCountInString(cell); n < width {
		err = writeSpaces(w, width-n)
	}

	if _, err = io.WriteString(w, cell); err != nil {
		return
	}

	return
}

func RenderCellCentered(w io.Writer, cell string, width int) (err error) {
	n := RuneCountInString(cell)
	l := 0
	r := 0

	if n < width {
		l = (width - n) / 2
		r = width - (n + l)
	}

	if err = writeSpaces(w, l); err != nil {
		return
	}

	if _, err = io.WriteString(w, cell); err != nil {
		return
	}

	if err = writeSpaces(w, r); err != nil {
		return
	}

	return
}

func RenderCell(w io.Writer, cell string, width int, align CellAlign) (err error) {
	switch align {
	case LeftAlign:
		return RenderCellLeftAlign(w, cell, width)
	case RightAlign:
		return RenderCellRightAlign(w, cell, width)
	default:
		return RenderCellCentered(w, cell, width)
	}
}

func writeSpaces(w io.Writer, n int) (err error) {
	if n != 0 {
		_, err = w.Write(makeSpaces(n))
	}
	return
}

func makeSpaces(n int) []byte {
	s := make([]byte, n)

	for i := range s {
		s[i] = ' '
	}

	return s
}
