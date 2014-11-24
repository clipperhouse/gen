package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func help(c config) error {
	cmd := filepath.Base(os.Args[0])
	spacer := strings.Repeat(" ", len(cmd))

	info := helpInfo{
		Name:       cmd,
		CustomName: c.customName,
		Spacer:     spacer,
	}

	if err := helpTmpl.Execute(c.out, info); err != nil {
		return err
	}

	return nil
}

type helpInfo struct {
	Name, CustomName, Spacer string
}

var helpTmpl = template.Must(template.New("help").Parse(`
Usage:
  {{.Name}}           Generate files for types marked with +gen.
  {{.Name}} get       Download and install typewriters (standard or custom). 
  {{.Spacer}}           Optional flags from go get: [-d] [-fix] [-t] [-u].
  {{.Name}} list      List available typewriters (standard or custom).
  {{.Name}} add       Add a third-party typewriter to the current package.
  {{.Name}} help      Print usage.

Further details are available at http://clipperhouse.github.io/gen

`))
