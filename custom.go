package main

import (
	"fmt"
	"os"
)

func custom(filename string) {
	w, err := os.Create(filename)

	if err != nil {
		// TODO: return err
		fmt.Println(err)
		return
	}

	defer w.Close()

	p := pkg{
		Name:    "main",
		Imports: stdImports,
		Main:    false,
	}

	tmpl.Execute(w, p)
}
