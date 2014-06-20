package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
)

func run() {
	if src, err := os.Open(customFilename); err == nil {
		// custom imports file exists, use it
		defer src.Close()
		runCustom(src)
	} else {
		// do it the regular way
		runStandard()
	}
}

func runStandard() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		panic(err)
	}

	if err := app.WriteAll(); err != nil {
		panic(err)
	}
}

func runCustom(src *os.File) {
	temp, err := getTempDir()
	if err != nil {
		// TODO return err?
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(temp)

	// set up imports file containing the custom typewriters (from _gen.go)
	imports, err := os.Create(filepath.Join(temp, "imports.go"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer imports.Close()

	io.Copy(imports, src)

	// set up main to be run
	main, err := os.Create(filepath.Join(temp, "main.go"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer main.Close()

	p := pkg{
		Name: "main",
		Imports: []string{
			`"github.com/clipperhouse/gen/typewriter"`,
		},
		Main: true,
	}
	tmpl.Execute(main, p)

	var out bytes.Buffer

	cmd := exec.Command("go", "run", main.Name(), imports.Name())
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err = cmd.Run(); err != nil {
		fmt.Println(err)
	}

	s := strings.TrimRight(out.String(), `
`)

	if len(s) > 0 {
		fmt.Println(s)
	}

	if strings.Contains(out.String(), "cannot find package") {
		fmt.Println("try running `go get` for individual imports in _gen.go")
	}
}
