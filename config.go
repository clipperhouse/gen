package main

import (
	"io"
	"os"

	"github.com/clipperhouse/typewriter"
)

type config struct {
	out        io.Writer
	customName string
	*typewriter.Config
}

var defaultConfig = config{
	out:        os.Stdout,
	customName: "_gen.go",
	Config:     &typewriter.Config{},
}

// keep in sync with imports.go
var stdImports = typewriter.NewImportSpecSet(
	typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/slice"},
	typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/stringer"},
)
