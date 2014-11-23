package main

import "github.com/clipperhouse/typewriter"

var stdImports = typewriter.NewImportSpecSet(
	typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/slicewriter"},
)

type pkg struct {
	Name    string
	Imports typewriter.ImportSpecSet
}
