package main

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	// use custom name so test won't interfere with a real _gen.go
	c := defaultConfig
	c.customName = "_gen_get_test.go"

	// clean up when done
	defer os.Remove(c.customName)

	// standard
	imps, err := getTypewriterImports(c)

	if err != nil {
		t.Error(err)
	}

	if len(imps) != 2 {
		t.Errorf("should return 2 imports, got %v", len(imps))
	}

	if err := add(c, "github.com/clipperhouse/foowriter", "github.com/clipperhouse/setwriter"); err != nil {
		t.Error(err)
	}

	// custom file now exists
	imps2, err := getTypewriterImports(c)

	if err != nil {
		t.Error(err)
	}

	if len(imps2) != 4 {
		t.Errorf("should return 4 custom imports, got %v", len(imps2))
	}

	// custom get
	if err := get(c); err != nil {
		t.Error(err)
	}

	// custom get with param
	if err := get(c, "-d"); err != nil {
		t.Error(err)
	}
}
