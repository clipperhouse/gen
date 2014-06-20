package main

import (
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestCustom(t *testing.T) {
	// ensure that generated _gen.go uses the same imports as built into this package

	filename := "_gen_test.go"
	custom(filename)

	defer os.Remove(filename)

	fset := token.NewFileSet()
	c, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)

	if err != nil {
		t.Error(err)
	}

	fset2 := token.NewFileSet()
	imp, err := parser.ParseFile(fset2, "imports.go", nil, parser.ImportsOnly)

	if err != nil {
		t.Error(err)
	}

	if len(c.Imports) != len(imp.Imports) {
		t.Errorf("generated imports should equal builtin imports; %v != %v", len(c.Imports), len(imp.Imports))
	}

	for i := range c.Imports {
		if c.Imports[i].Path.Value != imp.Imports[i].Path.Value {
			t.Errorf("generated imports should equal builtin imports")
		}
	}
}
