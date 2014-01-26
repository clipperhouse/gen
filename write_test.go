package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
)

func TestWrite(t *testing.T) {
	for _, typ := range typs {
		b := bytes.NewBufferString("")
		writeType(b, typ, options{})

		if _, err := parser.ParseFile(token.NewFileSet(), "dummy", b.String(), 0); err != nil {
			t.Error(err)
		}
	}
}
