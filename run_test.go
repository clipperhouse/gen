package main

import (
	"os"
	"testing"

	"github.com/clipperhouse/typewriter"
)

// +gen slice:"Where"
type dummy int

func TestRun(t *testing.T) {
	// use custom name so test won't interfere with a real _gen.go
	c := defaultConfig
	c.customName = "_gen_run_test.go"

	sliceName := "dummy_slice_test.go"
	fooName := "dummy_foo_test.go"

	// standard run
	if err := run(c); err != nil {
		t.Fatal(err)
	}

	// gen file should exist
	if _, err := os.Stat(sliceName); err != nil {
		t.Error(err)
	}

	// foo file should not exist, not a standard typewriter
	if _, err := os.Stat(fooName); err == nil {
		t.Errorf("%s should not have been generated", fooName)
	}

	// remove just-gen'd file
	if err := os.Remove(sliceName); err != nil {
		t.Fatal(err)
	}

	// create a custom typewriter import file
	imports := typewriter.NewImportSpecSet(
		typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/foowriter"},
	)

	if err := createCustomFile(c, imports); err != nil {
		t.Fatal(err)
	}

	// custom run
	if err := run(c); err != nil {
		t.Error(err)
	}

	// clean up custom file, no longer needed
	if err := os.Remove(c.customName); err != nil {
		t.Fatal(err)
	}

	// foo file should exist
	if _, err := os.Stat(fooName); err != nil {
		t.Error(err)
	}

	// clean up foo file
	if err := os.Remove(fooName); err != nil {
		t.Fatal(err)
	}

	// gen file should not exist, because it was not included in the custom file
	if _, err := os.Stat(sliceName); err == nil {
		t.Errorf("%s should not have been generated", sliceName)
	}
}
