package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
)

func TestWrite(t *testing.T) {
	for _, typ := range typs { // from package_test.go
		var b bytes.Buffer
		writeType(&b, typ, options{})

		if _, err := parser.ParseFile(token.NewFileSet(), "dummy", b.String(), 0); err != nil {
			t.Error(err)
		}

		if _, err := formatToBytes(&b); err != nil {
			t.Error(err)
		}
	}
}

func TestFormat(t *testing.T) {
	in := `package dummy 
 import   "fmt"
   type MyType  int`

	out := `package dummy

import "fmt"

type MyType int
`

	b := bytes.NewBufferString(in)
	byts, err := formatToBytes(b)

	if err != nil {
		t.Error(err)
	}

	s := string(byts)

	if s != out {
		t.Errorf("format failed, expected %v, got %v:", []byte(out), []byte(s))
	}
}

func TestDeletions(t *testing.T) {
	packages := getPackages()
	existing := map[string]bool{
		"thing1_gen.go":  true, // gen'd type, see package_test.go
		"thing2_gen.go":  true,
		"another_gen.go": true,
	}
	deletions := []string{
		"thing2_gen.go",
		"another_gen.go",
	}
	none := map[string]bool{}

	n := "n"
	y := "y"
	dummy := "dummy"

	var input, output bytes.Buffer

	input.WriteString(n)
	if _, ok := promptDeletions(packages, existing, &input, &output); ok {
		t.Errorf("promptDeletions should not be ok with '%v' as input", n)
	}

	input.WriteString(dummy)
	if _, ok := promptDeletions(packages, existing, &input, &output); ok {
		t.Errorf("promptDeletions should not be ok with '%v' as input", dummy)
	}

	input.WriteString(y)
	if _, ok := promptDeletions(packages, existing, &input, &output); !ok {
		t.Errorf("promptDeletions should be ok with '%v' as input", y)
	}

	input.WriteString(y)
	if d, _ := promptDeletions(packages, existing, &input, &output); !sliceEqual(d, deletions) {
		t.Errorf("promptDeletions should return %v, got %v", deletions, d)
	}

	if _, ok := promptDeletions(packages, none, &input, &output); !ok {
		t.Errorf("promptDeletions should return true when no files to delete")
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
