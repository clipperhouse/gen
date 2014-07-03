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

	if err := runMain([]string{"gen", "help"}); err != nil {
		t.Error(err)
	}

	if b.Len() == 0 {
		t.Errorf("help() should have resulted in output")
	}

	// grab the text for later comparison
	text := b.String()
	b.Reset()

	if err := runMain([]string{"gen", "foo"}); err != nil {
		t.Error(err)
	}

	if text != b.String() {
		t.Errorf("running an unknown command should return help text; returned %s", b.String())
	}
}
