package cli

import (
	"bytes"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Writer interface {
	io.Closer

	io.Writer

	Flush() error
}

func newWriter(term *terminal.Terminal, output *os.File) (writer Writer, err error) {
	if !terminal.IsTerminal(int(output.Fd())) {
		writer = newFileWriter(output)
		return
	}
	return newTermWriter(term, output)
}

type termWriter struct {
	t *terminal.Terminal
	s *terminal.State
	f *os.File
	c chan os.Signal
}

func newTermWriter(t *terminal.Terminal, f *os.File) (writer Writer, err error) {
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

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGWINCH)

	go func() {
		for _ = range sigchan {
			if w, h, err := terminal.GetSize(0); err != nil {
				t.SetSize(w, h)
			}
		}
	}()

	writer = termWriter{
		t: t,
		s: s,
		f: f,
		c: sigchan,
	}
	return
}

func (w termWriter) Close() (err error) {
	err = terminal.Restore(int(w.f.Fd()), w.s)
	w.f.Close()

	defer func() { recover() }()
	close(w.c)
	return
}

func (w termWriter) Write(b []byte) (n int, err error) {
	for err == nil && len(b) != 0 {
		var c int

		if newline := bytes.IndexByte(b, '\n'); newline >= 0 {
			if newline != 0 && b[newline-1] == '\r' {
				newline--
			}

			c, err = w.t.Write(b[:newline])
			n += c

			if err == nil {
				c, err = w.t.Write([]byte{'\r', '\n'})
				n += c
			}

			if err != nil {
				return
			}

			b = b[newline+1:]
			continue
		}

		c, err = w.t.Write(b)
		n += c
		b = nil
	}

	return
}

func (w termWriter) Flush() error {
	return nil
}

type fileWriter struct {
	b *bytes.Buffer
	f *os.File
}

func newFileWriter(f *os.File) Writer {
	b := &bytes.Buffer{}
	b.Grow(4096)
	return fileWriter{
		b: b,
		f: f,
	}
}

func (w fileWriter) Close() (err error) {
	err = w.Flush()
	w.f.Close()
	return
}

func (w fileWriter) Write(b []byte) (n int, err error) {
	if n, err = w.b.Write(b); err == nil {
		for {
			var ok bool
			if ok, err = w.flushLine(); !ok {
				break
			}
		}
	}
	return
}

func (w fileWriter) Flush() (err error) {
	_, err = io.Copy(w.f, w.b)
	return
}

func (w fileWriter) flushLine() (ok bool, err error) {
	b := w.b.Bytes()

	if newline := bytes.IndexByte(b, '\n'); newline >= 0 {
		if err = w.write(w.b.Next(newline + 1)); err == nil {
			ok = true
		}
	}

	return
}

func (w fileWriter) write(b []byte) (err error) {
	if off := bytes.IndexByte(b, '\033'); off >= 0 {
		b = StripStyles(b)
	}
	_, err = w.f.Write(b)
	return
}
