package main

import (
	"os"
	"testing"
)

// +gen methods:"Where"
type dummy int

func TestRunCustom(t *testing.T) {
	customName := "_gen_test.go"
	genName := "dummy_gen.go"
	fooName := "dummy_foo.go"

	// remove existing files, start fresh
	os.Remove(customName)
	os.Remove(genName)
	os.Remove(fooName)

	// standard run
	run(customName)

	// gen file should exist
	if _, err := os.Open(genName); err != nil {
		t.Error(err)
	}

	// foo file should not exist, not a standard typewriter
	if _, err := os.Open(fooName); err == nil {
		t.Errorf("%s should not have been generated", fooName)
	}

	// remove just-gen'd file
	os.Remove(genName)

	// create a custom typewriter import file
	w, err := os.Create(customName)

	if err != nil {
		t.Error(err)
	}

	p := pkg{
		Name: "main",
		Imports: []string{
			// non-standard typewriter
			`_ "github.com/clipperhouse/gen/typewriters/foowriter"`,
		},
		Main: false,
	}

	if err := tmpl.Execute(w, p); err != nil {
		t.Error(err)
	}

	// custom run
	run(customName)

	// foo file should exist
	if _, err := os.Open(fooName); err != nil {
		t.Error(err)
	}

	// gen file should not exist, because it was not included in the custom file
	if _, err := os.Open(genName); err == nil {
		t.Errorf("%s should not have been generated", genName)
	}

	// remove just-gen'd file
	os.Remove(fooName)

	// remove custom file
	os.Remove(customName)
}
