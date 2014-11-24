package main

import (
	"bytes"
	"testing"
)

func TestHelp(t *testing.T) {
	// use buffer instead of Stdout so help doesn't write to console
	var b bytes.Buffer
	c := defaultConfig
	c.out = &b

	if err := help(c); err != nil {
		t.Error(err)
	}

	if b.Len() == 0 {
		t.Errorf("help() should have resulted in output")
	}
}
