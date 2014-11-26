package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/clipperhouse/typewriter"
)

// add adds a new typewriter import to the current package, by creating (or appending) a _gen.go file.
func add(c config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify the import path of the typewriter you wish to add")
	}

	imports, err := getTypewriterImports(c)

	if err != nil {
		return err
	}

	for _, arg := range args {
		imp := typewriter.ImportSpec{Name: "_", Path: arg}

		// try to go get it
		cmd := exec.Command("go", "get", imp.Path)
		cmd.Stdout = c.out
		cmd.Stderr = c.out

		if err := cmd.Run(); err != nil {
			return err
		}

		imports.Add(imp)
	}

	if createCustomFile(c, imports); err != nil {
		return err
	}

	return nil
}

func createCustomFile(c config, imports typewriter.ImportSpecSet) error {
	w, err := os.Create(c.customName)

	if err != nil {
		return err
	}

	defer w.Close()

	p := pkg{
		Name:    "main",
		Imports: imports,
	}

	if err := tmpl.Execute(w, p); err != nil {
		return err
	}

	return nil
}
