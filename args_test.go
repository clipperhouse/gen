package main

import (
	"strings"
	"testing"
)

func testTypeArg(t *testing.T, s string, num int) {
	args := strings.Split(s, " ")
	opts, err := parseArgs(args)

	defaultOpts := options{}

	if opts != defaultOpts {
		t.Errorf("expected default options '%v', got '%v'", defaultOpts, opts)
	}

	if err != nil {
		t.Errorf("expected no errors, got '%v'", errs)
	}
}

func TestSynonymousFlags(t *testing.T) {
	args7 := strings.Split("-f", " ")
	opts7, _ := parseArgs(args7)

	args8 := strings.Split("-force", " ")
	opts8, _ := parseArgs(args8)

	if opts7 != opts8 {
		t.Errorf("-f and -force should by synonymous")
	}
}

func TestUnknownArgs(t *testing.T) {
	args := strings.Split("-clown", " ")
	_, err := parseArgs(args)

	if err == nil {
		t.Errorf("expected error for passing invalid flag, got none")
	}

	args2 := strings.Split("-b", " ")
	_, err2 := parseArgs(args2)

	if err2 == nil {
		t.Errorf("expected error for passing invalid flag, got none")
	}
}

func TestDeprecated(t *testing.T) {
	args := strings.Split("*package.Type", " ")
	_, err := parseArgs(args)

	if err == nil {
		t.Errorf("expected error for deprecated type argument, got none")
	}

	args2 := strings.Split("-a", " ")
	_, err2 := parseArgs(args2)

	if err2 == nil {
		t.Errorf("expected error for deprecated -all flag, got none")
	}

	args3 := strings.Split("-e", " ")
	_, err3 := parseArgs(args3)

	if err3 == nil {
		t.Errorf("expected error for deprecated -exported flag, got none")
	}
}
