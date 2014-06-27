package main

import (
	"os"
	"testing"
)

// +gen methods:"Where"
type dummy int

func TestRun(t *testing.T) {
	// use custom name so test won't interfere with a real _gen.go
	setCustomName("_gen_test.go")
	defer revertCustomName()

	genName := "dummy_gen_test.go"
	fooName := "dummy_foo_test.go"

	// remove existing files, start fresh
	os.Remove(customName)
	os.Remove(genName)
	os.Remove(fooName)

	// standard run
	if err := run(); err != nil {
		t.Error(err)
	}

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
	}

	if err := tmpl.Execute(w, p); err != nil {
		t.Error(err)
	}

	// custom run
	if err := run(); err != nil {
		t.Error(err)
	}

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
