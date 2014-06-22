package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// execute runs a gen command by first determining whether a custom imports file (typically _gen.go) exists
//
// If no custom file exists, it executes the passed 'standard' func.
//
// If the custom file exists, new files are written to a temp directory and executed via `go run` in the shell.
func execute(standard func() error, customFilename string, imports []string, body string) error {
	if src, err := os.Open(customFilename); err == nil {
		defer src.Close()

		// custom imports file exists, use it
		return executeCustom(src, imports, body)
	} else {
		// do it the regular way
		return standard()
	}
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

	var stdout, stderr bytes.Buffer

	// call `go run` on these files & send back output/err
	cmd := exec.Command("go", "run", main.Name(), imps.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmderr := cmd.Run()

	if s := trimBuffer(stdout); len(s) > 0 {
		fmt.Println(s)
	}

	if s := trimBuffer(stderr); len(s) > 0 {
		return errors.New(s)
	}

	if cmderr != nil {
		return cmderr
	}

	return nil
}

func trimBuffer(b bytes.Buffer) string {
	return strings.Trim(b.String(), "\n")
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
import ({{range .Imports}}
	{{.}}{{end}}
)
{{end}}`))
