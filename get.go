package main

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

// get runs `go get` for required typewriters, either default or specified in _gen.go
func get(args []string) error {
	imports, err := getTypewriterImports()

	if err != nil {
		return err
	}

	get := []string{"get"}
	get = append(get, args...)
	get = append(get, imports...)

	cmd := exec.Command("go", get...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil
	}

	return nil
}

func getTypewriterImports() ([]string, error) {
	imports := make([]string, 0)

	// check for existence of custom file
	if src, err := os.Open(customName); err == nil {
		defer src.Close()

		// custom file exists, parse its imports
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
		if err != nil {
			return imports, err
		}
		for _, v := range f.Imports {
			imports = append(imports, v.Path.Value)
		}
	} else {
		// doesn't exist, use standard
		imports = append(imports, stdImports...)
	}

	// clean `em up
	// TODO: a better way than strings to express imports
	for i := range imports {
		imports[i] = strings.Trim(imports[i], `_ "`)
	}

	return imports, nil
}
