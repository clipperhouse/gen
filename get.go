package main

import (
	"bytes"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

// get runs `go get` for required typewriters, either default or specified in _gen.go
func get(u bool) error {
	imports := make([]string, 0)

	// check for existence of custom file
	if src, err := os.Open(customName); err == nil {
		defer src.Close()

		// custom file exists, parse its imports
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, v := range f.Imports {
			imports = append(imports, v.Path.Value)
		}
	} else {
		// doesn't exist, use standard
		imports = append(imports, stdImports...)
	}

	// clean `em up
	// TODO: a better way to express imports
	for i := range imports {
		imports[i] = strings.Trim(imports[i], `_ "`)
	}

	get := []string{"get"}
	if u {
		get = append(get, "-u")
	}

	get = append(get, imports...)

	var outerr bytes.Buffer

	cmd := exec.Command("go", get...)
	cmd.Stdout = out
	cmd.Stderr = &outerr

	if err := cmd.Run(); err != nil {
		return err
	}

	if outerr.Len() > 0 {
		return errors.New(outerr.String())
	}

	return nil
}
