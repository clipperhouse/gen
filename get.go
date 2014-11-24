package main

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/clipperhouse/typewriter"
)

// get runs `go get` for required typewriters, either default or specified in _gen.go
func get(c config, args ...string) error {
	imports, err := getTypewriterImports(c)

	if err != nil {
		return err
	}

	// we just want the paths
	var imps []string
	for imp := range imports {
		imps = append(imps, imp.Path)
	}

	get := []string{"get"}
	get = append(get, args...)
	get = append(get, imps...)

	cmd := exec.Command("go", get...)
	cmd.Stdout = c.out
	cmd.Stderr = c.out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getTypewriterImports(c config) (typewriter.ImportSpecSet, error) {
	imports := typewriter.NewImportSpecSet()

	// check for existence of custom file
	if src, err := os.Open(c.customName); err == nil {
		defer src.Close()

		// custom file exists, parse its imports
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
		if err != nil {
			return imports, err
		}

		// convert ast imports into ImportSpecs
		for _, v := range f.Imports {
			imp := typewriter.ImportSpec{
				Name: v.Name.Name,
				Path: strings.Trim(v.Path.Value, `"`), // lose the quotes
			}
			imports.Add(imp)
		}
	} else {
		// doesn't exist, use standard (clone it)
		imports = stdImports.Clone()
	}

	return imports, nil
}
