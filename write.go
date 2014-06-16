package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// set up temp directory under current directory
// make sure to defer os.RemoveAll() in caller
func getTempDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

func write(w io.Writer, p pkg) {
	tmpl.Execute(w, p)
}

var tmpl = template.Must(template.New("package").Parse(`package {{.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	{{.}}{{end}}
)
{{end}}
{{if .Main}}
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		panic(err)
	}

	if err := app.WriteAll(); err != nil {
		panic(err)
	}
}
{{end}}
`))
