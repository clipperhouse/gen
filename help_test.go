package main

import (
	"bytes"
	"testing"
)

func TestHelp(t *testing.T) {
	// use buffer instead of Stdout so help doesn't write to console
	var b bytes.Buffer
	setOutput(&b)
	defer revertOutput()

	// just run it; not a great test but at least covers the code
	if err := help(); err != nil {
		t.Error(err)
	}

	if b.Len() == 0 {
		t.Errorf("help() should have resulted in output")
	}
}
