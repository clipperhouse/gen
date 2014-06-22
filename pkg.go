package main

var stdImports = []string{
	`_ "github.com/clipperhouse/gen/typewriters/container"`,
	`_ "github.com/clipperhouse/gen/typewriters/genwriter"`,
}

type pkg struct {
	Name    string
	Imports []string
}
