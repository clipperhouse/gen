package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

func get(u bool) {
	imports := make([]string, 0)

	if _, err := os.Open("_gen.go"); err == nil {
		// custom file exists, parse its imports
		fset := token.NewFileSet()
		if f, err := parser.ParseFile(fset, "_gen.go", nil, parser.ImportsOnly); err == nil {
			for _, v := range f.Imports {
				imports = append(imports, v.Path.Value)
			}
		}
	} else {
		// doesn't exist, use standard
		imports = append(imports, stdImports...)
	}

	// clean `em up, hacky; TODO: a better way
	for i := range imports {
		imports[i] = strings.Trim(imports[i], `_ "`)
	}

	get := []string{"get"}
	if u {
		get = append(get, "-u")
	}

	get = append(get, imports...)

	var out bytes.Buffer

	cmd := exec.Command("go", get...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	if out.Len() > 0 {
		fmt.Println(out.String())
	}
}
