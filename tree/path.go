package tree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/achille-roussel/cli"
)

type PathConfig struct {
	DirStyle      cli.StyleSet
	RegFileStyle  cli.StyleSet
	ExecFileStyle cli.StyleSet
	ShowHidden    bool
}

var (
	DefaultPathConfig PathConfig = PathConfig{
		DirStyle:      cli.Style(cli.Bold, cli.Blue),
		RegFileStyle:  cli.Normal,
		ExecFileStyle: cli.Style(cli.Bold, cli.Green),
		ShowHidden:    false,
	}
)

func Path(info os.FileInfo, path string) cli.TreeView {
	return PathWithConfig(info, path, DefaultPathConfig)
}

func PathWithConfig(info os.FileInfo, path string, config PathConfig) cli.TreeView {
	return makeFile(info, path, path, &config)
}

type file struct {
	path   string
	name   string
	info   os.FileInfo
	config *PathConfig
}

func makeFile(info os.FileInfo, name string, path string, config *PathConfig) file {
	return file{
		path:   path,
		name:   name,
		info:   info,
		config: config,
	}
}

func (f file) Cell() string {
	style := f.config.RegFileStyle

	if f.info.Mode().IsDir() {
		style = f.config.DirStyle
	} else if (f.info.Mode().Perm() & 0111) != 0 {
		style = f.config.ExecFileStyle
	}

	return style.S(f.name)
}

func (f file) Nodes() (nodes []cli.TreeView) {
	if f.info.IsDir() {
		files, _ := ioutil.ReadDir(f.path)
		nodes = make([]cli.TreeView, 0, len(files))

		for _, info := range files {
			name := info.Name()
			path := filepath.Join(f.path, name)

			if f.config.ShowHidden || !strings.HasPrefix(name, ".") {
				nodes = append(nodes, makeFile(info, name, path, f.config))
			}
		}
	}

	return
}
