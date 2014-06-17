package main

import (
	"fmt"
	"os"
)

func custom() {
	w, err := os.Create("_gen.go")

	if err != nil {
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
