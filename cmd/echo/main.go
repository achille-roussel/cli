package main

import (
	"io"
	"os"

	"github.com/achille-roussel/cli"
)

func main() {
	cli.Init()
	defer cli.Close()

	for {
		var line string
		var err error

		if line, err = cli.ReadLine("> "); err != nil {
			if err != io.EOF {
				cli.Println(err)
				os.Exit(1)
			}
			return
		}

		if len(line) != 0 {
			cli.Println(line)
		}
	}
}
