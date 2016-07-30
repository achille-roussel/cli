package cli

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	term   ReadWriter
	Input  Reader
	Output Writer
)

type ReadWriter interface {
	io.Closer

	io.Writer

	Flush() error

	ReadLine(prompt string) (line string, err error)

	ReadPassword(prompt string) (line string, err error)
}

type readWriter struct {
	Reader
	Writer
}

func (rw readWriter) Close() (err error) {
	rw.Writer.Close()
	rw.Reader.Close()
	return
}

func Close() error {
	return term.Close()
}

func Flush() error {
	return term.Flush()
}

func Print(args ...interface{}) (int, error) {
	return fmt.Fprint(term, args...)
}

func Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(term, args...)
}

func Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(term, format, args...)
}

func ReadLine(prompt string) (line string, err error) {
	return Input.ReadLine(prompt)
}

func ReadPassowrd(prompt string) (line string, err error) {
	return Input.ReadPassword(prompt)
}

func New(input *os.File, output *os.File) (rw ReadWriter, err error) {
	var reader Reader
	var writer Writer

	if terminal.IsTerminal(int(output.Fd())) {
		term := terminal.NewTerminal(struct {
			io.Reader
			io.Writer
		}{input, output}, "")

		if reader, err = newReader(term, input); err != nil {
			return
		}

		if writer, err = newWriter(term, output); err != nil {
			return
		}
	} else {
		reader = newFileReader(input)
		writer = newFileWriter(output)
	}

	rw = readWriter{reader, writer}
	return
}

func Init() (err error) {
	term, err = New(os.Stdin, os.Stdout)
	Input = term
	Output = term
	return
}
