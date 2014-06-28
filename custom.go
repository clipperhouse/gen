package main

import "os"

// custom generates a _gen.go file with default typewriters
func custom() error {
	w, err := os.Create(customName)

	if err != nil {
		return err
	}

	defer w.Close()

	p := pkg{
		Name:    "main",
		Imports: stdImports,
	}

	if err := tmpl.Execute(w, p); err != nil {
		return err
	}

	return nil
}
