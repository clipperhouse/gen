package main

import (
	"os"
	"testing"

	"github.com/clipperhouse/typewriter"
)

func TestAdd(t *testing.T) {
	// use custom name so test won't interfere with a real _gen.go
	c := defaultConfig
	c.customName = "_gen_add_test.go"
	defer os.Remove(c.customName)

	if err := add(c); err == nil {
		t.Error("add with no arguments should be an error")
	}

	foo := typewriter.ImportSpec{Name: "_", Path: "github.com/clipperhouse/foowriter"}
	before, err := getTypewriterImports(c)

	if err != nil {
		t.Error(err)
	}

	// ensure that the custom import is not in the default set
	if before.Contains(foo) {
		t.Errorf("default imports should not include %s", foo.Path)
	}

	// adding import which exists should succeed
	if err := add(c, foo.Path); err != nil {
		t.Error(err)
	}

	after, err := getTypewriterImports(c)

	// the new import should be reflected in imports
	if !after.Contains(foo) {
		t.Errorf("imports should include %s", foo.Path)
	}
}
