package main

import "github.com/clipperhouse/typewriter"

var stdImports = typewriter.ImportSpecSlice{
	{Name: "_", Path: "github.com/clipperhouse/slicewriter"},
}

type pkg struct {
	Name    string
	Imports typewriter.ImportSpecSlice
}
