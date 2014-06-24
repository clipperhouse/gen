package main

import (
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestCustom(t *testing.T) {
	// use custom name so test won't interfere with a real _gen.go
	setCustomName("_gen_test.go")
	defer revertCustomName()

	// ensure that generated _gen.go uses the same imports as built into this package
	custom()
	defer os.Remove(customName)

	fset := token.NewFileSet()
	c, err := parser.ParseFile(fset, customName, nil, parser.ImportsOnly)

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
