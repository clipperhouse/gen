package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/clipperhouse/typewriter"
)

// execute runs a gen command by first determining whether a custom imports file (typically _gen.go) exists
//
// If no custom file exists, it executes the passed 'standard' func.
//
// If the custom file exists, new files are written to a temp directory and executed via `go run` in the shell.
func execute(standard func(c config) error, c config, imports typewriter.ImportSpecSet, body *template.Template) error {
	if importsSrc, err := os.Open(c.customName); err == nil {
		defer importsSrc.Close()

		// custom imports file exists, use it
		return executeCustom(importsSrc, c, imports, body)
	}

	// do it the regular way
	return standard(c)
}

// executeCustom creates a temp directory, copies importsSrc into it and generates a main() using the passed imports and body.
//
// `go run` is then called on those files via os.Command.
func executeCustom(importsSrc io.Reader, c config, imports typewriter.ImportSpecSet, body *template.Template) error {
	temp, err := getTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	// set up imports file containing the custom typewriters (from _gen.go)
	importsDst, err := os.Create(filepath.Join(temp, "imports.go"))
	if err != nil {
		return err
	}
	defer importsDst.Close()

	io.Copy(importsDst, importsSrc)

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
	if err := body.Execute(main, c); err != nil {
		return err
	}

	// call `go run` on these files & send back output/err
	cmd := exec.Command("go", "run", main.Name(), importsDst.Name())
	cmd.Stdout = c.out
	cmd.Stderr = c.out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// set up temp directory under current directory
// make sure to defer os.RemoveAll() in caller
func getTempDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

var tmpl = template.Must(template.New("package").Parse(`package {{.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports.ToSlice}}
	{{.Name}} "{{.Path}}"{{end}}
)
{{end}}`))
