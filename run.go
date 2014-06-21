package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
)

func run() error {
	if src, err := os.Open(customFilename); err == nil {
		// custom imports file exists, use it
		defer src.Close()
		if err := runCustom(src); err != nil {
			return err
		}
	} else {
		// do it the regular way
		if err := runStandard(); err != nil {
			return err
		}
	}
	return nil
}

func runStandard() error {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	if err := app.WriteAll(); err != nil {
		return err
	}

	return nil
}

func runCustom(src *os.File) error {
	temp, err := getTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	// set up imports file containing the custom typewriters (from _gen.go)
	imports, err := os.Create(filepath.Join(temp, "imports.go"))
	if err != nil {
		return err
	}
	defer imports.Close()

	io.Copy(imports, src)

	// set up main to be run
	main, err := os.Create(filepath.Join(temp, "main.go"))
	if err != nil {
		return err
	}
	defer main.Close()

	p := pkg{
		Name: "main",
		Imports: []string{
			`"github.com/clipperhouse/gen/typewriter"`,
		},
		Main: true,
	}

	if err := tmpl.Execute(main, p); err != nil {
		return err
	}

	var out, outerr bytes.Buffer

	cmd := exec.Command("go", "run", main.Name(), imports.Name())
	cmd.Stdout = &out
	cmd.Stderr = &outerr

	if err = cmd.Run(); err != nil {
		return err
	}

	if outerr.Len() > 0 {
		return errors.New(outerr.String())
	}

	s := strings.TrimRight(out.String(), `
`)

	if len(s) > 0 {
		fmt.Println(s)
	}

	return nil
}
