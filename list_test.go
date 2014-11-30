package main

import (
	"bytes"
	"os"
	"testing"
)

func TestList(t *testing.T) {
	// use buffer instead of Stdout so we can inspect the results
	var b bytes.Buffer
	c := defaultConfig
	c.out = &b
	c.customName = "_gen_list_test.go"

	// clean up when done
	defer os.Remove(c.customName)

	// standard
	if err := list(c); err != nil {
		t.Error(err)
	}

	// 1 line for title + 1 standard typewriter (see imports.go)
	if lines := bytes.Count(b.Bytes(), []byte("\n")); lines != 2 {
		t.Errorf("standard list should output 2 lines, got %v", lines)
	}

	// clear out the buffer
	b.Reset()

	// create a custom typewriter import file
	if err := add(c, "github.com/clipperhouse/foowriter"); err != nil {
		t.Error(err)
	}

	// custom file now exists
	if err := list(c); err != nil {
		t.Error(err)
	}

	// 1 line for title + 2 custom typewriters
	if lines := bytes.Count(b.Bytes(), []byte("\n")); lines != 3 {
		t.Errorf("standard list should output 3 lines, got %v", lines)
	}
}
