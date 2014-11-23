package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/clipperhouse/typewriter"
)

func TestList(t *testing.T) {
	// use buffer instead of Stdout so we can inspect the results
	var b bytes.Buffer
	setOutput(&b)
	defer revertOutput()

	// use custom name so test won't interfere with a real _gen.go
	setCustomName("_gen_test.go")
	defer revertCustomName()

	// remove existing files, start fresh
	os.Remove(customName)

	// standard
	if err := runMain([]string{"gen", "list"}); err != nil {
		t.Error(err)
	}

	// one line for title, 2 standard typewriters
	if lines := bytes.Count(b.Bytes(), []byte("\n")); lines != 2 {
		t.Errorf("standard list should output 2 lines, got %v", lines)
	}

	// clear out the output buffer
	b.Reset()

	// create a custom typewriter import file
	w, err := os.Create(customName)

	if err != nil {
		t.Error(err)
	}

	defer os.Remove(customName)

	p := pkg{
		Name: "main",
		Imports: typewriter.NewImportSpecSet(
			typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/foowriter"},
			typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/slicewriter"},
		),
	}

	if err := tmpl.Execute(w, p); err != nil {
		t.Error(err)
	}

	// custom file now exists
	if err := runMain([]string{"gen", "list"}); err != nil {
		t.Error(err)
	}

	// one line for title, 3 custom typewriters
	if lines := bytes.Count(b.Bytes(), []byte("\n")); lines != 3 {
		t.Errorf("standard list should output 3 lines, got %v", lines)
	}
}
