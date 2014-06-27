package main

import (
	"os"
	"text/template"
)

func help() error {
	info := helpInfo{
		Name:       os.Args[0],
		CustomName: customName,
	}

	helpTmpl.Execute(out, info)
	return nil
}

type helpInfo struct {
	Name, CustomName string
}

var helpTmpl = template.Must(template.New("help").Parse(`
Usage:
  {{.Name}}           Generate files for types marked with +gen
  {{.Name}} get       Download and install typewriters (standard or custom)
  {{.Name}} list      List available typewriters (standard or custom)
  {{.Name}} custom    Create a standard {{.CustomName}} file for importing custom typewriters
  {{.Name}} help      Print usage

Further details are available at http://clipperhouse.github.io/gen
`))
