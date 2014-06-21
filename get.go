package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

func get(u bool) error {
	imports := make([]string, 0)

	if src, err := os.Open(customFilename); err == nil {
		// custom file exists, parse its imports
		defer src.Close()

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

	var out, outerr bytes.Buffer

	cmd := exec.Command("go", get...)
	cmd.Stdout = &out
	cmd.Stderr = &outerr

	if err := cmd.Run(); err != nil {
		return err
	}

	if out.Len() > 0 {
		fmt.Println(out.String())
	}

	if outerr.Len() > 0 {
		return errors.New(outerr.String())
	}

	return nil
}
