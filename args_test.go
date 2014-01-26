package main

import (
	"strings"
	"testing"
)

func TestFlags(t *testing.T) {
	args1 := strings.Split("-f", " ")
	opts1, _ := parseArgs(args1)

	args2 := strings.Split("-force", " ")
	opts2, _ := parseArgs(args2)

	if !opts1.Force || !opts2.Force {
		t.Errorf("should have detected force flag")
	}

	if opts1 != opts2 {
		t.Errorf("-f and -force should by synonymous")
	}

	args3 := strings.Split("-h", " ")
	opts3, _ := parseArgs(args3)

	args4 := strings.Split("-help", " ")
	opts4, _ := parseArgs(args4)

	if !opts3.Help || !opts3.Help {
		t.Errorf("should have detected help flag")
	}

	if opts3 != opts4 {
		t.Errorf("-h and -help should by synonymous")
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
