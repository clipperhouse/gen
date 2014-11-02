package main

var stdImports = []string{
	`_ "github.com/clipperhouse/containerwriter"`,
	`_ "github.com/clipperhouse/slicewriter"`,
}

type pkg struct {
	Name    string
	Imports []string
}
