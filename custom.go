package main

import "os"

func custom(filename string) error {
	w, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer w.Close()

	p := pkg{
		Name:    "main",
		Imports: stdImports,
		Main:    false,
	}

	if err := tmpl.Execute(w, p); err != nil {
		return err
	}

	return nil
}
