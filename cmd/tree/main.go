package main

import (
	"flag"
	"os"

	"github.com/achille-roussel/cli"
	"github.com/achille-roussel/cli/tree"
)

func main() {
	var config = tree.DefaultPathConfig
	var paths []string

	flag.BoolVar(&config.ShowHidden, "a", false, "show hidden files")
	flag.Parse()

	if paths = flag.Args(); len(paths) == 0 {
		paths = []string{"."}
	}

	cli.Init()
	defer cli.Close()

	for _, path := range paths {
		if dir, err := os.Open(path); err != nil {
			cli.Printf("%s: %s\n", path, err)
		} else {
			info, _ := dir.Stat()
			dir.Close()
			cli.RenderTreeView(cli.Output, tree.PathWithConfig(info, path, config))
		}
	}
}
