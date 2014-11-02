package main

var stdImports = []string{
	`_ "github.com/clipperhouse/gen/typewriters/container"`,
	`_ "github.com/clipperhouse/slicewriter"`,
}

type pkg struct {
	Name    string
	Imports []string
}
