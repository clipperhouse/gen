package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
)

func TestWrite(t *testing.T) {
	for _, typ := range typs { // from package_test.go
		b := bytes.NewBufferString("")
		writeType(b, typ, options{})

		if _, err := parser.ParseFile(token.NewFileSet(), "dummy", b.String(), 0); err != nil {
			t.Error(err)
		}

		if _, err := formatToBytes(b); err != nil {
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
