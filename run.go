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

// execute runs a gen command by first determining whether a custom imports file (typically _gen.go) exists
//
// If no custom file exists, it executes the passed 'standard' func.
//
// If the custom file exists, new files are written to a temp directory and executed via `go run` in the shell.
func execute(standard func() error, customFilename string, imports []string, body string) error {
	if src, err := os.Open(customFilename); err == nil {
		// custom imports file exists, use it
		defer src.Close()

		if err := executeCustom(src, imports, body); err != nil {
			return err
		}
	} else {
		// do it the regular way
		if err := standard(); err != nil {
			return err
		}
	}
	return nil
}

// executeCustom creates a temp directory, copies src into it and generates a main() using the passed imports and body.
//
// `go run` is then called on those files via os.Command.
func executeCustom(src io.Reader, imports []string, body string) error {
	temp, err := getTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	// set up imports file containing the custom typewriters (from _gen.go)
	imps, err := os.Create(filepath.Join(temp, "imports.go"))
	if err != nil {
		return err
	}
	defer imps.Close()

	io.Copy(imps, src)

	// set up main to be run
	main, err := os.Create(filepath.Join(temp, "main.go"))
	if err != nil {
		return err
	}
	defer main.Close()

	p := pkg{
		Name:    "main",
		Imports: imports,
	}

	// execute the package declaration and imports
	if err := tmpl.Execute(main, p); err != nil {
		return err
	}

	// write the body, usually a main()
	if _, err := main.WriteString(body); err != nil {
		return err
	}

	// call `go run` on these files & send back output/err
	var out, outerr bytes.Buffer

	cmd := exec.Command("go", "run", main.Name(), imps.Name())
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

func run(customFilename string) error {
	imports := []string{
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(runStandard, customFilename, imports, runBody)
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

func list(customFilename string) error {
	imports := []string{
		`"fmt"`,
		`"log"`,
		`"github.com/clipperhouse/gen/typewriter"`,
	}

	return execute(listStandard, customFilename, imports, listBody)
}

func listStandard() error {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		return err
	}

	fmt.Println("Installed typewriters:")
	for _, tw := range app.TypeWriters {
		fmt.Println("  " + tw.Name())
	}

	return nil
}
