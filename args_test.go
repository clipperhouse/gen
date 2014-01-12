package main

import (
	"strings"
	"testing"
)

func TestPlainTypeArg(t *testing.T) {
	testTypeArg(t, "pkg.Type", 1)
	testTypeArg(t, "*pkg.Type", 1)
	testTypeArg(t, "pkg.白鵬翔", 1)
	testTypeArg(t, "*αβ.Type", 1)
	testTypeArg(t, "pkg.Type pkg.Another", 2)
}

func testTypeArg(t *testing.T, s string, num int) {
	args := strings.Split(s, " ")
	typeArgs, opts, err := parseArgs(args)

	defaultOpts := options{}

	if len(typeArgs) != num {
		t.Errorf("expected %d typeArg(s), got %d", num, len(typeArgs))
	}

	if opts != defaultOpts {
		t.Errorf("expected default options '%v', got '%v'", defaultOpts, opts)
	}

	if err != nil {
		t.Errorf("expected no error, got '%v'", err)
	}
}

func TestPointerArgs(t *testing.T) {
	args := strings.Split("-a", " ")
	_, opts, _ := parseArgs(args)

	if len(opts.AllPointer) != 0 {
		t.Errorf("-a should not result in pointer")
	}

	args2 := strings.Split("-*a", " ")
	_, opts2, _ := parseArgs(args2)

	if opts2.AllPointer != "*" {
		t.Errorf("-a should result in pointer")
	}

	args3 := strings.Split("pkg.Type", " ")
	typeArgs3, _, _ := parseArgs(args3)

	if len(typeArgs3[0].Pointer) != 0 {
		t.Errorf("pkg.Type should not result in pointer")
	}

	args4 := strings.Split("*pkg.Type", " ")
	typeArgs4, _, _ := parseArgs(args4)

	if typeArgs4[0].Pointer != "*" {
		t.Errorf("*pkg.Type should result in pointer")
	}
}

func TestSynonymousFlags(t *testing.T) {
	args := strings.Split("-a", " ")
	_, opts, _ := parseArgs(args)

	args2 := strings.Split("-all", " ")
	_, opts2, _ := parseArgs(args2)

	if opts != opts2 {
		t.Errorf("-a and -all should by synonymous")
	}

	args3 := strings.Split("-*a", " ")
	_, opts3, _ := parseArgs(args3)

	args4 := strings.Split("-*all", " ")
	_, opts4, _ := parseArgs(args4)

	if opts3 != opts4 {
		t.Errorf("-*a and -*all should by synonymous")
	}

	args5 := strings.Split("-e", " ")
	_, opts5, _ := parseArgs(args5)

	args6 := strings.Split("-exported", " ")
	_, opts6, _ := parseArgs(args6)

	if opts5 != opts6 {
		t.Errorf("-e and -exported should by synonymous")
	}

	args7 := strings.Split("-f", " ")
	_, opts7, _ := parseArgs(args7)

	args8 := strings.Split("-force", " ")
	_, opts8, _ := parseArgs(args8)

	if opts7 != opts8 {
		t.Errorf("-f and -force should by synonymous")
	}
}

func TestUnknownArgs(t *testing.T) {
	args := strings.Split("pkg.Type -clown", " ")
	_, _, err := parseArgs(args)

	if err == nil {
		t.Errorf("expected error for passing invalid flag, got none")
	}

	args2 := strings.Split("-b", " ")
	_, _, err2 := parseArgs(args2)

	if err2 == nil {
		t.Errorf("expected error for passing invalid flag, got none")
	}

	args3 := strings.Split("*pkgType", " ")
	_, _, err3 := parseArgs(args3)

	if err3 == nil {
		t.Errorf("expected error for passing invalid type, got none")
	}

	args4 := strings.Split("*pkgType -f", " ")
	_, _, err4 := parseArgs(args4)

	if err4 == nil {
		t.Errorf("expected error for passing invalid type, got none")
	}
}

func TestConflictingArgs(t *testing.T) {
	args := strings.Split("pkg.Type -all", " ")
	_, _, err := parseArgs(args)

	if err == nil {
		t.Errorf("passing both type and -all should result in error")
	}

	args2 := strings.Split("-e", " ")
	_, _, err2 := parseArgs(args2)

	if err2 == nil {
		t.Errorf("passing -e(xported) without -a(ll) should result in error")
	}

	args3 := strings.Split("pkg.Type -e", " ")
	_, _, err3 := parseArgs(args3)

	if err3 == nil {
		t.Errorf("passing -e(xported) with type should result in error")
	}

	args4 := strings.Split("-a -e", " ")
	_, _, err4 := parseArgs(args4)

	if err4 != nil {
		t.Errorf("passing -a(ll) and -e(exported) should be ok")
	}
}
