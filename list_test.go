package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestList(t *testing.T) {
	customName := "_gen_test.go"
	// remove existing files, start fresh
	os.Remove(customName)

	// standard
	out, err := list(customName)

	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadAll(out)

	if err != nil {
		t.Error(err)
	}

	// one line for title, 2 standard typewriters
	if lines := bytes.Count(b, []byte("\n")); lines != 3 {
		t.Errorf("standard list should output 3 lines, got %v", lines)
	}

	// create a custom typewriter import file
	w, err := os.Create(customName)

	if err != nil {
		t.Error(err)
	}

	defer os.Remove(customName)

	p := pkg{
		Name: "main",
		Imports: []string{
			// non-standard typewriter
			`_ "github.com/clipperhouse/gen/typewriters/foowriter"`,
			`_ "github.com/clipperhouse/gen/typewriters/genwriter"`,
			`_ "github.com/clipperhouse/gen/typewriters/container"`,
		},
	}

	if err := tmpl.Execute(w, p); err != nil {
		t.Error(err)
	}

	// custom
	out2, err := list(customName)

	if err != nil {
		t.Error(err)
	}

	b2, err := ioutil.ReadAll(out2)

	if err != nil {
		t.Error(err)
	}

	// one line for title, 3 custom typewriters
	if lines := bytes.Count(b2, []byte("\n")); lines != 4 {
		t.Errorf("standard list should output 4 lines, got %v", lines)
	}
}
