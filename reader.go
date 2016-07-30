package cli

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type Reader interface {
	io.Closer

	ReadLine(prompt string) (line string, err error)

	ReadPassword(prompt string) (line string, err error)
}

func newReader(term *terminal.Terminal, input *os.File) (reader Reader, err error) {
	if !terminal.IsTerminal(int(input.Fd())) {
		reader = newFileReader(input)
		return
	}
	return newTermReader(term, input)
}

type termReader struct {
	t *terminal.Terminal
	s *terminal.State
	f *os.File
}

func newTermReader(t *terminal.Terminal, f *os.File) (reader Reader, err error) {
	var s *terminal.State
	var w int
	var h int
	var fd = int(f.Fd())

	if s, err = terminal.MakeRaw(fd); err != nil {
		return
	}

	if w, h, err = terminal.GetSize(fd); err != nil {
		return
	}

	if err = t.SetSize(w, h); err != nil {
		return
	}

	reader = termReader{
		t: t,
		s: s,
		f: f,
	}
	return
}

func (r termReader) Close() (err error) {
	err = terminal.Restore(int(r.f.Fd()), r.s)
	r.f.Close()
	return
}

func (r termReader) ReadLine(prompt string) (line string, err error) {
	r.t.SetPrompt(prompt)

	if line, err = r.t.ReadLine(); err == terminal.ErrPasteIndicator {
		err = nil
	}

	line = trimLine(line)

	if err == io.EOF {
		r.t.SetPrompt("")
		r.t.Write(append([]byte(prompt), '\r', '\n'))
	}

	return
}

func (r termReader) ReadPassword(prompt string) (string, error) {
	return r.t.ReadPassword(prompt)
}

type fileReader struct {
	b *bufio.Reader
	f *os.File
}

func newFileReader(f *os.File) Reader {
	return fileReader{
		b: bufio.NewReaderSize(f, 4096),
		f: f,
	}
}

func (r fileReader) Close() error {
	return r.f.Close()
}

func (r fileReader) ReadLine(prompt string) (line string, err error) {
	if line, err = r.b.ReadString('\n'); err != nil {
		if err == io.EOF && len(line) != 0 {
			err = nil
		}
	}
	line = trimLine(line)
	return
}

func (r fileReader) ReadPassword(prompt string) (string, error) {
	return r.ReadLine(prompt)
}

func trimLine(line string) string {
	if n := len(line); n != 0 && line[n-1] == '\n' {
		line = line[:n-1]
	}

	if n := len(line); n != 0 && line[n-1] == '\r' {
		line = line[:n-1]
	}

	return line
}
